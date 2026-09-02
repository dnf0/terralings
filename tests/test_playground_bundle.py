import json
import os
import sys
import pytest

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from docs.assets.playground.hcl_validator import validate_exercise
from scripts.build_playground_bundle import BUNDLE_OUTPUT_PATH, generate_bundle


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
        data = json.load(f)
        assert len(data["exercises"]) == 56


def test_all_56_reference_solutions_pass_validation():
    bundle = generate_bundle()
    failures = []
    for ex_id, ex in bundle["exercises"].items():
        sol = ex["solutionCode"]
        assert sol, f"Exercise {ex_id} has empty solution code"
        res = validate_exercise(sol, ex_id, ex.get("rules", {}))
        if not res["passed"]:
            failures.append(f"{ex_id}: {res.get('error')}")

    assert len(failures) == 0, f"The following reference solutions failed validation: {failures}"


def test_all_starter_templates_fail_validation_if_incomplete():
    bundle = generate_bundle()
    failed_starters = 0
    for ex_id, ex in bundle["exercises"].items():
        starter = ex["starterCode"]
        res = validate_exercise(starter, ex_id, ex.get("rules", {}))
        if not res["passed"]:
            failed_starters += 1

    # Most starter templates contain TODOs or placeholders and must fail
    assert failed_starters >= 50, f"Expected most starter templates to fail due to markers, but only {failed_starters} failed"
