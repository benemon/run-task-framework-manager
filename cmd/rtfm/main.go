package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
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

type TemplateData struct {
	RunTaskName    string
	RuntimeVersion string
}

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

func generateScaffold(runTaskName, workingDir, version string, templates embed.FS) error {
	targetDir, err := createTargetDir(workingDir, runTaskName)
	if err != nil {
		return err
	}

	tmplParams := TemplateData{
		RunTaskName:    runTaskName,
		RuntimeVersion: version,
	}

	var templateFiles []string
	var rootDir string
	err = fs.WalkDir(templates, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if rootDir == "" && !d.IsDir() {
			rootDir = filepath.Dir(path)
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".tmpl") {
			templateFiles = append(templateFiles, strings.TrimSuffix(path, ".tmpl"))
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to discover template files: %w", err)
	}

	for _, templatePath := range templateFiles {
		templatePath = filepath.FromSlash(templatePath)
		templateName := filepath.Base(templatePath)
		templateContent, err := templates.ReadFile(fmt.Sprintf("%s.tmpl", templatePath))
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

		outputFile := filepath.Join(targetDir, processTemplatePath(templatePath, rootDir, ".tmpl"))
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

		err = tmpl.Execute(f, tmplParams)
		if err != nil {
			return fmt.Errorf("failed to execute %s template: %w", templateName, err)
		}
	}

	return nil
}

func processTemplatePath(path, prefix, suffix string) string {
	return strings.TrimSuffix(strings.TrimPrefix(path, prefix), suffix)
}

func main() {
	flag.Parse()

	if *name == "" || *workingDir == "" || *language == "" {
		flag.PrintDefaults()
		return
	}

	switch *language {
	case "go":
		err := generateScaffold(*name, *workingDir, strings.TrimPrefix(runtime.Version(), "go"), goTemplates)
		if err != nil {
			fmt.Printf("Failed to generate Go scaffold: %v\n", err)
			return
		}
	case "python":
		err := generateScaffold(*name, *workingDir, "", pythonTemplates)
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
