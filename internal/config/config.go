package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey        string   `mapstructure:"api_key"`
	BaseURL       string   `mapstructure:"base_url"`
	Model         string   `mapstructure:"model"`
	MaxTokens     int      `mapstructure:"max_tokens"`
	MaxToolCalls  int      `mapstructure:"max_tool_calls"` // 最大工具调用次数
	AutoApprove   bool     `mapstructure:"auto_approve"`
	BlockList     []string `mapstructure:"block_list"`
	WorkDir       string   `mapstructure:"workdir"`
}

var DefaultConfig = Config{
	BaseURL:       "https://api.openai.com/v1",
	Model:         "gpt-4o",
	MaxTokens:     4096,
	MaxToolCalls:  30, // 增加到30次
	AutoApprove:   false,
	BlockList: []string{
		"rm -rf /",
		"rm -rf /*",
		"mkfs",
		"dd if=",
		":(){ :|:& };:",
	},
}

func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".accil")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("api_key", "")
	viper.SetDefault("base_url", DefaultConfig.BaseURL)
	viper.SetDefault("model", DefaultConfig.Model)
	viper.SetDefault("max_tokens", DefaultConfig.MaxTokens)
	viper.SetDefault("auto_approve", DefaultConfig.AutoApprove)
	viper.SetDefault("block_list", DefaultConfig.BlockList)

	// Environment variables
	viper.SetEnvPrefix("AI")
	viper.AutomaticEnv()
	viper.BindEnv("api_key", "AI_API_KEY")
	viper.BindEnv("base_url", "AI_BASE_URL")

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}

func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Override with environment variables
	if apiKey := os.Getenv("AI_API_KEY"); apiKey != "" {
		cfg.APIKey = apiKey
	}
	if baseURL := os.Getenv("AI_BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	viper.Set("api_key", cfg.APIKey)
	viper.Set("base_url", cfg.BaseURL)
	viper.Set("model", cfg.Model)
	viper.Set("max_tokens", cfg.MaxTokens)
	viper.Set("max_tool_calls", cfg.MaxToolCalls)
	viper.Set("auto_approve", cfg.AutoApprove)
	viper.Set("block_list", cfg.BlockList)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".accil")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")
	return viper.WriteConfigAs(configPath)
}

func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".accil"), nil
}

func GetSessionsDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	sessionsDir := filepath.Join(configDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return "", err
	}
	return sessionsDir, nil
}
