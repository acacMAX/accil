package quest

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/accil/accil/internal/ai"
	"github.com/accil/accil/internal/tools"
)

// Status represents quest status
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusPaused    Status = "paused"
)

// Step represents a single step in the quest
type Step struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Status      Status `json:"status"`
	Tool        string `json:"tool,omitempty"`
	Arguments   string `json:"arguments,omitempty"`
	Result      string `json:"result,omitempty"`
	Error       string `json:"error,omitempty"`
}

// Quest represents an autonomous programming task
type Quest struct {
	ID          string     `json:"id"`
	Goal        string     `json:"goal"`
	Status      Status     `json:"status"`
	Steps       []Step     `json:"steps"`
	CurrentStep int        `json:"current_step"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Context     string     `json:"context,omitempty"`
}

// Planner plans and executes autonomous tasks
type Planner struct {
	client   *ai.Client
	executor *tools.Executor
	maxSteps int
}

// NewPlanner creates a new quest planner
func NewPlanner(client *ai.Client, executor *tools.Executor) *Planner {
	return &Planner{
		client:   client,
		executor: executor,
		maxSteps: 20,
	}
}

// CreateQuest creates a new quest from a goal
func (p *Planner) CreateQuest(goal string) *Quest {
	return &Quest{
		ID:        fmt.Sprintf("quest-%d", time.Now().UnixNano()),
		Goal:      goal,
		Status:    StatusPending,
		Steps:     []Step{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Plan generates a plan for the quest
func (p *Planner) Plan(ctx context.Context, quest *Quest) error {
	planningPrompt := fmt.Sprintf(`You are an autonomous coding agent. Analyze the following goal and create a step-by-step plan.

Goal: %s

Create a detailed plan with numbered steps. Each step should be a specific action.
Format your response as a JSON array of steps, where each step has:
- "description": what the step does
- "tool": the tool to use (read_file, write_file, edit_file, run_command, list_dir, search_code, glob)
- "arguments": JSON object with tool arguments

Example format:
[
  {"description": "Read the main.go file to understand the project", "tool": "read_file", "arguments": "{\"path\": \"main.go\"}"},
  {"description": "Search for the function to modify", "tool": "search_code", "arguments": "{\"pattern\": \"func.*Handler\"}"}
]

Only respond with the JSON array, no other text.`, quest.Goal)

	messages := []ai.Message{
		{Role: "system", Content: p.getSystemPrompt()},
		{Role: "user", Content: planningPrompt},
	}

	resp, err := p.client.Chat(messages, nil)
	if err != nil {
		return fmt.Errorf("failed to generate plan: %w", err)
	}

	content := resp.Choices[0].Message.Content

	// Parse the plan
	var steps []Step
	// Try to extract JSON from the response
	jsonStart := strings.Index(content, "[")
	jsonEnd := strings.LastIndex(content, "]")
	if jsonStart != -1 && jsonEnd != -1 {
		jsonContent := content[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonContent), &steps); err != nil {
			// If parsing fails, create a single step with the goal
			steps = []Step{
				{ID: "step-1", Description: quest.Goal, Status: StatusPending},
			}
		}
	} else {
		steps = []Step{
			{ID: "step-1", Description: quest.Goal, Status: StatusPending},
		}
	}

	// Assign IDs
	for i := range steps {
		steps[i].ID = fmt.Sprintf("step-%d", i+1)
		steps[i].Status = StatusPending
	}

	quest.Steps = steps
	quest.Status = StatusPending
	quest.UpdatedAt = time.Now()

	return nil
}

// ExecuteStep executes a single step
func (p *Planner) ExecuteStep(ctx context.Context, quest *Quest, stepIndex int, autoApprove bool, approver func(desc string) bool) error {
	if stepIndex >= len(quest.Steps) {
		return fmt.Errorf("step index out of range")
	}

	step := &quest.Steps[stepIndex]
	step.Status = StatusRunning
	quest.UpdatedAt = time.Now()

	// Check if approval needed
	if step.Tool != "" && !autoApprove {
		needsConfirm, desc, _ := p.executor.NeedsConfirmation(step.Tool, step.Arguments)
		if needsConfirm && approver != nil {
			if !approver(desc) {
				step.Status = StatusFailed
				step.Error = "User declined"
				return fmt.Errorf("user declined operation")
			}
		}
	}

	// Execute the tool
	if step.Tool != "" {
		result := p.executor.Execute(step.Tool, step.Arguments)
		if result.Success {
			step.Status = StatusCompleted
			step.Result = result.Output
		} else {
			step.Status = StatusFailed
			step.Error = result.Error
			step.Result = result.Output
		}
	} else {
		// No tool specified, let AI handle it
		result, err := p.executeWithAI(ctx, quest, step)
		if err != nil {
			step.Status = StatusFailed
			step.Error = err.Error()
		} else {
			step.Status = StatusCompleted
			step.Result = result
		}
	}

	quest.CurrentStep = stepIndex + 1
	quest.UpdatedAt = time.Now()

	return nil
}

// Execute runs the entire quest autonomously
func (p *Planner) Execute(ctx context.Context, quest *Quest, autoApprove bool, approver func(desc string) bool, progress func(step Step, total int)) error {
	quest.Status = StatusRunning

	for i := range quest.Steps {
		select {
		case <-ctx.Done():
			quest.Status = StatusPaused
			return ctx.Err()
		default:
		}

		if progress != nil {
			progress(quest.Steps[i], len(quest.Steps))
		}

		if err := p.ExecuteStep(ctx, quest, i, autoApprove, approver); err != nil {
			quest.Status = StatusFailed
			return err
		}
	}

	// Check if all steps completed
	allCompleted := true
	for _, step := range quest.Steps {
		if step.Status != StatusCompleted {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		quest.Status = StatusCompleted
	} else {
		quest.Status = StatusFailed
	}

	return nil
}

// executeWithAI lets AI execute a step without a specific tool
func (p *Planner) executeWithAI(ctx context.Context, quest *Quest, step *Step) (string, error) {
	prompt := fmt.Sprintf(`You are executing a step in an autonomous coding task.

Goal: %s
Current Step: %s

Previous steps completed:
%s

Execute this step and report the result. If you need to use tools, specify them.`, 
		quest.Goal, 
		step.Description,
		p.getCompletedStepsSummary(quest))

	messages := []ai.Message{
		{Role: "system", Content: p.getSystemPrompt()},
		{Role: "user", Content: prompt},
	}

	resp, err := p.client.Chat(messages, ai.GetDefaultTools())
	if err != nil {
		return "", err
	}

	msg := resp.Choices[0].Message

	// Handle tool calls
	for len(msg.ToolCalls) > 0 {
		for _, tc := range msg.ToolCalls {
			result := p.executor.Execute(tc.Function.Name, tc.Function.Arguments)
			messages = append(messages, ai.Message{
				Role:       "tool",
				Content:    fmt.Sprintf("Result: %s", result.Output),
				ToolCallID: tc.ID,
				Name:       tc.Function.Name,
			})
		}

		resp, err = p.client.Chat(messages, ai.GetDefaultTools())
		if err != nil {
			return "", err
		}
		msg = resp.Choices[0].Message
	}

	return msg.Content, nil
}

func (p *Planner) getSystemPrompt() string {
	return `You are an autonomous coding agent called ACCIL. You can:
- Read and write files
- Execute shell commands
- Search code
- Analyze and modify code

Be methodical and thorough. Always verify your changes work correctly.`
}

func (p *Planner) getCompletedStepsSummary(quest *Quest) string {
	var sb strings.Builder
	for i, step := range quest.Steps {
		if step.Status == StatusCompleted {
			sb.WriteString(fmt.Sprintf("%d. %s - Done\n", i+1, step.Description))
			if step.Result != "" {
				sb.WriteString(fmt.Sprintf("   Result: %s\n", truncate(step.Result, 200)))
			}
		}
	}
	return sb.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// ContinueQuest continues a paused quest
func (p *Planner) ContinueQuest(ctx context.Context, quest *Quest, autoApprove bool, approver func(desc string) bool, progress func(step Step, total int)) error {
	return p.Execute(ctx, quest, autoApprove, approver, progress)
}

// RefinePlan allows AI to refine the plan based on results
func (p *Planner) RefinePlan(ctx context.Context, quest *Quest) error {
	refinementPrompt := fmt.Sprintf(`You are an autonomous coding agent. Review the current progress and refine the plan if needed.

Goal: %s

Completed steps:
%s

Current plan:
%s

Should we adjust the remaining steps? If yes, provide updated steps in JSON array format.
If the plan is still valid, respond with "PLAN_OK".`, 
		quest.Goal,
		p.getCompletedStepsSummary(quest),
		p.getRemainingSteps(quest))

	messages := []ai.Message{
		{Role: "system", Content: p.getSystemPrompt()},
		{Role: "user", Content: refinementPrompt},
	}

	resp, err := p.client.Chat(messages, nil)
	if err != nil {
		return err
	}

	content := resp.Choices[0].Message.Content
	if content == "PLAN_OK" || strings.Contains(content, "PLAN_OK") {
		return nil
	}

	// Try to parse new steps
	jsonStart := strings.Index(content, "[")
	jsonEnd := strings.LastIndex(content, "]")
	if jsonStart != -1 && jsonEnd != -1 {
		var newSteps []Step
		jsonContent := content[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonContent), &newSteps); err == nil {
			// Replace remaining steps
			quest.Steps = append(quest.Steps[:quest.CurrentStep], newSteps...)
			for i := quest.CurrentStep; i < len(quest.Steps); i++ {
				quest.Steps[i].ID = fmt.Sprintf("step-%d", i+1)
				quest.Steps[i].Status = StatusPending
			}
		}
	}

	return nil
}

func (p *Planner) getRemainingSteps(quest *Quest) string {
	var sb strings.Builder
	for i := quest.CurrentStep; i < len(quest.Steps); i++ {
		sb.WriteString(fmt.Sprintf("%d. %s [%s]\n", i+1, quest.Steps[i].Description, quest.Steps[i].Status))
	}
	return sb.String()
}
