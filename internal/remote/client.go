package remote

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

// Client represents an SSH remote client
type Client struct {
	conn       *ssh.Client
	session    *ssh.Session
	Host       string
	Port       string
	User       string
	WorkDir    string
	connected  bool
	lastError  error
}

// Config represents SSH connection configuration
type Config struct {
	Host       string
	Port       string
	User       string
	Password   string
	KeyPath    string
	WorkDir    string
	UseAgent   bool
}

// NewClient creates a new remote client
func NewClient(config Config) (*Client, error) {
	if config.Port == "" {
		config.Port = "22"
	}
	if config.WorkDir == "" {
		config.WorkDir = "."
	}

	client := &Client{
		Host:    config.Host,
		Port:    config.Port,
		User:    config.User,
		WorkDir: config.WorkDir,
	}

	if err := client.connect(config); err != nil {
		return nil, err
	}

	return client, nil
}

// connect establishes SSH connection
func (c *Client) connect(config Config) error {
	var authMethods []ssh.AuthMethod

	// 1. Try SSH agent first if enabled
	if config.UseAgent {
		if agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
			agentClient := agent.NewClient(agentConn)
			authMethods = append(authMethods, ssh.PublicKeysCallback(agentClient.Signers))
		}
	}

	// 2. Try private key
	if config.KeyPath != "" {
		key, err := os.ReadFile(config.KeyPath)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
	}

	// 3. Try default keys
	defaultKeys := []string{
		filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"),
		filepath.Join(os.Getenv("HOME"), ".ssh", "id_ed25519"),
		filepath.Join(os.Getenv("HOME"), ".ssh", "id_ecdsa"),
	}
	for _, keyPath := range defaultKeys {
		key, err := os.ReadFile(keyPath)
		if err == nil {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
				break
			}
		}
	}

	// 4. Use password if provided
	if config.Password != "" {
		authMethods = append(authMethods, ssh.Password(config.Password))
	}

	hostKeyCallback := c.makeHostKeyCallback(config)

	sshConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback,
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		c.lastError = err
		return fmt.Errorf("failed to connect to %s: %w", addr, err)
	}

	c.conn = conn
	c.connected = true
	return nil
}

func (c *Client) makeHostKeyCallback(config Config) ssh.HostKeyCallback {
	home, _ := os.UserHomeDir()
	knownHostsPath := filepath.Join(home, ".ssh", "known_hosts")

	// Ensure known_hosts exists (empty is ok; it will then reject unknown hosts).
	if err := os.MkdirAll(filepath.Dir(knownHostsPath), 0700); err == nil {
		if _, err := os.Stat(knownHostsPath); os.IsNotExist(err) {
			_ = os.WriteFile(knownHostsPath, []byte(""), 0600)
		}
	}

	cb, err := knownhosts.New(knownHostsPath)
	if err != nil {
		// As a last resort, fail closed with a clear error.
		return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return fmt.Errorf("无法加载 known_hosts(%s): %v", knownHostsPath, err)
		}
	}

	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		if err := cb(hostname, remote, key); err == nil {
			return nil
		}

		// Provide a helpful message on unknown/mismatch host keys.
		fp := ssh.FingerprintSHA256(key)

		// knownhosts.Line wants host patterns; include both raw host and [host]:port.
		port := config.Port
		if port == "" {
			port = "22"
		}
		host1 := config.Host
		host2 := fmt.Sprintf("[%s]:%s", config.Host, port)

		// Create a deterministic comment so users can later identify the line.
		sum := sha256.Sum256(key.Marshal())
		comment := fmt.Sprintf("accil-added-%x", sum[:6])
		line := knownhosts.Line([]string{host1, host2}, key) + " " + comment

		return fmt.Errorf(
			"SSH 主机密钥校验失败/未知主机。\n主机: %s\n指纹: %s\n请确认该指纹可信后，将以下一行追加到 %s：\n%s",
			config.Host, fp, knownHostsPath, line,
		)
	}
}

// IsConnected returns connection status
func (c *Client) IsConnected() bool {
	return c.connected && c.conn != nil
}

