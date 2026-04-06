package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/accil/accil/internal/agent"
	"github.com/accil/accil/internal/ai"
	"github.com/accil/accil/internal/config"
	appcontext "github.com/accil/accil/internal/context"
	"github.com/accil/accil/internal/memory"
	"github.com/accil/accil/internal/quest"
	"github.com/accil/accil/internal/review"
	"github.com/accil/accil/internal/session"
	"github.com/accil/accil/internal/tools"
	"github.com/accil/accil/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	Version      = "dev"
	flagWorkDir  string
	flagModel    string
	flagYolo     bool
	flagHeadless bool
	flagSession  string
	flagContinue bool
	flagSetup    bool
)

var rootCmd = &cobra.Command{
	Use:   "accil [prompt]",
	Short: "AI驱动的自主编程助手",
	Long: `Accil 是一个强大的终端AI编程助手：
  • 交互对话模式 - 对话式编程辅助
  • 任务模式 - 自主多步骤任务执行
  • 代码审查 - 安全性、性能和质量分析
  • 子代理 - 针对不同任务的专业代理
  • 上下文记忆 - 项目感知辅助`,
	Run: runRoot,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本号",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("accil 版本 %s %s/%s\n", Version, runtime.GOOS, runtime.GOARCH)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "编辑配置",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.EditConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
	},
}

var questCmd = &cobra.Command{
	Use:   "quest <目标>",
	Short: "启动自主任务",
	Args:  cobra.MinimumNArgs(1),
	Run:   runQuest,
}

var reviewCmd = &cobra.Command{
	Use:   "review [文件...]",
	Short: "审查代码问题",
	Run:   runReview,
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "管理和使用子代理",
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出可用代理",
	Run:   runAgentList,
}

var agentRunCmd = &cobra.Command{
	Use:   "run <代理> <任务>",
	Short: "运行指定代理",
	Args:  cobra.MinimumNArgs(2),
	Run:   runAgentTask,
}

var memoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "管理项目记忆",
}

var memoryInitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化项目记忆 (AGENTS.md)",
	Run:   runMemoryInit,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagWorkDir, "workdir", "w", "", "工作目录")
	rootCmd.PersistentFlags().StringVarP(&flagModel, "model", "m", "", "AI模型")
	rootCmd.PersistentFlags().BoolVarP(&flagYolo, "yolo", "y", false, "自动批准所有操作")
	rootCmd.PersistentFlags().BoolVar(&flagHeadless, "headless", false, "无头模式")
	rootCmd.PersistentFlags().StringVarP(&flagSession, "session", "s", "", "会话ID")
	rootCmd.PersistentFlags().BoolVarP(&flagContinue, "continue", "c", false, "继续上次会话")
	rootCmd.PersistentFlags().BoolVar(&flagSetup, "setup", false, "运行设置向导")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(questCmd)
	rootCmd.AddCommand(reviewCmd)

	agentCmd.AddCommand(agentListCmd)
	agentCmd.AddCommand(agentRunCmd)
	rootCmd.AddCommand(agentCmd)

	memoryCmd.AddCommand(memoryInitCmd)
	rootCmd.AddCommand(memoryCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	if flagSetup || !config.IsConfigured() {
		if err := config.SetupWizard(); err != nil {
			fmt.Fprintf(os.Stderr, "设置失败: %v\n", err)
			os.Exit(1)
		}
		if !flagSetup {
			fmt.Println("\n首次设置完成。运行 'accil' 启动。")
			return
		}
	}

	if err := config.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "初始化配置错误: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置错误: %v\n", err)
		os.Exit(1)
	}

	workDir := flagWorkDir
	if workDir == "" {
		workDir, _ = os.Getwd()
	}
	cfg.WorkDir = workDir

	if flagModel != "" {
		cfg.Model = flagModel
	}

	if flagYolo {
		cfg.AutoApprove = true
	}

	if cfg.APIKey == "" {
		fmt.Fprintln(os.Stderr, "错误: API密钥未配置")
		fmt.Fprintln(os.Stderr, "运行 'accil --setup' 配置或设置 AI_API_KEY 环境变量")
		os.Exit(1)
	}

	if len(args) > 0 {
		prompt := strings.Join(args, " ")
		runSingleShot(cfg, prompt)
		return
	}

	runInteractive(cfg)
}

