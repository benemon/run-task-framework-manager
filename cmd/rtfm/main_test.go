package main

import (
	"os"
	"testing"
)

func TestMainFunctionGo(t *testing.T) {
	// Backup command-line arguments and restore them at the end of the test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set command-line arguments for the test
	os.Args = []string{"cmd", "-name=go-run-task", "-dir=./testGo", "-language=go"}

	// Run the main function
	main()

	// Check if the expected files were created
	_, err := os.Stat("./testGo/go-run-task/main.go")
	if os.IsNotExist(err) {
		t.Errorf("main.go was not created for Go")
	}

	// Clean up: delete the test directory
	os.RemoveAll("./testGo")
}

func TestMainFunctionPython(t *testing.T) {
	// Backup command-line arguments and restore them at the end of the test
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set command-line arguments for the test
	os.Args = []string{"cmd", "-name=python-run-task", "-dir=./testPython", "-language=python"}

	// Run the main function
	main()

	// Check if the expected files were created
	_, err := os.Stat("./testPython/python-run-task/python-run-task.py")
	if os.IsNotExist(err) {
		t.Errorf("python-run-task.py was not created for Python")
	}

	// Clean up: delete the test directory
	os.RemoveAll("./testPython")
}
