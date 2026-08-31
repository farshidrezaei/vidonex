package spec

import (
	"path/filepath"
	"strings"
)

// ResolveAssetPaths resolves all relative media and subtitle file paths within the spec
// against the provided base directory (e.g., the directory where the spec file is located).
func ResolveAssetPaths(specification *VideoSpec, baseDirectory string) {
	if specification == nil || baseDirectory == "" {
		return
	}

	// Resolve Output path if relative
	if specification.Output != "" && !filepath.IsAbs(specification.Output) {
		specification.Output = filepath.Join(baseDirectory, specification.Output)
	}

	for trackIndex := range specification.Tracks {
		track := &specification.Tracks[trackIndex]

		// Resolve Subtitle file path if present
		if track.Subtitles != nil && track.Subtitles.File != "" && !filepath.IsAbs(track.Subtitles.File) {
			track.Subtitles.File = filepath.Join(baseDirectory, track.Subtitles.File)
		}

		// Resolve Clip sources
		for clipIndex := range track.Clips {
			clip := &track.Clips[clipIndex]
			if clip.Source != "" && !filepath.IsAbs(clip.Source) && !isURLLike(clip.Source) {
				clip.Source = filepath.Join(baseDirectory, clip.Source)
			}
		}
	}
}

func isURLLike(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "rtmp://") ||
		strings.HasPrefix(path, "rtsp://")
}
