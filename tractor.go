package main

import (
	"log"
	"time"
)

const ChannelSize = 64
const MaxSpeed int32 = 1000
const (
	OBSTACLE_NO   = iota
	OBSTACLE_FAR  = iota
	OBSTACLE_NEAR = iota
)

// Implements Robot interface.
type Tractor struct {
	AutopilotEnabled bool
	Control          chan func()
	QuitAutopilot    chan bool
	LeftMotor        *Motor
	RightMotor       *Motor
	Speaker          *Speaker
	Beeper           *Beeper
	IrSensor         *IrSensor
}

func MakeTractor() (*Tractor, error) {
	var tractor Tractor
	tractor.AutopilotEnabled = false
	tractor.Control = make(chan func(), ChannelSize)
	tractor.QuitAutopilot = make(chan bool)
	tractor.LeftMotor = MakeMotor("outB")
	tractor.RightMotor = MakeMotor("outC")
	tractor.Speaker = MakeSpeaker()
	tractor.Beeper = MakeBeeper(tractor.Speaker)
	tractor.IrSensor = MakeIrSensor("0")
	go tractor.ControlLoop()
	tractor.Control <- func() {
		tractor.Beeper.Beep()
	}
	return &tractor, nil
}

func (tractor *Tractor) Stop() {
	tractor.LeftMotor.SetSpeed(0)
	tractor.RightMotor.SetSpeed(0)
}

func (tractor *Tractor) Shutdown() {
	tractor.Stop()
	tractor.LeftMotor.Close()
	tractor.RightMotor.Close()
	tractor.Speaker.Close()
}

func (tractor *Tractor) HandleButton(id int, pressed bool) {
	var status = "pressed"
	if !pressed {
		status = "released"
	}
	if Verbose() {
		log.Printf("Button %d %s", id, status)
	}
	if pressed {
		if id == 0 { // A.
			tractor.Control <- func() {
				tractor.Beeper.Beep()
			}
		} else if id == 3 { // Y.
			tractor.Control <- func() {
				tractor.Beeper.Beep()
			}
		} else if id == 9 { // Start.
			tractor.Control <- func() {
				tractor.ToggleAutopilot()
			}
		}
	}
}

func (tractor *Tractor) HandleAxisChange(id int, value float64) {
	if Verbose() {
		log.Printf("Axis %d = %f\n", id, value)
	}
	speed := int32(-float64(MaxSpeed) * value)
	if id == 1 {
		tractor.Control <- func() {
			tractor.LeftMotor.SetSpeed(speed)
		}
	} else if id == 3 {
		tractor.Control <- func() {
			tractor.RightMotor.SetSpeed(speed)
		}
	}
}

func (tractor *Tractor) ToggleAutopilot() {
	tractor.AutopilotEnabled = !tractor.AutopilotEnabled
	if tractor.AutopilotEnabled {
		log.Printf("Autopilot enabled")
		go tractor.AutoPilotLoop()
	} else {
		log.Printf("Autopilot disabled")
		tractor.QuitAutopilot <- true
		tractor.Stop()
	}
}

func (tractor *Tractor) ControlLoop() {
	for command := range tractor.Control {
		command()
	}
}

func (tractor *Tractor) AutoPilotLoop() {
	for {
		select {
		case <-tractor.QuitAutopilot:
			return
		default:
		}
		tractor.Control <- func() {
			tractor.LeftMotor.SetSpeed(300)
		}
		tractor.Control <- func() {
			tractor.RightMotor.SetSpeed(-300)
		}
		time.Sleep(time.Second)
		tractor.Control <- func() {
			tractor.LeftMotor.SetSpeed(-300)
		}
		tractor.Control <- func() {
			tractor.RightMotor.SetSpeed(300)
		}
		time.Sleep(time.Second)
	}
}
