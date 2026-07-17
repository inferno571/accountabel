import os

def replace_in_file(filepath):
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        return

    orig_content = content
    content = content.replace('Accountabel AI', 'Accountabel AI')
    content = content.replace('accountabel_bot', 'accountabel_bot')
    content = content.replace('accountabel', 'accountabel')
    content = content.replace('accountabel_user', 'accountabel_user')
    content = content.replace('accountabel_session', 'accountabel_session')
    content = content.replace('accountabel_signup', 'accountabel_signup')

    if content != orig_content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Updated {filepath}")

def main():
    root_dir = r"c:\Users\saksham\Downloads\Accountabel AI-main\Accountabel AI-main"
    for dirpath, dirnames, filenames in os.walk(root_dir):
        if '.git' in dirpath:
            continue
        for filename in filenames:
            if filename.endswith('.exe') or filename.endswith('.png') or filename.endswith('.jpg'):
                continue
            filepath = os.path.join(dirpath, filename)
            replace_in_file(filepath)

if __name__ == '__main__':
    main()
