import pytest
import sys
import os

# Add root directory to sys.path so we can import docs.assets.playground.hcl_validator
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from docs.assets.playground.hcl_validator import (
    check_markers,
    parse_hcl,
    validate_exercise,
)


def test_check_markers_detects_not_done_and_todo():
    code_with_not_done = '# I AM NOT DONE\nterraform {\n  required_version = ">= 1.6.0"\n}'
    assert check_markers(code_with_not_done) is True

    code_with_todo = 'terraform {\n  # TODO: specify version\n  required_version = "___"\n}'
    assert check_markers(code_with_todo) is True

    code_with_comment_todo_only = '# TODO (What): complete\nterraform {}'
    assert check_markers(code_with_comment_todo_only) is True

    clean_code = 'terraform {\n  required_version = ">= 1.6.0"\n}'
    assert check_markers(clean_code) is False


def test_parse_hcl_extracts_blocks_and_attributes():
    hcl = '''
    terraform {
      required_version = ">= 1.6.0"
      required_providers {
        local = {
          source  = "hashicorp/local"
          version = "~> 2.0"
        }
      }
    }
    variable "environment" {
      type    = string
      default = "dev"
    }
    locals {
      service_name = "auth"
    }
    '''
    ast = parse_hcl(hcl)
    assert "terraform" in ast
    assert ast["terraform"][0]["required_version"] == ">= 1.6.0"
    assert "variable" in ast
    assert ast["variable"][0]["_label"] == "environment"
    assert ast["variable"][0]["default"] == "dev"
    assert "locals" in ast
    assert ast["locals"][0]["service_name"] == "auth"


def test_validate_exercise_primitives01_success_and_failure():
    starter_fail = 'terraform {\n  # I AM NOT DONE\n  required_version = ">= 1.6.0"\n}'
    res_fail = validate_exercise(starter_fail, "primitives01", {})
    assert res_fail["passed"] is False
    assert "NOT DONE" in res_fail["error"] or "incomplete" in res_fail["error"].lower()

    valid_sol = '''
    terraform {
      required_version = ">= 1.6.0"
      required_providers {
        local = {
          source  = "hashicorp/local"
          version = "~> 2.0"
        }
      }
    }
    '''
    res_pass = validate_exercise(valid_sol, "primitives01", {})
    assert res_pass["passed"] is True
    assert res_pass["error"] is None
