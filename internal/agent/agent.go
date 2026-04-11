package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/accil/accil/internal/ai"
	"github.com/accil/accil/internal/tools"
)

// AgentType represents different types of agents
type AgentType string

const (
	AgentGeneral    AgentType = "general"
	AgentCoder      AgentType = "coder"
	AgentReviewer   AgentType = "reviewer"
	AgentArchitect  AgentType = "architect"
	AgentTester     AgentType = "tester"
	AgentDebugger   AgentType = "debugger"
	AgentResearcher AgentType = "researcher"
)

// Agent represents a specialized sub-agent
type Agent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        AgentType `json:"type"`
	Description string    `json:"description"`
	SystemPrompt string   `json:"system_prompt"`
	Tasks       []Task    `json:"tasks"`
}

// Task represents a task assigned to an agent
type Task struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Result      string `json:"result,omitempty"`
	Error       string `json:"error,omitempty"`
}

// Manager manages sub-agents
type Manager struct {
	client   *ai.Client
	executor *tools.Executor
	agents   map[string]*Agent
	mu       sync.RWMutex
}

// NewManager creates a new agent manager
func NewManager(client *ai.Client, executor *tools.Executor) *Manager {
	m := &Manager{
		client:   client,
		executor: executor,
		agents:   make(map[string]*Agent),
	}

	// Initialize default agents
	m.initDefaultAgents()

	return m
}

// initDefaultAgents creates the default specialized agents
func (m *Manager) initDefaultAgents() {
	defaultAgents := []Agent{
		{
			ID:   "coder",
			Name: "Code Generator",
			Type: AgentCoder,
			Description: "Specialized in writing clean, efficient code. " +
				"Follows best practices and coding standards.",
			SystemPrompt: `You are an expert code generator. Your responsibilities:
- Write clean, maintainable, and efficient code
- Follow language-specific best practices
- Include appropriate error handling
- Write self-documenting code with clear variable names
- Consider edge cases and potential bugs`,
		},
		{
			ID:   "reviewer",
			Name: "Code Reviewer",
			Type: AgentReviewer,
			Description: "Reviews code for quality, security, and performance. " +
				"Provides detailed feedback and suggestions.",
			SystemPrompt: `You are an expert code reviewer. Your responsibilities:
- Identify potential bugs and errors
- Check for security vulnerabilities (SQL injection, XSS, etc.)
- Evaluate code performance and suggest optimizations
- Ensure code follows style guides and best practices
- Provide constructive, actionable feedback`,
		},
		{
			ID:   "architect",
			Name: "Software Architect",
			Type: AgentArchitect,
			Description: "Designs software architecture and system structure. " +
				"Ensures scalability and maintainability.",
			SystemPrompt: `You are an expert software architect. Your responsibilities:
- Design scalable and maintainable system architectures
- Choose appropriate design patterns
- Consider performance, security, and reliability
- Create clear documentation and diagrams
- Evaluate trade-offs between different approaches`,
		},
		{
			ID:   "tester",
			Name: "Test Engineer",
			Type: AgentTester,
			Description: "Creates comprehensive test suites and ensures code quality. " +
				"Identifies edge cases and regression risks.",
			SystemPrompt: `You are an expert test engineer. Your responsibilities:
- Write comprehensive unit tests
- Create integration and end-to-end tests
- Identify edge cases and boundary conditions
- Ensure good test coverage
- Write clear test descriptions`,
		},
		{
			ID:   "debugger",
			Name: "Debug Specialist",
			Type: AgentDebugger,
			Description: "Analyzes and fixes bugs. Uses systematic debugging approaches.",
			SystemPrompt: `You are an expert debugger. Your responsibilities:
- Systematically identify root causes of bugs
- Use debugging tools and techniques effectively
- Fix bugs without introducing new issues
- Document the bug and its fix
- Suggest preventive measures`,
		},
		{
			ID:   "researcher",
			Name: "Research Agent",
			Type: AgentResearcher,
			Description: "Researches best practices, libraries, and solutions. " +
				"Provides informed recommendations.",
			SystemPrompt: `You are an expert research agent. Your responsibilities:
- Research and evaluate libraries and frameworks
- Find best practices and design patterns
- Analyze documentation and examples
- Provide well-researched recommendations
- Compare alternatives objectively`,
		},
	}

	for _, agent := range defaultAgents {
		m.agents[agent.ID] = &agent
	}
}

