package anvilutils

import (
	"fmt"
	"io"
	"math/big"
	"os/exec"
	"strings"
)

type ForkCmd struct {
	cmd *exec.Cmd

	// listening port
	port string
}

type ForkOpts struct {
	ForkUrl     string
	Port        string
	BlockNumber *big.Int
}

func NewForkCommand(opt ForkOpts) (*ForkCmd, error) {
	_, err := exec.LookPath("anvil")
	if err != nil {
		return nil, err
	}

	if opt.Port == "" {
		opt.Port = "8545"
	}
	args := []string{
		"--port", opt.Port,
	}
	if opt.ForkUrl != "" {
		args = append(args, []string{
			"--fork-url", opt.ForkUrl,
		}...)
	}
	if opt.BlockNumber != nil {
		args = append(args, []string{
			"--fork-block-number", opt.BlockNumber.String(),
		}...)
	}
	return &ForkCmd{
		cmd:  exec.Command("anvil", args...), // #nosec
		port: opt.Port,
	}, nil
}

// start fork a network for testing purposes
func (f *ForkCmd) Start() (string, error) {
	out, err := f.cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := f.cmd.Start(); err != nil {
		return "", err
	}
	if err := waitForAnvilReady(out, f.port); err != nil {
		return "", err
	}
	return fmt.Sprintf("http://localhost:%s", f.port), nil
}

// stop current fork process
func (f *ForkCmd) Stop() error {
	if err := f.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("kill old fork process failed: %v", err.Error())
	}
	if err := f.cmd.Wait(); err != nil {
		return fmt.Errorf("stop old fork process failed: %v", err.Error())
	}
	return nil
}

func waitForAnvilReady(out io.ReadCloser, port string) error {
	readyString := "Listening on 127.0.0.1:" + port
	for {
		buff := make([]byte, 1024)
		if _, err := out.Read(buff); err != nil {
			return err
		}
		if strings.Contains(string(buff), readyString) {
			return nil
		}
	}
}
