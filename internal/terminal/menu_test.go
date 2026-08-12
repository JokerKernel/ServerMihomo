package terminal

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestSelectReprintsAfterEmptyAndInvalidInput(t *testing.T) {
	withInput(t, "\n9\n1\n")

	var action int
	var err error
	output := captureStdout(t, func() {
		action, err = Select("测试菜单:", "[0-1]", []MenuOption[int]{
			{Number: 1, Label: "继续", Value: 42},
			{Number: 0, Label: "返回", Value: 0},
		})
	})
	if err != nil {
		t.Fatalf("Select() error = %v", err)
	}
	if action != 42 {
		t.Fatalf("action = %d, want 42", action)
	}
	if count := strings.Count(output, "ServerMihomo  ›  测试菜单"); count != 3 {
		t.Fatalf("menu prints = %d, want 3\n%s", count, output)
	}
	for _, want := range []string{"输入不能为空，请输入菜单编号。", "输入无效，请重新输入。"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}

func TestSelectAcceptsReturnAliases(t *testing.T) {
	for _, input := range []string{"q\n", "Q\n", "exit\n", " EXIT \n"} {
		t.Run(strings.TrimSpace(input), func(t *testing.T) {
			withInput(t, input)

			action, err := Select("测试菜单", "[0-1]", []MenuOption[int]{
				{Number: 1, Label: "继续", Value: 1},
				{Number: 0, Label: "返回", Value: 0},
			})
			if err != nil {
				t.Fatalf("Select() error = %v", err)
			}
			if action != 0 {
				t.Fatalf("action = %d, want return action", action)
			}
		})
	}
}

func TestSelectRendersDescriptionAndDistinctReturn(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	withInput(t, "0\n")

	output := captureStdout(t, func() {
		_, _ = Select("安装与更新", "[0-1]", []MenuOption[int]{
			{Number: 1, Label: "本地安装", Description: "使用内嵌资源", Value: 1},
			{Number: 0, Label: "返回", Value: 0},
		})
	})

	for _, want := range []string{
		"ServerMihomo  ›  安装与更新",
		"1 本地安装 — 使用内嵌资源",
		"0/q 返回",
		"❯ 输入选项 [0-1]:",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
	if strings.Contains(output, "\033[") {
		t.Fatalf("NO_COLOR output contains ANSI escape sequence: %q", output)
	}
}

func TestHomeTitleIncludesVersionAndPlatform(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	output := captureStdout(t, func() {
		HomeTitle("v1.2.3", "linux/arm64")
	})
	for _, want := range []string{"ServerMihomo", "[版本 v1.2.3]", "[linux/arm64]"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q:\n%s", want, output)
		}
	}
}

func withInput(t *testing.T, input string) {
	t.Helper()
	restore := SetInput(strings.NewReader(input))
	t.Cleanup(restore)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = writer
	defer func() {
		os.Stdout = originalStdout
	}()

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout writer error = %v", err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read stdout error = %v", err)
	}
	return string(data)
}
