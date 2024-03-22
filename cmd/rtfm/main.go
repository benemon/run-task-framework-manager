package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

//go:embed resources/go/*
var goTemplates embed.FS

//go:embed resources/python/*
var pythonTemplates embed.FS

var (
	name       *string
	workingDir *string
	language   *string
)

func init() {
	name = flag.String("name", "", "Run Task name (required)")
	workingDir = flag.String("dir", "", "Target directory (required)")
	language = flag.String("language", "", "Language (required). Supported languages: go, python")
}

func createTargetDir(workingDir, runTaskName string) (string, error) {
	targetDir := filepath.Join(workingDir, runTaskName)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create target directory: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("failed to check if target directory exists: %w", err)
	}

	return targetDir, nil
}

func generateGoScaffold(runTaskName, workingDir string) error {
	goVersion := strings.TrimPrefix(runtime.Version(), "go")

	targetDir, err := createTargetDir(workingDir, runTaskName)
	if err != nil {
		return err
	}

	templates := map[string]interface{}{
		"go.mod":                                          struct{ RunTaskName, GoVersion string }{RunTaskName: runTaskName, GoVersion: goVersion},
		"cmd/main.go":                                     struct{ RunTaskName string }{RunTaskName: runTaskName},
		"internal/api/run_task_request.go":                struct{}{},
		"internal/api/run_task_response.go":               struct{}{},
		"internal/controller/run_task_controller.go":      struct{}{},
		"internal/controller/run_task_controller_test.go": struct{}{},
		"Containerfile":                                   struct{ RunTaskName, GoVersion string }{RunTaskName: runTaskName, GoVersion: goVersion},
		// Add more templates here...
	}

	for templatePath, data := range templates {
		templatePath = filepath.FromSlash(templatePath)
		templateName := filepath.Base(templatePath)
		templateContent, err := goTemplates.ReadFile(fmt.Sprintf("resources/go/%s.tmpl", templatePath))
		if err != nil {
			return fmt.Errorf("failed to read %s template: %w", templateName, err)
		}

		tmpl, err := template.New(templateName).Parse(string(templateContent))
		if err != nil {
			return fmt.Errorf("failed to parse %s template: %w", templateName, err)
		}

		outputFile := filepath.Join(targetDir, templatePath)
		outputDir := filepath.Dir(outputFile)

		// Ensure the directory exists
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", templatePath, err)
		}

		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create %s file: %w", templatePath, err)
		}
		defer f.Close()

		err = tmpl.Execute(f, data)
		if err != nil {
			return fmt.Errorf("failed to execute %s template: %w", templateName, err)
		}
	}

	return nil
}

func generatePythonScaffold(runTaskName, workingDir string) error {

	targetDir, err := createTargetDir(workingDir, runTaskName)
	if err != nil {
		return err
	}

	templates := map[string]interface{}{
		"main.py":          struct{}{},
		"requirements.txt": struct{}{},
		"Containerfile":    struct{ RunTaskName string }{RunTaskName: runTaskName},
		// Add more templates here...
	}

	for templatePath, data := range templates {
		templatePath = filepath.FromSlash(templatePath)
		templateName := filepath.Base(templatePath)
		templateContent, err := pythonTemplates.ReadFile(fmt.Sprintf("resources/python/%s.tmpl", templatePath))
		if err != nil {
			return fmt.Errorf("failed to read %s template: %w", templateName, err)
		}

		tmpl, err := template.New(templateName).Parse(string(templateContent))
		if err != nil {
			return fmt.Errorf("failed to parse %s template: %w", templateName, err)
		}

		// If the template is main.py, replace it with runTaskName
		if templateName == "main.py" {
			templatePath = filepath.Join(filepath.Dir(templatePath), runTaskName+".py")
		}

		outputFile := filepath.Join(targetDir, templatePath)
		outputDir := filepath.Dir(outputFile)

		// Ensure the directory exists
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", templatePath, err)
		}

		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create %s file: %w", templatePath, err)
		}
		defer f.Close()

		err = tmpl.Execute(f, data)
		if err != nil {
			return fmt.Errorf("failed to execute %s template: %w", templateName, err)
		}

	}

	return nil
}

func main() {
	flag.Parse()

	if *name == "" || *workingDir == "" || *language == "" {
		flag.PrintDefaults()
		return
	}

	switch *language {
	case "go":
		err := generateGoScaffold(*name, *workingDir)
		if err != nil {
			fmt.Printf("Failed to generate Go scaffold: %v\n", err)
			return
		}
	case "python":
		err := generatePythonScaffold(*name, *workingDir)
		if err != nil {
			fmt.Printf("Failed to generate Python scaffold: %v\n", err)
			return
		}
	default:
		fmt.Printf("Unsupported language: %s", *language)
		return
	}

	fmt.Println("Scaffold generated successfully")
}
