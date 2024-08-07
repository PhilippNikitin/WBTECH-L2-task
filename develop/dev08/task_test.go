package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestCd(t *testing.T) {
	initialDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get initial directory: %s", err)
	}

	tempDir := os.TempDir()
	cd([]string{tempDir})
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %s", err)
	}

	if currentDir != tempDir {
		t.Errorf("Expected directory %s, but got %s", tempDir, currentDir)
	}

	// Change back to the initial directory
	cd([]string{initialDir})
}

func TestPwd(t *testing.T) {
	initialDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get initial directory: %s", err)
	}

	output := captureOutput(func() {
		pwd()
	})

	if strings.TrimSpace(output) != initialDir {
		t.Errorf("Expected directory %s, but got %s", initialDir, strings.TrimSpace(output))
	}
}

func TestEcho(t *testing.T) {
	output := captureOutput(func() {
		echo([]string{"Hello", "World"})
	})

	expected := "Hello World"
	if output != expected {
		t.Errorf("Expected output %s, but got %s", expected, output)
	}
}

func TestKill(t *testing.T) {
	// Start a long-running process
	cmd := exec.Command("sleep", "100")
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start sleep command: %s", err)
	}

	// Kill the process
	pid := cmd.Process.Pid
	kill([]string{strconv.Itoa(pid)})

	// Check if the process is killed
	err = cmd.Wait()
	if err == nil || !strings.Contains(err.Error(), "signal: killed") {
		t.Errorf("Expected process to be killed, but got error: %s", err)
	}
}

func TestPs(t *testing.T) {
	output := captureOutput(func() {
		ps()
	})

	// Check for PID and CMD headers, which are commonly used
	if !strings.Contains(output, "PID") || !strings.Contains(output, "CMD") {
		t.Errorf("Expected ps output to contain headers, but got %s", output)
	}
}

func TestPipeline(t *testing.T) {
	output := captureOutput(func() {
		executePipeline([][]string{
			{"echo", "Hello"},
			{"awk", "{print $1 \" World\"}"},
		})
	})

	expected := "Hello World"
	if output != expected {
		t.Errorf("Expected output %s, but got %s", expected, output)
	}
}

// Helper function to capture output
// Function to capture the output of a function
func captureOutput(f func()) string {
	var buf bytes.Buffer

	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	stderr := os.Stderr

	// Redirect stdout and stderr to the pipe
	os.Stdout = w
	os.Stderr = w

	// Read output in a separate goroutine
	done := make(chan struct{})
	go func() {
		buf.ReadFrom(r)
		close(done)
	}()

	// Execute the function
	f()

	// Close the pipe and restore stdout/stderr
	w.Close()
	os.Stdout = stdout
	os.Stderr = stderr
	<-done // Ensure the goroutine has finished

	return strings.TrimSpace(buf.String())
}
