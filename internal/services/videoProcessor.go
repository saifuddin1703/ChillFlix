package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type VideoProcessorService struct {
}

func NewVideoProcessorService() *VideoProcessorService {
	return &VideoProcessorService{}
}

func (v *VideoProcessorService) ProcessVideo(videoPath string) error {
	base := filepath.Base(videoPath)
	videoID := base[:len(base)-len(filepath.Ext(base))]
	outputDir := filepath.Join("hls", videoID)

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Step 1: Detect if audio exists using ffprobe
	hasAudio := true
	ffprobeCmd := exec.Command("ffprobe", "-i", videoPath, "-show_streams", "-select_streams", "a", "-loglevel", "error")
	if err := ffprobeCmd.Run(); err != nil {
		hasAudio = false
	}

	// Step 2: Build FFmpeg command dynamically
	args := []string{
		"-i", videoPath,
		"-filter_complex",
		`[0:v]split=4[v1][v2][v3][v4];` +
			`[v1]scale=w=1920:h=1080[v1out];` +
			`[v2]scale=w=1280:h=720[v2out];` +
			`[v3]scale=w=854:h=480[v3out];` +
			`[v4]scale=w=640:h=360[v4out]`,
		"-map", "[v1out]",
		"-map", "[v2out]",
		"-map", "[v3out]",
		"-map", "[v4out]",
		"-c:v:0", "libx264", "-b:v:0", "5000k",
		"-c:v:1", "libx264", "-b:v:1", "3000k",
		"-c:v:2", "libx264", "-b:v:2", "1500k",
		"-c:v:3", "libx264", "-b:v:3", "800k",
		"-preset", "veryfast",
		"-g", "48", "-sc_threshold", "0",
		"-master_pl_name", "master.m3u8",
		"-f", "hls",
		"-hls_time", "6",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(outputDir, "v%v", "segment%d.ts"),
		filepath.Join(outputDir, "v%v", "playlist.m3u8"),
	}

	if hasAudio {
		// Map audio once and reuse it across all
		args = append(args,
			"-map", "0:a:0",
			"-map", "0:a:0",
			"-map", "0:a:0",
			"-map", "0:a:0",
			"-c:a", "aac", "-ar", "48000", "-b:a", "128k",
			"-var_stream_map", "v:0,a:0 v:1,a:1 v:2,a:2 v:3,a:3",
		)
	} else {
		args = append(args,
			"-var_stream_map", "v:0 v:1 v:2 v:3",
		)
	}

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Processing video with FFmpeg...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg processing failed: %w", err)
	}

	fmt.Println("Video successfully processed at:", outputDir)
	return nil
}
