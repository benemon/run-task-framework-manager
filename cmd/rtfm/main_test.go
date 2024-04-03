package main

import (
	"bytes"
	"os"
	"testing"
)

func TestMainFunctionGo(t *testing.T) {
	// Backup command-line arguments and restore them at the end of the test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set command-line arguments for the test
	os.Args = []string{"cmd", "-name=go-run-task", "-dir=./testMainFunctionGo", "-language=go"}

	// Run the main function
	main()

	// Check if the expected files were created
	_, err := os.Stat("./testMainFunctionGo/go-run-task/main.go")
	if os.IsNotExist(err) {
		t.Errorf("main.go was not created for Go")
	}

	// Clean up: delete the test directory
	os.RemoveAll("./testMainFunctionGo")
}

func TestMainFunctionPython(t *testing.T) {
	// Backup command-line arguments and restore them at the end of the test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set command-line arguments for the test
	os.Args = []string{"cmd", "-name=python-run-task", "-dir=./testMainFunctionPython", "-language=python"}

	// Run the main function
	main()

	// Check if the expected files were created
	_, err := os.Stat("./testMainFunctionPython/python-run-task/python-run-task.py")
	if os.IsNotExist(err) {
		t.Errorf("python-run-task.py was not created for Python")
	}

	// Clean up: delete the test directory
	os.RemoveAll("./testMainFunctionPython")
}

func TestMainFunctionUnsupportedLanguage(t *testing.T) {
	// Backup command-line arguments and restore them at the end of the test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set command-line arguments for the test
	os.Args = []string{"cmd", "-name=unsupported-run-task", "-dir=./testMainFunctionUnsupportedLanguage", "-language=madeuplanguage"}

	// Backup and redirect output streams
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	// Run the main function
	main()

	// Stop capturing the output
	w.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Read the captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check the output
	expectedOutput := "Unsupported language: madeuplanguage"
	if output != expectedOutput {
		t.Errorf("Expected output '%s', got '%s'", expectedOutput, output)
	}
}
