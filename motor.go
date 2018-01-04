package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Motor struct {
	Port                string
	CommandFile         *os.File
	SpeedFile           *os.File
	SpeedRegulationFile *os.File
}

func MakeMotor(port string) *Motor {
	var motor Motor
	motor.Port = port
	path, err := FindMotor(port)
	if err != nil {
		log.Printf("Didn't find motor %s", port)
		motor.CommandFile = nil
		motor.SpeedFile = nil
		motor.SpeedRegulationFile = nil
	} else {
		log.Printf("Found motor %s: %s", port, path)
		motor.CommandFile = OpenOrDie(path + "/command")
		motor.SpeedFile = OpenOrDie(path + "/speed_sp")
		motor.SpeedRegulationFile = OpenOrDie(path + "/speed_regulation")
	}
	motor.Reset()
	return &motor
}

func FindMotor(port string) (string, error) {
	motorsDir := "/sys/class/tacho-motor"
	subdirs, err := ioutil.ReadDir(motorsDir)
	if err != nil {
		return "", err
	}
	for _, subdir := range subdirs {
		path := motorsDir + "/" + subdir.Name()
		portName, err := ioutil.ReadFile(path + "/port_name")
		if err != nil {
			log.Printf("Error reading port file: %v", err)
			continue
		}
		if port == strings.TrimRight(string(portName), " \t\n\r") {
			return path, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Motor not found: %s", port))
}

func (motor *Motor) Close() {
	motor.Reset()
	if motor.SpeedFile != nil {
		motor.SpeedFile.Close()
	}
	if motor.CommandFile != nil {
		motor.CommandFile.Close()
	}
	if motor.SpeedRegulationFile != nil {
		motor.SpeedRegulationFile.Close()
	}
}

func (motor *Motor) Reset() {
	WriteCommand(motor.CommandFile, "reset")
	WriteCommand(motor.SpeedRegulationFile, "on")
}

func (motor *Motor) Stop() {
	WriteCommand(motor.CommandFile, "stop")
}

// 0..1000
func (motor *Motor) SetSpeed(speed int32) {
	if Verbose() {
		log.Printf("Motor %s speed %d", motor.Port, speed)
	}
	if speed == 0 {
		motor.Stop()
		return
	}
	WriteInt(motor.SpeedFile, speed)
	WriteCommand(motor.CommandFile, "run-forever")
}
