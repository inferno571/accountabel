import os

def fix_imports(filepath):
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        return

    orig_content = content
    content = content.replace('github.com/yincongcyincong/MuseBot', 'github.com/yincongcyincong/MuseBot')
    content = content.replace('muse_bot.log', 'muse_bot.log')

    if content != orig_content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Fixed {filepath}")

def main():
    root_dir = r"c:\Users\saksham\Downloads\MuseBot-main\MuseBot-main"
    for dirpath, dirnames, filenames in os.walk(root_dir):
        if '.git' in dirpath:
            continue
        for filename in filenames:
            if filename.endswith('.exe') or filename.endswith('.png') or filename.endswith('.jpg'):
                continue
            filepath = os.path.join(dirpath, filename)
            fix_imports(filepath)

if __name__ == '__main__':
    main()
