#!/usr/bin/env python3
"""Build-time Playground Bundle Generator for Terralings.

Extracts all 13 chapters, 56 exercises, starter files, reference solutions,
progressive hints, and the in-memory HCL AST validator into a single
self-contained JSON asset for zero-install Pyodide browser execution.
"""

from __future__ import annotations

import json
import re
from pathlib import Path
from typing import Any

REPO_ROOT = Path(__file__).resolve().parent.parent
MANIFEST_GO_PATH = REPO_ROOT / "internal" / "manifest" / "manifest.go"
VALIDATOR_PY_PATH = REPO_ROOT / "docs" / "assets" / "playground" / "hcl_validator.py"
BUNDLE_OUTPUT_PATH = REPO_ROOT / "docs" / "assets" / "playground" / "playground-bundle.json"


def parse_go_hints(hints_raw: str) -> list[str]:
    """Parse Go slice string literals, preserving inner quotes correctly."""
    hints = []
    for line in hints_raw.strip().splitlines():
        line = line.strip().rstrip(",")
        if line.startswith('"') and line.endswith('"') and len(line) >= 2:
            inner = line[1:-1].replace(r'\"', '"').replace(r"\\", "\\")
            if inner:
                hints.append(inner)
    return hints


def parse_manifest_go(manifest_path: Path) -> list[dict[str, Any]]:
    """Parse internal/manifest/manifest.go to extract chapters and exercises."""
    content = manifest_path.read_text(encoding="utf-8")

    # Match Chapter structs
    # { Number: 1, Name: "01_primitives", Title: "...", Description: "...", Exercises: [...] }
    chapter_blocks = re.findall(
        r'{\s*Number:\s*(\d+),\s*Name:\s*"([^"]+)",\s*Title:\s*"([^"]+)",\s*Description:\s*"([^"]+)",\s*Exercises:\s*\[\]models\.Exercise\{([\s\S]*?)\n\s*\},?\s*\}',
        content,
    )

    chapters = []

    for ch_num, ch_name, ch_title, ch_desc, ex_block in chapter_blocks:
        chapter_dict: dict[str, Any] = {
            "number": int(ch_num),
            "name": ch_name,
            "title": ch_title,
            "description": ch_desc,
            "exercise_ids": [],
            "exercises": [],
        }

        # Match Exercise structs inside chapter
        # { Name: "...", Title: "...", Path: "...", ChapterName: "...", Hints: []string{ ... }, Mode: ... }
        ex_matches = re.finditer(
            r'{\s*Name:\s*"([^"]+)",\s*Title:\s*"([^"]+)",\s*Path:\s*"([^"]+)",\s*ChapterName:\s*"([^"]+)",\s*Hints:\s*\[\]string\{([\s\S]*?)\},',
            ex_block,
        )

        for ex_m in ex_matches:
            name, title, rel_path, chapter_name, hints_raw = ex_m.groups()
            hints = parse_go_hints(hints_raw)

            chapter_dict["exercise_ids"].append(name)
            chapter_dict["exercises"].append(
                {
                    "id": name,
                    "name": name,
                    "title": title,
                    "path": rel_path,
                    "chapterName": chapter_name,
                    "hints": hints,
                }
            )

        chapters.append(chapter_dict)

    return chapters


def read_file_or_dir(base_path: Path) -> str:
    """Read a single file or concatenate files in a directory."""
    if not base_path.exists():
        return ""
    if base_path.is_file():
        return base_path.read_text(encoding="utf-8")
    if base_path.is_dir():
        # Concatenate .tf, .hcl files in directory
        parts = []
        for file in sorted(base_path.glob("**/*")):
            if file.is_file() and (file.suffix in (".tf", ".hcl", ".json")):
                rel = file.relative_to(base_path)
                parts.append(f"# --- File: {rel} ---\n" + file.read_text(encoding="utf-8"))
        return "\n\n".join(parts)
    return ""


def generate_bundle() -> dict[str, Any]:
    """Generate complete curriculum bundle structure."""
    chapters = parse_manifest_go(MANIFEST_GO_PATH)
    exercises_map: dict[str, Any] = {}

    total_exercises = 0

    for ch in chapters:
        for ex in ch["exercises"]:
            total_exercises += 1
            ex_name = ex["name"]
            starter_path = REPO_ROOT / ex["path"]
            # Solution path mirrors exercise path replacing exercises/ with solutions/
            sol_rel_path = ex["path"].replace("exercises/", "solutions/")
            sol_path = REPO_ROOT / sol_rel_path

            starter_code = read_file_or_dir(starter_path)
            solution_code = read_file_or_dir(sol_path)

            ex_data = {
                "id": ex_name,
                "name": ex_name,
                "title": ex["title"],
                "path": ex["path"],
                "chapter": ex["chapterName"],
                "chapterName": ex["chapterName"],
                "chapter_number": ch["number"],
                "chapterNumber": ch["number"],
                "chapter_title": ch["title"],
                "chapterTitle": ch["title"],
                "filename": Path(ex["path"]).name,
                "hints": ex["hints"],
                "starter_code": starter_code,
                "starterCode": starter_code,
                "solution_code": solution_code,
                "solutionCode": solution_code,
            }
            exercises_map[ex_name] = ex_data

    validator_code = VALIDATOR_PY_PATH.read_text(encoding="utf-8") if VALIDATOR_PY_PATH.exists() else ""

    bundle = {
        "version": 1,
        "stats": {
            "totalChapters": len(chapters),
            "totalExercises": total_exercises,
        },
        "total_chapters": len(chapters),
        "total_exercises": total_exercises,
        "chapters": chapters,
        "exercises": exercises_map,
        "validator_code": validator_code,
    }

    return bundle


def write_bundle() -> None:
    """Generate and write bundle to playground-bundle.json."""
    bundle = generate_bundle()
    BUNDLE_OUTPUT_PATH.parent.mkdir(parents=True, exist_ok=True)
    with open(BUNDLE_OUTPUT_PATH, "w", encoding="utf-8") as f:
        json.dump(bundle, f, indent=2, ensure_ascii=False)
    print(f"✓ Generated playground bundle with {bundle['total_exercises']} exercises at {BUNDLE_OUTPUT_PATH}")


if __name__ == "__main__":
    write_bundle()
