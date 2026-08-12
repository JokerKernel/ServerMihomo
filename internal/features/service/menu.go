package service

import "snailproxy/internal/terminal"

type Action int

const (
	ActionReturn Action = iota
	ActionStart
	ActionRestart
	ActionStop
	ActionWriteProxyEnv
	ActionClearProxyEnv
)

func Select() (Action, error) {
	return terminal.Select("mihomo 服务与代理", "[0-5]", []terminal.MenuOption[Action]{
		{Number: 1, Label: "启动 mihomo 服务", Description: "启动已安装的 systemd 服务", Value: ActionStart},
		{Number: 2, Label: "重启 mihomo 服务", Description: "重新加载当前配置", Value: ActionRestart},
		{Number: 3, Label: "停止 mihomo 服务", Description: "停止代理服务", Value: ActionStop},
		{Number: 4, Label: "写入代理环境变量", Description: "配置当前用户 Shell 代理", Value: ActionWriteProxyEnv},
		{Number: 5, Label: "清除代理环境变量", Description: "移除本工具写入的代理配置", Value: ActionClearProxyEnv},
		{Number: 0, Label: "返回", Value: ActionReturn},
	})
}
