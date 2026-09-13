package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/farshidrezaei/vidonex/executor"
)

// ProgressBar manages interactive live terminal progress visualization.
type ProgressBar struct {
	Writer         io.Writer
	TotalDuration  time.Duration
	StartTime      time.Time
	IsTTY          bool
	Width          int
	lastPercentage int
	mutex          sync.Mutex
}

// NewProgressBar creates a new progress tracker.
func NewProgressBar(writer io.Writer, totalDuration time.Duration) *ProgressBar {
	if writer == nil {
		writer = os.Stdout
	}

	isTTY := false
	if file, ok := writer.(*os.File); ok {
		fileInfo, err := file.Stat()
		if err == nil && (fileInfo.Mode()&os.ModeCharDevice) != 0 {
			isTTY = true
		}
	}

	return &ProgressBar{
		Writer:        writer,
		TotalDuration: totalDuration,
		StartTime:     time.Now(),
		IsTTY:         isTTY,
		Width:         30,
	}
}

// Update renders the current progress state.
func (bar *ProgressBar) Update(event executor.ProgressEvent) {
	bar.mutex.Lock()
	defer bar.mutex.Unlock()

	percentage := event.Percentage
	if percentage > 100.0 {
		percentage = 100.0
	}

	elapsed := time.Since(bar.StartTime)
	var eta time.Duration
	if percentage > 1.0 {
		totalEstimated := time.Duration(float64(elapsed) / (percentage / 100.0))
		eta = totalEstimated - elapsed
		if eta < 0 {
			eta = 0
		}
	}

	if bar.IsTTY {
		filledWidth := int((percentage / 100.0) * float64(bar.Width))
		if filledWidth > bar.Width {
			filledWidth = bar.Width
		}
		emptyWidth := bar.Width - filledWidth
		if emptyWidth < 0 {
			emptyWidth = 0
		}

		barString := strings.Repeat("█", filledWidth) + strings.Repeat("░", emptyWidth)

		_, _ = fmt.Fprintf(bar.Writer, "\r\033[K  [%s] %5.1f%% | Speed: %4.2fx | FPS: %4.1f | Elapsed: %s | ETA: %s",
			barString,
			percentage,
			event.Speed,
			event.FPS,
			FormatDuration(elapsed),
			FormatDuration(eta),
		)
	} else {
		// Non-TTY (CI / Cloud logs): print milestone percentages
		currentMilestone := int(percentage / 10.0) * 10
		if currentMilestone > bar.lastPercentage && currentMilestone <= 100 {
			bar.lastPercentage = currentMilestone
			_, _ = fmt.Fprintf(bar.Writer, "Rendering: %3d%% (Speed: %.2fx, FPS: %.1f, Elapsed: %s)\n",
				currentMilestone,
				event.Speed,
				event.FPS,
				FormatDuration(elapsed),
			)
		}
	}
}

// Finish concludes progress output.
func (bar *ProgressBar) Finish() {
	bar.mutex.Lock()
	defer bar.mutex.Unlock()

	if bar.IsTTY {
		_, _ = fmt.Fprintln(bar.Writer, "")
	}
}

// FormatDuration formats a duration into MM:SS format.
func FormatDuration(duration time.Duration) string {
	duration = duration.Round(time.Second)
	totalSeconds := int(duration.Seconds())
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
