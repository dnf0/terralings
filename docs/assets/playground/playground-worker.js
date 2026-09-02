/**
 * Web Worker for Terralings Pyodide WebAssembly Runtime.
 *
 * Loads Pyodide v0.26+ in a background Web Worker, mounts the pure-Python
 * in-memory HCL AST and rule validator, and provides sandboxed execution
 * with captured stdout and millisecond timing.
 */

/* global loadPyodide, importScripts */
importScripts("https://cdn.jsdelivr.net/pyodide/v0.26.2/full/pyodide.js");

let pyodide = null;
let bundleData = null;

/**
 * Initialize Pyodide WebAssembly runtime and mount Terralings virtual modules.
 * @param {Object} bundle - Playground bundle containing validator_code.
 */
async function initPyodide(bundle) {
  bundleData = bundle;
  self.postMessage({
    type: "STATUS",
    stage: "loading_pyodide",
    message: "⚡ Initializing Python 3.12 WebAssembly Runtime...",
  });

  pyodide = await loadPyodide({
    indexURL: "https://cdn.jsdelivr.net/pyodide/v0.26.2/full/",
  });

  self.postMessage({
    type: "STATUS",
    stage: "mounting_bundle",
    message: "🔧 Mounting Terralings HCL AST Validator...",
  });

  // Create /lib/terralings virtual package in Pyodide FS
  pyodide.FS.mkdirTree("/lib/terralings");
  pyodide.FS.writeFile("/lib/terralings/__init__.py", "");
  pyodide.FS.writeFile(
    "/lib/terralings/hcl_validator.py",
    (bundle && bundle.validator_code) || ""
  );

  // Setup Python sys.path and validation function wrapper
  await pyodide.runPythonAsync(`
import sys
import io
import time
import traceback

if "/lib" not in sys.path:
    sys.path.insert(0, "/lib")

import terralings.hcl_validator as validator

def run_hcl_validation(user_code_str, exercise_id_str, rules_dict=None):
    try:
        res = validator.validate_exercise(user_code_str, exercise_id_str, rules_dict)
        return {
            "passed": bool(res.get("passed", False)),
            "error": res.get("error"),
            "output": res.get("output", ""),
            "durationMs": float(res.get("duration_ms", 0.0)),
            "line": res.get("line")
        }
    except Exception as e:
        tb = traceback.format_exc()
        return {
            "passed": False,
            "error": f"{type(e).__name__}: {e}",
            "traceback": tb,
            "output": f"✕ Unhandled Validator Exception: {e}",
            "durationMs": 0.0,
            "line": None
        }
`);

  self.postMessage({
    type: "STATUS",
    stage: "ready",
    message: "✅ Ready! WebAssembly HCL engine loaded.",
  });
}

self.onmessage = async function (e) {
  const msg = e.data;
  if (!msg || !msg.type) return;

  if (msg.type === "INIT") {
    try {
      await initPyodide(msg.bundle || {});
    } catch (err) {
      self.postMessage({
        type: "STATUS",
        stage: "error",
        message:
          "Error initializing Pyodide: " +
          (err && err.message ? err.message : String(err)),
      });
    }
  } else if (msg.type === "RUN_EXERCISE") {
    if (!pyodide) {
      self.postMessage({
        type: "RUN_RESULT",
        exerciseId: msg.exerciseId,
        passed: false,
        error: "WebAssembly runtime is still initializing...",
        output: "Please wait for Pyodide to finish loading.",
        durationMs: 0,
        line: null,
      });
      return;
    }

    try {
      pyodide.globals.set("temp_code_str", msg.code || "");
      pyodide.globals.set("temp_exercise_id", msg.exerciseId || "");
      pyodide.globals.set(
        "temp_rules",
        pyodide.toPy(msg.rules || {})
      );

      const resProxy = await pyodide.runPythonAsync(
        "run_hcl_validation(temp_code_str, temp_exercise_id, temp_rules)"
      );
      const resultObj = resProxy.toJs({ dict_converter: Object.fromEntries });

      self.postMessage({
        type: "RUN_RESULT",
        exerciseId: msg.exerciseId,
        passed: resultObj.passed,
        error: resultObj.error,
        output: resultObj.output,
        durationMs: resultObj.durationMs,
        line: resultObj.line,
      });
    } catch (err) {
      self.postMessage({
        type: "RUN_RESULT",
        exerciseId: msg.exerciseId,
        passed: false,
        error: String(err),
        output: "✕ Evaluation runtime error: " + String(err),
        durationMs: 0,
        line: null,
      });
    }
  }
};