// Disconnect closes the SSH connection
func (c *Client) Disconnect() error {
	c.connected = false
	if c.session != nil {
		c.session.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Execute runs a command on the remote server
func (c *Client) Execute(command string) (string, string, error) {
	if !c.IsConnected() {
		return "", "", fmt.Errorf("not connected")
	}

	session, err := c.conn.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Change to working directory
	if c.WorkDir != "" && c.WorkDir != "." {
		command = fmt.Sprintf("cd %s && %s", c.WorkDir, command)
	}

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)
	return stdout.String(), stderr.String(), err
}

// ReadFile reads a file from the remote server
func (c *Client) ReadFile(path string) (string, error) {
	if !c.IsConnected() {
		return "", fmt.Errorf("not connected")
	}

	// Use cat to read file
	cmd := fmt.Sprintf("cat '%s' 2>&1", escapePath(path))
	stdout, stderr, err := c.Execute(cmd)
	if err != nil {
		if stderr != "" {
			return "", fmt.Errorf("%s", stderr)
		}
		return "", err
	}
	return stdout, nil
}

// WriteFile writes content to a file on the remote server
func (c *Client) WriteFile(path string, content string) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	// Use tee to write file
	session, err := c.conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		mkdirCmd := fmt.Sprintf("mkdir -p '%s'", escapePath(dir))
		c.Execute(mkdirCmd)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}

	cmd := fmt.Sprintf("cat > '%s'", escapePath(path))
	if err := session.Start(cmd); err != nil {
		return err
	}

	io.WriteString(stdin, content)
	stdin.Close()

	return session.Wait()
}

// EditFile edits a file on the remote server
func (c *Client) EditFile(path string, oldString, newString string) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	content, err := c.ReadFile(path)
	if err != nil {
		return err
	}

	newContent := strings.Replace(content, oldString, newString, 1)
	if newContent == content {
		return fmt.Errorf("old string not found in file")
	}

	return c.WriteFile(path, newContent)
}

// ListDir lists directory contents on the remote server
func (c *Client) ListDir(path string) (string, error) {
	if !c.IsConnected() {
		return "", fmt.Errorf("not connected")
	}

	cmd := fmt.Sprintf("ls -la '%s' 2>&1", escapePath(path))
	stdout, _, err := c.Execute(cmd)
	return stdout, err
}

// Glob finds files matching a pattern on the remote server
func (c *Client) Glob(pattern string) (string, error) {
	if !c.IsConnected() {
		return "", fmt.Errorf("not connected")
	}

	cmd := fmt.Sprintf("find . -type f -name '%s' 2>/dev/null | head -50", escapePath(pattern))
	stdout, _, err := c.Execute(cmd)
	return stdout, err
}

// SearchCode searches for code patterns on the remote server
func (c *Client) SearchCode(pattern string) (string, error) {
	if !c.IsConnected() {
		return "", fmt.Errorf("not connected")
	}

	cmd := fmt.Sprintf("grep -r '%s' . --include='*.go' --include='*.js' --include='*.ts' --include='*.py' --include='*.java' 2>/dev/null | head -30", escapePath(pattern))
	stdout, _, err := c.Execute(cmd)
	return stdout, err
}

// GetInfo returns remote server information
func (c *Client) GetInfo() (map[string]string, error) {
	info := make(map[string]string)

	if !c.IsConnected() {
		return info, fmt.Errorf("not connected")
	}

	// Get hostname
	if out, _, err := c.Execute("hostname"); err == nil {
		info["hostname"] = strings.TrimSpace(out)
	}

	// Get OS info
	if out, _, err := c.Execute("uname -a"); err == nil {
		info["os"] = strings.TrimSpace(out)
	}

	// Get current directory
	if out, _, err := c.Execute("pwd"); err == nil {
		info["pwd"] = strings.TrimSpace(out)
	}

	// Get user
	if out, _, err := c.Execute("whoami"); err == nil {
		info["user"] = strings.TrimSpace(out)
	}

	return info, nil
}

// escapePath escapes single quotes in path
func escapePath(path string) string {
	return strings.ReplaceAll(path, "'", "'\"'\"'")
}

// TestConnection tests if connection is alive
func (c *Client) TestConnection() error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected")
	}

	_, _, err := c.Execute("echo 'ping'")
	return err
}
