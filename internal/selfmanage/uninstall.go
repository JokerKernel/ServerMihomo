package selfmanage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	InstallScriptURL = "https://raw.githubusercontent.com/JokerKernel/ServerMihomo/main/scripts/install.sh"
	maxScriptSize    = 1024 * 1024
)

// Update downloads the repository installer and delegates updating to it.
func Update(ctx context.Context) error {
	client := &http.Client{Timeout: 30 * time.Second}
	fmt.Printf("安装脚本来源: %s\n", InstallScriptURL)
	return run(ctx, client, InstallScriptURL, os.Stdin, os.Stdout, os.Stderr)
}

// Uninstall downloads the repository installer and delegates self-removal to it.
func Uninstall(ctx context.Context) error {
	client := &http.Client{Timeout: 30 * time.Second}
	fmt.Printf("安装脚本来源: %s\n", InstallScriptURL)
	return run(ctx, client, InstallScriptURL, os.Stdin, os.Stdout, os.Stderr, "uninstall")
}

func run(ctx context.Context, client *http.Client, scriptURL string, stdin io.Reader, stdout, stderr io.Writer, scriptArgs ...string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, scriptURL, nil)
	if err != nil {
		return fmt.Errorf("创建安装脚本请求失败: %w", err)
	}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("下载安装脚本失败: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("下载安装脚本失败: HTTP %d %s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	data, err := io.ReadAll(io.LimitReader(response.Body, maxScriptSize+1))
	if err != nil {
		return fmt.Errorf("读取安装脚本失败: %w", err)
	}
	if len(data) > maxScriptSize {
		return fmt.Errorf("安装脚本超过 %d 字节限制", maxScriptSize)
	}
	if !bytes.HasPrefix(data, []byte("#!/bin/sh")) {
		return fmt.Errorf("下载内容不是有效的 ServerMihomo 安装脚本")
	}

	temporary, err := os.CreateTemp("", "servermihomo-installer-*.sh")
	if err != nil {
		return fmt.Errorf("创建安装脚本临时文件失败: %w", err)
	}
	path := temporary.Name()
	defer os.Remove(path)
	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return fmt.Errorf("设置安装脚本权限失败: %w", err)
	}
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return fmt.Errorf("保存安装脚本失败: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("保存安装脚本失败: %w", err)
	}

	commandArgs := append([]string{path}, scriptArgs...)
	command := exec.CommandContext(ctx, "sh", commandArgs...)
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("执行安装脚本失败: %w", err)
	}
	return nil
}
