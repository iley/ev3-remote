package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func WriteCommand(file *os.File, command string) {
	if file != nil {
		_, err := file.Write([]byte(command + "\n"))
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Printf("Simulation: command %s", command)
	}
}

func WriteInt(file *os.File, value int32) {
	if file != nil {
		stringValue := strconv.Itoa(int(value)) + "\n"
		_, err := file.Write([]byte(stringValue))
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Printf("Simulation: write int %d", value)
	}
}

func TryOpen(path string) *os.File {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0775)
	if err != nil {
		log.Printf(err.Error())
		return nil
	}
	log.Printf("Opened %s", path)
	return file
}

func OpenOrDie(path string) *os.File {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0775)
	if err != nil {
		panic(err)
	}
	log.Printf("Opened %s", path)
	return file
}

func TryWrite(path string, value string) {
	err := ioutil.WriteFile(path, []byte(value), 0775)
	if err != nil {
		log.Printf(err.Error())
	}
}

func WriteOrDie(path string, value string) {
	err := ioutil.WriteFile(path, []byte(value), 0775)
	if err != nil {
		panic(err)
	}
}

func PlaySoundFile(path string) {
	exec.Command("aplay", path).Run()
}

func MaxInt32(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func MinInt32(x, y int32) int32 {
	if x > y {
		return y
	}
	return x
}