// GetAgent returns an agent by ID
func (m *Manager) GetAgent(id string) (*Agent, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	agent, ok := m.agents[id]
	return agent, ok
}

// ListAgents returns all available agents
func (m *Manager) ListAgents() []*Agent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	agents := make([]*Agent, 0, len(m.agents))
	for _, agent := range m.agents {
		agents = append(agents, agent)
	}
	return agents
}

// AssignTask assigns a task to an agent
func (m *Manager) AssignTask(ctx context.Context, agentID string, task Task, autoApprove bool, approver func(desc string) bool) (*Task, error) {
	agent, ok := m.GetAgent(agentID)
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", agentID)
	}

	task.Status = "running"
	task.ID = fmt.Sprintf("task-%d", len(agent.Tasks)+1)

	messages := []ai.Message{
		{Role: "system", Content: agent.SystemPrompt + "\n\n" + m.getToolInstructions()},
		{Role: "user", Content: task.Description},
	}

	// Execute with tool support - 无限循环直到完成
	for i := 0; ; i++ {
		resp, err := m.client.Chat(messages, ai.GetDefaultTools())
		if err != nil {
			task.Status = "failed"
			task.Error = err.Error()
			return &task, err
		}

		msg := resp.Choices[0].Message
		messages = append(messages, msg)

		if len(msg.ToolCalls) == 0 {
			task.Status = "completed"
			task.Result = msg.Content
			break
		}

		// Execute tools
		for _, tc := range msg.ToolCalls {
			needsConfirm, desc, _ := m.executor.NeedsConfirmation(tc.Function.Name, tc.Function.Arguments)
			if needsConfirm && !autoApprove && approver != nil {
				if !approver(desc) {
					messages = append(messages, ai.Message{
						Role:       "tool",
						Content:    "Operation cancelled by user",
						ToolCallID: tc.ID,
						Name:       tc.Function.Name,
					})
					continue
				}
			}

			result := m.executor.Execute(tc.Function.Name, tc.Function.Arguments)
			messages = append(messages, ai.Message{
				Role:       "tool",
				Content:    formatResult(result),
				ToolCallID: tc.ID,
				Name:       tc.Function.Name,
			})
		}
	}

	agent.Tasks = append(agent.Tasks, task)
	return &task, nil
}

// Collaborate runs multiple agents on a complex task
func (m *Manager) Collaborate(ctx context.Context, task string, agentTypes []AgentType, autoApprove bool, approver func(desc string) bool) (map[string]*Task, error) {
	results := make(map[string]*Task)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, agentType := range agentTypes {
		wg.Add(1)
		go func(at AgentType) {
			defer wg.Done()

			// Find agent by type
			var agent *Agent
			for _, a := range m.ListAgents() {
				if a.Type == at {
					agent = a
					break
				}
			}

			if agent == nil {
				return
			}

			// Adjust task for agent type
			adjustedTask := m.adjustTaskForAgent(task, at)

			result, err := m.AssignTask(ctx, agent.ID, Task{
				Description: adjustedTask,
			}, autoApprove, approver)
			if err != nil {
				return
			}

			mu.Lock()
			results[agent.ID] = result
			mu.Unlock()
		}(agentType)
	}

	wg.Wait()
	return results, nil
}

// adjustTaskForAgent adjusts the task description based on agent type
func (m *Manager) adjustTaskForAgent(task string, agentType AgentType) string {
	switch agentType {
	case AgentReviewer:
		return fmt.Sprintf("Review the following code/task and provide detailed feedback:\n\n%s", task)
	case AgentTester:
		return fmt.Sprintf("Create tests for the following:\n\n%s", task)
	case AgentArchitect:
		return fmt.Sprintf("Design the architecture for:\n\n%s", task)
	case AgentDebugger:
		return fmt.Sprintf("Debug and fix issues in:\n\n%s", task)
	case AgentResearcher:
		return fmt.Sprintf("Research and provide information about:\n\n%s", task)
	default:
		return task
	}
}

