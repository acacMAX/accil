package tui

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// MarkdownRenderer handles markdown rendering
type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
}

// NewMarkdownRenderer creates a new markdown renderer
func NewMarkdownRenderer() *MarkdownRenderer {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	return &MarkdownRenderer{renderer: r}
}

// Render renders markdown content
func (r *MarkdownRenderer) Render(content string, width int) string {
	// Configure width
	if r.renderer == nil {
		return content
	}

	// Render markdown
	rendered, err := r.renderer.Render(content)
	if err != nil {
		return content
	}

	return strings.TrimSpace(rendered)
}

// CodeBlockStyle styles for code blocks
var codeBlockStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("252")).
	Background(lipgloss.Color("236")).
	Padding(0, 1)

// InlineCodeStyle style for inline code
var inlineCodeStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("203")).
	Background(lipgloss.Color("236")).
	Padding(0, 1)

// HighlightCode highlights code syntax
func HighlightCode(code, language string) string {
	// Simple syntax highlighting for common languages
	switch language {
	case "go":
		return highlightGo(code)
	case "python", "py":
		return highlightPython(code)
	case "javascript", "js", "typescript", "ts":
		return highlightJS(code)
	default:
		return codeBlockStyle.Render(code)
	}
}

func highlightGo(code string) string {
	// Simple Go highlighting
	keywords := []string{"func", "package", "import", "var", "const", "type", "struct", "interface", "if", "else", "for", "range", "return", "go", "defer", "chan", "select", "case", "default", "switch", "break", "continue", "goto", "fallthrough"}
	return highlightKeywords(code, keywords)
}

func highlightPython(code string) string {
	keywords := []string{"def", "class", "import", "from", "if", "else", "elif", "for", "while", "return", "yield", "lambda", "with", "as", "try", "except", "finally", "raise", "assert", "pass", "break", "continue"}
	return highlightKeywords(code, keywords)
}

func highlightJS(code string) string {
	keywords := []string{"function", "const", "let", "var", "class", "if", "else", "for", "while", "return", "import", "export", "from", "async", "await", "try", "catch", "finally", "throw", "new", "this", "super", "extends", "static"}
	return highlightKeywords(code, keywords)
}

func highlightKeywords(code string, keywords []string) string {
	result := code
	keywordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("81")).Bold(true)

	for _, kw := range keywords {
		result = strings.ReplaceAll(result, kw, keywordStyle.Render(kw))
	}
	return result
}
