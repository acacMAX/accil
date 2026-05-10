package remote

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/accil/accil/internal/tools"
)

// RemoteExecutor executes tools on remote server
type RemoteExecutor struct {
	client    *Client
	blockList []string
}

// NewRemoteExecutor creates a new remote executor
func NewRemoteExecutor(client *Client, blockList []string) *RemoteExecutor {
	return &RemoteExecutor{client: client, blockList: blockList}
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
	if err := tools.ParseArgumentsInto(args, &params); err != nil {
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
	if err := tools.ParseArgumentsInto(args, &params); err != nil {
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
	if err := tools.ParseArgumentsInto(args, &params); err != nil {
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
	if err := tools.ParseArgumentsInto(args, &params); err != nil {
		return &tools.ToolResult{Success: false, Error: err.Error()}
	}

	if e.isBlocked(params.Command) {
		return &tools.ToolResult{
			Success: false,
			Error:   fmt.Sprintf("Command is blocked: %s", params.Command),
		}
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
	if err := tools.ParseArgumentsInto(args, &params); err != nil {
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
	if err := tools.ParseArgumentsInto(args, &params); err != nil {
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
	if err := tools.ParseArgumentsInto(args, &params); err != nil {
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
		if err := tools.ParseArgumentsInto(arguments, &params); err != nil {
			return false, "", err
		}
		return true, fmt.Sprintf("%s %s", toolName, params.Path), nil
	case "run_command":
		var params struct {
			Command string `json:"command"`
		}
		if err := tools.ParseArgumentsInto(arguments, &params); err != nil {
			return false, "", err
		}
		if e.isBlocked(params.Command) {
			return true, fmt.Sprintf("blocked command: %s", params.Command), nil
		}

		// Extra "dangerous" heuristics (regex-based).
		dangerousPatterns := []string{
			`(?i)\brm\s+-rf\b`,
			`(?i)\bmkfs\b`,
			`(?i)\bdd\s+if=`,
			`(?i)\b:\s*\(\)\s*\{\s*:\|\:&\s*\};:\b`, // fork bomb
			`(?i)\b(curl|wget)\b[^|]*\|\s*(sh|bash)\b`,
			`(?i)>\s*/dev/sd[a-z]\b`,
		}
		for _, pattern := range dangerousPatterns {
			if matched, _ := matchRegex(params.Command, pattern); matched {
				return true, fmt.Sprintf("potentially dangerous command: %s", params.Command), nil
			}
		}

		return true, fmt.Sprintf("execute: %s", params.Command), nil
	}
	return false, "", nil
}

func matchRegex(s, pattern string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchString(s), nil
}

func (e *RemoteExecutor) isBlocked(command string) bool {
	for _, blocked := range e.blockList {
		if blocked == "" {
			continue
		}
		if strings.Contains(command, blocked) {
			return true
		}
	}
	return false
}

// GetClient returns the underlying SSH client
func (e *RemoteExecutor) GetClient() *Client {
	return e.client
}
