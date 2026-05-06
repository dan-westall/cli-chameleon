package executor

import (
	"bytes"
	"io"
	"os/exec"
	"sync"
)

type Result struct {
	Output string
	Err    error
}

func Run(commands []string, parallel bool) Result {
	if parallel {
		return runParallel(commands)
	}
	return runSequential(commands)
}

func runSequential(commands []string) Result {
	var output bytes.Buffer
	for _, cmd := range commands {
		c := exec.Command("sh", "-c", cmd)
		c.Stdout = &output
		c.Stderr = &output
		if err := c.Run(); err != nil {
			return Result{Output: output.String(), Err: err}
		}
	}
	return Result{Output: output.String()}
}

func runParallel(commands []string) Result {
	var (
		mu     sync.Mutex
		output bytes.Buffer
		wg     sync.WaitGroup
		errs   []error
	)

	for _, cmd := range commands {
		wg.Add(1)
		go func(cmd string) {
			defer wg.Done()
			c := exec.Command("sh", "-c", cmd)
			var buf bytes.Buffer
			c.Stdout = &buf
			c.Stderr = &buf
			err := c.Run()
			mu.Lock()
			output.Write(buf.Bytes())
			if err != nil {
				errs = append(errs, err)
			}
			mu.Unlock()
		}(cmd)
	}

	wg.Wait()
	if len(errs) > 0 {
		return Result{Output: output.String(), Err: errs[0]}
	}
	return Result{Output: output.String()}
}

// StreamCmd starts a command and returns a reader for its combined output.
func StreamCmd(commands []string) (*exec.Cmd, io.ReadCloser, error) {
	// For streaming, join commands with &&
	joined := ""
	for i, c := range commands {
		if i > 0 {
			joined += " && "
		}
		joined += c
	}

	cmd := exec.Command("sh", "-c", joined)
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw

	if err := cmd.Start(); err != nil {
		pw.Close()
		pr.Close()
		return nil, nil, err
	}

	// Close writer when process exits so reader gets EOF.
	go func() {
		cmd.Wait()
		pw.Close()
	}()

	return cmd, pr, nil
}
