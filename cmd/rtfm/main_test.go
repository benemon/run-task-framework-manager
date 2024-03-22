package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
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
	_, err := os.Stat("./testMainFunctionGo/go-run-task/cmd/main.go")
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

func TestGenerateGoScaffold(t *testing.T) {
	runTaskName := "testGenerateGoScaffold"
	workingDir := "."

	err := generateGoScaffold(runTaskName, workingDir)
	require.NoError(t, err)

	// Verify that the expected files are generated
	expectedFiles := []string{
		"go.mod",
		"cmd/main.go",
		"internal/api/run_task_request.go",
		"internal/api/run_task_response.go",
		"internal/controller/run_task_controller.go",
		"internal/controller/run_task_controller_test.go",
		"Containerfile",
		"README.md",
		// Add more expected files here...
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(workingDir, runTaskName, file)
		_, err := os.Stat(filePath)
		require.NoErrorf(t, err, "expected file %s to be generated", filePath)
	}

	// Clean up: delete the test directory
	os.RemoveAll("./testGenerateGoScaffold")
}

func TestGeneratePythonScaffold(t *testing.T) {
	runTaskName := "testGeneratePythonScaffold"
	workingDir := "."

	err := generatePythonScaffold(runTaskName, workingDir)
	require.NoError(t, err)

	// Verify the existence of the generated files
	expectedFiles := []string{
		"testGeneratePythonScaffold.py",
		"requirements.txt",
		"Containerfile",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(workingDir, runTaskName, file)
		_, err := os.Stat(filePath)
		require.NoErrorf(t, err, "expected file %s to be generated", filePath)
	}

	// Clean up: delete the test directory
	os.RemoveAll("./testGeneratePythonScaffold")
}
