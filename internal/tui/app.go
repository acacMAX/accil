package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ═══════════════════════════════════════════════════════════════
//  复古极简终端风格 - Retro Minimal Terminal Aesthetic
// ═══════════════════════════════════════════════════════════════

// 核心配色方案 - 温暖琥珀色调
const (
	colorBg         = "#1a1a1a" // 深灰黑背景
	colorFg         = "#e8e6e3" // 温暖白前景
	colorAmber      = "#ffb347" // 琥珀色强调
	colorAmberDim   = "#cc8a3c" // 暗琥珀色
	colorGreen      = "#98c379" // 柔和绿
	colorRed        = "#e06c75" // 柔和红
	colorBlue       = "#61afef" // 柔和蓝
	colorPurple     = "#c678dd" // 柔和紫
	colorGray       = "#5c6370" // 中灰
	colorGrayLight  = "#abb2bf" // 浅灰
	colorCursor     = "#ffb347" // 光标色
)

// 样式定义 - 极简风格
var (
	// 主容器样式
	containerStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(colorBg))

	// 标题栏 - 简洁线条风格
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAmber)).
		Background(lipgloss.Color(colorBg)).
		Bold(true).
		Padding(0, 1).
		BorderStyle(lipgloss.Border{
			Top:    "─",
			Bottom: "─",
			Left:   "│",
			Right:  "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "├",
			BottomRight: "┤",
		}).
		BorderForeground(lipgloss.Color(colorGray))

	// 模式标签 - 胶囊形状
	modeCapsuleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBg)).
		Background(lipgloss.Color(colorAmber)).
		Bold(true).
		Padding(0, 2).
		Margin(0, 1)

	modeChatStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBg)).
		Background(lipgloss.Color(colorGreen)).
		Bold(true).
		Padding(0, 2).
		Margin(0, 1)

	modeQuestStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBg)).
		Background(lipgloss.Color(colorPurple)).
		Bold(true).
		Padding(0, 2).
		Margin(0, 1)

	modeReviewStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBg)).
		Background(lipgloss.Color(colorBlue)).
		Bold(true).
		Padding(0, 2).
		Margin(0, 1)

	modeRemoteStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBg)).
		Background(lipgloss.Color(colorRed)).
		Bold(true).
		Padding(0, 2).
		Margin(0, 1)

	// 用户信息 - 右对齐标签
	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGrayLight)).
		Italic(true)

	// 消息样式 - 极简边框
	userMsgBoxStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAmber)).
		BorderStyle(lipgloss.Border{
			Left: "▸",
		}).
		BorderForeground(lipgloss.Color(colorAmber)).
		PaddingLeft(1).
		MarginLeft(2)

	userContentStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorFg)).
		MarginLeft(4)

	assistantMsgBoxStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGreen)).
		BorderStyle(lipgloss.Border{
			Left: "◆",
		}).
		BorderForeground(lipgloss.Color(colorGreen)).
		PaddingLeft(1).
		MarginLeft(2)

	assistantContentStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorFg)).
		MarginLeft(4)

	systemMsgStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGrayLight)).
		Italic(true).
		MarginLeft(2)

	toolMsgStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBlue)).
		MarginLeft(4)

	errorMsgStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorRed)).
		Bold(true).
		MarginLeft(2)

	successMsgStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGreen)).
		MarginLeft(2)

	// 处理中动画样式
	processingStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAmber)).
		Bold(true).
		Blink(true)

	// 输入区域 - 底部固定
	inputBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.Border{
			Top:    "─",
			Left:   "│",
			Right:  "│",
			Bottom: "─",
			TopLeft:     "├",
			TopRight:    "┤",
			BottomLeft:  "╰",
			BottomRight: "╯",
		}).
		BorderForeground(lipgloss.Color(colorGray)).
		Padding(0, 1).
		Background(lipgloss.Color(colorBg))

	inputPromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAmber)).
		Bold(true)

	inputCursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorCursor)).
		Background(lipgloss.Color(colorAmberDim))

	// 状态栏 - 极简底部
	statusBarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGray)).
		Background(lipgloss.Color(colorBg)).
		Padding(0, 1)

	statusActiveStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAmber)).
		Bold(true)

	// 滚动条样式
	scrollbarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGray))

	// 帮助面板
	helpBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.Border{
			Top:    "─",
			Bottom: "─",
			Left:   "│",
			Right:  "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "╰",
			BottomRight: "╯",
		}).
		BorderForeground(lipgloss.Color(colorGray)).
		Padding(1, 2).
		Background(lipgloss.Color(colorBg))

	helpTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAmber)).
		Bold(true).
		Underline(true)

	helpKeyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAmber)).
		Bold(true)

	helpDescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGrayLight))

	// 确认对话框
	confirmBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(colorAmber)).
		Padding(1, 3).
		Background(lipgloss.Color(colorBg))

	confirmTitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAmber)).
		Bold(true)

	// 时间戳样式
	timestampStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGray)).
		Italic(true)
)

