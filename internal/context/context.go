package context

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ContextType represents different types of context
type ContextType string

const (
	ContextFile     ContextType = "file"
	ContextCode     ContextType = "code"
	ContextCommand  ContextType = "command"
	ContextDecision ContextType = "decision"
	ContextError    ContextType = "error"
	ContextLearned  ContextType = "learned"
	// 新增上下文类型
	ContextFunction   ContextType = "function"
	ContextDependency ContextType = "dependency"
	ContextPattern    ContextType = "pattern"
	ContextTest       ContextType = "test"
)

// Entry represents a context entry
type Entry struct {
	ID           string      `json:"id"`
	Type         ContextType `json:"type"`
	Content      string      `json:"content"`
	File         string      `json:"file,omitempty"`
	Line         int         `json:"line,omitempty"`
	Tags         []string    `json:"tags,omitempty"`
	Importance   int         `json:"importance"` // 1-10
	CreatedAt    time.Time   `json:"created_at"`
	AccessCount  int         `json:"access_count"`
	LastAccessed time.Time   `json:"last_accessed"`
}

// Context manages conversation and project context
type Context struct {
	WorkDir     string        `json:"work_dir"`
	Entries     []Entry       `json:"entries"`
	ProjectInfo *ProjectInfo  `json:"project_info,omitempty"`
	RecentFiles []string      `json:"recent_files"`
	Decisions   []Decision    `json:"decisions"`
	Patterns    []CodePattern `json:"patterns"`
	// 新增：代码关系图谱
	CodeGraph *CodeGraph `json:"code_graph,omitempty"`
	// 新增：会话历史摘要
	SessionSummary []string  `json:"session_summary,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ProjectInfo contains analyzed project information
type ProjectInfo struct {
	Name          string            `json:"name"`
	Type          string            `json:"type"`
	Language      string            `json:"language"`
	Framework     string            `json:"framework,omitempty"`
	Structure     map[string]string `json:"structure"`
	EntryPoints   []string          `json:"entry_points"`
	Dependencies  []string          `json:"dependencies"`
	ConfigFiles   []string          `json:"config_files"`
	TestFramework string            `json:"test_framework,omitempty"`
	CodeStyle     string            `json:"code_style,omitempty"`
}

// Decision represents a design decision
type Decision struct {
	ID           string    `json:"id"`
	Topic        string    `json:"topic"`
	Decision     string    `json:"decision"`
	Rationale    string    `json:"rationale"`
	Alternatives []string  `json:"alternatives,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// CodePattern represents a discovered code pattern
type CodePattern struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Examples    []string `json:"examples"`
	Frequency   int      `json:"frequency"`
}

// FunctionInfo 函数信息
type FunctionInfo struct {
	Name       string   `json:"name"`
	File       string   `json:"file"`
	Line       int      `json:"line"`
	Signature  string   `json:"signature"`
	Doc        string   `json:"doc"`
	Calls      []string `json:"calls"`      // 调用的函数
	CalledBy   []string `json:"called_by"`  // 被谁调用
	Complexity int      `json:"complexity"` // 复杂度
}

// CodeGraph 代码关系图谱
type CodeGraph struct {
	Functions    map[string]*FunctionInfo `json:"functions"`
	Files        map[string]*FileInfo     `json:"files"`
	Dependencies []Dependency             `json:"dependencies"`
}

// FileInfo 文件信息
type FileInfo struct {
	Path      string   `json:"path"`
	Package   string   `json:"package"`
	Imports   []string `json:"imports"`
	Functions []string `json:"functions"`
	Types     []string `json:"types"`
	LineCount int      `json:"line_count"`
	TestFile  bool     `json:"test_file"`
}

// Dependency 依赖关系
type Dependency struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"` // import, call, inherit
}

// Manager manages context
type Manager struct {
	workDir    string
	contextDir string
	context    *Context
}

// NewManager creates a new context manager
func NewManager(workDir string) (*Manager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	contextDir := filepath.Join(home, ".accil", "contexts")
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		return nil, err
	}

	m := &Manager{
		workDir:    workDir,
		contextDir: contextDir,
	}

	m.context, _ = m.Load()

	return m, nil
}

