package main

type Device interface {
	IsPoweredOn() (bool, error)
	GetPairedDevices() ([]string, error)
	PowerOn() error
	PromptAvailable() string
	ConnectTo(string) error
}

type blueman struct {
	cliBin string
	opts   Options
}

type Options struct {
	listPaired  string
	powerStatus string
	connectTo   string
}
