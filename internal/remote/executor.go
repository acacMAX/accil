package remote

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/accil/accil/internal/tools"
)

// RemoteExecutor executes tools on remote server
type RemoteExecutor struct {
	client *Client
}

// NewRemoteExecutor creates a new remote executor
func NewRemoteExecutor(client *Client) *RemoteExecutor {
	return &RemoteExecutor{client: client}
}

// Execute executes a tool on the remote server
func (e *RemoteExecutor) Execute(toolName string, arguments string) *tools.ToolResult {
	switch toolName {
	case "read_file":
		return e.readFile(arguments)
	case "write_file":
		return e.writeFile(arguments)
	case "edit_file":
		return e.editFile(arguments)
	case "run_command":
		return e.runCommand(arguments)
	case "list_dir":
		return e.listDir(arguments)
	case "search_code":
		return e.searchCode(arguments)
	case "glob":
		return e.glob(arguments)
	default:
		return &tools.ToolResult{
			Success: false,
			Error:   fmt.Sprintf("unknown tool: %s", toolName),
		}
	}
}

func (e *RemoteExecutor) readFile(args string) *tools.ToolResult {
	var params struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	content, err := e.client.ReadFile(params.Path)
	if err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	return &tools.ToolResult{Success: true, Output: content}
}

func (e *RemoteExecutor) writeFile(args string) *tools.ToolResult {
	var params struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	if err := e.client.WriteFile(params.Path, params.Content); err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	return &tools.ToolResult{Success: true, Output: fmt.Sprintf("File written: %s", params.Path)}
}

func (e *RemoteExecutor) editFile(args string) *tools.ToolResult {
	var params struct {
		Path      string `json:"path"`
		OldString string `json:"old_string"`
		NewString string `json:"new_string"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	if err := e.client.EditFile(params.Path, params.OldString, params.NewString); err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	return &tools.ToolResult{Success: true, Output: fmt.Sprintf("File edited: %s", params.Path)}
}

func (e *RemoteExecutor) runCommand(args string) *tools.ToolResult {
	var params struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	stdout, stderr, err := e.client.Execute(params.Command)
	result := &tools.ToolResult{
		Success: err == nil,
		Output:  stdout,
		Error:   stderr,
	}
	if err != nil {
		result.Error = fmt.Sprintf("%s: %v", stderr, err)
	}
	return result
}

func (e *RemoteExecutor) listDir(args string) *tools.ToolResult {
	var params struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	if params.Path == "" {
		params.Path = "."
	}

	content, err := e.client.ListDir(params.Path)
	if err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	return &tools.ToolResult{Success: true, Output: content}
}

func (e *RemoteExecutor) searchCode(args string) *tools.ToolResult {
	var params struct {
		Pattern string `json:"pattern"`
		Path    string `json:"path"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	content, err := e.client.SearchCode(params.Pattern)
	if err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	return &tools.ToolResult{Success: true, Output: content}
}

func (e *RemoteExecutor) glob(args string) *tools.ToolResult {
	var params struct {
		Pattern string `json:"pattern"`
		Path    string `json:"path"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	content, err := e.client.Glob(params.Pattern)
	if err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	return &tools.ToolResult{Success: true, Output: content}
}

// NeedsConfirmation checks if a tool needs user confirmation
func (e *RemoteExecutor) NeedsConfirmation(toolName string, arguments string) (bool, string, error) {
	// Same as local executor
	switch toolName {
	case "write_file", "edit_file":
		var params struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return false, "", err
		}
		return true, fmt.Sprintf("%s %s", toolName, params.Path), nil
	case "run_command":
		var params struct {
			Command string `json:"command"`
		}
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return false, "", err
		}
		// Check for dangerous commands
		dangerous := []string{"rm -rf", "> /dev", "mkfs", "dd if", "curl.*|.*sh", "wget.*|.*sh"}
		for _, pattern := range dangerous {
			if matched, _ := matchPattern(params.Command, pattern); matched {
				return true, fmt.Sprintf("potentially dangerous command: %s", params.Command), nil
			}
		}
		return true, fmt.Sprintf("execute: %s", params.Command), nil
	}
	return false, "", nil
}

func matchPattern(s, pattern string) (bool, error) {
	// Simple pattern matching
	return strings.Contains(s, strings.TrimSuffix(pattern, ".*")), nil
}

// GetClient returns the underlying SSH client
func (e *RemoteExecutor) GetClient() *Client {
	return e.client
}
