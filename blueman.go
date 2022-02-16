package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/manifoldco/promptui"
)

var (
	cliBin = "blueutil"
)

func (b *blueman) GetPairedDevices() ([]string, error) {
	res, err := runCmd(b.cliBin, b.opts.listPaired)

	fmt.Printf("\n\nAVAILABLE DEVICES\n %s\n\n", strings.Repeat("-", 100))
	fmt.Println(res)

	paired := []string{}
	for _, line := range strings.Split(res, "\n") {
		if len(line) == 0 {
			continue
		}
		paired = append(paired, strings.ReplaceAll(strings.Split(line, " ")[1], ",", ""))
	}

	return paired, err

}

func NewDevice() Device {
	if runtime.GOOS != "darwin" {
		must("ONLY MAC is supported as of now", fmt.Errorf("unsupported operating system"))
	}
	_, err := exec.LookPath(cliBin)
	must("blueutil is not installed. `brew install blueutil` to get started", err)
	flaggen := func(arg string) string {
		return "--" + arg
	}
	return &blueman{
		cliBin: cliBin,
		opts: Options{
			listPaired:  flaggen("paired"),
			powerStatus: flaggen("power"),
			connectTo:   flaggen("connect"),
		},
	}
}

func (b *blueman) IsPoweredOn() (bool, error) {
	val, err := runCmd(b.cliBin, b.opts.powerStatus)
	return strings.TrimSpace(val) == "1", err
}

func (b *blueman) PowerOn() error {
	_, err := runCmd(b.cliBin, b.opts.powerStatus, "1")
	return err
}

func (b *blueman) PromptAvailable() string {
	availableDevs, err := b.GetPairedDevices()
	must("Listing available devices", err)
	if len(availableDevs) == 0 {
		log.Println("There are NO devices available to connect to. Exiting...")
		os.Exit(0)
	}
	prompt := promptui.Select{
		Label: "Select Device to connect to",
		Items: availableDevs,
	}

	_, result, err := prompt.Run()
	must("Prompt failed %v\n", err)

	return result
}

func (b *blueman) ConnectTo(connTo string) error {
	_, err := runCmd(b.cliBin, b.opts.connectTo, connTo)
	return err
}

func runCmd(cmds ...string) (string, error) {
	bin, opts := cmds[0], cmds[1:]
	cmd := exec.Command(bin, opts...)
	stdout, err := cmd.Output()
	execCmd := fmt.Sprintf("Running OS command: %s", cmds)
	must(execCmd, err)
	return strings.Trim(string(stdout), " "), err
}

func must(msg string, err error) {
	if err != nil {
		log.Println("FATAL::", msg, ":", err.Error())
		os.Exit(-1)
	}
}