// Load loads context from disk
func (m *Manager) Load() (*Context, error) {
	contextFile := filepath.Join(m.workDir, ".accil-context.json")
	data, err := os.ReadFile(contextFile)
	if err != nil {
		// Create new context
		return &Context{
			WorkDir:     m.workDir,
			Entries:     []Entry{},
			RecentFiles: []string{},
			Decisions:   []Decision{},
			Patterns:    []CodePattern{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}

	var ctx Context
	if err := json.Unmarshal(data, &ctx); err != nil {
		return nil, err
	}

	return &ctx, nil
}

// Save saves context to disk
func (m *Manager) Save() error {
	m.context.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(m.context, "", "  ")
	if err != nil {
		return err
	}

	contextFile := filepath.Join(m.workDir, ".accil-context.json")
	return os.WriteFile(contextFile, data, 0644)
}

// AddEntry adds a new context entry
func (m *Manager) AddEntry(entryType ContextType, content string, file string, tags []string) *Entry {
	entry := Entry{
		ID:         fmt.Sprintf("ctx-%d", time.Now().UnixNano()),
		Type:       entryType,
		Content:    content,
		File:       file,
		Tags:       tags,
		Importance: 5,
		CreatedAt:  time.Now(),
	}

	m.context.Entries = append(m.context.Entries, entry)

	// Keep only recent 100 entries
	if len(m.context.Entries) > 100 {
		m.context.Entries = m.context.Entries[len(m.context.Entries)-100:]
	}

	return &entry
}

// AddDecision records a design decision
func (m *Manager) AddDecision(topic, decision, rationale string, alternatives []string) {
	d := Decision{
		ID:           fmt.Sprintf("dec-%d", time.Now().UnixNano()),
		Topic:        topic,
		Decision:     decision,
		Rationale:    rationale,
		Alternatives: alternatives,
		CreatedAt:    time.Now(),
	}

	m.context.Decisions = append(m.context.Decisions, d)
}

// TrackFile tracks a recently accessed file
func (m *Manager) TrackFile(filePath string) {
	// Remove if exists
	for i, f := range m.context.RecentFiles {
		if f == filePath {
			m.context.RecentFiles = append(m.context.RecentFiles[:i], m.context.RecentFiles[i+1:]...)
			break
		}
	}

	// Add to front
	m.context.RecentFiles = append([]string{filePath}, m.context.RecentFiles...)

	// Keep only 20 recent files
	if len(m.context.RecentFiles) > 20 {
		m.context.RecentFiles = m.context.RecentFiles[:20]
	}
}

// GetRelevantContext retrieves context relevant to a query
func (m *Manager) GetRelevantContext(query string, maxTokens int) string {
	var sb strings.Builder
	tokenCount := 0

	// Add project info
	if m.context.ProjectInfo != nil {
		sb.WriteString("## Project Information\n")
		sb.WriteString(fmt.Sprintf("Language: %s\n", m.context.ProjectInfo.Language))
		if m.context.ProjectInfo.Framework != "" {
			sb.WriteString(fmt.Sprintf("Framework: %s\n", m.context.ProjectInfo.Framework))
		}
		tokenCount += 50
	}

	// Add recent decisions
	if len(m.context.Decisions) > 0 {
		sb.WriteString("\n## Recent Decisions\n")
		for i := len(m.context.Decisions) - 1; i >= 0 && i >= len(m.context.Decisions)-5; i-- {
			d := m.context.Decisions[i]
			sb.WriteString(fmt.Sprintf("- %s: %s\n", d.Topic, d.Decision))
			tokenCount += 20
		}
	}

	// Add relevant entries based on query keywords
	keywords := extractKeywords(query)
	for _, entry := range m.context.Entries {
		if tokenCount >= maxTokens {
			break
		}

		relevance := calculateRelevance(entry, keywords)
		if relevance > 0.3 {
			sb.WriteString(fmt.Sprintf("\n[%s] %s\n", entry.Type, entry.Content))
			if entry.File != "" {
				sb.WriteString(fmt.Sprintf("File: %s\n", entry.File))
			}
			tokenCount += len(entry.Content) / 4
		}
	}

	return sb.String()
}

// AnalyzeProject analyzes the project structure
func (m *Manager) AnalyzeProject() error {
	info := &ProjectInfo{
		Structure:    make(map[string]string),
		EntryPoints:  []string{},
		Dependencies: []string{},
		ConfigFiles:  []string{},
	}

	// Detect project type
	entries, err := os.ReadDir(m.workDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		switch {
		case name == "go.mod":
			info.Language = "Go"
			info.Type = "Go Module"
			m.parseGoProject(info)
		case name == "package.json":
			info.Language = "JavaScript/TypeScript"
			info.Type = "Node.js Project"
			m.parseJSProject(info)
		case name == "requirements.txt" || name == "pyproject.toml":
			info.Language = "Python"
			info.Type = "Python Project"
		case name == "Cargo.toml":
			info.Language = "Rust"
			info.Type = "Rust Project"
		}
	}

	// Find entry points
	m.findEntryPoints(info)

	// Find config files
	m.findConfigFiles(info)

	m.context.ProjectInfo = info
	return nil
}

func (m *Manager) parseGoProject(info *ProjectInfo) {
	goMod, err := os.ReadFile(filepath.Join(m.workDir, "go.mod"))
	if err != nil {
		return
	}

	content := string(goMod)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			info.Name = strings.TrimPrefix(line, "module ")
		}
		if !strings.HasPrefix(line, "require") && !strings.HasPrefix(line, ")") && strings.Contains(line, " ") {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				info.Dependencies = append(info.Dependencies, parts[0])
			}
		}
	}
}

