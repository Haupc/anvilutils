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

// start fork a network for testing purposes
func StartFork(forkUrl string, port string, blockNumber *big.Int) (string, error) {
	_, err := exec.LookPath("anvil")
	if err != nil {
		return "", err
	}
	if port == "" {
		port = "8545"
	}
	mu.Lock()
	if forkProcess != nil {
		forkProcess.Process.Kill()
		if err := forkProcess.Wait(); err != nil {
			return "", fmt.Errorf("stop old fork process failed: %v", err.Error())
		}
	}
	args := []string{
		"--fork-url", forkUrl,
		"--port", port,
	}
	if blockNumber != nil {
		args = append(args, []string{
			"--fork-block-number", blockNumber.String(),
		}...)
	}
	forkProcess = exec.Command("anvil", args...)
	defer mu.Unlock()
	if err := forkProcess.Start(); err != nil {
		return "", err
	}
	return fmt.Sprintf("http://localhost:%s", port), nil
}

// stop current fork process
func StopFork() error {
	mu.Lock()
	defer mu.Unlock()
	if forkProcess != nil {
		forkProcess.Process.Kill()
		if err := forkProcess.Wait(); err != nil {
			return fmt.Errorf("stop old fork process failed: %v", err.Error())
		}
	}
	return nil
}
