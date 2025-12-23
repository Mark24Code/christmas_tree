# 🎄 圣诞魔法 - 炫酷命令行动画

一个充满魔法的命令行圣诞动画，使用 Go 语言编写，包含三幕精彩场景！

## 🚀 快速开始

### 系统要求
- Go 1.21+
- 最小终端：50x25
- 支持 ANSI 颜色和 Unicode 的终端

### 安装依赖

```bash
# 如果遇到网络问题，使用国内镜像
export GOPROXY=https://goproxy.cn,direct

# 下载依赖
go mod tidy

# 编译
go build -o christmas_tree
```

### 运行

```bash
# 默认名字运行
./christmas_tree

# 自定义名字
./christmas_tree -name "圣诞快乐 2024"
./christmas_tree -name "Happy Holidays!"

# 或直接运行
go run main.go -name "你的名字"
```

### 控制
- **退出**：按 `ESC`、`Ctrl+C` 或 `q`
- **自动退出**：动画播放完毕（约 26 秒）后自动退出


## 🎯 设计灵感

灵感来源于：
- Ruby Trick Contest 的创意编程艺术

## 🛠️ 依赖

- [tcell v2](https://github.com/gdamore/tcell) - 强大的终端控制库

## 💡 提示

1. **最佳观看体验**：
   - 使用支持 UTF-8 的现代终端
   - 深色主题效果更好
   - 较大的终端窗口能看到完整效果

2. **录制分享**：
   - 可以使用 `asciinema` 录制动画
   - 分享给朋友观看

3. **自定义名字**：
   - 适合作为节日问候
   - 可以放名字、祝福语等

## 🐛 故障排除

**网络超时**
```bash
export GOPROXY=https://goproxy.cn,direct
go mod tidy
```

**终端太小**
```bash
# 调整终端窗口到至少 50x25
# 推荐 80x30 以上获得最佳效果
```

**字符显示异常**
```bash
# 确保终端支持 UTF-8
export LANG=en_US.UTF-8

# 使用现代终端：
# - macOS: iTerm2
# - Windows: Windows Terminal
# - Linux: GNOME Terminal, Konsole
```

**颜色异常**
```bash
# 检查终端支持 256 色
export TERM=xterm-256color
```

## 📜 许可证

MIT License


---

**祝你圣诞快乐！🎄✨**

*Merry Christmas and Happy Holidays!*
