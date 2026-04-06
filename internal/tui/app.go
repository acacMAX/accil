package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// 样式定义
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Padding(0, 2)

	modeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			Background(lipgloss.Color("22")).
			Padding(0, 1)

	questModeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("226")).
			Background(lipgloss.Color("130")).
			Padding(0, 1)

	reviewModeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("63")).
			Padding(0, 1)

	userMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true)

	assistantMsgStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("141")).
				Bold(true)

	systemMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true)

	toolMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("51"))

	errorMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	successMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("82"))

	processingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226"))

	inputPromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("62")).
			Bold(true)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))

	statusBarStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("241")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// Mode 运行模式
type Mode string

const (
	ModeChat   Mode = "chat"
	ModeQuest  Mode = "quest"
	ModeReview Mode = "review"
	ModeAgent  Mode = "agent"
)

// DisplayMessage 显示的消息
type DisplayMessage struct {
	Role      string
	Content   string
}

// Model TUI模型
type Model struct {
	Messages      []DisplayMessage
	viewport      viewport.Model
	Input         textinput.Model
	ready         bool
	width         int
	height        int
	Err           error
	Mode          Mode
	AwaitingConfirm bool
	ConfirmDesc     string
	ConfirmCallback func(bool)
	IsStreaming    bool
	history        []string
	historyIndex   int
	renderer       *MarkdownRenderer
	ModelName      string
	Provider       string
	QuestStatus    string
	QuestProgress  string
	CurrentAgent   string
	ProcessingMsg  string // 当前处理信息
}

// NewModel 创建新的TUI模型
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "输入消息或 /help..."
	ti.Prompt = "│ "
	ti.Focus()
	ti.Width = 40

	vp := viewport.New(40, 10)

	return Model{
		Messages:     []DisplayMessage{},
		Input:        ti,
		viewport:     vp,
		ready:        false,
		history:      []string{},
		historyIndex: -1,
		renderer:     NewMarkdownRenderer(),
		Mode:         ModeChat,
		ModelName:    "unknown",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnableMouseCellMotion)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlL:
			m.Messages = []DisplayMessage{}
			m.updateViewport()
			return m, nil
		case tea.KeyUp:
			if !m.Input.Focused() {
				m.viewport.LineUp(1)
				return m, nil
			}
			if len(m.history) > 0 && m.historyIndex < len(m.history)-1 {
				m.historyIndex++
				m.Input.SetValue(m.history[len(m.history)-1-m.historyIndex])
			}
			return m, nil
		case tea.KeyDown:
			if !m.Input.Focused() {
				m.viewport.LineDown(1)
				return m, nil
			}
			if m.historyIndex > 0 {
				m.historyIndex--
				m.Input.SetValue(m.history[len(m.history)-1-m.historyIndex])
			} else if m.historyIndex == 0 {
				m.historyIndex = -1
				m.Input.SetValue("")
			}
			return m, nil
		case tea.KeyPgUp:
			m.viewport.HalfViewUp()
			return m, nil
		case tea.KeyPgDown:
			m.viewport.HalfViewDown()
			return m, nil
		case tea.KeyEnter:
			if m.AwaitingConfirm {
				return m, nil
			}
			content := strings.TrimSpace(m.Input.Value())
			if content == "" {
				return m, nil
			}

			if strings.HasPrefix(content, "/") {
				return m.handleSlashCommand(content)
			}

			m.history = append(m.history, content)
			m.historyIndex = -1

			m.Messages = append(m.Messages, DisplayMessage{
				Role:    "user",
				Content: content,
			})
			m.updateViewport()

			m.Input.SetValue("")
			return m, m.sendMessage(content)
		case tea.KeyRunes:
			if m.AwaitingConfirm {
				char := string(msg.Runes)
				if char == "y" || char == "Y" {
					if m.ConfirmCallback != nil {
						m.ConfirmCallback(true)
					}
					m.AwaitingConfirm = false
					return m, nil
				} else if char == "n" || char == "N" {
					if m.ConfirmCallback != nil {
						m.ConfirmCallback(false)
					}
					m.AwaitingConfirm = false
					return m, nil
				}
			}
		}

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.viewport.LineUp(3)
		case tea.MouseWheelDown:
			m.viewport.LineDown(3)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		viewportHeight := m.height - 10
		if viewportHeight < 5 {
			viewportHeight = 5
		}
		
		viewportWidth := m.width - 2
		if viewportWidth < 10 {
			viewportWidth = 10
		}
		inputWidth := m.width - 8
		if inputWidth < 10 {
			inputWidth = 10
		}
		
		m.viewport.Width = viewportWidth
		m.viewport.Height = viewportHeight
		m.Input.Width = inputWidth
		m.ready = true
		m.updateViewport()

	case ErrorMessage:
		m.Err = msg.Error
		m.IsStreaming = false
		m.ProcessingMsg = ""
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "error",
			Content: msg.Error.Error(),
		})
		m.updateViewport()

	case AssistantMessage:
		m.IsStreaming = false
		m.ProcessingMsg = ""
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "assistant",
			Content: msg.Content,
		})
		m.updateViewport()
		m.viewport.GotoBottom()

	case ToolCallMessage:
		m.ProcessingMsg = fmt.Sprintf("🔧 执行: %s(%s)", msg.Tool, truncateStr(msg.Args, 30))
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "tool",
			Content: m.ProcessingMsg,
		})
		m.updateViewport()
		m.viewport.GotoBottom()

	case ToolResultMessage:
		m.ProcessingMsg = ""
		var resultMsg string
		if msg.Success {
			resultMsg = fmt.Sprintf("✅ 完成: %s", truncateStr(msg.Result, 100))
		} else {
			resultMsg = fmt.Sprintf("❌ 失败: %s", truncateStr(msg.Result, 100))
		}
		// 更新最后一条消息
		if len(m.Messages) > 0 {
			m.Messages[len(m.Messages)-1].Content += "\n" + resultMsg
		}
		m.updateViewport()

	case ProcessingUpdate:
		m.ProcessingMsg = msg.Message
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "system",
			Content: "⏳ " + msg.Message,
		})
		m.updateViewport()
		m.viewport.GotoBottom()

	case QuestStatusMessage:
		m.QuestStatus = msg.Status
		m.QuestProgress = msg.Progress

	case ModeChangeMessage:
		m.Mode = msg.Mode
	}

	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, vpCmd)

	var inCmd tea.Cmd
	m.Input, inCmd = m.Input.Update(msg)
	cmds = append(cmds, inCmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) updateViewport() {
	var content strings.Builder

	for _, msg := range m.Messages {
		rendered := m.renderMessage(msg)
		content.WriteString(rendered)
		content.WriteString("\n")
	}

	m.viewport.SetContent(content.String())
}