// Mode 运行模式
type Mode string

const (
	ModeChat   Mode = "chat"
	ModeQuest  Mode = "quest"
	ModeReview Mode = "review"
	ModeAgent  Mode = "agent"
	ModeRemote Mode = "remote"
)

// DisplayMessage 显示的消息
type DisplayMessage struct {
	Role      string
	Content   string
	Timestamp time.Time
}

// Model TUI模型
type Model struct {
	Messages        []DisplayMessage
	viewport        viewport.Model
	Input           textarea.Model
	ready           bool
	width           int
	height          int
	Err             error
	Mode            Mode
	AwaitingConfirm bool
	ConfirmDesc     string
	ConfirmCallback func(bool)
	IsStreaming     bool
	StopRequested   bool
	history         []string
	historyIndex    int
	renderer        *MarkdownRenderer
	ModelName       string
	Provider        string
	QuestStatus     string
	QuestProgress   string
	CurrentAgent    string
	ProcessingMsg   string
	ShowQuitConfirm bool
	QuitConfirmMsg  string
	ShowHelp        bool
	frameCount      int    // 动画帧计数
	RemoteHost      string // 远程服务器地址
	RemoteConnected bool   // 远程连接状态
}

// NewModel 创建新的TUI模型
func NewModel() Model {
	ta := textarea.New()
	ta.Placeholder = "在此输入消息..."
	ta.SetWidth(40)
	ta.SetHeight(3)
	ta.Focus()

	vp := viewport.New(40, 10)

	return Model{
		Messages:     []DisplayMessage{},
		Input:        ta,
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
	return tea.Batch(
		m.Input.Focus(),
		tea.EnableMouseCellMotion,
		m.tick(),
	)
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg{t}
	})
}

