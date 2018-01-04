package main

type Robot interface {
	Stop()
	Shutdown()
	HandleButton(id int, pressed bool)
	HandleAxisChange(id int, value float64)
}
