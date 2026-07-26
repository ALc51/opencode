#!/usr/bin/env python3
"""手动翻译映射,替换 opencode TUI Go 源码 UI 字符串(不依赖 API)。"""
import re, os, sys

MANUAL = {
    "Session deleted successfully": "会话已删除",
    "Share URL copied to clipboard!": "分享链接已复制到剪贴板!",
    "Session unshared successfully": "会话已取消分享",
    "Message copied to clipboard": "消息已复制到剪贴板",
    "Provider error: ": "提供商错误: ",
    "Failed to open session": "打开会话失败",
    "Failed to read file": "读取文件失败",
    "No EDITOR set, can't open editor": "未设置 EDITOR,无法打开编辑器",
    "Something went wrong, couldn't open editor": "出错了,无法打开编辑器",
    "Failed to share session": "分享会话失败",
    "Failed to unshare session": "取消分享会话失败",
    "opencode updated to ": "opencode 已更新到 ",
    "New version installed": "新版本已安装",
    "Tool details are now visible": "工具详情已显示",
    "Tool details are now hidden": "工具详情已隐藏",
    "Select Model": "选择模型",
    "Failed to delete session: ": "删除会话失败: ",
    "Switch Session": "切换会话",
    "Help": "帮助",
    "Select Theme": "选择主题",
    "Find Files": "查找文件",
    "show help": "显示帮助",
    "open editor": "打开编辑器",
    "new session": "新建会话",
    "list sessions": "会话列表",
    "share session": "分享会话",
    "unshare session": "取消分享",
    "interrupt session": "中断会话",
    "compact the session": "压缩会话",
    "toggle tool details": "切换工具详情",
    "list models": "模型列表",
    "list themes": "主题列表",
    "list files": "文件列表",
    "close file": "关闭文件",
    "search file": "搜索文件",
    "split/unified diff": "分屏/统一差异视图",
    "create/update AGENTS.md": "创建/更新 AGENTS.md",
    "clear input": "清空输入",
    "paste content": "粘贴内容",
    "submit message": "发送消息",
    "insert newline": "插入换行",
    "page up": "向上翻页",
    "page down": "向下翻页",
    "half page up": "向上翻半页",
    "half page down": "向下翻半页",
    "previous message": "上一条消息",
    "next message": "下一条消息",
    "first message": "第一条消息",
    "last message": "最后一条消息",
    "toggle layout": "切换布局",
    "copy message": "复制消息",
    "revert message": "回退消息",
    "exit the app": "退出应用",
}

UI_PATTERNS = [
    r'toast\.NewSuccessToast\("([^"]+)"',
    r'toast\.NewErrorToast\("([^"]+)"',
    r'toast\.New\w*\([^,]*?"([^"]+)"',
    r'WithTitle\("([^"]+)"',
    r'WithDescription\("([^"]+)"',
    r'WithLabel\("([^"]+)"',
    r'Title:\s*"([^"]+)"',
    r'Label:\s*"([^"]+)"',
    r'Placeholder:\s*"([^"]+)"',
    r'Description:\s*"([^"]+)"',
    r'ButtonText:\s*"([^"]+)"',
    r'HelpText:\s*"([^"]+)"',
    r'message\s*:?=\s*"([^"]+)"',
    r'msg\s*:?=\s*"([^"]+)"',
]


def process_file(path, translations):
    with open(path, encoding='utf-8') as f:
        content = f.read()
    original = content
    for pattern in UI_PATTERNS:
        def replacer(m):
            s = m.group(1)
            t = translations.get(s)
            if t and t != s:
                return m.group(0).replace('"' + s + '"', '"' + t + '"', 1)
            return m.group(0)
        content = re.sub(pattern, replacer, content)
    if content != original:
        with open(path, 'w', encoding='utf-8') as f:
            f.write(content)
        return True
    return False


def main():
    tui_dir = sys.argv[1] if len(sys.argv) > 1 else 'internal'
    changed = 0
    for root, _, fnames in os.walk(tui_dir):
        for fn in fnames:
            if not fn.endswith('.go') or fn.endswith('_test.go'):
                continue
            fp = os.path.join(root, fn)
            if process_file(fp, MANUAL):
                changed += 1
                print(f"  已更新: {fp}")
    print(f"共更新 {changed} 个文件")


if __name__ == '__main__':
    main()