type tickMsg struct {
	time.Time
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		m.frameCount++
		return m, m.tick()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			if !m.ShowQuitConfirm {
				m.ShowQuitConfirm = true
				m.QuitConfirmMsg = "再次按下 Ctrl+C 确认退出"
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyEsc:
			if m.ShowHelp {
				m.ShowHelp = false
				return m, nil
			}
			if m.ShowQuitConfirm {
				m.ShowQuitConfirm = false
				m.QuitConfirmMsg = ""
				return m, nil
			}
			if m.IsStreaming || m.ProcessingMsg != "" {
				m.StopRequested = true
				m.ProcessingMsg = ""
				m.Messages = append(m.Messages, DisplayMessage{
					Role:      "system",
					Content:   "◼ 已停止",
					Timestamp: time.Now(),
				})
				m.updateViewport()
				return m, nil
			}

		case tea.KeyCtrlL:
			m.Messages = []DisplayMessage{}
			m.updateViewport()
			return m, nil

		case tea.KeyCtrlH:
			m.ShowHelp = !m.ShowHelp
			return m, nil

		case tea.KeyCtrlS, tea.KeyF5:
			if !m.AwaitingConfirm && !m.ShowHelp {
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
					Role:      "user",
					Content:   content,
					Timestamp: time.Now(),
				})
				m.updateViewport()

				m.Input.SetValue("")
				return m, m.sendMessage(content)
			}

		case tea.KeyPgUp:
			m.viewport.HalfViewUp()
			return m, nil

		case tea.KeyPgDown:
			m.viewport.HalfViewDown()
			return m, nil

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

		viewportWidth := m.width - 4
		if viewportWidth < 10 {
			viewportWidth = 10
		}
		inputWidth := m.width - 6
		if inputWidth < 10 {
			inputWidth = 10
		}

		m.viewport.Width = viewportWidth
		m.viewport.Height = viewportHeight
		m.Input.SetWidth(inputWidth)
		m.Input.SetHeight(3)
		m.ready = true
		m.updateViewport()

	case ErrorMessage:
		m.Err = msg.Error
		m.IsStreaming = false
		m.ProcessingMsg = ""
		m.Messages = append(m.Messages, DisplayMessage{
			Role:      "error",
			Content:   msg.Error.Error(),
			Timestamp: time.Now(),
		})
		m.updateViewport()

	case AssistantMessage:
		m.IsStreaming = false
		m.ProcessingMsg = ""
		m.Messages = append(m.Messages, DisplayMessage{
			Role:      "assistant",
			Content:   msg.Content,
			Timestamp: time.Now(),
		})
		m.updateViewport()
		m.viewport.GotoBottom()

	case ToolCallMessage:
		m.ProcessingMsg = fmt.Sprintf("⟳ %s", msg.Tool)
		m.Messages = append(m.Messages, DisplayMessage{
			Role:      "tool",
			Content:   fmt.Sprintf("⟳ %s %s", msg.Tool, msg.Args),
			Timestamp: time.Now(),
		})
		m.updateViewport()
		m.viewport.GotoBottom()

	case ToolResultMessage:
		m.ProcessingMsg = ""
		var icon string
		if msg.Success {
			icon = "✓"
		} else {
			icon = "✗"
		}
		if len(m.Messages) > 0 {
			lastMsg := &m.Messages[len(m.Messages)-1]
			if lastMsg.Role == "tool" {
				lastMsg.Content = fmt.Sprintf("%s %s", icon, truncateStr(msg.Result, 60))
			}
		}
		m.updateViewport()

	case ProcessingUpdate:
		m.ProcessingMsg = msg.Message
		m.Messages = append(m.Messages, DisplayMessage{
			Role:      "system",
			Content:   m.getSpinner() + " " + msg.Message,
			Timestamp: time.Now(),
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

// getSpinner 返回当前动画帧
func (m Model) getSpinner() string {
	spinners := []string{"◐", "◓", "◑", "◒"}
	return spinners[m.frameCount%len(spinners)]
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
	timestamp := timestampStyle.Render(fmt.Sprintf("[%02d:%02d] ", msg.Timestamp.Hour(), msg.Timestamp.Minute()))

	switch msg.Role {
	case "user":
		return fmt.Sprintf("%s%s\n%s",
			timestamp,
			userMsgBoxStyle.Render("YOU"),
			userContentStyle.Render(msg.Content))

	case "assistant":
		rendered := m.renderer.Render(msg.Content, m.width-8)
		return fmt.Sprintf("%s%s\n%s",
			timestamp,
			assistantMsgBoxStyle.Render("AI"),
			assistantContentStyle.Render(rendered))

	case "system":
		return fmt.Sprintf("%s%s", timestamp, systemMsgStyle.Render(msg.Content))

	case "tool":
		return toolMsgStyle.Render(msg.Content)

	case "error":
		return fmt.Sprintf("%s%s", timestamp, errorMsgStyle.Render("ERROR: "+msg.Content))

	case "success":
		return fmt.Sprintf("%s%s", timestamp, successMsgStyle.Render(msg.Content))

	default:
		return fmt.Sprintf("  %s", msg.Content)
	}
}

func (m Model) View() string {
	if !m.ready {
		return m.renderSplash()
	}

	if m.ShowHelp {
		return m.renderHelp()
	}

	var sb strings.Builder

	// 标题栏
	sb.WriteString(m.renderHeader())
	sb.WriteString("\n")

	// 消息区域
	sb.WriteString(m.renderViewport())
	sb.WriteString("\n")

	// 状态栏
	sb.WriteString(m.renderStatusBar())
	sb.WriteString("\n")

	// 确认对话框
	if m.ShowQuitConfirm {
		sb.WriteString(m.renderQuitConfirm())
		sb.WriteString("\n")
	}

	// 输入区域
	sb.WriteString(m.renderInput())

	return containerStyle.Render(sb.String())
}

func (m Model) renderSplash() string {
	logo := `
    ╭────────────────────────────────────╮
    │                                    │
    │    ▗▄▖ ▗▄▄▖ ▗▄▄▖ ▗▖  ▗▖ ▗▄▖       │
    │   ▐▌ ▐▌▐▌   ▐▌   ▐▌  ▐▌▐▌ ▐▌      │
    │   ▐▛▀▜▌▐▛▀▀▘▐▌   ▐▌  ▐▌▐▛▀▜▌      │
    │   ▐▌ ▐▌▐▙▄▄▖▐▙▄▄▖▐▙▄▄▞▘▐▌ ▐▌      │
    │                                    │
    │         智能编程助手               │
    │                                    │
    ╰────────────────────────────────────╯
`
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorAmber)).
		Render(logo)
}