func (m *Manager) parseJSProject(info *ProjectInfo) {
	packageJSON, err := os.ReadFile(filepath.Join(m.workDir, "package.json"))
	if err != nil {
		return
	}

	var pkg struct {
		Name         string                 `json:"name"`
		Dependencies map[string]interface{} `json:"dependencies"`
	}
	if err := json.Unmarshal(packageJSON, &pkg); err != nil {
		return
	}

	info.Name = pkg.Name
	for dep := range pkg.Dependencies {
		info.Dependencies = append(info.Dependencies, dep)
	}
}

func (m *Manager) findEntryPoints(info *ProjectInfo) {
	entryPatterns := []string{"main.go", "main.py", "index.js", "index.ts", "app.py", "server.go"}
	for _, pattern := range entryPatterns {
		if _, err := os.Stat(filepath.Join(m.workDir, pattern)); err == nil {
			info.EntryPoints = append(info.EntryPoints, pattern)
		}
	}
}

func (m *Manager) findConfigFiles(info *ProjectInfo) {
	configPatterns := []string{
		".gitignore", "Dockerfile", "docker-compose.yml",
		".env.example", "Makefile", "README.md",
		"tsconfig.json", "eslint.config.js", ".prettierrc",
	}

	for _, pattern := range configPatterns {
		if _, err := os.Stat(filepath.Join(m.workDir, pattern)); err == nil {
			info.ConfigFiles = append(info.ConfigFiles, pattern)
		}
	}
}

// LearnPattern learns a code pattern
func (m *Manager) LearnPattern(name, description string, example string) {
	// Check if pattern exists
	for i, p := range m.context.Patterns {
		if p.Name == name {
			m.context.Patterns[i].Frequency++
			m.context.Patterns[i].Examples = append(p.Examples, example)
			return
		}
	}

	// Add new pattern
	m.context.Patterns = append(m.context.Patterns, CodePattern{
		Name:        name,
		Description: description,
		Examples:    []string{example},
		Frequency:   1,
	})
}

// GetPromptContext builds context for AI prompts
func (m *Manager) GetPromptContext() string {
	var sb strings.Builder

	if m.context.ProjectInfo != nil {
		sb.WriteString(fmt.Sprintf("Project: %s (%s)\n",
			m.context.ProjectInfo.Name,
			m.context.ProjectInfo.Language))
	}

	if len(m.context.RecentFiles) > 0 {
		sb.WriteString("Recent files:\n")
		for _, f := range m.context.RecentFiles[:5] {
			sb.WriteString(fmt.Sprintf("- %s\n", f))
		}
	}

	if len(m.context.Decisions) > 0 {
		sb.WriteString("Design decisions:\n")
		for _, d := range m.context.Decisions {
			sb.WriteString(fmt.Sprintf("- %s\n", d.Decision))
		}
	}

	return sb.String()
}

