package terminal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type MenuOption[T any] struct {
	Number      int
	Label       string
	Description string
	Value       T
}

type Terminal interface {
	ReadLine() (string, error)
	ConfirmNoDefault(prompt string) (bool, error)
}

type stdTerminal struct{}

var stdinReader = bufio.NewReader(os.Stdin)

func Default() Terminal {
	return stdTerminal{}
}

func SetInput(reader io.Reader) func() {
	originalReader := stdinReader
	stdinReader = bufio.NewReader(reader)
	return func() {
		stdinReader = originalReader
	}
}

func Select[T any](title string, promptRange string, options []MenuOption[T]) (T, error) {
	return selectMenu(func() {
		printMenuHeading(title)
	}, promptRange, options)
}

func SelectHome[T any](buildVersion string, platform string, promptRange string, options []MenuOption[T]) (T, error) {
	return selectMenu(func() {
		HomeTitle(buildVersion, platform)
	}, promptRange, options)
}

func selectMenu[T any](printHeading func(), promptRange string, options []MenuOption[T]) (T, error) {
	var zero T
	if len(options) == 0 {
		return zero, fmt.Errorf("菜单选项不能为空")
	}

	actions := make(map[int]T, len(options))
	for _, option := range options {
		actions[option.Number] = option.Value
	}

	ClearScreen()
	for {
		printHeading()
		for _, option := range options {
			key := strconv.Itoa(option.Number)
			if option.Number == 0 {
				key = "0/q"
				PrintMenuExit(key, option.Label)
				continue
			}
			PrintMenuOption(key, option.Label, option.Description)
		}
		fmt.Println()
		value, err := Ask(fmt.Sprintf("输入选项 %s: ", promptRange))
		if err != nil {
			return zero, fmt.Errorf("读取用户输入失败: %w", err)
		}

		if value == "" {
			Warning("输入不能为空，请输入菜单编号。")
			fmt.Println()
			continue
		}
		if isReturnChoice(value) {
			if action, ok := actions[0]; ok {
				return action, nil
			}
		}

		number, err := strconv.Atoi(value)
		if err == nil {
			if action, ok := actions[number]; ok {
				return action, nil
			}
		}
		Warning("输入无效，请重新输入。")
		fmt.Println()
	}
}

func isReturnChoice(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "q", "exit":
		return true
	default:
		return false
	}
}

func Ask(prompt string) (string, error) {
	fmt.Print(paint(bold+cyan, "❯ ") + paint(bold, prompt))
	line, err := ReadLine()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func ConfirmNoDefault(prompt string) (bool, error) {
	return Default().ConfirmNoDefault(prompt)
}

func (stdTerminal) ConfirmNoDefault(prompt string) (bool, error) {
	answer, err := Ask(prompt)
	if err != nil {
		return false, fmt.Errorf("读取用户输入失败: %w", err)
	}

	answer = strings.ToLower(answer)
	return answer == "y" || answer == "yes", nil
}

func Pause(prompt string) error {
	if strings.TrimSpace(prompt) == "" {
		prompt = "按 Enter 继续..."
	}
	fmt.Print(paint(dim, prompt))
	if _, err := ReadLine(); err != nil {
		return fmt.Errorf("读取用户输入失败: %w", err)
	}
	return nil
}

func ReadLine() (string, error) {
	return Default().ReadLine()
}

func (stdTerminal) ReadLine() (string, error) {
	return stdinReader.ReadString('\n')
}

func ClearScreen() {
	if !terminalOutput() || strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return
	}
	fmt.Print("\033[2J\033[H\033[3J")
}

func printMenuHeading(title string) {
	lines := strings.Split(strings.TrimSpace(title), "\n")
	path := make([]string, 0, 2)
	if len(lines) > 0 {
		for _, part := range strings.Split(strings.TrimSuffix(strings.TrimSpace(lines[0]), ":"), "›") {
			if part = strings.TrimSpace(part); part != "" {
				path = append(path, part)
			}
		}
	}
	MenuTitle(path...)
	for _, line := range lines[1:] {
		fmt.Println(paint(dim, line))
	}
	if len(lines) > 1 {
		fmt.Println()
	}
}

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	cyan   = "\033[36m"
	blue   = "\033[34m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
)

func HomeTitle(buildVersion string, platform string) {
	buildVersion = strings.TrimSpace(buildVersion)
	if buildVersion == "" {
		buildVersion = "dev"
	}
	platform = strings.TrimSpace(platform)

	fmt.Println(paint(cyan, "╭─") + " " + paint(bold+cyan, "ServerMihomo"))
	fmt.Print(paint(cyan, "│") + " " + paint(dim, "mihomo 安装与配置工具") + "  " + Badge("版本 "+buildVersion, true))
	if platform != "" {
		fmt.Print(" " + Badge(platform, true))
	}
	fmt.Println()
	fmt.Println(paint(cyan, "╰"+strings.Repeat("─", 52)))
	fmt.Println()
}

func MenuTitle(parts ...string) {
	clean := make([]string, 0, len(parts)+1)
	clean = append(clean, "ServerMihomo")
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" && part != "主菜单" {
			clean = append(clean, part)
		}
	}
	title := strings.Join(clean, "  ›  ")
	fmt.Println(paint(cyan, "╭─") + " " + paint(bold+cyan, title))
	fmt.Println(paint(cyan, "╰"+strings.Repeat("─", 52)))
	fmt.Println()
}

func PrintMenuOption(key, label, description string) {
	fmt.Printf("  %s %s", paint(bold+blue, key), label)
	if description = strings.TrimSpace(description); description != "" {
		fmt.Print(paint(dim, " — "+description))
	}
	fmt.Println()
}

func PrintMenuExit(key, label string) {
	fmt.Printf("  %s %s\n", paint(bold+yellow, key), paint(dim, label))
}

func Badge(text string, positive bool) string {
	color := yellow
	if positive {
		color = green
	}
	return paint(color, "["+text+"]")
}

func Info(message string) {
	fmt.Println(paint(cyan, "•") + " " + message)
}

func Success(message string) {
	fmt.Println(paint(green, "✓") + " " + message)
}

func Warning(message string) {
	fmt.Println(paint(yellow, "!") + " " + message)
}

func Error(message string) {
	fmt.Println(paint(red, "✗") + " " + message)
}

func paint(style, text string) string {
	if !colorEnabled() {
		return text
	}
	return style + text + reset
}

func colorEnabled() bool {
	if _, disabled := os.LookupEnv("NO_COLOR"); disabled || strings.EqualFold(os.Getenv("TERM"), "dumb") {
		return false
	}
	if forced := os.Getenv("CLICOLOR_FORCE"); forced != "" && forced != "0" {
		return true
	}
	return terminalOutput()
}

func terminalOutput() bool {
	info, err := os.Stdout.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
