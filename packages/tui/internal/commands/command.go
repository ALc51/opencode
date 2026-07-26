package commands

import (
	"encoding/json"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/sst/opencode-sdk-go"
)

type ExecuteCommandMsg Command
type ExecuteCommandsMsg []Command
type CommandExecutedMsg Command

type Keybinding struct {
	RequiresLeader bool
	Key            string
}

func (k Keybinding) Matches(msg tea.KeyPressMsg, leader bool) bool {
	key := k.Key
	key = strings.TrimSpace(key)
	return key == msg.String() && (k.RequiresLeader == leader)
}

type CommandName string
type Command struct {
	Name        CommandName
	Description string
	Keybindings []Keybinding
	Trigger     string
}

func (c Command) Keys() []string {
	var keys []string
	for _, k := range c.Keybindings {
		keys = append(keys, k.Key)
	}
	return keys
}

type CommandRegistry map[CommandName]Command

func (r CommandRegistry) Sorted() []Command {
	var commands []Command
	for _, command := range r {
		commands = append(commands, command)
	}
	slices.SortFunc(commands, func(a, b Command) int {
		if a.Name == AppExitCommand {
			return 1
		}
		if b.Name == AppExitCommand {
			return -1
		}
		return strings.Compare(string(a.Name), string(b.Name))
	})
	return commands
}

func (r CommandRegistry) Matches(msg tea.KeyPressMsg, leader bool) []Command {
	var matched []Command
	for _, command := range r.Sorted() {
		if command.Matches(msg, leader) {
			matched = append(matched, command)
		}
	}
	return matched
}

const (
	AppHelpCommand              CommandName = "app_help"
	EditorOpenCommand           CommandName = "editor_open"
	SessionNewCommand           CommandName = "session_new"
	SessionListCommand          CommandName = "session_list"
	SessionShareCommand         CommandName = "session_share"
	SessionUnshareCommand       CommandName = "session_unshare"
	SessionInterruptCommand     CommandName = "session_interrupt"
	SessionCompactCommand       CommandName = "session_compact"
	ToolDetailsCommand          CommandName = "tool_details"
	ModelListCommand            CommandName = "model_list"
	ThemeListCommand            CommandName = "theme_list"
	FileListCommand             CommandName = "file_list"
	FileCloseCommand            CommandName = "file_close"
	FileSearchCommand           CommandName = "file_search"
	FileDiffToggleCommand       CommandName = "file_diff_toggle"
	ProjectInitCommand          CommandName = "project_init"
	InputClearCommand           CommandName = "input_clear"
	InputPasteCommand           CommandName = "input_paste"
	InputSubmitCommand          CommandName = "input_submit"
	InputNewlineCommand         CommandName = "input_newline"
	MessagesPageUpCommand       CommandName = "messages_page_up"
	MessagesPageDownCommand     CommandName = "messages_page_down"
	MessagesHalfPageUpCommand   CommandName = "messages_half_page_up"
	MessagesHalfPageDownCommand CommandName = "messages_half_page_down"
	MessagesPreviousCommand     CommandName = "messages_previous"
	MessagesNextCommand         CommandName = "messages_next"
	MessagesFirstCommand        CommandName = "messages_first"
	MessagesLastCommand         CommandName = "messages_last"
	MessagesLayoutToggleCommand CommandName = "messages_layout_toggle"
	MessagesCopyCommand         CommandName = "messages_copy"
	MessagesRevertCommand       CommandName = "messages_revert"
	AppExitCommand              CommandName = "app_exit"
)

func (k Command) Matches(msg tea.KeyPressMsg, leader bool) bool {
	for _, binding := range k.Keybindings {
		if binding.Matches(msg, leader) {
			return true
		}
	}
	return false
}

func parseBindings(bindings ...string) []Keybinding {
	var parsedBindings []Keybinding
	for _, binding := range bindings {
		for p := range strings.SplitSeq(binding, ",") {
			requireLeader := strings.HasPrefix(p, "<leader>")
			keybinding := strings.ReplaceAll(p, "<leader>", "")
			keybinding = strings.TrimSpace(keybinding)
			parsedBindings = append(parsedBindings, Keybinding{
				RequiresLeader: requireLeader,
				Key:            keybinding,
			})
		}
	}
	return parsedBindings
}