// Helper functions
func extractKeywords(query string) []string {
	words := strings.Fields(strings.ToLower(query))
	keywords := []string{}
	stopWords := map[string]bool{"the": true, "a": true, "an": true, "is": true, "are": true, "was": true, "were": true, "be": true, "been": true, "being": true, "have": true, "has": true, "had": true, "do": true, "does": true, "did": true, "will": true, "would": true, "could": true, "should": true, "may": true, "might": true, "must": true, "shall": true, "can": true, "need": true, "to": true, "of": true, "in": true, "for": true, "on": true, "with": true, "at": true, "by": true, "from": true, "as": true, "into": true, "through": true, "during": true, "before": true, "after": true, "above": true, "below": true, "between": true, "under": true, "again": true, "further": true, "then": true, "once": true}

	for _, word := range words {
		if !stopWords[word] && len(word) > 2 {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

func calculateRelevance(entry Entry, keywords []string) float64 {
	if len(keywords) == 0 {
		return 0
	}

	content := strings.ToLower(entry.Content)
	matches := 0
	for _, kw := range keywords {
		if strings.Contains(content, kw) {
			matches++
		}
	}

	return float64(matches) / float64(len(keywords))
}

// ==================== 新增代码图谱功能 ====================

// BuildCodeGraph 构建代码关系图谱
func (m *Manager) BuildCodeGraph() error {
	if m.context.CodeGraph == nil {
		m.context.CodeGraph = &CodeGraph{
			Functions: make(map[string]*FunctionInfo),
			Files:     make(map[string]*FileInfo),
		}
	}

	// 遍历项目文件
	err := filepath.Walk(m.workDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// 跳过隐藏文件和依赖目录
		if strings.HasPrefix(info.Name(), ".") || strings.Contains(path, "node_modules") || strings.Contains(path, "vendor") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 只分析代码文件
		if !info.IsDir() && isCodeFile(path) {
			m.analyzeFile(path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	m.context.UpdatedAt = time.Now()
	return m.Save()
}

// analyzeFile 分析单个文件
func (m *Manager) analyzeFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	relPath, _ := filepath.Rel(m.workDir, filePath)
	lines := strings.Split(string(content), "\n")

	fileInfo := &FileInfo{
		Path:      relPath,
		LineCount: len(lines),
		TestFile:  strings.HasSuffix(filePath, "_test.go") || strings.Contains(filePath, ".test."),
	}

	// 简单分析：提取包名、导入、函数定义
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// 包名
		if strings.HasPrefix(line, "package ") {
			fileInfo.Package = strings.TrimPrefix(line, "package ")
		}

		// 导入
		if strings.HasPrefix(line, "import ") || strings.HasPrefix(line, "\"") {
			// 简化处理
			if strings.Contains(line, "\"") {
				parts := strings.Split(line, "\"")
				if len(parts) >= 2 {
					fileInfo.Imports = append(fileInfo.Imports, parts[1])
				}
			}
		}

		// 函数定义 (Go 风格)
		if strings.HasPrefix(line, "func ") {
			funcName := extractFuncName(line)
			if funcName != "" {
				fileInfo.Functions = append(fileInfo.Functions, funcName)

				// 添加到函数图谱
				m.context.CodeGraph.Functions[funcName] = &FunctionInfo{
					Name:      funcName,
					File:      relPath,
					Line:      i + 1,
					Signature: line,
				}
			}
		}
	}

	m.context.CodeGraph.Files[relPath] = fileInfo
}

// extractFuncName 提取函数名
func extractFuncName(line string) string {
	// 简化处理：func Name(...) 或 func (r *Receiver) Name(...)
	line = strings.TrimPrefix(line, "func ")
	if idx := strings.Index(line, "("); idx != -1 {
		before := line[:idx]
		// 检查是否有接收器
		if strings.HasPrefix(before, "(") {
			// 有接收器，找接收器后的名称
			if idx2 := strings.Index(line, ")"); idx2 != -1 {
				after := strings.TrimSpace(line[idx2+1:])
				if idx3 := strings.Index(after, "("); idx3 != -1 {
					return strings.TrimSpace(after[:idx3])
				}
			}
		} else {
			return strings.TrimSpace(before)
		}
	}
	return ""
}

// isCodeFile 检查是否为代码文件
func isCodeFile(path string) bool {
	extensions := []string{".go", ".js", ".ts", ".py", ".java", ".rs", ".cpp", ".c", ".h"}
	for _, ext := range extensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}

// GetFunctionContext 获取函数上下文
func (m *Manager) GetFunctionContext(funcName string) *FunctionInfo {
	if m.context.CodeGraph == nil {
		return nil
	}
	return m.context.CodeGraph.Functions[funcName]
}

// FindRelatedFunctions 查找相关函数
func (m *Manager) FindRelatedFunctions(funcName string) []string {
	if m.context.CodeGraph == nil {
		return nil
	}

	funcInfo, exists := m.context.CodeGraph.Functions[funcName]
	if !exists {
		return nil
	}

	// 合并调用和被调用
	related := make(map[string]bool)
	for _, f := range funcInfo.Calls {
		related[f] = true
	}
	for _, f := range funcInfo.CalledBy {
		related[f] = true
	}

	var result []string
	for f := range related {
		result = append(result, f)
	}
	return result
}

// AddSessionSummary 添加会话摘要
func (m *Manager) AddSessionSummary(summary string) {
	m.context.SessionSummary = append(m.context.SessionSummary, summary)
	// 只保留最近20条摘要
	if len(m.context.SessionSummary) > 20 {
		m.context.SessionSummary = m.context.SessionSummary[len(m.context.SessionSummary)-20:]
	}
	m.context.UpdatedAt = time.Now()
	m.Save()
}

// GetSessionContext 获取会话上下文
func (m *Manager) GetSessionContext() string {
	if len(m.context.SessionSummary) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("## 会话历史\n")
	for _, summary := range m.context.SessionSummary {
		sb.WriteString(fmt.Sprintf("- %s\n", summary))
	}
	return sb.String()
}
