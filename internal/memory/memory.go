package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Memory represents project memory/context
type Memory struct {
	ProjectType string   `json:"project_type"`
	Framework   string   `json:"framework"`
	Languages   []string `json:"languages"`
	Structure   string   `json:"structure"`
	Rules       []string `json:"rules"`
	Notes       string   `json:"notes"`
	// 新增：代码语义记忆
	CodeSemantics *CodeSemantics `json:"code_semantics,omitempty"`
	// 新增：学习历史
	LearningHistory []LearningEntry `json:"learning_history,omitempty"`
	// 新增：错误记忆
	ErrorPatterns []ErrorPattern `json:"error_patterns,omitempty"`
	// 新增：API 使用模式
	APIPatterns []APIPattern `json:"api_patterns,omitempty"`
	// 新增：文件关系图谱
	FileRelations map[string][]string `json:"file_relations,omitempty"`
	// 元数据
	LastUpdated time.Time `json:"last_updated"`
	Version     string    `json:"version"`
}

// CodeSemantics 代码语义理解
type CodeSemantics struct {
	// 关键函数/方法及其用途
	KeyFunctions map[string]string `json:"key_functions"`
	// 数据结构定义
	DataStructures map[string]string `json:"data_structures"`
	// 接口和抽象
	Interfaces map[string]string `json:"interfaces"`
	// 业务逻辑关键点
	BusinessLogic []LogicPoint `json:"business_logic"`
}

// LogicPoint 业务逻辑点
type LogicPoint struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
	Importance  int      `json:"importance"`
}

// LearningEntry 学习历史条目
type LearningEntry struct {
	ID          string    `json:"id"`
	Topic       string    `json:"topic"`
	Content     string    `json:"content"`
	Source      string    `json:"source"`
	LearnedAt   time.Time `json:"learned_at"`
	AccessCount int       `json:"access_count"`
}

