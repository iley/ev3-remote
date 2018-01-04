package main

import (
	"os"
	"time"
)

const SpeakerPath = "/sys/devices/platform/snd-legoev3"
const MaxVolume = 100

type Speaker struct {
	VolumeFile *os.File
	ToneFile   *os.File
}

func MakeSpeaker() *Speaker {
	var speaker Speaker
	speaker.VolumeFile = TryOpen(SpeakerPath + "/volume")
	speaker.ToneFile = TryOpen(SpeakerPath + "/tone")
	return &speaker
}

func (speaker *Speaker) Close() {
	if speaker.VolumeFile != nil {
		speaker.VolumeFile.Close()
	}
	if speaker.ToneFile != nil {
		speaker.ToneFile.Close()
	}
}

func (speaker *Speaker) SetVolume(volume int32) {
	if volume < 0 {
		volume = 0
	} else if volume > MaxVolume {
		volume = MaxVolume
	}
	WriteInt(speaker.VolumeFile, volume)
}

func (speaker *Speaker) Play(frequency int32, d time.Duration) {
	WriteInt(speaker.ToneFile, frequency)
	time.Sleep(d)
	WriteInt(speaker.ToneFile, 0)
}
