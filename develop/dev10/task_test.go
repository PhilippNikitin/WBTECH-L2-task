package main

import (
	"bytes"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Тестовый сервер для обработки соединений
func startTestServer(t *testing.T, port string, response string) net.Listener {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				io.Copy(c, bytes.NewBufferString(response))
			}(conn)
		}
	}()

	return ln
}

// Компилируем `main.go` в исполняемый файл перед тестами
func buildClient(t *testing.T) string {
	tmpDir := os.TempDir()
	output := filepath.Join(tmpDir, "telnet_client")
	cmd := exec.Command("go", "build", "-o", output, "task.go")

	// Захватываем стандартный вывод и ошибки команды
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build client: %v, stderr: %v", err, stderr.String())
	}
	return output
}

// Тест подключения клиента к серверу и обмена данными
// Тест подключения клиента к серверу и обмена данными
func TestTelnetClient(t *testing.T) {
	port := "12345"
	response := "Hello, World!\n"
	ln := startTestServer(t, port, response)
	defer ln.Close()

	clientPath := buildClient(t)

	// Аргументы командной строки для теста
	args := []string{"--timeout=2s", "--host=localhost", "--port=" + port}

	cmd := exec.Command(clientPath, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to get stdin pipe: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

	// Закрыть stdin, чтобы завершить работу клиента
	stdin.Close()

	// Ожидание завершения работы клиента
	err = cmd.Wait()
	if err != nil {
		t.Fatalf("Command finished with error: %v", err)
	}

	// Проверка ответа от сервера
	if out.String() != response {
		t.Fatalf("Expected response %q but got %q", response, out.String())
	}
}

// Тест подключения к несуществующему серверу
func TestTelnetClientTimeout(t *testing.T) {
	port := "12346"

	clientPath := buildClient(t)

	cmd := exec.Command(clientPath, "--timeout=2s", "localhost", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

	// Ожидание завершения работы клиента
	err = cmd.Wait()
	if err == nil {
		t.Fatalf("Expected command to fail due to timeout, but it succeeded")
	}
}
