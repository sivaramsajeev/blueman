package main

import "log"

func main() {
	device := NewDevice()
	isDevOn, err := device.IsPoweredOn()
	must("Fetching power status details", err)
	if isDevOn {
		log.Println("Bluetooth service is ON. Continuing...")
	} else {
		log.Println("Bluetooth service is DOWN. Trying to power-on ...")
		err := device.PowerOn()
		must("Powering ON", err)
	}

	connTo := device.PromptAvailable()
	err = device.ConnectTo(connTo)
	must("Connecting to device", err)
}
