# ACCIL 发布说明

## 适用版本

本文档对应 `v1.4.6` 之后的发布流程。

## 发布前检查

- 更新 `VERSION`
- 更新 [CHANGELOG.md](CHANGELOG.md)
- 更新 [README.md](README.md)
- 更新 [README_zh.md](README_zh.md)
- 更新 [QUICKSTART.md](QUICKSTART.md)
- 确认安装脚本中的默认版本号一致

## 版本文件

需要保持一致的文件：

- `VERSION`
- `cmd/root.go`
- `build.bat`
- `install.bat`
- `install.ps1`
- `install.sh`

## 配置命令说明

从 `v1.4.6` 开始，运行时配置统一使用：

```text
/config show
/config provider <name>
/config model <name>
/config baseurl <url>
```

旧命令 `/provider`、`/model`、`/baseurl` 仅作为兼容别名保留。

## 文档要求

发布时需要确认：

- 用户文档优先使用官网链接
- 不再把 GitHub 作为主要安装或社区入口
- 中英文文档中的命令示例保持一致
- 更新说明准确描述当前版本变化

## 打包建议

### Windows

```powershell
build.bat
```

### 通用源码构建

```bash
go mod tidy
go build -o accil .
```

## 发布检查清单

- 版本号已更新为目标版本
- `CHANGELOG.md` 已新增对应条目
- README 中的 “What’s New” 已同步
- 中文 README 已同步
- Quick Start 已同步
- `/config` 命令说明已覆盖主流程

## 备注

如果后续还要继续收口配置入口，优先扩展 `/config` 子命令，不再新增新的顶层命令。
