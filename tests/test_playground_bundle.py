import os
import sys
import pytest

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from scripts.build_playground_bundle import generate_bundle, BUNDLE_OUTPUT_PATH


def test_generate_bundle_structure():
    bundle = generate_bundle()
    assert "chapters" in bundle
    assert "exercises" in bundle
    assert "validator_code" in bundle
    assert "stats" in bundle

    assert len(bundle["chapters"]) == 13
    assert bundle["stats"]["totalChapters"] == 13
    assert bundle["stats"]["totalExercises"] == 56
    assert len(bundle["exercises"]) == 56


def test_sample_exercises_have_required_fields():
    bundle = generate_bundle()
    sample = bundle["exercises"]["primitives01"]
    assert sample["name"] == "primitives01"
    assert sample["title"] == "Terraform Configuration Block"
    assert sample["chapterName"] == "01_primitives"
    assert "required_version" in sample["starterCode"]
    assert "required_version" in sample["solutionCode"]
    assert len(sample["hints"]) >= 2


def test_bundle_file_exists_and_valid_json():
    assert os.path.exists(BUNDLE_OUTPUT_PATH)
    with open(BUNDLE_OUTPUT_PATH, "r", encoding="utf-8") as f:
        import json
        data = json.load(f)
        assert len(data["exercises"]) == 56
