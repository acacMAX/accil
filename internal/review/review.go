package review

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/accil/accil/internal/ai"
	"github.com/accil/accil/internal/tools"
)

// Severity represents issue severity
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// Category represents issue category
type Category string

const (
	CategorySecurity  Category = "security"
	CategoryPerformance Category = "performance"
	CategoryStyle     Category = "style"
	CategoryBug       Category = "bug"
	CategoryDesign    Category = "design"
	CategoryTest      Category = "test"
	CategoryDoc       Category = "documentation"
)

// Issue represents a code review issue
type Issue struct {
	File     string   `json:"file"`
	Line     int      `json:"line,omitempty"`
	Column   int      `json:"column,omitempty"`
	Severity Severity `json:"severity"`
	Category Category `json:"category"`
	Message  string   `json:"message"`
	Suggestion string  `json:"suggestion,omitempty"`
	Code     string   `json:"code,omitempty"`
}

// Report represents a code review report
type Report struct {
	Files     []string `json:"files"`
	Issues    []Issue  `json:"issues"`
	Summary   string   `json:"summary"`
	Score     int      `json:"score"` // 0-100
	Duration  string   `json:"duration"`
}

// Reviewer performs code reviews
type Reviewer struct {
	client   *ai.Client
	executor *tools.Executor
}

// NewReviewer creates a new code reviewer
func NewReviewer(client *ai.Client, executor *tools.Executor) *Reviewer {
	return &Reviewer{
		client:   client,
		executor: executor,
	}
}

// ReviewFile reviews a single file
func (r *Reviewer) ReviewFile(ctx context.Context, filePath string) (*Report, error) {
	// Read file content
	result := r.executor.Execute("read_file", fmt.Sprintf(`{"path": "%s"}`, filePath))
	if !result.Success {
		return nil, fmt.Errorf("failed to read file: %s", result.Error)
	}

	content := result.Output

	// Perform AI review
	reviewPrompt := fmt.Sprintf(`Review the following code file and identify issues.

File: %s

Code:
%s

Analyze for:
1. Security vulnerabilities (SQL injection, XSS, hardcoded secrets, etc.)
2. Bugs and potential errors
3. Performance issues
4. Code style and best practices
5. Design issues
6. Missing tests or documentation

For each issue, provide:
- Line number (approximate if unclear)
- Severity: critical, high, medium, low, info
- Category: security, performance, style, bug, design, test, documentation
- Description of the issue
- Suggested fix

Format your response as JSON:
{
  "issues": [
    {"line": 10, "severity": "high", "category": "security", "message": "...", "suggestion": "..."}
  ],
  "summary": "Overall assessment",
  "score": 75
}`, filePath, content)

	messages := []ai.Message{
		{Role: "system", Content: r.getSystemPrompt()},
		{Role: "user", Content: reviewPrompt},
	}

	resp, err := r.client.Chat(messages, nil)
	if err != nil {
		return nil, err
	}

	return r.parseReviewResponse(filePath, resp.Choices[0].Message.Content)
}

// ReviewFiles reviews multiple files
func (r *Reviewer) ReviewFiles(ctx context.Context, filePaths []string) (*Report, error) {
	report := &Report{
		Files:  filePaths,
		Issues: []Issue{},
		Score:  100,
	}

	totalScore := 0
	for _, path := range filePaths {
		fileReport, err := r.ReviewFile(ctx, path)
		if err != nil {
			continue
		}
		report.Issues = append(report.Issues, fileReport.Issues...)
		totalScore += fileReport.Score
	}

	if len(filePaths) > 0 {
		report.Score = totalScore / len(filePaths)
	}

	report.Summary = r.generateSummary(report)

	return report, nil
}

// ReviewChanges reviews git changes
func (r *Reviewer) ReviewChanges(ctx context.Context) (*Report, error) {
	// Get git diff
	result := r.executor.Execute("run_command", `{"command": "git diff HEAD"}`)
	if !result.Success {
		return nil, fmt.Errorf("failed to get git diff: %s", result.Error)
	}

	if result.Output == "" {
		return &Report{
			Summary: "No changes to review",
			Score:   100,
		}, nil
	}

	reviewPrompt := fmt.Sprintf(`Review the following git diff and identify issues.

Diff:
%s

Focus on:
1. Security issues in new/modified code
2. Potential bugs
3. Performance concerns
4. Code quality
5. Breaking changes

Format your response as JSON:
{
  "issues": [
    {"file": "path/to/file.go", "line": 10, "severity": "high", "category": "security", "message": "...", "suggestion": "..."}
  ],
  "summary": "Overall assessment",
  "score": 75
}`, result.Output)

	messages := []ai.Message{
		{Role: "system", Content: r.getSystemPrompt()},
		{Role: "user", Content: reviewPrompt},
	}

	resp, err := r.client.Chat(messages, nil)
	if err != nil {
		return nil, err
	}

	return r.parseReviewResponse("", resp.Choices[0].Message.Content)
}

// ReviewProject performs a comprehensive project review
func (r *Reviewer) ReviewProject(ctx context.Context, workDir string) (*Report, error) {
	// Get project structure
	structureResult := r.executor.Execute("run_command", `{"command": "find . -type f -name "*.go" -o -name "*.js" -o -name "*.ts" -o -name "*.py" -o -name "*.java" | head -50"}`)
	
	// Get file list
	var files []string
	if structureResult.Success {
		files = strings.Split(strings.TrimSpace(structureResult.Output), "\n")
	}

	// Review each file
	report := &Report{
		Files:  files,
		Issues: []Issue{},
	}

	// Sample review of key files
	keyFiles := r.identifyKeyFiles(files)
	for _, file := range keyFiles {
		fileReport, err := r.ReviewFile(ctx, file)
		if err != nil {
			continue
		}
		report.Issues = append(report.Issues, fileReport.Issues...)
	}

	report.Score = r.calculateScore(report)
	report.Summary = r.generateSummary(report)

	return report, nil
}