func (m Model) renderHeader() string {
	// 标题
	title := titleStyle.Render(" ACCIL ")

	// 模式胶囊
	var modeStr string
	switch m.Mode {
	case ModeChat:
		modeStr = modeChatStyle.Render(" CHAT ")
	case ModeQuest:
		modeStr = modeQuestStyle.Render(" QUEST ")
	case ModeReview:
		modeStr = modeReviewStyle.Render(" REVIEW ")
	case ModeAgent:
		modeStr = modeCapsuleStyle.Render(" AGENT ")
	case ModeRemote:
		modeStr = modeRemoteStyle.Render(" REMOTE ")
	}

	// 模型信息或远程连接信息
	var infoStr string
	if m.Mode == ModeRemote && m.RemoteConnected {
		infoStr = infoStyle.Render(fmt.Sprintf("ssh: %s", m.RemoteHost))
	} else {
		infoStr = infoStyle.Render(fmt.Sprintf("model: %s", m.ModelName))
	}

	// 计算间距
	usedWidth := lipgloss.Width(title) + lipgloss.Width(modeStr) + lipgloss.Width(infoStr) + 4
	spacing := m.width - usedWidth
	if spacing < 0 {
		spacing = 0
	}

	return title + modeStr + strings.Repeat(" ", spacing) + infoStr
}

func (m Model) renderViewport() string {
	// 顶部边框
	topBorder := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGray)).
		Render("╭" + strings.Repeat("─", m.viewport.Width) + "╮")

	// 内容区域
	content := m.viewport.View()

	// 底部边框
	bottomBorder := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGray)).
		Render("╰" + strings.Repeat("─", m.viewport.Width) + "╯")

	return topBorder + "\n" + content + "\n" + bottomBorder
}

func (m Model) renderStatusBar() string {
	var left, right string

	if m.Mode == ModeQuest && m.QuestStatus != "" {
		left = fmt.Sprintf("● %s", m.QuestProgress)
	} else if m.CurrentAgent != "" {
		left = fmt.Sprintf("● agent: %s", m.CurrentAgent)
	} else {
		left = fmt.Sprintf("● messages: %d", len(m.Messages))
	}

	if m.IsStreaming || m.ProcessingMsg != "" {
		right = "[ESC]停止 [Ctrl+C]退出 [Ctrl+S]发送 [Ctrl+H]帮助"
	} else {
		right = "[Ctrl+C]退出 [Ctrl+S]发送 [Ctrl+H]帮助 [Ctrl+L]清屏"
	}

	leftStyled := statusActiveStyle.Render(left)
	rightStyled := statusBarStyle.Render(right)

	spacing := m.width - lipgloss.Width(leftStyled) - lipgloss.Width(rightStyled)
	if spacing < 0 {
		spacing = 0
	}

	return leftStyled + strings.Repeat(" ", spacing) + rightStyled
}

func (m Model) renderQuitConfirm() string {
	msg := confirmTitleStyle.Render(m.QuitConfirmMsg)
	return lipgloss.NewStyle().
		MarginLeft(2).
		Render(msg)
}

func (m Model) renderInput() string {
	prompt := inputPromptStyle.Render("❯ ")
	inputView := m.Input.View()

	// 包装输入框
	wrapped := inputBoxStyle.
		Width(m.width - 4).
		Render(prompt + inputView)

	return wrapped
}

func (m Model) renderHelp() string {
	helpText := `
  命令
  ─────
  /help, /?      显示帮助
  /clear, /cls   清除对话
  /quit, /exit   退出程序
  /chat          对话模式
  /quest         任务模式
  /review        审查模式
  /agent         代理模式
  /remote        远程开发模式
  /model <name>  切换模型
  /context       显示上下文

  快捷键
  ──────
  Ctrl+S, F5     发送消息
  Ctrl+C (x2)    退出
  Ctrl+L         清屏
  Ctrl+H         显示/隐藏帮助
  ESC            停止思考
  PgUp/PgDn      翻页
  鼠标滚轮        滚动消息
`
	return helpBoxStyle.
		Width(m.width - 10).
		Height(m.height - 6).
		Render(helpTitleStyle.Render("  帮助  ") + "\n" + helpDescStyle.Render(helpText))
}

