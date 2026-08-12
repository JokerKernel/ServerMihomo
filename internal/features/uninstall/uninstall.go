package uninstall

import (
	"context"

	"snailproxy/internal/feature"
	"snailproxy/internal/terminal"
)

type Feature struct{}

func (Feature) ID() string {
	return "uninstall"
}

func (Feature) Label() string {
	return "卸载"
}

func (Feature) Description() string {
	return "停止服务并移除 mihomo 运行文件"
}

func (Feature) Order() int {
	return 40
}

func (Feature) Run(ctx context.Context, runtime feature.Runtime) error {
	terminal.ClearScreen()
	terminal.MenuTitle("卸载")
	terminal.Warning("此操作将停止 mihomo 服务，并删除 /opt/mihomo 下的运行文件。")
	confirmed, err := terminal.ConfirmNoDefault("确认继续卸载? [y/N]: ")
	if err != nil {
		return err
	}
	if !confirmed {
		terminal.Info("已取消卸载。")
		return feature.ErrReturn
	}

	manager, err := runtime.NewPlatformManager()
	if err != nil {
		return err
	}
	return manager.Uninstall(ctx)
}
