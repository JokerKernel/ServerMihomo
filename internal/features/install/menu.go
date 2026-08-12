package install

import (
	"fmt"

	"snailproxy/internal/infra/github"
	"snailproxy/internal/terminal"
)

type Action int

const (
	ActionReturn Action = iota
	ActionLocal
	ActionOnline
	ActionService
)

func Select() (Action, error) {
	return terminal.Select("安装与更新", "[0-3]", []terminal.MenuOption[Action]{
		{Number: 1, Label: "本地安装", Description: "使用程序内嵌资源，无需联网", Value: ActionLocal},
		{Number: 2, Label: "在线安装 mihomo", Description: "从 GitHub 获取当前架构版本", Value: ActionOnline},
		{Number: 3, Label: "安装/更新 systemd 服务", Description: "写入服务配置，默认不启动", Value: ActionService},
		{Number: 0, Label: "返回", Value: ActionReturn},
	})
}

func SelectAsset(assets []github.Asset) (github.Asset, error) {
	options := make([]terminal.MenuOption[github.Asset], 0, len(assets)+1)
	for i, asset := range assets {
		options = append(options, terminal.MenuOption[github.Asset]{
			Number:      i + 1,
			Label:       asset.Name,
			Description: formatAssetSize(asset.Size),
			Value:       asset,
		})
	}
	options = append(options, terminal.MenuOption[github.Asset]{Number: 0, Label: "返回"})

	asset, err := terminal.Select("安装与更新 › 选择安装包", fmt.Sprintf("[0-%d]", len(assets)), options)
	if err != nil {
		return github.Asset{}, err
	}
	if asset.Name == "" {
		return github.Asset{}, errReturnToInstallMenu
	}
	return asset, nil
}

func formatAssetSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for value := size / unit; value >= unit; value /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}

func ConfirmOverwrite(path string) (bool, error) {
	fmt.Printf("本地文件已存在: %s\n", path)
	return terminal.ConfirmNoDefault("是否覆盖重新下载? [y/N]: ")
}

func ConfirmOverwriteInstall(path string) (bool, error) {
	fmt.Printf("安装目录中已存在程序文件: %s\n", path)
	return terminal.ConfirmNoDefault("是否覆盖安装? [y/N]: ")
}

func ConfirmOverwriteLocalInstall(targetDir string) (bool, error) {
	fmt.Printf("本地安装目录: %s\n", targetDir)
	return terminal.ConfirmNoDefault("遇到同名文件时是否覆盖安装? [y/N]: ")
}