func (m *Manager) getToolInstructions() string {
	return `You have access to these tools:
- read_file(path): Read file contents
- write_file(path, content): Write to file
- edit_file(path, old_string, new_string): Edit file
- run_command(command): Execute shell command
- list_dir(path): List directory contents
- search_code(pattern): Search code with regex
- glob(pattern): Find files matching pattern

Use these tools to accomplish your tasks effectively.`
}

func formatResult(result *tools.ToolResult) string {
	if result.Success {
		return result.Output
	}
	return fmt.Sprintf("Error: %s\nOutput: %s", result.Error, result.Output)
}

// CreateCustomAgent creates a new custom agent
func (m *Manager) CreateCustomAgent(id, name, description, systemPrompt string) *Agent {
	m.mu.Lock()
	defer m.mu.Unlock()

	agent := &Agent{
		ID:           id,
		Name:         name,
		Type:         AgentGeneral,
		Description:  description,
		SystemPrompt: systemPrompt,
		Tasks:        []Task{},
	}

	m.agents[id] = agent
	return agent
}

// GetAgentPrompt returns the system prompt for an agent
func (m *Manager) GetAgentPrompt(agentType AgentType) string {
	prompts := map[AgentType]string{
		AgentCoder: `You are a coding specialist. Focus on:
- Writing clean, efficient code
- Following language conventions
- Implementing features correctly
- Handling errors appropriately`,
		AgentReviewer: `You are a code review specialist. Focus on:
- Code quality and maintainability
- Security vulnerabilities
- Performance issues
- Best practices violations`,
		AgentArchitect: `You are a software architecture specialist. Focus on:
- System design and structure
- Scalability considerations
- Design patterns
- Technology choices`,
		AgentTester: `You are a testing specialist. Focus on:
- Unit test coverage
- Integration tests
- Edge cases
- Test reliability`,
		AgentDebugger: `You are a debugging specialist. Focus on:
- Root cause analysis
- Systematic debugging
- Fix verification
- Prevention strategies`,
	}

	if prompt, ok := prompts[agentType]; ok {
		return prompt
	}
	return "You are a helpful AI assistant specialized in software development."
}

// AnalyzeAndAssign analyzes a task and assigns it to the most appropriate agent
func (m *Manager) AnalyzeAndAssign(ctx context.Context, task string) (*Task, *Agent, error) {
	// Use AI to determine the best agent
	analysisPrompt := fmt.Sprintf(`Analyze this task and determine which type of agent is best suited:

Task: %s

Available agent types:
- coder: For writing and modifying code
- reviewer: For reviewing code quality and security
- architect: For system design and architecture decisions
- tester: For creating tests
- debugger: For fixing bugs and errors
- researcher: For research and information gathering

Respond with only the agent type name (one word).`, task)

	messages := []ai.Message{
		{Role: "user", Content: analysisPrompt},
	}

	resp, err := m.client.Chat(messages, nil)
	if err != nil {
		return nil, nil, err
	}

	agentTypeStr := strings.TrimSpace(strings.ToLower(resp.Choices[0].Message.Content))
	var agentType AgentType
	switch agentTypeStr {
	case "coder":
		agentType = AgentCoder
	case "reviewer":
		agentType = AgentReviewer
	case "architect":
		agentType = AgentArchitect
	case "tester":
		agentType = AgentTester
	case "debugger":
		agentType = AgentDebugger
	case "researcher":
		agentType = AgentResearcher
	default:
		agentType = AgentCoder
	}

	// Find the agent
	var agent *Agent
	for _, a := range m.ListAgents() {
		if a.Type == agentType {
			agent = a
			break
		}
	}

	if agent == nil {
		agent = m.agents["coder"]
	}

	// Assign the task
	result, err := m.AssignTask(ctx, agent.ID, Task{Description: task}, false, nil)
	return result, agent, err
}