func LoadFromConfig(config *opencode.Config) CommandRegistry {
	defaults := []Command{
		{
			Name:        AppHelpCommand,
			Description: "显示帮助",
			Keybindings: parseBindings("<leader>h"),
			Trigger:     "help",
		},
		{
			Name:        EditorOpenCommand,
			Description: "打开编辑器",
			Keybindings: parseBindings("<leader>e"),
			Trigger:     "editor",
		},
		{
			Name:        SessionNewCommand,
			Description: "新建会话",
			Keybindings: parseBindings("<leader>n"),
			Trigger:     "new",
		},
		{
			Name:        SessionListCommand,
			Description: "会话列表",
			Keybindings: parseBindings("<leader>l"),
			Trigger:     "sessions",
		},
		{
			Name:        SessionShareCommand,
			Description: "分享会话",
			Keybindings: parseBindings("<leader>s"),
			Trigger:     "share",
		},
		{
			Name:        SessionUnshareCommand,
			Description: "取消分享",
			Keybindings: parseBindings("<leader>u"),
			Trigger:     "unshare",
		},
		{
			Name:        SessionInterruptCommand,
			Description: "中断会话",
			Keybindings: parseBindings("esc"),
		},
		{
			Name:        SessionCompactCommand,
			Description: "压缩会话",
			Keybindings: parseBindings("<leader>c"),
			Trigger:     "compact",
		},
		{
			Name:        ToolDetailsCommand,
			Description: "切换工具详情",
			Keybindings: parseBindings("<leader>d"),
			Trigger:     "details",
		},
		{
			Name:        ModelListCommand,
			Description: "模型列表",
			Keybindings: parseBindings("<leader>m"),
			Trigger:     "models",
		},
		{
			Name:        ThemeListCommand,
			Description: "主题列表",
			Keybindings: parseBindings("<leader>t"),
			Trigger:     "themes",
		},
		{
			Name:        FileListCommand,
			Description: "文件列表",
			Keybindings: parseBindings("<leader>f"),
			Trigger:     "files",
		},
		{
			Name:        FileCloseCommand,
			Description: "关闭文件",
			Keybindings: parseBindings("esc"),
		},
		{
			Name:        FileSearchCommand,
			Description: "搜索文件",
			Keybindings: parseBindings("<leader>/"),
		},
		{
			Name:        FileDiffToggleCommand,
			Description: "分屏/统一差异视图",
			Keybindings: parseBindings("<leader>v"),
		},
		{
			Name:        ProjectInitCommand,
			Description: "创建/更新 AGENTS.md",
			Keybindings: parseBindings("<leader>i"),
			Trigger:     "init",
		},
		{
			Name:        InputClearCommand,
			Description: "清空输入",
			Keybindings: parseBindings("ctrl+c"),
		},
		{
			Name:        InputPasteCommand,
			Description: "粘贴内容",
			Keybindings: parseBindings("ctrl+v"),
		},
		{
			Name:        InputSubmitCommand,
			Description: "发送消息",
			Keybindings: parseBindings("enter"),
		},
		{
			Name:        InputNewlineCommand,
			Description: "插入换行",
			Keybindings: parseBindings("shift+enter", "ctrl+j"),
		},
		{
			Name:        MessagesPageUpCommand,
			Description: "向上翻页",
			Keybindings: parseBindings("pgup"),
		},
		{
			Name:        MessagesPageDownCommand,
			Description: "向下翻页",
			Keybindings: parseBindings("pgdown"),
		},
		{
			Name:        MessagesHalfPageUpCommand,
			Description: "向上翻半页",
			Keybindings: parseBindings("ctrl+alt+u"),
		},
		{
			Name:        MessagesHalfPageDownCommand,
			Description: "向下翻半页",
			Keybindings: parseBindings("ctrl+alt+d"),
		},
		{
			Name:        MessagesPreviousCommand,
			Description: "上一条消息",
			Keybindings: parseBindings("ctrl+up"),
		},
		{
			Name:        MessagesNextCommand,
			Description: "下一条消息",
			Keybindings: parseBindings("ctrl+down"),
		},
		{
			Name:        MessagesFirstCommand,
			Description: "第一条消息",
			Keybindings: parseBindings("ctrl+g"),
		},
		{
			Name:        MessagesLastCommand,
			Description: "最后一条消息",
			Keybindings: parseBindings("ctrl+alt+g"),
		},
		{
			Name:        MessagesLayoutToggleCommand,
			Description: "切换布局",
			Keybindings: parseBindings("<leader>p"),
		},
		{
			Name:        MessagesCopyCommand,
			Description: "复制消息",
			Keybindings: parseBindings("<leader>y"),
		},
		{
			Name:        MessagesRevertCommand,
			Description: "回退消息",
			Keybindings: parseBindings("<leader>r"),
		},
		{
			Name:        AppExitCommand,
			Description: "退出应用",
			Keybindings: parseBindings("ctrl+c", "<leader>q"),
			Trigger:     "exit",
		},
	}
	registry := make(CommandRegistry)
	keybinds := map[string]string{}
	marshalled, _ := json.Marshal(config.Keybinds)
	json.Unmarshal(marshalled, &keybinds)
	for _, command := range defaults {
		if keybind, ok := keybinds[string(command.Name)]; ok && keybind != "" {
			command.Keybindings = parseBindings(keybind)
		}
		registry[command.Name] = command
	}
	return registry
}
