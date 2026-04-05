package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Memory represents project memory/context
type Memory struct {
	ProjectType string   `json:"project_type"`
	Framework   string   `json:"framework"`
	Languages   []string `json:"languages"`
	Structure   string   `json:"structure"`
	Rules       []string `json:"rules"`
	Notes       string   `json:"notes"`
}

const AgentsFileName = "AGENTS.md"

// Manager manages project memory
type Manager struct {
	workDir string
}

// NewManager creates a new memory manager
func NewManager(workDir string) *Manager {
	return &Manager{workDir: workDir}
}

// Exists checks if AGENTS.md exists
func (m *Manager) Exists() bool {
	path := filepath.Join(m.workDir, AgentsFileName)
	_, err := os.Stat(path)
	return err == nil
}

// Load loads the memory file
func (m *Manager) Load() (*Memory, error) {
	path := filepath.Join(m.workDir, AgentsFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse markdown file to extract memory
	content := string(data)
	memory := &Memory{}

	// Simple markdown parsing
	lines := strings.Split(content, "\n")
	var currentSection string
	var rules []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "## Project Type") {
			currentSection = "project_type"
		} else if strings.HasPrefix(line, "## Framework") {
			currentSection = "framework"
		} else if strings.HasPrefix(line, "## Languages") {
			currentSection = "languages"
		} else if strings.HasPrefix(line, "## Directory Structure") {
			currentSection = "structure"
		} else if strings.HasPrefix(line, "## Coding Rules") {
			currentSection = "rules"
		} else if strings.HasPrefix(line, "## Notes") {
			currentSection = "notes"
		} else if line != "" && !strings.HasPrefix(line, "#") {
			switch currentSection {
			case "project_type":
				memory.ProjectType = line
			case "framework":
				memory.Framework = line
			case "languages":
				if strings.HasPrefix(line, "- ") {
					memory.Languages = append(memory.Languages, strings.TrimPrefix(line, "- "))
				}
			case "structure":
				memory.Structure += line + "\n"
			case "rules":
				if strings.HasPrefix(line, "- ") {
					rules = append(rules, strings.TrimPrefix(line, "- "))
				}
			case "notes":
				memory.Notes += line + "\n"
			}
		}
	}

	memory.Rules = rules
	return memory, nil
}

// LoadRaw loads the raw content of AGENTS.md
func (m *Manager) LoadRaw() (string, error) {
	path := filepath.Join(m.workDir, AgentsFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Save saves the memory file
func (m *Manager) Save(memory *Memory) error {
	var sb strings.Builder

	sb.WriteString("# Project Memory\n\n")
	sb.WriteString("This file contains project context for the AI assistant.\n\n")

	sb.WriteString("## Project Type\n")
	sb.WriteString(memory.ProjectType + "\n\n")

	sb.WriteString("## Framework\n")
	sb.WriteString(memory.Framework + "\n\n")

	sb.WriteString("## Languages\n")
	for _, lang := range memory.Languages {
		sb.WriteString("- " + lang + "\n")
	}
	sb.WriteString("\n")

	sb.WriteString("## Directory Structure\n")
	sb.WriteString("```\n")
	sb.WriteString(memory.Structure)
	sb.WriteString("```\n\n")

	sb.WriteString("## Coding Rules\n")
	for _, rule := range memory.Rules {
		sb.WriteString("- " + rule + "\n")
	}
	sb.WriteString("\n")

	sb.WriteString("## Notes\n")
	sb.WriteString(memory.Notes + "\n")

	path := filepath.Join(m.workDir, AgentsFileName)
	return os.WriteFile(path, []byte(sb.String()), 0644)
}

// Generate generates memory by analyzing the project
func (m *Manager) Generate() (*Memory, error) {
	memory := &Memory{
		Languages: []string{},
		Rules:     []string{},
	}

	// Analyze project structure
	entries, err := os.ReadDir(m.workDir)
	if err != nil {
		return nil, err
	}

	// Detect languages and frameworks
	for _, entry := range entries {
		name := entry.Name()

		switch {
		case name == "go.mod":
			memory.Languages = append(memory.Languages, "Go")
			memory.ProjectType = "Go Module"
			memory.Framework = detectGoFramework(m.workDir)
		case name == "package.json":
			memory.Languages = append(memory.Languages, "JavaScript/TypeScript")
			memory.ProjectType = "Node.js Project"
			memory.Framework = detectJSFramework(m.workDir)
		case name == "requirements.txt" || name == "pyproject.toml":
			memory.Languages = append(memory.Languages, "Python")
			memory.ProjectType = "Python Project"
		case name == "Cargo.toml":
			memory.Languages = append(memory.Languages, "Rust")
			memory.ProjectType = "Rust Project"
		case name == "pom.xml":
			memory.Languages = append(memory.Languages, "Java")
			memory.ProjectType = "Java/Maven Project"
		case name == "build.gradle":
			memory.Languages = append(memory.Languages, "Java/Kotlin")
			memory.ProjectType = "Gradle Project"
		}
	}

	// Generate directory structure
	structure, err := generateStructure(m.workDir, 2)
	if err == nil {
		memory.Structure = structure
	}

	return memory, nil
}

func detectGoFramework(workDir string) string {
	goMod, err := os.ReadFile(filepath.Join(workDir, "go.mod"))
	if err != nil {
		return "Standard Go"
	}

	content := string(goMod)
	switch {
	case strings.Contains(content, "gin-gonic"):
		return "Gin"
	case strings.Contains(content, "echo"):
		return "Echo"
	case strings.Contains(content, "fiber"):
		return "Fiber"
	case strings.Contains(content, "chi"):
		return "Chi"
	default:
		return "Standard Go"
	}
}

func detectJSFramework(workDir string) string {
	packageJSON, err := os.ReadFile(filepath.Join(workDir, "package.json"))
	if err != nil {
		return "Node.js"
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(packageJSON, &pkg); err != nil {
		return "Node.js"
	}

	deps, ok := pkg["dependencies"].(map[string]interface{})
	if !ok {
		return "Node.js"
	}

	switch {
	case deps["react"] != nil:
		return "React"
	case deps["vue"] != nil:
		return "Vue"
	case deps["angular"] != nil || deps["@angular/core"] != nil:
		return "Angular"
	case deps["svelte"] != nil:
		return "Svelte"
	case deps["next"] != nil:
		return "Next.js"
	case deps["express"] != nil:
		return "Express"
	case deps["fastify"] != nil:
		return "Fastify"
	default:
		return "Node.js"
	}
}

func generateStructure(workDir string, maxDepth int) (string, error) {
	var sb strings.Builder

	err := filepath.Walk(workDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(workDir, path)
		if err != nil {
			return nil
		}

		// Skip hidden files and common ignore patterns
		parts := strings.Split(relPath, string(os.PathSeparator))
		for _, part := range parts {
			if strings.HasPrefix(part, ".") || part == "node_modules" || part == "vendor" {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Limit depth
		depth := len(parts)
		if depth > maxDepth {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Create tree structure
		indent := strings.Repeat("  ", depth-1)
		if info.IsDir() {
			sb.WriteString(fmt.Sprintf("%s%s/\n", indent, info.Name()))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s\n", indent, info.Name()))
		}

		return nil
	})

	return sb.String(), err
}
