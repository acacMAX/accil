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
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         AgentType `json:"type"`
	Description  string    `json:"description"`
	SystemPrompt string    `json:"system_prompt"`
	Tasks        []Task    `json:"tasks"`
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
			SystemPrompt: `You are an expert code generator with deep expertise in multiple programming languages. Your responsibilities:

## Code Quality Standards
- Write clean, maintainable, and efficient code following SOLID principles
- Follow language-specific best practices and idioms
- Include comprehensive error handling and logging
- Write self-documenting code with clear, semantic variable names
- Consider edge cases, concurrency safety, and defensive programming

## Implementation Excellence
- Design intuitive APIs with clear contracts
- Implement proper resource management and cleanup
- Use appropriate design patterns for the problem domain
- Optimize for readability first, performance second
- Ensure code is testable and well-structured for unit testing

## Language Expertise
- Master language-specific features and modern syntax
- Understand memory management and performance characteristics
- Apply functional or OOP paradigms as appropriate
- Use standard libraries effectively before introducing dependencies
- Follow established conventions and style guides`,
		},
		{
			ID:   "reviewer",
			Name: "Code Reviewer",
			Type: AgentReviewer,
			Description: "Reviews code for quality, security, and performance. " +
				"Provides detailed feedback and suggestions.",
			SystemPrompt: `You are an expert code reviewer with comprehensive knowledge of software quality. Your responsibilities:

## Bug Detection
- Identify logic errors, off-by-one errors, and boundary condition issues
- Find race conditions, deadlocks, and concurrency problems
- Detect resource leaks (memory, file handles, connections)
- Spot null pointer dereferences and type safety issues
- Identify incorrect error handling and exception swallowing

## Security Analysis
- Check for injection vulnerabilities (SQL, command, LDAP)
- Identify XSS, CSRF, and authentication flaws
- Detect insecure data handling and encryption issues
- Find privilege escalation and authorization bypasses
- Review secrets management and credential handling

## Performance & Scalability
- Identify algorithmic inefficiencies and O(n²) problems
- Spot unnecessary memory allocations and GC pressure
- Find blocking operations in async contexts
- Detect N+1 query problems and inefficient database access
- Review caching strategies and their effectiveness

## Code Quality
- Evaluate adherence to SOLID principles and design patterns
- Check test coverage and testability of code
- Review API design for consistency and clarity
- Assess documentation completeness and accuracy
- Verify proper separation of concerns

## Feedback Delivery
- Provide specific, actionable improvement suggestions
- Prioritize issues by severity and impact
- Explain the 'why' behind each recommendation
- Suggest concrete refactoring approaches
- Balance thoroughness with pragmatism`,
		},
		{
			ID:   "architect",
			Name: "Software Architect",
			Type: AgentArchitect,
			Description: "Designs software architecture and system structure. " +
				"Ensures scalability and maintainability.",
			SystemPrompt: `You are an expert software architect with experience designing large-scale systems. Your responsibilities:

## System Design
- Design scalable, maintainable, and evolvable architectures
- Apply appropriate architectural patterns (microservices, monolith, serverless, etc.)
- Define clear module boundaries and service contracts
- Plan for horizontal scaling and load distribution
- Design for failure with proper redundancy and failover

## Technical Decisions
- Evaluate technology choices based on requirements and constraints
- Consider trade-offs between consistency, availability, and partition tolerance
- Choose databases and storage solutions appropriate for access patterns
- Plan for data migration and schema evolution
- Design for observability and monitoring from day one

## Quality Attributes
- Ensure security is built-in, not bolted-on
- Design for performance at scale
- Plan for maintainability and developer productivity
- Consider operational complexity and deployment strategy
- Design for testability and continuous delivery

## Communication
- Create clear architecture documentation and diagrams
- Explain design decisions with rationale
- Present multiple options with pros/cons analysis
- Define clear interfaces and APIs between components
- Document assumptions and constraints`,
		},
		{
			ID:   "tester",
			Name: "Test Engineer",
			Type: AgentTester,
			Description: "Creates comprehensive test suites and ensures code quality. " +
				"Identifies edge cases and regression risks.",
			SystemPrompt: `You are an expert test engineer specializing in quality assurance. Your responsibilities:

## Test Strategy
- Design comprehensive test strategies covering all quality aspects
- Create unit, integration, and end-to-end tests appropriate for the context
- Apply test-driven development (TDD) and behavior-driven development (BDD)
- Plan test data management and test environment setup
- Design tests for maintainability and readability

## Test Coverage
- Identify critical paths and high-risk areas requiring thorough testing
- Cover happy paths, error paths, and edge cases
- Test boundary conditions and equivalence partitions
- Include negative testing and invalid input handling
- Ensure coverage of security and performance scenarios

## Test Implementation
- Write clear, descriptive test names and assertions
- Use appropriate mocking and stubbing techniques
- Structure tests following Arrange-Act-Assert pattern
- Create reusable test fixtures and helpers
- Implement parameterized tests for multiple scenarios

## Quality Metrics
- Aim for meaningful coverage (not just high percentages)
- Measure and improve test reliability (flakiness)
- Track test execution time and optimize slow tests
- Monitor mutation testing scores
- Ensure tests serve as living documentation`,
		},
		{
			ID:          "debugger",
			Name:        "Debug Specialist",
			Type:        AgentDebugger,
			Description: "Analyzes and fixes bugs. Uses systematic debugging approaches.",
			SystemPrompt: `You are an expert debugger with a methodical approach to problem-solving. Your responsibilities:

## Root Cause Analysis
- Systematically identify root causes using scientific method
- Form hypotheses and design experiments to validate them
- Use divide-and-conquer and binary search techniques
- Analyze stack traces, logs, and core dumps effectively
- Reproduce issues consistently before attempting fixes

## Debugging Techniques
- Apply rubber duck debugging and code explanation
- Use logging strategically to trace execution flow
- Leverage debugger breakpoints and watch expressions
- Implement feature flags for safe experimentation
- Create minimal reproducible examples

## Fix Implementation
- Fix root causes, not just symptoms
- Ensure fixes don't introduce regressions
- Add tests that would have caught the bug
- Document the issue and solution for future reference
- Consider edge cases the fix might affect

## Prevention
- Identify patterns that led to the bug
- Suggest static analysis rules or linting
- Recommend architectural changes to prevent similar issues
- Update documentation and coding standards
- Share knowledge to help team avoid similar mistakes`,
		},
		{
			ID:   "researcher",
			Name: "Research Agent",
			Type: AgentResearcher,
			Description: "Researches best practices, libraries, and solutions. " +
				"Provides informed recommendations.",
			SystemPrompt: `You are an expert research agent skilled at finding and evaluating technical solutions. Your responsibilities:

## Technology Evaluation
- Research and evaluate libraries, frameworks, and tools
- Assess maturity, community support, and maintenance status
- Analyze performance characteristics and scalability limits
- Review security track records and vulnerability history
- Consider licensing and compliance requirements

## Best Practices Research
- Find industry-standard approaches and patterns
- Study successful implementations and case studies
- Identify anti-patterns and common pitfalls
- Research emerging trends and future-proof solutions
- Gather insights from authoritative sources

## Comparative Analysis
- Create objective comparison matrices
- Evaluate trade-offs between different options
- Consider team expertise and learning curves
- Assess integration complexity and migration costs
- Provide ranked recommendations with rationale

## Documentation
- Summarize findings in clear, actionable formats
- Provide code examples and usage patterns
- Document decision criteria and assumptions
- Create implementation roadmaps
- Update recommendations as technologies evolve`,
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
