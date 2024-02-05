package foundryutils

import (
	"fmt"
	"math/big"
	"os/exec"
	"sync"
)

var (
	forkProcess *exec.Cmd
	mu          sync.Mutex
)

type ForkOpts struct {
	Port        string
	BlockNumber *big.Int
}

// start fork a network for testing purposes
func StartFork(forkUrl string, opt ForkOpts) (string, error) {
	_, err := exec.LookPath("anvil")
	if err != nil {
		return "", err
	}
	if forkProcess != nil {
		return "", fmt.Errorf("fork process already exists")
	}
	if opt.Port == "" {
		opt.Port = "8545"
	}
	args := []string{
		"--fork-url", forkUrl,
		"--port", opt.Port,
	}
	if opt.BlockNumber != nil {
		args = append(args, []string{
			"--fork-block-number", opt.BlockNumber.String(),
		}...)
	}
	mu.Lock()
	forkProcess = exec.Command("anvil", args...) // #nosec
	defer mu.Unlock()
	if err := forkProcess.Start(); err != nil {
		return "", err
	}
	return fmt.Sprintf("http://localhost:%s", opt.Port), nil
}

// stop current fork process
func StopFork() error {
	mu.Lock()
	defer mu.Unlock()
	if forkProcess != nil {
		if err := forkProcess.Process.Kill(); err != nil {
			return fmt.Errorf("kill old fork process failed: %v", err.Error())
		}
		if err := forkProcess.Wait(); err != nil {
			return fmt.Errorf("stop old fork process failed: %v", err.Error())
		}
	}
	return nil
}