func runSingleShot(cfg *config.Config, prompt string) {
	client := ai.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model)
	executor := newExecutor(cfg)

	messages := []ai.Message{
		{Role: "system", Content: getSystemPrompt(cfg.WorkDir)},
		{Role: "user", Content: prompt},
	}

	maxCalls := cfg.MaxToolCalls
	if maxCalls <= 0 {
		maxCalls = 30
	}

	for i := 0; i < maxCalls; i++ {
		fmt.Printf("\n[思考中...]\n")
		
		resp, err := client.Chat(messages, ai.GetDefaultTools())
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}

		msg := resp.Choices[0].Message
		messages = append(messages, msg)

		if len(msg.ToolCalls) == 0 {
			fmt.Println("\n[结果]")
			fmt.Println(msg.Content)
			return
		}

		// 显示AI的想法
		if msg.Content != "" {
			fmt.Printf("\n[AI] %s\n", truncateString(msg.Content, 200))
		}

		// 执行工具调用
		for _, tc := range msg.ToolCalls {
			fmt.Printf("\n[执行] %s\n", tc.Function.Name)
			fmt.Printf("  参数: %s\n", truncateString(tc.Function.Arguments, 50))
			
			result := executor.Execute(tc.Function.Name, tc.Function.Arguments)
			
			if result.Success {
				fmt.Printf("[完成] %s\n", truncateString(result.Output, 100))
			} else {
				fmt.Printf("[错误] %s\n", result.Error)
			}

			messages = append(messages, ai.Message{
				Role:       "tool",
				Content:    formatToolResult(result),
				ToolCallID: tc.ID,
				Name:       tc.Function.Name,
			})
		}
	}

	fmt.Fprintln(os.Stderr, "警告: 达到最大工具调用次数")
}

func runInteractive(cfg *config.Config) {
	sessionMgr, err := session.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化会话管理器错误: %v\n", err)
		os.Exit(1)
	}

	var sess *session.Session
	if flagSession != "" {
		sess, err = sessionMgr.Load(flagSession)
		if err != nil {
			fmt.Fprintf(os.Stderr, "加载会话错误: %v\n", err)
			os.Exit(1)
		}
	} else if flagContinue {
		sess, err = sessionMgr.GetLastSession()
	}

	if sess == nil {
		sess = sessionMgr.NewSession("新会话")
	}

	client := ai.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model)
	executor := newExecutor(cfg)
	contextMgr, _ := appcontext.NewManager(cfg.WorkDir)

	app := NewApp(cfg, client, executor, sessionMgr, sess, contextMgr)

	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		p.Quit()
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "运行TUI错误: %v\n", err)
		os.Exit(1)
	}

	sessionMgr.Save(sess)
}

