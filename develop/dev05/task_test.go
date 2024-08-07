package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"testing"
)

func TestGrep(t *testing.T) {
	content := `line 1
hello world
Line with hello in between
another line
last line`
	tempFile := createTempFile(content)
	defer os.Remove(tempFile.Name())

	testCases := []struct {
		desc     string
		args     []string
		expected string
	}{
		{
			desc:     "Basic case with pattern 'hello'",
			args:     []string{"grep", "-n", "hello", tempFile.Name()},
			expected: "2:hello world\n3:Line with hello in between\n",
		},
		{
			desc:     "Case insensitive search with pattern 'line'",
			args:     []string{"grep", "-i", "-n", "line", tempFile.Name()},
			expected: "1:line 1\n3:Line with hello in between\n4:another line\n5:last line\n",
		},
		{
			desc:     "Printing lines after match with pattern 'line'",
			args:     []string{"grep", "-n", "-A", "1", "line", tempFile.Name()},
			expected: "1:line 1\n2:hello world\n4:another line\n5:last line\n",
		},
		{
			desc:     "Printing lines before match with pattern 'line'",
			args:     []string{"grep", "-n", "-B", "1", "line", tempFile.Name()},
			expected: "1:line 1\n3:Line with hello in between\n4:another line\n5:last line\n",
		},
		{
			desc:     "Printing lines around match with pattern 'line'",
			args:     []string{"grep", "-n", "-C", "1", "line", tempFile.Name()},
			expected: "1:line 1\n2:hello world\n3:Line with hello in between\n4:another line\n5:last line\n",
		},
		{
			desc:     "Inverting match with pattern 'line'",
			args:     []string{"grep", "-n", "-v", "line", tempFile.Name()},
			expected: "2:hello world\n3:Line with hello in between\n",
		},
		{
			desc:     "Fixed string match with pattern 'hello world'",
			args:     []string{"grep", "-n", "-F", "hello world", tempFile.Name()},
			expected: "2:hello world\n",
		},
		{
			desc:     "Counting lines with pattern 'line'",
			args:     []string{"grep", "-c", "line", tempFile.Name()},
			expected: "3\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			output := runGrep(tc.args...)
			if output != tc.expected {
				t.Errorf("Expected output:\n%s\nGot:\n%s\n", tc.expected, output)
			}
		})
	}
}

func createTempFile(content string) *os.File {
	file, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		panic(err)
	}
	file.WriteString(content)
	file.Close()
	return file
}

func runGrep(args ...string) string {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Redirect stdout to buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a new flag set and parse the args
	flag.CommandLine = flag.NewFlagSet(args[0], flag.ExitOnError)
	os.Args = args
	main()

	// Capture output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = old

	return buf.String()
}