func (m Model) renderMessage(msg DisplayMessage) string {
	var sb strings.Builder

	switch msg.Role {
	case "user":
		sb.WriteString(userMsgStyle.Render("┌─ 用户"))
		sb.WriteString("\n")
		sb.WriteString(m.indentContent(msg.Content, "│ "))
		sb.WriteString("\n└─")
	case "assistant":
		sb.WriteString(assistantMsgStyle.Render("┌─ 助手"))
		sb.WriteString("\n")
		rendered := m.renderer.Render(msg.Content, m.width-8)
		sb.WriteString(m.indentContent(rendered, "│ "))
		sb.WriteString("\n└─")
	case "system":
		sb.WriteString(systemMsgStyle.Render("◆ "))
		sb.WriteString(msg.Content)
	case "tool":
		sb.WriteString(toolMsgStyle.Render("│ "))
		sb.WriteString(msg.Content)
	case "error":
		sb.WriteString(errorMsgStyle.Render("✗ "))
		sb.WriteString(msg.Content)
	case "success":
		sb.WriteString(successMsgStyle.Render("✓ "))
		sb.WriteString(msg.Content)
	default:
		sb.WriteString("  ")
		sb.WriteString(msg.Content)
	}

	return sb.String()
}

func (m Model) View() string {
	if !m.ready {
		return m.renderSplash()
	}

	var sb strings.Builder

	sb.WriteString(m.renderHeader())
	sb.WriteString("\n")
	sb.WriteString(m.renderViewportBorder())
	sb.WriteString("\n")
	sb.WriteString(m.renderStatusBar())
	sb.WriteString("\n")

	if m.AwaitingConfirm {
		sb.WriteString(m.renderConfirmPrompt())
		return sb.String()
	}

	// 显示处理中状态
	if m.ProcessingMsg != "" {
		sb.WriteString(processingStyle.Render("  ⏳ " + m.ProcessingMsg))
		sb.WriteString("\n")
	}

	sb.WriteString(m.renderInput())

	return sb.String()
}

func (m Model) renderSplash() string {
	return `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║   █████╗ ██████╗ ██████╗  ██████╗██╗  ██╗██╗     ███████╗   ║
║  ██╔══██╗██╔══██╗██╔══██╗██╔════╝██║  ██║██║     ██╔════╝   ║
║  ███████║██████╔╝██████╔╝██║     ███████║██║     █████╗     ║
║  ██╔══██║██╔══██╗██╔══██╗██║     ██╔══██║██║     ██╔══╝     ║
║  ██║  ██║██████╔╝██████╔╝╚██████╗██║  ██║███████╗███████╗   ║
║  ╚═╝  ╚═╝╚═════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝   ║
║                                                              ║
║            AI驱动的自主编程助手                               ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝

  正在初始化...
`
}

