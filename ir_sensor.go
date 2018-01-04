package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const IrSensorPath = "/sys/class/lego-sensor/sensor"
const MaxProximity = 100

type IrSensor struct {
	Id        string
	ValuePath string
}

func MakeIrSensor(id string) *IrSensor {
	var sensor IrSensor
	sensor.Id = id
	sensor.ValuePath = IrSensorPath + id + "/value0"
	TryWrite(IrSensorPath+id+"/mode", "IR-PROX")
	return &sensor
}

// 0..100
func (sensor *IrSensor) Proximity() int32 {
	contents, err := ioutil.ReadFile(sensor.ValuePath)
	if err != nil {
		return MaxProximity
	}
	contentsString := strings.TrimRight(string(contents[:]), " \t\n\r")
	proximity, err := strconv.Atoi(contentsString)
	if err != nil {
		log.Printf("Error parsing IR sensor data '%s': %s", contentsString,
			err.Error())
		return MaxProximity
	}
	return int32(proximity)
}
