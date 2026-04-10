package output

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"twitch-redeem-trigger/src/config"
	"twitch-redeem-trigger/src/logger"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type SoundDevice struct {
	config config.Media
	logger *logger.Logger
}

type MediaType int

const (
	MediaTypeUnknown MediaType = 0
	MediaTypeAudio   MediaType = 1
	MediaTypeVideo   MediaType = 2
)

var (
	speakerInitialized = false
	speakerMutex       sync.Mutex
	defaultSampleRate  = beep.SampleRate(44100)
)

func BuildMedia(cfgMedia config.Media, l *logger.Logger) Device {
	l.Info("test")
	err := testSound(l)
	if err != nil {
		l.Error("Failed to initialize speaker: %v", err)
	}

	return SoundDevice{
		config: cfgMedia,
		logger: l,
	}
}

func (d SoundDevice) Toggle(input ToggleInput) error {
	if d.config.Port == 0 {
		return nil
	}

	d.logger.Info("Media: Toggle")

	for redeem, sound := range d.config.Mappings {
		if redeem == input.RedeemName {
			d.logger.Info("User: %s", input.User)
			d.logger.Info("RedeemName: %s", redeem)
			d.logger.Info("Media: %s", sound)

			return playMedia(sound)
		}
	}

	return nil
}

func testSound(l *logger.Logger) error {
	l.Debug("Testing speaker...")

	speakerMutex.Lock()
	defer speakerMutex.Unlock()

	if speakerInitialized {
		return nil
	}

	// Initialisiere mit Standard-Samplerate (44.1kHz)
	err := speaker.Init(defaultSampleRate, defaultSampleRate.N(time.Second/10))
	if err != nil {
		return fmt.Errorf("failed to initialize speaker: %w", err)
	}

	speakerInitialized = true
	return nil
}

func playMedia(file string) error {
	mediaType := detectMediaType(file)

	switch mediaType {
	case MediaTypeAudio:
		log.Printf("Playing audio: %s", file)
		return playSound(file)

	case MediaTypeVideo:
		log.Printf("Triggering video: %s", file)
		return playVideo(file) // write the video file into a temp file for webserver

	default:
		return fmt.Errorf("unknown media type for file: %s", file)
	}
}

func detectMediaType(filePath string) MediaType {
	ext := strings.ToLower(filepath.Ext(filePath))

	audioExts := []string{".mp3", ".wav", ".ogg", ".flac"}
	videoExts := []string{".mp4", ".webm", ".mov", ".avi", ".mkv"}

	for _, audioExt := range audioExts {
		if ext == audioExt {
			return MediaTypeAudio
		}
	}

	for _, videoExt := range videoExts {
		if ext == videoExt {
			return MediaTypeVideo
		}
	}

	return MediaTypeUnknown
}

func playSound(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open sound file: %w", err)
	}
	defer f.Close()

	var streamer beep.StreamSeekCloser
	var format beep.Format

	// Unterstütze verschiedene Formate
	ext := filepath.Ext(file)
	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	default:
		return fmt.Errorf("unsupported audio format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to decode audio: %w", err)
	}
	defer streamer.Close()

	// Initialisiere Speaker beim ersten Aufruf
	if !speakerInitialized {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		speakerInitialized = true
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}

func playVideo(videoPath string) error {
	fileName := "current_media.txt"
	content := []byte(videoPath)

	err := os.WriteFile(fileName, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