// identifyKeyFiles identifies important files to review
func (r *Reviewer) identifyKeyFiles(files []string) []string {
	var keyFiles []string
	priorityPatterns := []string{
		"main.", "handler", "controller", "service", "model", "auth", "api",
	}

	for _, file := range files {
		fileLower := strings.ToLower(file)
		for _, pattern := range priorityPatterns {
			if strings.Contains(fileLower, pattern) {
				keyFiles = append(keyFiles, file)
				break
			}
		}
		if len(keyFiles) >= 10 {
			break
		}
	}

	return keyFiles
}

// calculateScore calculates an overall score
func (r *Reviewer) calculateScore(report *Report) int {
	score := 100
	for _, issue := range report.Issues {
		switch issue.Severity {
		case SeverityCritical:
			score -= 20
		case SeverityHigh:
			score -= 10
		case SeverityMedium:
			score -= 5
		case SeverityLow:
			score -= 2
		}
	}
	if score < 0 {
		score = 0
	}
	return score
}

// generateSummary generates a summary of the review
func (r *Reviewer) generateSummary(report *Report) string {
	issueCounts := make(map[Severity]int)
	categoryCounts := make(map[Category]int)

	for _, issue := range report.Issues {
		issueCounts[issue.Severity]++
		categoryCounts[issue.Category]++
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d issues across %d files.\n", len(report.Issues), len(report.Files)))

	if critical := issueCounts[SeverityCritical]; critical > 0 {
		sb.WriteString(fmt.Sprintf("- %d critical issues require immediate attention\n", critical))
	}
	if high := issueCounts[SeverityHigh]; high > 0 {
		sb.WriteString(fmt.Sprintf("- %d high priority issues\n", high))
	}
	if security := categoryCounts[CategorySecurity]; security > 0 {
		sb.WriteString(fmt.Sprintf("- %d security concerns identified\n", security))
	}

	sb.WriteString(fmt.Sprintf("\nOverall code health score: %d/100", report.Score))

	return sb.String()
}

// parseReviewResponse parses the AI review response
func (r *Reviewer) parseReviewResponse(filePath, response string) (*Report, error) {
	report := &Report{
		Files:  []string{filePath},
		Issues: []Issue{},
		Score:  100,
	}

	// Try to extract JSON
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")
	if jsonStart == -1 || jsonEnd == -1 {
		report.Summary = response
		return report, nil
	}

	jsonContent := response[jsonStart : jsonEnd+1]

	var parsed struct {
		Issues []struct {
			File       string `json:"file"`
			Line       int    `json:"line"`
			Severity   string `json:"severity"`
			Category   string `json:"category"`
			Message    string `json:"message"`
			Suggestion string `json:"suggestion"`
		} `json:"issues"`
		Summary string `json:"summary"`
		Score   int    `json:"score"`
	}

	if err := parseJSON(jsonContent, &parsed); err != nil {
		report.Summary = response
		return report, nil
	}

	for _, issue := range parsed.Issues {
		file := issue.File
		if file == "" {
			file = filePath
		}
		report.Issues = append(report.Issues, Issue{
			File:       file,
			Line:       issue.Line,
			Severity:   Severity(issue.Severity),
			Category:   Category(issue.Category),
			Message:    issue.Message,
			Suggestion: issue.Suggestion,
		})
	}

	report.Summary = parsed.Summary
	report.Score = parsed.Score
	if report.Score == 0 {
		report.Score = 100 - len(report.Issues)*5
		if report.Score < 0 {
			report.Score = 0
		}
	}

	return report, nil
}

func parseJSON(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func (r *Reviewer) getSystemPrompt() string {
	return `You are an expert code reviewer with deep knowledge in:
- Security vulnerabilities (OWASP Top 10)
- Performance optimization
- Clean code principles
- Design patterns
- Language-specific best practices

Provide thorough but constructive feedback. Always suggest specific fixes.`
}

// FormatReport formats a report for display
func FormatReport(report *Report) string {
	var sb strings.Builder

	sb.WriteString("╔══════════════════════════════════════════════════════════════╗\n")
	sb.WriteString("║                    CODE REVIEW REPORT                        ║\n")
	sb.WriteString("╚══════════════════════════════════════════════════════════════╝\n\n")

	if len(report.Files) > 0 {
		sb.WriteString("Files Reviewed:\n")
		for _, f := range report.Files {
			sb.WriteString(fmt.Sprintf("  • %s\n", f))
		}
		sb.WriteString("\n")
	}

	if len(report.Issues) > 0 {
		sb.WriteString("Issues Found:\n")
		sb.WriteString("─────────────\n")
		for i, issue := range report.Issues {
			sb.WriteString(fmt.Sprintf("\n%d. [%s][%s] %s\n", i+1, 
				strings.ToUpper(string(issue.Severity)), 
				strings.ToUpper(string(issue.Category)), 
				issue.Message))
			if issue.File != "" {
				sb.WriteString(fmt.Sprintf("   File: %s", issue.File))
				if issue.Line > 0 {
					sb.WriteString(fmt.Sprintf(":%d", issue.Line))
				}
				sb.WriteString("\n")
			}
			if issue.Suggestion != "" {
				sb.WriteString(fmt.Sprintf("   Suggestion: %s\n", issue.Suggestion))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("─────────────\n")
	sb.WriteString(fmt.Sprintf("Score: %d/100\n\n", report.Score))
	sb.WriteString(report.Summary)

	return sb.String()
}
