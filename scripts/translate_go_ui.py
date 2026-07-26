#!/usr/bin/env python3
"""扫描 opencode TUI Go 源码,提取 UI 字符串,AI 逐个翻译,替换源码。
只翻译明确的 UI 模式(toast/Title/Label/Placeholder 等),避免误改命令名/日志/key。"""
import re, os, json, urllib.request, sys, time

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


def translate_single(text):
    if not isinstance(text, str) or not text.strip():
        return text
    prompt = (
        "Translate the following UI string to Simplified Chinese for a terminal app. "
        "Reply with ONLY the Chinese translation. Preserve {placeholders} like {name} unchanged. "
        "Keep concise, natural UI Chinese. No explanation, no quotes.\n\n" + text
    )
    body = json.dumps({
        "model": "mimo-v2.5-free",
        "messages": [{"role": "user", "content": prompt}],
        "max_tokens": 2000,
    }).encode()
    req = urllib.request.Request(
        "https://opencode.ai/zen/v1/chat/completions",
        data=body, headers={"Content-Type": "application/json"})
    for attempt in range(3):
        try:
            with urllib.request.urlopen(req, timeout=60) as r:
                d = json.loads(r.read())
            content = d['choices'][0]['message'].get('content')
            return content.strip() if content else text
        except Exception as e:
            if attempt < 2:
                time.sleep(3)
                continue
            print(f"  翻译失败: {e}", file=sys.stderr)
            return text


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
    all_strings = []
    files = []
    for root, _, fnames in os.walk(tui_dir):
        for fn in fnames:
            if not fn.endswith('.go') or fn.endswith('_test.go'):
                continue
            fp = os.path.join(root, fn)
            try:
                with open(fp, encoding='utf-8') as f:
                    content = f.read()
            except Exception:
                continue
            for pattern in UI_PATTERNS:
                for m in re.finditer(pattern, content):
                    all_strings.append(m.group(1))
            files.append(fp)
    unique = list(dict.fromkeys(all_strings))
    print(f"提取 {len(all_strings)} 个候选(去重 {len(unique)} 个)")
    if not unique:
        print("无 UI 字符串,退出")
        return
    print("调 AI 逐个翻译...")
    translations = {}
    for i, s in enumerate(unique):
        t = translate_single(s)
        if t and t != s:
            translations[s] = t
            print(f"  [{i+1}/{len(unique)}] {s!r} -> {t!r}")
        else:
            print(f"  [{i+1}/{len(unique)}] {s!r} (保留原文)")
        time.sleep(0.5)
    print(f"\n获得 {len(translations)} 个翻译")
    changed = 0
    for fp in files:
        if process_file(fp, translations):
            changed += 1
            print(f"  已更新: {fp}")
    print(f"共更新 {changed} 个文件")


if __name__ == '__main__':
    main()
