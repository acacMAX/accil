# Quick Start Guide

## 🚀 Get Started in 5 Minutes

### Step 1: Install ACCIL

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/accil/accil/main/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/accil/accil/main/install.ps1 | iex
```

### Step 2: Configure API

Run the setup wizard:
```bash
accil --setup
```

You'll be prompted to:
1. Choose an API provider (OpenAI, DeepSeek, Anthropic, Ollama, etc.)
2. Enter your API key
3. Select a model

**Or set environment variables:**
```bash
export AI_API_KEY="your-api-key"
export AI_BASE_URL="https://api.openai.com/v1"
```

### Step 3: Start Using!

**Interactive mode:**
```bash
accil
```

**Single task:**
```bash
accil "Create a Python hello world script"
```

---

## 💡 Common Use Cases

### Create Files
```bash
accil "Create a README.md with project description"
```

### Read and Explain Code
```bash
accil "Read main.go and explain what it does"
```

### Execute Commands
```bash
accil "List all Go files and count them"
```

### Multi-step Tasks (Quest Mode)
```bash
accil --quest "Build a simple web server with health check endpoint"
```

### Code Review
```bash
accil review ./main.go
```

---

## 🎯 Tips for Best Results

1. **Be Specific**: Clear instructions get better results
   - ✅ "Create a Python function that sorts a list"
   - ❌ "Make something"

2. **Use Quest Mode for Complex Tasks**: 
   ```bash
   accil /quest "Create a REST API with CRUD operations"
   ```

3. **Review Generated Code**: Always review before executing
   - Don't use `--yolo` unless you trust the output

4. **Leverage Context**: ACCIL remembers your project
   - It will automatically read AGENTS.md if present

5. **Use Sub-agents for Specialized Tasks**:
   ```bash
   accil agent run reviewer "Check for security issues"
   accil agent run tester "Write unit tests"
   ```

---

## 🔧 Troubleshooting

### "Command not found"
- Restart your terminal after installation
- Or run: `source ~/.bashrc` (Linux) or refresh PATH (Windows)

### "API Key not configured"
- Run `accil --setup`
- Or set `AI_API_KEY` environment variable

### Network Timeout
- Check your internet connection
- The tool has automatic retry (up to 3 attempts)
- Consider using a proxy if needed

### Tool Execution Not Visible
- Make sure you're in interactive mode (`accil`)
- Real-time display only works in TUI mode

---

## 📚 Next Steps

- Read the full [README](README.md) for detailed documentation
- Explore [examples](examples/) directory for more use cases
- Join our [community discussions](https://github.com/accil/accil/discussions)
- Contribute! See [CONTRIBUTING.md](CONTRIBUTING.md)

---

<div align="center">

**Happy Coding with ACCIL! 🎉**

</div>
