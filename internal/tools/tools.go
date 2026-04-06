package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

// Executor handles tool execution
type Executor struct {
	workDir   string
	blockList []string
}

// NewExecutor creates a new tool executor
func NewExecutor(workDir string, blockList []string) *Executor {
	return &Executor{
		workDir:   workDir,
		blockList: blockList,
	}
}

// Execute executes a tool call
func (e *Executor) Execute(name string, arguments string) *ToolResult {
	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return &ToolResult{
			Success: false,
			Error:   fmt.Sprintf("Failed to parse arguments: %v", err),
		}
	}

	switch name {
	case "read_file":
		return e.readFile(params)
	case "write_file":
		return e.writeFile(params)
	case "edit_file":
		return e.editFile(params)
	case "run_command":
		return e.runCommand(params)
	case "list_dir":
		return e.listDir(params)
	case "search_code":
		return e.searchCode(params)
	case "glob":
		return e.glob(params)
	case "web_search":
		return e.webSearch(params)
	case "web_fetch":
		return e.webFetch(params)
	default:
		return &ToolResult{
			Success: false,
			Error:   fmt.Sprintf("Unknown tool: %s", name),
		}
	}
}

// NeedsConfirmation returns true if the tool needs user confirmation
func (e *Executor) NeedsConfirmation(name string, arguments string) (bool, string, error) {
	switch name {
	case "read_file", "list_dir", "search_code", "glob", "web_search", "web_fetch":
		return false, "", nil
	case "write_file", "edit_file", "run_command":
		params := make(map[string]interface{})
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return true, "", err
		}

		var desc string
		switch name {
		case "write_file":
			path, _ := params["path"].(string)
			desc = fmt.Sprintf("Write to file: %s", path)
		case "edit_file":
			path, _ := params["path"].(string)
			desc = fmt.Sprintf("Edit file: %s", path)
		case "run_command":
			cmd, _ := params["command"].(string)
			desc = fmt.Sprintf("Execute command: %s", cmd)
		}

		return true, desc, nil
	default:
		return true, "", nil
	}
}

// IsBlocked checks if a command is blocked
func (e *Executor) IsBlocked(command string) bool {
	for _, blocked := range e.blockList {
		if strings.Contains(command, blocked) {
			return true
		}
	}
	return false
}

func (e *Executor) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(e.workDir, path)
}

func (e *Executor) readFile(params map[string]interface{}) *ToolResult {
	path, ok := params["path"].(string)
	if !ok {
		return &ToolResult{Success: false, Error: "path parameter is required"}
	}

	fullPath := e.resolvePath(path)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return &ToolResult{Success: false, Error: err.Error()}
	}

	return &ToolResult{
		Success: true,
		Output:  string(content),
	}
}

func (e *Executor) writeFile(params map[string]interface{}) *ToolResult {
	path, ok := params["path"].(string)
	if !ok {
		return &ToolResult{Success: false, Error: "path parameter is required"}
	}

	content, ok := params["content"].(string)
	if !ok {
		return &ToolResult{Success: false, Error: "content parameter is required"}
	}

	fullPath := e.resolvePath(path)

	// Create directory if needed
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &ToolResult{Success: false, Error: err.Error()}
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return &ToolResult{Success: false, Error: err.Error()}
	}

	return &ToolResult{
		Success: true,
		Output:  fmt.Sprintf("Successfully wrote to %s", path),
	}
}

func (e *Executor) editFile(params map[string]interface{}) *ToolResult {
	path, ok := params["path"].(string)
	if !ok {
		return &ToolResult{Success: false, Error: "path parameter is required"}
	}

	oldString, ok := params["old_string"].(string)
	if !ok {
		return &ToolResult{Success: false, Error: "old_string parameter is required"}
	}

	newString, ok := params["new_string"].(string)
	if !ok {
		return &ToolResult{Success: false, Error: "new_string parameter is required"}
	}

	fullPath := e.resolvePath(path)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return &ToolResult{Success: false, Error: err.Error()}
	}

	// Check if old_string exists
	if !bytes.Contains(content, []byte(oldString)) {
		return &ToolResult{
			Success: false,
			Error:   "old_string not found in file",
		}
	}

	// Check for multiple occurrences
	if bytes.Count(content, []byte(oldString)) > 1 {
		return &ToolResult{
			Success: false,
			Error:   "old_string appears multiple times in file, please provide more context",
		}
	}

	newContent := bytes.Replace(content, []byte(oldString), []byte(newString), 1)
	if err := os.WriteFile(fullPath, newContent, 0644); err != nil {
		return &ToolResult{Success: false, Error: err.Error()}
	}

	return &ToolResult{
		Success: true,
		Output:  fmt.Sprintf("Successfully edited %s", path),
	}
}

func (e *Executor) runCommand(params map[string]interface{}) *ToolResult {
	command, ok := params["command"].(string)
	if !ok {
		return &ToolResult{Success: false, Error: "command parameter is required"}
	}

	// Check block list
	if e.IsBlocked(command) {
		return &ToolResult{
			Success: false,
			Error:   fmt.Sprintf("Command is blocked: %s", command),
		}
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	cmd.Dir = e.workDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return &ToolResult{
			Success: false,
			Output:  string(output),
			Error:   err.Error(),
		}
	}

	return &ToolResult{
		Success: true,
		Output:  string(output),
	}
}

