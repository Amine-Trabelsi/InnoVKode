import os
import sys

def safe_print(text):
    sys.stdout.buffer.write((text + "\n").encode("utf-8"))

def print_tree_only(
    startpath, 
    ignore=("__pycache__", ".git", "env", ".venv", ".mypy_cache", ".pytest_cache"), 
    max_files_preview=1
):
    safe_print("=== PROJECT STRUCTURE (no file contents) ===")
    for root, dirs, files in os.walk(startpath):
        dirs[:] = [d for d in dirs if d not in ignore]

        level = root.replace(startpath, "").count(os.sep)
        indent = " " * 4 * level
        safe_print(f"{indent}{os.path.basename(root)}/")

        subindent = " " * 4 * (level + 1)

        code_files = [f for f in files if f.endswith((".py", ".go"))]
        other_files = [f for f in files if not f.endswith((".py", ".go"))]

        if len(other_files) > max_files_preview:
            other_files = other_files[:max_files_preview] + ["..."]

        for f in code_files + other_files:
            safe_print(f"{subindent}{f}")

def print_tree_with_code_content(
    startpath, 
    ignore=("__pycache__", ".git", "env", ".venv", "venv", ".mypy_cache", ".pytest_cache"),
    max_files_preview=15,
    max_code_lines=500
):
    safe_print("\n=== PROJECT STRUCTURE WITH .py/.go CONTENTS (truncated) ===")
    for root, dirs, files in os.walk(startpath):
        dirs[:] = [d for d in dirs if d not in ignore]

        level = root.replace(startpath, "").count(os.sep)
        indent = " " * 4 * level
        safe_print(f"{indent}{os.path.basename(root)}/")

        subindent = " " * 4 * (level + 1)

        code_files = [f for f in files if f.endswith((".py", ".go"))]
        other_files = [f for f in files if not f.endswith((".py", ".go"))]

        if len(other_files) > max_files_preview:
            other_files = other_files[:max_files_preview] + ["..."]

        for f in code_files + other_files:
            safe_print(f"{subindent}{f}")
            if f.endswith((".py", ".go")) and f != "...":
                filepath = os.path.join(root, f)
                try:
                    with open(filepath, encoding="utf-8") as file:
                        for i, line in enumerate(file):
                            if i >= max_code_lines:
                                safe_print(f"{subindent}    ...")
                                break
                            safe_print(f"{subindent}    {line.rstrip()}")
                except Exception as e:
                    safe_print(f"{subindent}    [Could not read file: {e}]")

if __name__ == "__main__":
    root_path = "."
    print_tree_only(root_path)
    print_tree_with_code_content(root_path)