func (m Model) renderHeader() string {
	left := headerStyle.Render(" ACCIL ")

	var modeIndicator string
	switch m.Mode {
	case ModeChat:
		modeIndicator = modeStyle.Render(" 对话 ")
	case ModeQuest:
		modeIndicator = questModeStyle.Render(" 任务 ")
	case ModeReview:
		modeIndicator = reviewModeStyle.Render(" 审查 ")
	case ModeAgent:
		modeIndicator = modeStyle.Render(" 代理 ")
	}

	modelInfo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(fmt.Sprintf("模型: %s", m.ModelName))

	middle := strings.Repeat(" ", max(0, m.width-len(left)-len(modeIndicator)-len(modelInfo)-4))
	
	return left + middle + modeIndicator + " " + modelInfo
}

func (m Model) renderViewportBorder() string {
	content := m.viewport.View()
	
	scrollInfo := ""
	if m.viewport.TotalLineCount() > m.viewport.Height {
		totalLines := m.viewport.TotalLineCount() - m.viewport.Height
		if totalLines > 0 {
			percentage := float64(m.viewport.YOffset) / float64(totalLines) * 100
			if percentage < 0 {
				percentage = 0
			}
			if percentage > 100 {
				percentage = 100
			}
			scrollInfo = lipgloss.NewStyle().
				Foreground(lipgloss.Color("62")).
				Render(fmt.Sprintf(" [%.0f%% ↑↓滚动]", percentage))
		}
	}

	borderWidth := m.width - 20 - len(scrollInfo)
	if borderWidth < 2 {
		borderWidth = 2
	}
	topBorder := lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Render("╭─ 消息" + scrollInfo + strings.Repeat("─", borderWidth) + "╮")

	bottomBorderWidth := m.width - 2
	if bottomBorderWidth < 2 {
		bottomBorderWidth = 2
	}
	bottomBorder := lipgloss.NewStyle().
		Foreground(lipgloss.Color("62")).
		Render("╰" + strings.Repeat("─", bottomBorderWidth) + "╯")

	return topBorder + "\n" + content + "\n" + bottomBorder
}

func (m Model) renderStatusBar() string {
	left := ""

	if m.Mode == ModeQuest && m.QuestStatus != "" {
		left = fmt.Sprintf("任务: %s", m.QuestProgress)
	}

	if m.CurrentAgent != "" {
		left = fmt.Sprintf("代理: %s", m.CurrentAgent)
	}

	right := "Ctrl+C:退出 | Ctrl+L:清屏 | PgUp/PgDn:翻页 | /help:帮助"

	spacing := m.width - len(left) - len(right) - 2
	if spacing < 0 {
		spacing = 0
	}

	return statusBarStyle.Render(left + strings.Repeat(" ", spacing) + right)
}

func (m Model) renderInput() string {
	inputView := m.Input.View()
	if inputView == "" {
		inputView = "│ "
	}
	return inputPromptStyle.Render("╭─ 输入") + "\n" +
		inputStyle.Render(inputView) + "\n" +
		"╰─► "
}

func (m Model) renderConfirmPrompt() string {
	return lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("226")).
		Padding(1, 2).
		Render(fmt.Sprintf("\n⚠ %s\n\n  确认? (y/n): ", m.ConfirmDesc))
}