func (e *Executor) listDir(params map[string]interface{}) *ToolResult {
	path, _ := params["path"].(string)
	if path == "" {
		path = "."
	}

	fullPath := e.resolvePath(path)
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return &ToolResult{Success: false, Error: err.Error()}
	}

	var output strings.Builder
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		typeStr := "file"
		if entry.IsDir() {
			typeStr = "dir"
		}

		output.WriteString(fmt.Sprintf("%s\t%s\t%s\n", entry.Name(), typeStr, info.Mode().String()))
	}

	return &ToolResult{
		Success: true,
		Output:  output.String(),
	}
}

func (e *Executor) searchCode(params map[string]interface{}) *ToolResult {
	pattern, ok := params["pattern"].(string)
	if !ok {
		return &ToolResult{Success: false, Error: "pattern parameter is required"}
	}

	path, _ := params["path"].(string)
	if path == "" {
		path = "."
	}

	fullPath := e.resolvePath(path)

	re, err := regexp.Compile(pattern)
	if err != nil {
		return &ToolResult{Success: false, Error: fmt.Sprintf("Invalid regex: %v", err)}
	}

	var output strings.Builder
	err = filepath.WalkDir(fullPath, func(walkPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden directories and common ignore patterns
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip binary files
		content, err := os.ReadFile(walkPath)
		if err != nil {
			return nil
		}

		// Simple binary check
		if bytes.IndexByte(content, 0) != -1 {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if re.MatchString(line) {
				relPath, _ := filepath.Rel(fullPath, walkPath)
				output.WriteString(fmt.Sprintf("%s:%d: %s\n", relPath, i+1, strings.TrimSpace(line)))
			}
		}

		return nil
	})

	if err != nil {
		return &ToolResult{Success: false, Error: err.Error()}
	}

	return &ToolResult{
		Success: true,
		Output:  output.String(),
	}
}

func (e *Executor) glob(params map[string]interface{}) *ToolResult {
	pattern, ok := params["pattern"].(string)
	if !ok {
		return &ToolResult{Success: false, Error: "pattern parameter is required"}
	}

	path, _ := params["path"].(string)
	if path == "" {
		path = "."
	}

	fullPath := e.resolvePath(path)

	matches, err := filepath.Glob(filepath.Join(fullPath, pattern))
	if err != nil {
		return &ToolResult{Success: false, Error: err.Error()}
	}

	var output strings.Builder
	for _, match := range matches {
		relPath, _ := filepath.Rel(fullPath, match)
		output.WriteString(relPath + "\n")
	}

	return &ToolResult{
		Success: true,
		Output:  output.String(),
	}
}

// webSearch performs a web search using DuckDuckGo or similar service
func (e *Executor) webSearch(params map[string]interface{}) *ToolResult {
	query, ok := params["query"].(string)
	if !ok {
		return &ToolResult{Success: false, Error: "query parameter is required"}
	}

	// Use DuckDuckGo Instant Answer API (no API key required)
	escapedQuery := url.QueryEscape(query)
	searchURL := "https://api.duckduckgo.com/?q=" + escapedQuery + "&format=json&no_html=1"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(searchURL)
	if err != nil {
		return &ToolResult{Success: false, Error: fmt.Sprintf("Search failed: %v", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ToolResult{Success: false, Error: fmt.Sprintf("Failed to read response: %v", err)}
	}

	// Parse the response
	var ddgResponse struct {
		AbstractText   string `json:"AbstractText"`
		AbstractSource string `json:"AbstractSource"`
		AbstractURL    string `json:"AbstractURL"`
		Heading        string `json:"Heading"`
		RelatedTopics  []struct {
			Text string `json:"Text"`
			URL  string `json:"FirstURL"`
		} `json:"RelatedTopics"`
		Results []struct {
			Text string `json:"Text"`
			URL  string `json:"FirstURL"`
		} `json:"Results"`
	}

	if err := json.Unmarshal(body, &ddgResponse); err != nil {
		return &ToolResult{Success: false, Error: fmt.Sprintf("Failed to parse response: %v", err)}
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("Search results for: %s\n\n", query))

	if ddgResponse.Heading != "" && ddgResponse.AbstractText != "" {
		output.WriteString(fmt.Sprintf("📚 %s\n%s\nSource: %s\nURL: %s\n\n",
			ddgResponse.Heading,
			ddgResponse.AbstractText,
			ddgResponse.AbstractSource,
			ddgResponse.AbstractURL))
	}

	if len(ddgResponse.RelatedTopics) > 0 {
		output.WriteString("🔗 Related Topics:\n")
		for i, topic := range ddgResponse.RelatedTopics {
			if i >= 5 {
				break
			}
			if topic.Text != "" {
				output.WriteString(fmt.Sprintf("  • %s\n", topic.Text))
				if topic.URL != "" {
					output.WriteString(fmt.Sprintf("    URL: %s\n", topic.URL))
				}
			}
		}
	}

	if output.Len() == 0 {
		output.WriteString("No results found. Try a different search query.")
	}

	return &ToolResult{
		Success: true,
		Output:  output.String(),
	}
}

// webFetch fetches content from a URL
func (e *Executor) webFetch(params map[string]interface{}) *ToolResult {
	targetURL, ok := params["url"].(string)
	if !ok {
		return &ToolResult{Success: false, Error: "url parameter is required"}
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(targetURL)
	if err != nil {
		return &ToolResult{Success: false, Error: fmt.Sprintf("Fetch failed: %v", err)}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ToolResult{Success: false, Error: fmt.Sprintf("Failed to read response: %v", err)}
	}

	// Limit response size
	maxSize := 10000
	content := string(body)
	if len(content) > maxSize {
		content = content[:maxSize] + "\n... (truncated)"
	}

	return &ToolResult{
		Success: true,
		Output:  content,
	}
}