func (m Model) handleSlashCommand(content string) (tea.Model, tea.Cmd) {
	cmd := strings.Fields(content)
	if len(cmd) == 0 {
		return m, nil
	}

	switch cmd[0] {
	case "/help", "/?":
		m.ShowHelp = true

	case "/clear", "/cls":
		m.Messages = []DisplayMessage{}
		m.updateViewport()

	case "/quit", "/exit", "/q":
		return m, tea.Quit

	case "/quest":
		m.Mode = ModeQuest
		m.Messages = append(m.Messages, DisplayMessage{
			Role:      "system",
			Content:   "进入任务模式。描述你的目标，AI将自主规划和执行。",
			Timestamp: time.Now(),
		})
		m.updateViewport()

	case "/review":
		m.Mode = ModeReview
		m.Messages = append(m.Messages, DisplayMessage{
			Role:      "system",
			Content:   "进入审查模式。发送代码路径进行审查。",
			Timestamp: time.Now(),
		})
		m.updateViewport()

	case "/agent":
		m.Mode = ModeAgent
		m.Messages = append(m.Messages, DisplayMessage{
			Role:      "system",
			Content:   "进入代理模式。可用: coder, reviewer, architect, tester, debugger, researcher",
			Timestamp: time.Now(),
		})
		m.updateViewport()

	case "/chat":
		m.Mode = ModeChat
		m.Messages = append(m.Messages, DisplayMessage{
			Role:      "system",
			Content:   "切换到对话模式。",
			Timestamp: time.Now(),
		})
		m.updateViewport()

	case "/remote":
		m.Mode = ModeRemote
		if len(cmd) > 1 {
			switch cmd[1] {
			case "connect":
				// 建立连接
				if m.RemoteHost == "" {
					m.Messages = append(m.Messages, DisplayMessage{
						Role:      "error",
						Content:   "请先设置目标服务器: /remote <hostname>",
						Timestamp: time.Now(),
					})
					m.updateViewport()
				} else {
					// 发送连接请求消息
					return m, func() tea.Msg {
						return RemoteConnectMessage{Host: m.RemoteHost}
					}
				}
			case "disconnect":
				// 发送断开连接请求
				return m, func() tea.Msg {
					return RemoteDisconnectMessage{}
				}
			default:
				// 设置主机名
				m.RemoteHost = cmd[1]
				m.RemoteConnected = false
				m.Messages = append(m.Messages, DisplayMessage{
					Role:      "system",
					Content:   fmt.Sprintf("已设置目标服务器: %s\n使用 /remote connect 建立连接", cmd[1]),
					Timestamp: time.Now(),
				})
				m.updateViewport()
			}
		} else {
			m.Messages = append(m.Messages, DisplayMessage{
				Role:      "system",
				Content:   "远程开发模式\n\n用法:\n  /remote <host>       - 设置目标服务器\n  /remote connect      - 建立连接\n  /remote disconnect   - 断开连接",
				Timestamp: time.Now(),
			})
			m.updateViewport()
		}

	case "/model":
		if len(cmd) > 1 {
			m.ModelName = cmd[1]
			m.Messages = append(m.Messages, DisplayMessage{
				Role:      "success",
				Content:   fmt.Sprintf("模型已切换为: %s", cmd[1]),
				Timestamp: time.Now(),
			})
			m.updateViewport()
		}

	case "/context":
		m.Messages = append(m.Messages, DisplayMessage{
			Role:      "system",
			Content:   fmt.Sprintf("工作目录: %s\n模型: %s\n模式: %s", ".", m.ModelName, m.Mode),
			Timestamp: time.Now(),
		})
		m.updateViewport()

	default:
		m.Messages = append(m.Messages, DisplayMessage{
			Role:      "error",
			Content:   fmt.Sprintf("未知命令: %s", cmd[0]),
			Timestamp: time.Now(),
		})
		m.updateViewport()
	}

	return m, nil
}

func (m Model) sendMessage(content string) tea.Cmd {
	return func() tea.Msg {
		return UserMessage{Content: content}
	}
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

type ProcessingUpdate struct {
	Message string
}

// RemoteConnectMessage 远程连接请求
type RemoteConnectMessage struct {
	Host string
}

// RemoteDisconnectMessage 远程断开连接请求
type RemoteDisconnectMessage struct{}

// 设置方法
func (m *Model) SetModelName(name string) {
	m.ModelName = name
}

func (m *Model) SetProvider(provider string) {
	m.Provider = provider
}

func (m *Model) AddMessage(role, content string) {
	m.Messages = append(m.Messages, DisplayMessage{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})
	m.updateViewport()
}

func (m *Model) SetMode(mode Mode) {
	m.Mode = mode
}

func (m *Model) ShowConfirm(desc string, callback func(bool)) {
	m.AwaitingConfirm = true
	m.ConfirmDesc = desc
	m.ConfirmCallback = callback
}

// RefreshViewport 公开方法，更新视口内容
func (m *Model) RefreshViewport() {
	m.updateViewport()
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