func runQuest(cmd *cobra.Command, args []string) {
	cfg := loadConfig()
	client := ai.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model)
	executor := newExecutor(cfg)

	goal := strings.Join(args, " ")
	planner := quest.NewPlanner(client, executor)
	q := planner.CreateQuest(goal)

	fmt.Printf("\n╔══════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  任务: %-53s║\n", truncateString(goal, 53))
	fmt.Printf("╚══════════════════════════════════════════════════════════════╝\n\n")

	ctx := context.Background()
	if err := planner.Plan(ctx, q); err != nil {
		fmt.Fprintf(os.Stderr, "创建计划错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("已创建 %d 步计划:\n", len(q.Steps))
	for i, step := range q.Steps {
		fmt.Printf("  %d. %s\n", i+1, step.Description)
	}

	if !cfg.AutoApprove {
		fmt.Printf("\n开始执行? (y/n): ")
		var input string
		fmt.Scanln(&input)
		if strings.ToLower(input) != "y" {
			fmt.Println("任务已取消。")
			return
		}
	}

	progress := func(step quest.Step, total int) {
		fmt.Printf("\n[%d/%d] %s\n", total, total, step.Description)
	}

	approver := func(desc string) bool {
		fmt.Printf("确认: %s? (y/n): ", desc)
		var input string
		fmt.Scanln(&input)
		return strings.ToLower(input) == "y"
	}

	if err := planner.Execute(ctx, q, cfg.AutoApprove, approver, progress); err != nil {
		fmt.Fprintf(os.Stderr, "\n任务失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n╔══════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  任务 %s                                              ║\n", strings.ToUpper(string(q.Status)))
	fmt.Printf("╚══════════════════════════════════════════════════════════════╝\n")
}

func runReview(cmd *cobra.Command, args []string) {
	cfg := loadConfig()
	client := ai.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model)
	executor := newExecutor(cfg)

	reviewer := review.NewReviewer(client, executor)
	ctx := context.Background()

	var report *review.Report
	var err error

	if len(args) == 0 {
		fmt.Println("审查 git 变更...")
		report, err = reviewer.ReviewChanges(ctx)
	} else {
		fmt.Printf("审查 %d 个文件...\n", len(args))
		report, err = reviewer.ReviewFiles(ctx, args)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(review.FormatReport(report))
}

func runAgentList(cmd *cobra.Command, args []string) {
	fmt.Println("\n可用代理:")
	fmt.Println("─────────────────")

	agents := []struct {
		id, name, desc string
	}{
		{"coder", "代码生成器", "编写干净、高效的代码"},
		{"reviewer", "代码审查员", "安全性和质量分析"},
		{"architect", "软件架构师", "系统设计和结构"},
		{"tester", "测试工程师", "创建全面的测试"},
		{"debugger", "调试专家", "分析和修复错误"},
		{"researcher", "研究代理", "查找最佳实践"},
	}

	for _, a := range agents {
		fmt.Printf("\n  %-12s %s\n", a.id+":", a.name)
		fmt.Printf("  %-12s %s\n", "", a.desc)
	}

	fmt.Println("\n用法: accil agent run <代理> <任务>")
}

func runAgentTask(cmd *cobra.Command, args []string) {
	cfg := loadConfig()
	client := ai.NewClient(cfg.APIKey, cfg.BaseURL, cfg.Model)
	executor := newExecutor(cfg)

	agentID := args[0]
	task := strings.Join(args[1:], " ")

	mgr := agent.NewManager(client, executor)
	ctx := context.Background()
	
	result, err := mgr.AssignTask(ctx, agentID, agent.Task{Description: task}, cfg.AutoApprove, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("结果:")
	fmt.Println(result.Result)
}

func runMemoryInit(cmd *cobra.Command, args []string) {
	workDir := flagWorkDir
	if workDir == "" {
		workDir, _ = os.Getwd()
	}

	mgr := memory.NewManager(workDir)
	mem, err := mgr.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成记忆错误: %v\n", err)
		os.Exit(1)
	}

	if err := mgr.Save(mem); err != nil {
		fmt.Fprintf(os.Stderr, "保存记忆错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("已在 %s 创建 %s\n", workDir, memory.AgentsFileName)
}

func loadConfig() *config.Config {
	if err := config.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "初始化配置错误: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置错误: %v\n", err)
		os.Exit(1)
	}

	if cfg.APIKey == "" {
		fmt.Fprintln(os.Stderr, "错误: API密钥未配置。运行 'accil --setup'")
		os.Exit(1)
	}

	if flagWorkDir != "" {
		cfg.WorkDir = flagWorkDir
	} else {
		cfg.WorkDir, _ = os.Getwd()
	}

	if flagModel != "" {
		cfg.Model = flagModel
	}

	if flagYolo {
		cfg.AutoApprove = true
	}

	return cfg
}

// newExecutor 创建配置好的工具执行器
func newExecutor(cfg *config.Config) *tools.Executor {
	executor := newExecutor(cfg)
	if cfg.CommandTimeout > 0 {
		executor.SetCommandTimeout(time.Duration(cfg.CommandTimeout) * time.Second)
	}
	return executor
}

func getSystemPrompt(workDir string) string {
	basePrompt := `你是 Accil，一个 AI 驱动的自主编程助手。你可以：
- 读取、写入和编辑文件
- 执行 shell 命令
- 搜索和分析代码
- 规划和执行多步骤任务

总是有帮助、准确和彻底。做更改时，解释你在做什么。

重要：当用户给你一个任务时，你应该立即调用工具来完成任务，而不是只是说你"将要"做什么。直接执行操作！`

	mgr := memory.NewManager(workDir)
	if mgr.Exists() {
		content, err := mgr.LoadRaw()
		if err == nil {
			return basePrompt + "\n\n# 项目上下文\n\n" + content
		}
	}

	return basePrompt
}

func formatToolResult(result *tools.ToolResult) string {
	if result.Success {
		return result.Output
	}
	return fmt.Sprintf("错误: %s\n输出: %s", result.Error, result.Output)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// App 主应用
type App struct {
	model      tui.Model
	cfg        *config.Config
	client     *ai.Client
	executor   *tools.Executor
	sessionMgr *session.Manager
	session    *session.Session
	contextMgr *appcontext.Manager
	agentMgr   *agent.Manager
	planner    *quest.Planner
	reviewer   *review.Reviewer
	streaming  bool
	msgChan    chan tea.Msg // 用于持续接收流式消息
}

func NewApp(cfg *config.Config, client *ai.Client, executor *tools.Executor,
	sessionMgr *session.Manager, sess *session.Session, contextMgr *appcontext.Manager) App {

	model := tui.NewModel()
	model.SetModelName(cfg.Model)

	return App{
		model:      model,
		cfg:        cfg,
		client:     client,
		executor:   executor,
		sessionMgr: sessionMgr,
		session:    sess,
		contextMgr: contextMgr,
		agentMgr:   agent.NewManager(client, executor),
		planner:    quest.NewPlanner(client, executor),
		reviewer:   review.NewReviewer(client, executor),
	}
}

func (a App) Init() tea.Cmd {
	return a.model.Init()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 如果正在处理流式消息，持续从channel读取
	if a.msgChan != nil {
		select {
		case streamMsg, ok := <-a.msgChan:
			if !ok {
				// channel关闭，处理完成
				a.msgChan = nil
				// 更新模型处理该消息
				newModel, _ := a.model.Update(msg)
				a.model = newModel.(tui.Model)
				return a, nil
			}
			// 处理流式消息
			newModel, _ := a.model.Update(streamMsg)
			a.model = newModel.(tui.Model)
			// 继续等待下一条消息
			return a, func() tea.Msg {
				if nextMsg, nextOk := <-a.msgChan; nextOk {
					return nextMsg
				}
				return nil
			}
		default:
			// 没有立即的消息，继续等待
		}
	}

	switch msg := msg.(type) {
	case tui.UserMessage:
		cmd := a.processUserMessageWithTools(msg.Content)
		return a, cmd
	}

	newModel, cmd := a.model.Update(msg)
	a.model = newModel.(tui.Model)

	return a, cmd
}

func (a App) View() string {
	return a.model.View()
}

// processUserMessageWithTools 处理用户消息并实时显示执行过程
func (a *App) processUserMessageWithTools(content string) tea.Cmd {
	// 先添加用户消息到会话
	a.session.AddMessage("user", content)

	// 创建channel用于异步通信
	msgChan := make(chan tea.Msg, 100)
	a.msgChan = msgChan // 保存到App中以便Update方法使用

	// 启动goroutine进行异步处理
	go func() {
		defer close(msgChan)

		messages := []ai.Message{
			{Role: "system", Content: getSystemPrompt(a.cfg.WorkDir)},
		}
		messages = append(messages, a.session.Messages...)

		maxCalls := a.cfg.MaxToolCalls
		if maxCalls <= 0 {
			maxCalls = 30
		}

		var allOutput strings.Builder

		// 工具调用循环
		for i := 0; i < maxCalls; i++ {
			// 发送处理状态
			msgChan <- tui.ProcessingUpdate{
				Message: fmt.Sprintf("正在思考... (第%d轮)", i+1),
			}

			resp, err := a.client.Chat(messages, ai.GetDefaultTools())
			if err != nil {
				msgChan <- tui.ErrorMessage{Error: err}
				return
			}

			aiMsg := resp.Choices[0].Message
			messages = append(messages, aiMsg)

			// 如果AI有话要说，立即显示
			if aiMsg.Content != "" {
				allOutput.WriteString(aiMsg.Content)
				allOutput.WriteString("\n")
				msgChan <- tui.AssistantMessage{Content: aiMsg.Content}
			}

			// 没有工具调用时返回结果
			if len(aiMsg.ToolCalls) == 0 {
				a.session.AddMessage("assistant", allOutput.String())
				msgChan <- tui.ProcessingUpdate{Message: ""} // 清除处理提示
				return
			}

			// 执行工具调用，并实时显示
			for _, tc := range aiMsg.ToolCalls {
				// 构建工具调用日志
				toolLog := fmt.Sprintf("🔧 %s", tc.Function.Name)

				// 解析参数以显示更易读的信息
				var args map[string]interface{}
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err == nil {
					if path, ok := args["path"].(string); ok {
						toolLog += fmt.Sprintf(" → %s", path)
					}
					if cmd, ok := args["command"].(string); ok {
						toolLog += fmt.Sprintf(": %s", truncateString(cmd, 30))
					}
					if content, ok := args["content"].(string); ok {
						toolLog += fmt.Sprintf(" (%d 字符)", len(content))
					}
				}

				// 显示工具调用开始
				msgChan <- tui.ToolCallMessage{
					Tool: tc.Function.Name,
					Args: toolLog,
				}

				// 执行工具
				result := a.executor.Execute(tc.Function.Name, tc.Function.Arguments)

				// 显示执行结果
				if result.Success {
					msgChan <- tui.ToolResultMessage{
						Success: true,
						Result:  truncateString(result.Output, 200),
					}
				} else {
					msgChan <- tui.ToolResultMessage{
						Success: false,
						Result:  truncateString(result.Error, 200),
					}
				}

				// 累积输出
				if result.Success {
					allOutput.WriteString(fmt.Sprintf("%s ✅\n%s\n", toolLog, truncateString(result.Output, 100)))
				} else {
					allOutput.WriteString(fmt.Sprintf("%s ❌ %s\n", toolLog, truncateString(result.Error, 100)))
				}

				// 添加工具结果到消息历史
				messages = append(messages, ai.Message{
					Role:       "tool",
					Content:    formatToolResult(result),
					ToolCallID: tc.ID,
					Name:       tc.Function.Name,
				})
			}
		}

		// 达到最大调用次数
		finalOutput := allOutput.String() + "\n⚠️ [达到最大调用次数，任务可能未完成]"
		a.session.AddMessage("assistant", finalOutput)
		msgChan <- tui.ProcessingUpdate{Message: ""}
		msgChan <- tui.AssistantMessage{Content: finalOutput}
	}()

	// 立即返回第一条消息
	return func() tea.Msg {
		msg, ok := <-msgChan
		if !ok {
			return nil
		}
		return msg
	}
}
