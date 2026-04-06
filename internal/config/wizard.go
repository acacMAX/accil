package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SetupWizard runs the first-time setup wizard
func SetupWizard() error {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                              ║")
	fmt.Println("║   █████╗ ██████╗ ██████╗  ██████╗██╗  ██╗██╗     ███████╗   ║")
	fmt.Println("║  ██╔══██╗██╔══██╗██╔══██╗██╔════╝██║  ██║██║     ██╔════╝   ║")
	fmt.Println("║  ███████║██████╔╝██████╔╝██║     ███████║██║     █████╗     ║")
	fmt.Println("║  ██╔══██║██╔══██╗██╔══██╗██║     ██╔══██║██║     ██╔══╝     ║")
	fmt.Println("║  ██║  ██║██████╔╝██████╔╝╚██████╗██║  ██║███████╗███████╗   ║")
	fmt.Println("║  ╚═╝  ╚═╝╚═════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝╚══════╝╚══════╝   ║")
	fmt.Println("║                                                              ║")
	fmt.Println("║           AI-Powered Autonomous Coding Assistant             ║")
	fmt.Println("║                                                              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  Welcome to ACCIL! Let's set up your configuration.")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// API Key
	fmt.Print("  Enter your API Key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	// Base URL
	fmt.Println()
	fmt.Println("  Select API Provider:")
	fmt.Println("  [1] OpenAI          (https://api.openai.com/v1)")
	fmt.Println("  [2] DeepSeek        (https://api.deepseek.com/v1)")
	fmt.Println("  [3] Anthropic       (https://api.anthropic.com/v1)")
	fmt.Println("  [4] Azure OpenAI    (custom)")
	fmt.Println("  [5] Ollama (local)  (http://localhost:11434/v1)")
	fmt.Println("  [6] Custom URL")
	fmt.Println()
	fmt.Print("  Select [1-6]: ")
	providerChoice, _ := reader.ReadString('\n')
	providerChoice = strings.TrimSpace(providerChoice)

	var baseURL string
	switch providerChoice {
	case "1":
		baseURL = "https://api.openai.com/v1"
	case "2":
		baseURL = "https://api.deepseek.com/v1"
	case "3":
		baseURL = "https://api.anthropic.com/v1"
	case "4":
		fmt.Print("  Enter Azure endpoint: ")
		baseURL, _ = reader.ReadString('\n')
		baseURL = strings.TrimSpace(baseURL)
	case "5":
		baseURL = "http://localhost:11434/v1"
	case "6":
		fmt.Print("  Enter custom base URL: ")
		baseURL, _ = reader.ReadString('\n')
		baseURL = strings.TrimSpace(baseURL)
	default:
		baseURL = "https://api.openai.com/v1"
	}

	// Model selection
	fmt.Println()
	fmt.Println("  Common models:")
	fmt.Println("  OpenAI:     gpt-4o, gpt-4-turbo, gpt-3.5-turbo")
	fmt.Println("  DeepSeek:   deepseek-chat, deepseek-coder")
	fmt.Println("  Anthropic:  claude-3-opus, claude-3-sonnet")
	fmt.Println("  Ollama:     llama3, codellama, mistral")
	fmt.Println()
	fmt.Print("  Enter model name: ")
	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)

	if model == "" {
		model = "gpt-4o"
	}

	// Save configuration
	cfg := &Config{
		APIKey:      apiKey,
		BaseURL:     baseURL,
		Model:       model,
		MaxTokens:   4096,
		AutoApprove: false,
		BlockList:   DefaultConfig.BlockList,
	}

	if err := Save(cfg); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("  ╔════════════════════════════════════════╗")
	fmt.Println("  ║        Configuration saved!            ║")
	fmt.Println("  ╚════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("  API URL: %s\n", baseURL)
	fmt.Printf("  Model:   %s\n", model)
	fmt.Println()
	fmt.Println("  Run 'accil' to start, or 'accil config' to change settings.")
	fmt.Println()

	return nil
}

// IsConfigured checks if the application is configured
func IsConfigured() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	configPath := filepath.Join(home, ".accil", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	}

	// Read config file directly instead of relying on viper
	data, err := os.ReadFile(configPath)
	if err != nil {
		return false
	}

	// Simple check: look for api_key in the file content
	content := string(data)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "api_key:") {
			keyPart := strings.TrimPrefix(line, "api_key:")
			keyPart = strings.TrimSpace(keyPart)
			keyPart = strings.Trim(keyPart, "\"'")
			return keyPart != ""
		}
	}
	return false
}

// EditConfig opens the configuration editor
func EditConfig() error {
	cfg, err := Load()
	if err != nil {
		cfg = &DefaultConfig
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("  ╔════════════════════════════════════════╗")
	fmt.Println("  ║         Configuration Editor           ║")
	fmt.Println("  ╚════════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("  Current API Key: %s***\n", maskKey(cfg.APIKey))
	fmt.Print("  New API Key (press Enter to keep): ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)
	if apiKey != "" {
		cfg.APIKey = apiKey
	}

	fmt.Println()
	fmt.Printf("  Current Base URL: %s\n", cfg.BaseURL)
	fmt.Print("  New Base URL (press Enter to keep): ")
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}

	fmt.Println()
	fmt.Printf("  Current Model: %s\n", cfg.Model)
	fmt.Print("  New Model (press Enter to keep): ")
	model, _ := reader.ReadString('\n')
	model = strings.TrimSpace(model)
	if model != "" {
		cfg.Model = model
	}

	fmt.Println()
	fmt.Print("  Auto-approve operations? (y/N): ")
	autoApprove, _ := reader.ReadString('\n')
	cfg.AutoApprove = strings.ToLower(strings.TrimSpace(autoApprove)) == "y"

	if err := Save(cfg); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("  Configuration updated!")
	return nil
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return ""
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}