func (m Model) handleSlashCommand(content string) (tea.Model, tea.Cmd) {
	cmd := strings.Fields(content)
	if len(cmd) == 0 {
		return m, nil
	}

	switch cmd[0] {
	case "/help", "/帮助":
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "system",
			Content: m.getHelpText(),
		})
		m.updateViewport()
	case "/clear", "/清屏":
		m.Messages = []DisplayMessage{}
		m.updateViewport()
	case "/quit", "/exit", "/退出":
		return m, tea.Quit
	case "/quest", "/任务":
		m.Mode = ModeQuest
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "system",
			Content: "进入任务模式。描述你的目标，我将自主规划和执行。",
		})
		m.updateViewport()
	case "/review", "/审查":
		m.Mode = ModeReview
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "system",
			Content: "进入审查模式。使用 /review file <路径> 或 /review diff 审查代码。",
		})
		m.updateViewport()
	case "/agent", "/代理":
		m.Mode = ModeAgent
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "system",
			Content: "代理模式。可用代理: coder(编码), reviewer(审查), architect(架构), tester(测试), debugger(调试), researcher(研究)\n使用 /agent <类型> <任务> 分配任务。",
		})
		m.updateViewport()
	case "/chat", "/对话":
		m.Mode = ModeChat
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "system",
			Content: "切换到对话模式。",
		})
		m.updateViewport()
	case "/model", "/模型":
		if len(cmd) > 1 {
			m.ModelName = cmd[1]
			m.Messages = append(m.Messages, DisplayMessage{
				Role:    "success",
				Content: fmt.Sprintf("模型已更改为: %s", cmd[1]),
			})
			m.updateViewport()
		}
	case "/config", "/配置":
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "system",
			Content: "请在终端运行 'accil config' 命令来编辑配置。",
		})
		m.updateViewport()
	case "/memory":
		if len(cmd) > 1 && cmd[1] == "init" {
			m.Messages = append(m.Messages, DisplayMessage{
				Role:    "system",
				Content: "请在终端运行 'accil memory init' 命令来初始化项目记忆。",
			})
		} else {
			m.Messages = append(m.Messages, DisplayMessage{
				Role:    "system",
				Content: "用法: /memory init - 初始化项目记忆 (AGENTS.md)",
			})
		}
		m.updateViewport()
	case "/context", "/上下文":
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "system",
			Content: fmt.Sprintf("当前上下文:\n  工作目录: %s\n  模型: %s\n  模式: %s", ".", m.ModelName, m.Mode),
		})
		m.updateViewport()
	default:
		m.Messages = append(m.Messages, DisplayMessage{
			Role:    "system",
			Content: fmt.Sprintf("未知命令: %s。输入 /help 查看可用命令。", cmd[0]),
		})
		m.updateViewport()
	}

	return m, nil
}

func (m Model) getHelpText() string {
	return `
╭─────────────────────────────────────────────────────────────────╮
│                        ACCIL 命令帮助                            │
├─────────────────────────────────────────────────────────────────┤
│  /help, /帮助     显示此帮助信息                                  │
│  /clear, /清屏    清除对话历史                                    │
│  /quit, /退出     退出 ACCIL                                      │
│                                                                 │
│  模式切换:                                                        │
│  /chat, /对话     切换到对话模式 (默认)                            │
│  /quest, /任务    进入任务模式 (自主编程)                          │
│  /review, /审查   进入审查模式 (代码审查)                          │
│  /agent, /代理    进入代理模式 (专业子代理)                        │
│                                                                 │
│  设置:                                                           │
│  /model <名称>    更改 AI 模型                                    │
│  /config          打开配置编辑器                                  │
│                                                                 │
│  项目:                                                           │
│  /memory init     初始化项目记忆 (AGENTS.md)                      │
│  /context         显示当前上下文                                  │
│                                                                 │
│  快捷键:                                                          │
│  Ctrl+C           退出                                           │
│  Ctrl+L           清屏                                           │
│  ↑/↓              浏览历史 / 滚动消息                             │
│  PgUp/PgDn        翻页                                           │
│  鼠标滚轮          滚动消息                                       │
╰─────────────────────────────────────────────────────────────────╯`
}

func (m Model) sendMessage(content string) tea.Cmd {
	return func() tea.Msg {
		return UserMessage{Content: content}
	}
}

func (m Model) indentContent(content, prefix string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// 消息类型
type UserMessage struct {
	Content string
}

type ErrorMessage struct {
	Error error
}

type AssistantMessage struct {
	Content string
}

type ToolCallMessage struct {
	Tool string
	Args string
}

type ToolResultMessage struct {
	Success bool
	Result  string
}

type QuestStatusMessage struct {
	Status   string
	Progress string
}

type ModeChangeMessage struct {
	Mode Mode
}

// ProcessingUpdate 处理更新消息
type ProcessingUpdate struct {
	Message string
}

// SetModelName 设置模型名称
func (m *Model) SetModelName(name string) {
	m.ModelName = name
}

// SetProvider 设置提供商名称
func (m *Model) SetProvider(provider string) {
	m.Provider = provider
}

// AddMessage 添加消息
func (m *Model) AddMessage(role, content string) {
	m.Messages = append(m.Messages, DisplayMessage{
		Role:    role,
		Content: content,
	})
	m.updateViewport()
}

// SetMode 设置模式
func (m *Model) SetMode(mode Mode) {
	m.Mode = mode
}

// ShowConfirm 显示确认对话框
func (m *Model) ShowConfirm(desc string, callback func(bool)) {
	m.AwaitingConfirm = true
	m.ConfirmDesc = desc
	m.ConfirmCallback = callback
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