// ErrorPattern 错误模式记忆
type ErrorPattern struct {
	ID        string    `json:"id"`
	Pattern   string    `json:"pattern"`
	ErrorType string    `json:"error_type"`
	Solution  string    `json:"solution"`
	Context   string    `json:"context"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	Frequency int       `json:"frequency"`
}

// APIPattern API 使用模式
type APIPattern struct {
	Name           string   `json:"name"`
	Package        string   `json:"package"`
	Usage          string   `json:"usage"`
	Examples       []string `json:"examples"`
	CommonMistakes []string `json:"common_mistakes"`
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

// ==================== 新增记忆功能 ====================

// RecordError 记录错误模式
func (m *Manager) RecordError(errorType, pattern, solution, context string) {
	memory, _ := m.Load()
	if memory == nil {
		memory = &Memory{Version: "2.0"}
	}

	now := time.Now()

	// 检查是否已存在相似错误
	for i, ep := range memory.ErrorPatterns {
		if ep.Pattern == pattern && ep.ErrorType == errorType {
			memory.ErrorPatterns[i].Frequency++
			memory.ErrorPatterns[i].LastSeen = now
			memory.LastUpdated = now
			m.Save(memory)
			return
		}
	}

	// 添加新错误模式
	newError := ErrorPattern{
		ID:        fmt.Sprintf("err-%d", now.UnixNano()),
		Pattern:   pattern,
		ErrorType: errorType,
		Solution:  solution,
		Context:   context,
		FirstSeen: now,
		LastSeen:  now,
		Frequency: 1,
	}
	memory.ErrorPatterns = append(memory.ErrorPatterns, newError)
	memory.LastUpdated = now
	m.Save(memory)
}

// FindSimilarErrors 查找相似错误
func (m *Manager) FindSimilarErrors(errorMsg string) []ErrorPattern {
	memory, _ := m.Load()
	if memory == nil {
		return nil
	}

	var matches []ErrorPattern
	for _, ep := range memory.ErrorPatterns {
		// 简单匹配：检查错误消息是否包含模式
		if strings.Contains(strings.ToLower(errorMsg), strings.ToLower(ep.Pattern)) {
			matches = append(matches, ep)
		}
	}

	return matches
}

// Learn 记录学习到的知识
func (m *Manager) Learn(topic, content, source string) {
	memory, _ := m.Load()
	if memory == nil {
		memory = &Memory{Version: "2.0"}
	}

	entry := LearningEntry{
		ID:        fmt.Sprintf("learn-%d", time.Now().UnixNano()),
		Topic:     topic,
		Content:   content,
		Source:    source,
		LearnedAt: time.Now(),
	}

	memory.LearningHistory = append(memory.LearningHistory, entry)
	memory.LastUpdated = time.Now()

	// 只保留最近100条学习记录
	if len(memory.LearningHistory) > 100 {
		memory.LearningHistory = memory.LearningHistory[len(memory.LearningHistory)-100:]
	}

	m.Save(memory)
}

// GetRelevantLearning 获取相关学习记录
func (m *Manager) GetRelevantLearning(query string, limit int) []LearningEntry {
	memory, _ := m.Load()
	if memory == nil {
		return nil
	}

	keywords := strings.Fields(strings.ToLower(query))
	var scored []struct {
		entry LearningEntry
		score int
	}

	for _, entry := range memory.LearningHistory {
		score := 0
		content := strings.ToLower(entry.Topic + " " + entry.Content)
		for _, kw := range keywords {
			if strings.Contains(content, kw) {
				score++
			}
		}
		if score > 0 {
			scored = append(scored, struct {
				entry LearningEntry
				score int
			}{entry, score})
		}
	}

	// 按相关性排序
	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// 返回前 limit 个
	var result []LearningEntry
	for i := 0; i < len(scored) && i < limit; i++ {
		result = append(result, scored[i].entry)
	}

	return result
}

// AddFileRelation 添加文件关系
func (m *Manager) AddFileRelation(file string, relatedFiles []string) {
	memory, _ := m.Load()
	if memory == nil {
		memory = &Memory{Version: "2.0", FileRelations: make(map[string][]string)}
	}

	if memory.FileRelations == nil {
		memory.FileRelations = make(map[string][]string)
	}

	// 合并关系，去重
	existing := memory.FileRelations[file]
	relationMap := make(map[string]bool)
	for _, f := range existing {
		relationMap[f] = true
	}
	for _, f := range relatedFiles {
		relationMap[f] = true
	}

	var newRelations []string
	for f := range relationMap {
		newRelations = append(newRelations, f)
	}

	memory.FileRelations[file] = newRelations
	memory.LastUpdated = time.Now()
	m.Save(memory)
}

// GetRelatedFiles 获取相关文件
func (m *Manager) GetRelatedFiles(file string) []string {
	memory, _ := m.Load()
	if memory == nil || memory.FileRelations == nil {
		return nil
	}
	return memory.FileRelations[file]
}

// AddKeyFunction 记录关键函数
func (m *Manager) AddKeyFunction(name, description string) {
	memory, _ := m.Load()
	if memory == nil {
		memory = &Memory{Version: "2.0"}
	}

	if memory.CodeSemantics == nil {
		memory.CodeSemantics = &CodeSemantics{
			KeyFunctions:   make(map[string]string),
			DataStructures: make(map[string]string),
			Interfaces:     make(map[string]string),
		}
	}

	memory.CodeSemantics.KeyFunctions[name] = description
	memory.LastUpdated = time.Now()
	m.Save(memory)
}

// GetEnhancedPromptContext 获取增强的提示上下文
func (m *Manager) GetEnhancedPromptContext() string {
	memory, _ := m.Load()
	if memory == nil {
		return ""
	}

	var sb strings.Builder

	// 项目基本信息
	if memory.ProjectType != "" {
		sb.WriteString(fmt.Sprintf("项目类型: %s\n", memory.ProjectType))
	}
	if memory.Framework != "" {
		sb.WriteString(fmt.Sprintf("框架: %s\n", memory.Framework))
	}

	// 关键函数
	if memory.CodeSemantics != nil && len(memory.CodeSemantics.KeyFunctions) > 0 {
		sb.WriteString("\n关键函数:\n")
		for name, desc := range memory.CodeSemantics.KeyFunctions {
			sb.WriteString(fmt.Sprintf("  - %s: %s\n", name, desc))
		}
	}

	// 最近的错误模式
	if len(memory.ErrorPatterns) > 0 {
		sb.WriteString("\n常见错误及解决方案:\n")
		for i := len(memory.ErrorPatterns) - 1; i >= 0 && i >= len(memory.ErrorPatterns)-3; i-- {
			ep := memory.ErrorPatterns[i]
			sb.WriteString(fmt.Sprintf("  - %s: %s\n", ep.ErrorType, ep.Solution))
		}
	}

	// 编码规则
	if len(memory.Rules) > 0 {
		sb.WriteString("\n编码规则:\n")
		for _, rule := range memory.Rules {
			sb.WriteString(fmt.Sprintf("  - %s\n", rule))
		}
	}

	return sb.String()
}
