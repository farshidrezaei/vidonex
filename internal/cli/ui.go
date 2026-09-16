// Package cli provides private terminal UI, live progress telemetry, and structured logging helpers for the Vidonex CLI tool.
package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// ANSI color codes for rich terminal styling.
const (
	ColorReset     = "\033[0m"
	ColorBold      = "\033[1m"
	ColorDim       = "\033[2m"
	ColorCyan      = "\033[36m"
	ColorGreen     = "\033[32m"
	ColorYellow    = "\033[33m"
	ColorRed       = "\033[31m"
	ColorMagenta   = "\033[35m"
	ColorBlue      = "\033[34m"
	ColorWhiteBold = "\033[1;37m"
)

// UI manages stylized terminal output.
type UI struct {
	Writer   io.Writer
	IsTTY    bool
	NoColors bool
}

// NewUI creates a new UI manager.
func NewUI(writer io.Writer) *UI {
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
	return &UI{
		Writer:   writer,
		IsTTY:    isTTY,
		NoColors: os.Getenv("NO_COLOR") != "",
	}
}

// Colorize wraps text with ANSI color sequences if TTY and colors are enabled.
func (ui *UI) Colorize(colorCode, text string) string {
	if ui.NoColors || !ui.IsTTY {
		return text
	}
	return colorCode + text + ColorReset
}

// PrintBanner outputs the stylized Vidonex engine banner.
func (ui *UI) PrintBanner() {
	banner := `
  __   ___     __                  
  \ \ / (_)___/ /__  ___  __ _____ 
   \ V / / / _  / _ \/ _ \/ // /\ \ /
    \_/_/_/\_,_/\___/_//_/\_, //_\_\ 
                         /___/       
`
	_, _ = fmt.Fprintln(ui.Writer, ui.Colorize(ColorCyan+ColorBold, banner))
	_, _ = fmt.Fprintln(ui.Writer, ui.Colorize(ColorDim, "  Declarative Video Composition & FFmpeg Filtergraph Engine in Go"))
	_, _ = fmt.Fprintln(ui.Writer, ui.Colorize(ColorDim, "  -------------------------------------------------------------"))
}

// PrintStep prints a major pipeline step badge.
func (ui *UI) PrintStep(stepIndex, totalSteps int, symbol, description string) {
	badge := fmt.Sprintf("[%d/%d] %s", stepIndex, totalSteps, symbol)
	_, _ = fmt.Fprintf(ui.Writer, "\n%s %s\n", ui.Colorize(ColorCyan+ColorBold, badge), ui.Colorize(ColorWhiteBold, description))
}

// PrintSuccess prints a success message.
func (ui *UI) PrintSuccess(message string) {
	_, _ = fmt.Fprintf(ui.Writer, "%s %s\n", ui.Colorize(ColorGreen+ColorBold, "✓"), ui.Colorize(ColorGreen, message))
}

// PrintError prints an error message.
func (ui *UI) PrintError(err error) {
	_, _ = fmt.Fprintf(ui.Writer, "\n%s %s\n", ui.Colorize(ColorRed+ColorBold, "✖ Error:"), ui.Colorize(ColorRed, err.Error()))
}

// PrintSummaryCard prints an executive completion card.
func (ui *UI) PrintSummaryCard(outputPath string, fileSizeBytes int64, renderDuration, videoDuration time.Duration) {
	formattedSize := FormatFileSize(fileSizeBytes)
	speedRatio := 0.0
	if renderDuration > 0 {
		speedRatio = videoDuration.Seconds() / renderDuration.Seconds()
	}

	border := strings.Repeat("─", 60)
	_, _ = fmt.Fprintln(ui.Writer, "")
	_, _ = fmt.Fprintln(ui.Writer, ui.Colorize(ColorCyan, "┌"+border+"┐"))
	_, _ = fmt.Fprintf(ui.Writer, ui.Colorize(ColorCyan, "│")+"  %s\n", ui.Colorize(ColorGreen+ColorBold, "✨ Composition Rendered Successfully!"))
	_, _ = fmt.Fprintln(ui.Writer, ui.Colorize(ColorCyan, "├"+border+"┤"))
	_, _ = fmt.Fprintf(ui.Writer, ui.Colorize(ColorCyan, "│")+"  %-20s : %s\n", ui.Colorize(ColorDim, "Output File"), ui.Colorize(ColorBold, outputPath))
	_, _ = fmt.Fprintf(ui.Writer, ui.Colorize(ColorCyan, "│")+"  %-20s : %s\n", ui.Colorize(ColorDim, "File Size"), formattedSize)
	_, _ = fmt.Fprintf(ui.Writer, ui.Colorize(ColorCyan, "│")+"  %-20s : %s\n", ui.Colorize(ColorDim, "Video Duration"), videoDuration.Round(10*time.Millisecond))
	_, _ = fmt.Fprintf(ui.Writer, ui.Colorize(ColorCyan, "│")+"  %-20s : %s (Speedup: %.2fx)\n", ui.Colorize(ColorDim, "Render Elapsed"), renderDuration.Round(10*time.Millisecond), speedRatio)
	_, _ = fmt.Fprintln(ui.Writer, ui.Colorize(ColorCyan, "└"+border+"┘"))
}

// FormatFileSize formats raw bytes into human readable binary units.
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
