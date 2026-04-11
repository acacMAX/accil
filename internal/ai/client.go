package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Message represents a chat message
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

// ToolCall represents a tool call from the AI
type ToolCall struct {
	ID       string     `json:"id"`
	Type     string     `json:"type"`
	Function Function   `json:"function"`
}

// Function represents a function call
type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Tool represents a tool definition
type Tool struct {
	Type     string       `json:"type"`
	Function FunctionDef  `json:"function"`
}

// FunctionDef represents a function definition
type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Tools       []Tool    `json:"tools,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a response choice
type Choice struct {
	Index        int      `json:"index"`
	Message      Message  `json:"message"`
	FinishReason string   `json:"finish_reason"`
}

// Usage represents token usage
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamResponse represents a streaming response
type StreamResponse struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// StreamChoice represents a streaming choice
type StreamChoice struct {
	Index        int           `json:"index"`
	Delta        StreamDelta   `json:"delta"`
	FinishReason string        `json:"finish_reason"`
}

// StreamDelta represents streaming delta content
type StreamDelta struct {
	Role      string     `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Client represents an AI client
type Client struct {
	apiKey        string
	baseURL       string
	model         string
	client        *http.Client
	TotalTokens   int // 总用量统计
	PromptTokens  int // 输入token统计
	OutputTokens  int // 输出token统计
	RequestCount  int // 请求次数
}

// NewClient creates a new AI client
func NewClient(apiKey, baseURL, model string) *Client {
	// 支持代理
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		client: &http.Client{
			Timeout:   300 * time.Second, // 增加到5分钟
			Transport: transport,
		},
		TotalTokens:  0,
		PromptTokens: 0,
		OutputTokens: 0,
		RequestCount: 0,
	}
}

// Chat sends a chat completion request with retry
func (c *Client) Chat(messages []Message, tools []Tool) (*ChatResponse, error) {
	req := ChatRequest{
		Model:     c.model,
		Messages:  messages,
		MaxTokens: 4096,
		Tools:     tools,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 重试机制
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second) // 重试间隔
		}

		httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

		resp, err := c.client.Do(httpReq)
		if err != nil {
			lastErr = fmt.Errorf("network error (attempt %d/3): %w", attempt+1, err)
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			resp.Body.Close()
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("API error: %s - %s", resp.Status, string(respBody))
			continue
		}

		var chatResp ChatResponse
		if err := json.Unmarshal(respBody, &chatResp); err != nil {
			lastErr = fmt.Errorf("failed to unmarshal response: %w", err)
			continue
		}

		// 统计用量
		c.TotalTokens += chatResp.Usage.TotalTokens
		c.PromptTokens += chatResp.Usage.PromptTokens
		c.OutputTokens += chatResp.Usage.CompletionTokens
		c.RequestCount++

		return &chatResp, nil
	}

	return nil, fmt.Errorf("请求失败，已重试3次: %w", lastErr)
}

// GetUsageStats 返回用量统计
func (c *Client) GetUsageStats() (totalTokens, promptTokens, outputTokens, requestCount int) {
	return c.TotalTokens, c.PromptTokens, c.OutputTokens, c.RequestCount
}

// ResetUsageStats 重置用量统计
func (c *Client) ResetUsageStats() {
	c.TotalTokens = 0
	c.PromptTokens = 0
	c.OutputTokens = 0
	c.RequestCount = 0
}

// StreamChat sends a streaming chat completion request
func (c *Client) StreamChat(messages []Message, tools []Tool) (*http.Response, error) {
	req := ChatRequest{
		Model:     c.model,
		Messages:  messages,
		MaxTokens: 4096,
		Tools:     tools,
		Stream:    true,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(respBody))
	}

	return resp, nil
}

// ParseStreamResponse parses a streaming response line
func ParseStreamResponse(line string) (*StreamResponse, error) {
	if len(line) < 6 || line[:6] != "data: " {
		return nil, nil
	}

	data := line[6:]
	if data == "[DONE]" {
		return nil, nil
	}

	var resp StreamResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stream response: %w", err)
	}

	return &resp, nil
}

// GetDefaultTools returns the default tools for the AI
func GetDefaultTools() []Tool {
	return []Tool{
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "read_file",
				Description: "Read the contents of a file",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "The path to the file to read",
						},
					},
					"required": []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "write_file",
				Description: "Write content to a file, creating it if it doesn't exist",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "The path to the file to write",
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "The content to write to the file",
						},
					},
					"required": []string{"path", "content"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "edit_file",
				Description: "Edit a file by replacing exact text matches",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "The path to the file to edit",
						},
						"old_string": map[string]interface{}{
							"type":        "string",
							"description": "The exact text to find and replace",
						},
						"new_string": map[string]interface{}{
							"type":        "string",
							"description": "The new text to replace it with",
						},
					},
					"required": []string{"path", "old_string", "new_string"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "run_command",
				Description: "Execute a shell command",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"command": map[string]interface{}{
							"type":        "string",
							"description": "The shell command to execute",
						},
					},
					"required": []string{"command"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "list_dir",
				Description: "List the contents of a directory",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "The path to the directory to list",
						},
					},
					"required": []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "search_code",
				Description: "Search for code in the project using regex",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"pattern": map[string]interface{}{
							"type":        "string",
							"description": "The regex pattern to search for",
						},
						"path": map[string]interface{}{
							"type":        "string",
							"description": "The path to search in (optional, defaults to current directory)",
						},
					},
					"required": []string{"pattern"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "glob",
				Description: "Find files matching a glob pattern",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"pattern": map[string]interface{}{
							"type":        "string",
							"description": "The glob pattern to match files",
						},
						"path": map[string]interface{}{
							"type":        "string",
							"description": "The path to search in (optional)",
						},
					},
					"required": []string{"pattern"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "web_search",
				Description: "Search the web for information using DuckDuckGo. Returns relevant search results.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "The search query",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: FunctionDef{
				Name:        "web_fetch",
				Description: "Fetch content from a URL. Use this to read web pages or API responses.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"url": map[string]interface{}{
							"type":        "string",
							"description": "The URL to fetch content from",
						},
					},
					"required": []string{"url"},
				},
			},
		},
	}
}
