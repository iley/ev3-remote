package main

import (
	"time"
)

type Beeper struct {
	Speaker *Speaker
	Channel chan bool
}

func MakeBeeper(speaker *Speaker) *Beeper {
	var beeper Beeper
	beeper.Speaker = speaker
	beeper.Channel = make(chan bool, 1)
	go beeper.Loop()
	return &beeper
}

func (beeper *Beeper) Beep() {
	select {
	case beeper.Channel <- true:
	default:
	}
}

func (beeper *Beeper) Loop() {
	for {
		<-beeper.Channel
		beeper.Speaker.Play(100, 50*time.Millisecond)
		time.Sleep(time.Second)
	}
}
