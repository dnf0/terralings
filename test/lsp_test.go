package test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dnf0/terralings/internal/detector"
	"github.com/dnf0/terralings/internal/lsp"
	"github.com/dnf0/terralings/internal/manifest"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/state"
	"github.com/dnf0/terralings/internal/watcher"
)

// lspSafeBuffer is a thread-safe bytes.Buffer for concurrent test assertions
type lspSafeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *lspSafeBuffer) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *lspSafeBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

// lspClientHelper sends framed LSP messages and reads framed responses
type lspClientHelper struct {
	inPipe  io.WriteCloser
	outPipe *bufio.Reader
	mu      sync.Mutex
}

func newLSPClientHelper(serverIn io.WriteCloser, serverOut io.Reader) *lspClientHelper {
	return &lspClientHelper{
		inPipe:  serverIn,
		outPipe: bufio.NewReader(serverOut),
	}
}

func (c *lspClientHelper) Send(payload interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	if _, err := c.inPipe.Write([]byte(header)); err != nil {
		return err
	}
	_, err = c.inPipe.Write(data)
	return err
}

func (c *lspClientHelper) ReadMessage() ([]byte, error) {
	// Read headers
	contentLength := 0
	for {
		line, err := c.outPipe.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			// End of headers
			break
		}
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				_, _ = fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &contentLength)
			}
		}
	}

	if contentLength <= 0 {
		return nil, fmt.Errorf("invalid content length: %d", contentLength)
	}

	body := make([]byte, contentLength)
	_, err := io.ReadFull(c.outPipe, body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func TestLSP_Initialize(t *testing.T) {
	clientInR, serverInW := io.Pipe()
	serverInR, clientInW := io.Pipe()

	client := newLSPClientHelper(clientInW, clientInR)
	m := manifest.GetManifest()
	r := runner.NewRunner("")
	srv := lsp.NewServer(r, m, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverDone := make(chan error, 1)
	go func() {
		serverDone <- srv.RunWithContext(ctx, serverInR, serverInW)
	}()

	initReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"processId": os.Getpid(),
			"rootUri":   "file://" + t.TempDir(),
		},
	}

	if err := client.Send(initReq); err != nil {
		t.Fatalf("Failed to send initialize request: %v", err)
	}

	respBytes, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read initialize response: %v", err)
	}

	var resp struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			Capabilities struct {
				TextDocumentSync   int  `json:"textDocumentSync"`
				HoverProvider      bool `json:"hoverProvider"`
				CodeActionProvider bool `json:"codeActionProvider"`
			} `json:"capabilities"`
			ServerInfo struct {
				Name string `json:"name"`
			} `json:"serverInfo"`
		} `json:"result"`
	}

	if err := json.Unmarshal(respBytes, &resp); err != nil {
		t.Fatalf("Failed to parse initialize response: %v. Raw: %s", err, string(respBytes))
	}

	if resp.ID != 1 {
		t.Errorf("Expected response id 1, got %d", resp.ID)
	}
	if resp.Result.Capabilities.TextDocumentSync == 0 {
		t.Errorf("Expected TextDocumentSync capability > 0, got %d", resp.Result.Capabilities.TextDocumentSync)
	}
	if !resp.Result.Capabilities.HoverProvider {
		t.Error("Expected HoverProvider capability to be true")
	}
	if !resp.Result.Capabilities.CodeActionProvider {
		t.Error("Expected CodeActionProvider capability to be true")
	}
	if resp.Result.ServerInfo.Name == "" {
		t.Error("Expected serverInfo name to be populated")
	}

	// Clean shutdown
	cancel()
	_ = clientInW.Close()
	_ = serverInW.Close()
}

func TestLSP_DidOpenAndDidSave_PublishesDiagnostics(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on system PATH; skipping LSP test")
	}

	tmpDir := t.TempDir()
	exFile := filepath.Join(tmpDir, "exercises", "01_primitives", "primitives01.tf")
	_ = os.MkdirAll(filepath.Dir(exFile), 0755)

	// File has # I AM NOT DONE marker
	_ = os.WriteFile(exFile, []byte("# I AM NOT DONE\nterraform {\n  required_version = \">= 1.6.0\"\n}\n"), 0644)

	customEx := models.Exercise{
		Name:        "primitives01",
		Title:       "Terraform Configuration Block",
		Path:        exFile,
		ChapterName: "01_primitives",
		Mode:        models.ModeValidate,
	}

	customManifest := &models.Manifest{
		Chapters: []models.Chapter{
			{
				Number:    1,
				Name:      "01_primitives",
				Title:     "Primitives",
				Exercises: []models.Exercise{customEx},
			},
		},
	}

	statePath := filepath.Join(tmpDir, "state.json")
	store, err := state.NewStore(statePath)
	if err != nil {
		t.Fatalf("Failed to initialize state: %v", err)
	}

	clientInR, serverInW := io.Pipe()
	serverInR, clientInW := io.Pipe()

	client := newLSPClientHelper(clientInW, clientInR)
	r := runner.NewRunner(bin)
	srv := lsp.NewServer(r, customManifest, store)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = srv.RunWithContext(ctx, serverInR, serverInW)
	}()

	fileURI := "file://" + filepath.ToSlash(exFile)

	// 1. Send didOpen notification
	didOpenNotification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "textDocument/didOpen",
		"params": map[string]interface{}{
			"textDocument": map[string]interface{}{
				"uri":        fileURI,
				"languageId": "terraform",
				"version":    1,
				"text":       "# I AM NOT DONE\nterraform {\n  required_version = \">= 1.6.0\"\n}\n",
			},
		},
	}

	if err := client.Send(didOpenNotification); err != nil {
		t.Fatalf("Failed to send didOpen: %v", err)
	}

	msgBytes, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read didOpen diagnostics notification: %v", err)
	}

	var diagNotif struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  struct {
			URI         string `json:"uri"`
			Diagnostics []struct {
				Severity int    `json:"severity"`
				Message  string `json:"message"`
				Range    struct {
					Start struct {
						Line int `json:"line"`
					} `json:"start"`
				} `json:"range"`
			} `json:"diagnostics"`
		} `json:"params"`
	}

	if err := json.Unmarshal(msgBytes, &diagNotif); err != nil {
		t.Fatalf("Failed to parse diagnostic notification: %v. Raw: %s", err, string(msgBytes))
	}

	if diagNotif.Method != "textDocument/publishDiagnostics" {
		t.Fatalf("Expected textDocument/publishDiagnostics, got %s", diagNotif.Method)
	}

	if len(diagNotif.Params.Diagnostics) == 0 {
		t.Fatal("Expected at least 1 diagnostic for '# I AM NOT DONE' marker")
	}

	hasMarkerDiag := false
	for _, d := range diagNotif.Params.Diagnostics {
		if strings.Contains(strings.ToLower(d.Message), "not done") || strings.Contains(strings.ToLower(d.Message), "not finished") {
			hasMarkerDiag = true
			if d.Range.Start.Line != 0 {
				t.Errorf("Expected marker diagnostic at line 0 (0-indexed), got %d", d.Range.Start.Line)
			}
		}
	}

	if !hasMarkerDiag {
		t.Errorf("Expected marker diagnostic message, got %+v", diagNotif.Params.Diagnostics)
	}

	// 2. Send didSave notification
	didSaveNotification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "textDocument/didSave",
		"params": map[string]interface{}{
			"textDocument": map[string]interface{}{
				"uri": fileURI,
			},
		},
	}

	if err := client.Send(didSaveNotification); err != nil {
		t.Fatalf("Failed to send didSave: %v", err)
	}

	msgBytes2, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read didSave diagnostics notification: %v", err)
	}

	var diagNotif2 struct {
		JSONRPC string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  struct {
			URI         string `json:"uri"`
			Diagnostics []struct {
				Severity int `json:"severity"`
			} `json:"diagnostics"`
		} `json:"params"`
	}
	if err := json.Unmarshal(msgBytes2, &diagNotif2); err != nil {
		t.Fatalf("Failed to parse didSave diagnostic: %v", err)
	}
	if diagNotif2.Method != "textDocument/publishDiagnostics" {
		t.Fatalf("Expected publishDiagnostics on didSave, got %s", diagNotif2.Method)
	}

	// Verify attempt recorded in state
	exSt := store.GetExerciseState("primitives01")
	if exSt == nil || exSt.Attempts < 2 {
		t.Errorf("Expected at least 2 attempts in state store after didOpen and didSave, got %+v", exSt)
	}

	cancel()
	_ = clientInW.Close()
	_ = serverInW.Close()
}

func TestLSP_Hover(t *testing.T) {
	m := manifest.GetManifest()
	r := runner.NewRunner("")
	srv := lsp.NewServer(r, m, nil)

	clientInR, serverInW := io.Pipe()
	serverInR, clientInW := io.Pipe()

	client := newLSPClientHelper(clientInW, clientInR)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = srv.RunWithContext(ctx, serverInR, serverInW)
	}()

	all := m.AllExercises()
	if len(all) == 0 {
		t.Fatal("Manifest has no exercises")
	}
	targetEx := all[0]

	cwd, _ := os.Getwd()
	absPath := filepath.Join(cwd, targetEx.Path)
	fileURI := "file://" + filepath.ToSlash(absPath)

	hoverReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      42,
		"method":  "textDocument/hover",
		"params": map[string]interface{}{
			"textDocument": map[string]interface{}{
				"uri": fileURI,
			},
			"position": map[string]interface{}{
				"line":      0,
				"character": 0,
			},
		},
	}

	if err := client.Send(hoverReq); err != nil {
		t.Fatalf("Failed to send hover request: %v", err)
	}

	respBytes, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read hover response: %v", err)
	}

	var resp struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			Contents struct {
				Kind  string `json:"kind"`
				Value string `json:"value"`
			} `json:"contents"`
		} `json:"result"`
	}

	if err := json.Unmarshal(respBytes, &resp); err != nil {
		t.Fatalf("Failed to unmarshal hover response: %v. Raw: %s", err, string(respBytes))
	}

	if resp.ID != 42 {
		t.Errorf("Expected response id 42, got %d", resp.ID)
	}

	val := resp.Result.Contents.Value
	if !strings.Contains(val, targetEx.Title) && !strings.Contains(val, targetEx.Name) {
		t.Errorf("Expected hover content to contain exercise title or name, got: %s", val)
	}
	if len(targetEx.Hints) > 0 && !strings.Contains(val, targetEx.Hints[0]) {
		t.Errorf("Expected hover content to contain hint '%s', got: %s", targetEx.Hints[0], val)
	}

	cancel()
	_ = clientInW.Close()
	_ = serverInW.Close()
}

func TestLSP_CodeAction(t *testing.T) {
	m := manifest.GetManifest()
	r := runner.NewRunner("")
	srv := lsp.NewServer(r, m, nil)

	clientInR, serverInW := io.Pipe()
	serverInR, clientInW := io.Pipe()

	client := newLSPClientHelper(clientInW, clientInR)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = srv.RunWithContext(ctx, serverInR, serverInW)
	}()

	all := m.AllExercises()
	if len(all) == 0 {
		t.Fatal("Manifest has no exercises")
	}
	targetEx := all[0]
	cwd, _ := os.Getwd()
	absPath := filepath.Join(cwd, targetEx.Path)
	fileURI := "file://" + filepath.ToSlash(absPath)

	caReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      55,
		"method":  "textDocument/codeAction",
		"params": map[string]interface{}{
			"textDocument": map[string]interface{}{
				"uri": fileURI,
			},
			"range": map[string]interface{}{
				"start": map[string]int{"line": 0, "character": 0},
				"end":   map[string]int{"line": 0, "character": 10},
			},
			"context": map[string]interface{}{
				"diagnostics": []map[string]interface{}{
					{
						"severity": 2,
						"message":  "Exercise is not finished ('I AM NOT DONE' marker present)",
						"range": map[string]interface{}{
							"start": map[string]int{"line": 0, "character": 0},
							"end":   map[string]int{"line": 0, "character": 15},
						},
					},
				},
			},
		},
	}

	if err := client.Send(caReq); err != nil {
		t.Fatalf("Failed to send codeAction request: %v", err)
	}

	respBytes, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read codeAction response: %v", err)
	}

	var resp struct {
		JSONRPC string           `json:"jsonrpc"`
		ID      int              `json:"id"`
		Result  []lsp.CodeAction `json:"result"`
	}

	if err := json.Unmarshal(respBytes, &resp); err != nil {
		t.Fatalf("Failed to unmarshal codeAction response: %v", err)
	}

	if resp.ID != 55 {
		t.Errorf("Expected id 55, got %d", resp.ID)
	}

	if len(resp.Result) == 0 {
		t.Error("Expected code action suggestions, got 0")
	}

	hasQuickFix := false
	for _, action := range resp.Result {
		if strings.Contains(strings.ToLower(action.Title), "remove") {
			hasQuickFix = true
		}
	}
	if !hasQuickFix {
		t.Errorf("Expected quickfix action to remove marker, got %+v", resp.Result)
	}

	cancel()
	_ = clientInW.Close()
	_ = serverInW.Close()
}

func TestLSP_MethodNotFound(t *testing.T) {
	m := manifest.GetManifest()
	r := runner.NewRunner("")
	srv := lsp.NewServer(r, m, nil)

	clientInR, serverInW := io.Pipe()
	serverInR, clientInW := io.Pipe()

	client := newLSPClientHelper(clientInW, clientInR)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = srv.RunWithContext(ctx, serverInR, serverInW)
	}()

	unknownReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      999,
		"method":  "workspace/unknownMethod",
	}

	if err := client.Send(unknownReq); err != nil {
		t.Fatalf("Failed to send unknown method request: %v", err)
	}

	respBytes, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read error response: %v", err)
	}

	var resp struct {
		JSONRPC string        `json:"jsonrpc"`
		ID      int           `json:"id"`
		Error   *lsp.RPCError `json:"error"`
	}

	if err := json.Unmarshal(respBytes, &resp); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("Expected RPC error for unknown method")
	}
	if resp.Error.Code != lsp.MethodNotFound {
		t.Errorf("Expected error code %d (MethodNotFound), got %d", lsp.MethodNotFound, resp.Error.Code)
	}

	cancel()
	_ = clientInW.Close()
	_ = serverInW.Close()
}

func TestLSP_Shutdown(t *testing.T) {
	m := manifest.GetManifest()
	r := runner.NewRunner("")
	srv := lsp.NewServer(r, m, nil)

	clientInR, serverInW := io.Pipe()
	serverInR, clientInW := io.Pipe()

	client := newLSPClientHelper(clientInW, clientInR)

	ctx := context.Background()
	serverDone := make(chan error, 1)

	go func() {
		serverDone <- srv.RunWithContext(ctx, serverInR, serverInW)
	}()

	// 1. Send shutdown request
	shutdownReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      100,
		"method":  "shutdown",
		"params":  nil,
	}

	if err := client.Send(shutdownReq); err != nil {
		t.Fatalf("Failed to send shutdown request: %v", err)
	}

	respBytes, err := client.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read shutdown response: %v", err)
	}

	var resp struct {
		JSONRPC string      `json:"jsonrpc"`
		ID      int         `json:"id"`
		Result  interface{} `json:"result"`
	}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		t.Fatalf("Failed to unmarshal shutdown response: %v", err)
	}
	if resp.ID != 100 {
		t.Errorf("Expected response id 100, got %d", resp.ID)
	}

	// 2. Send exit notification
	exitNotif := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "exit",
	}
	if err := client.Send(exitNotif); err != nil {
		t.Fatalf("Failed to send exit notification: %v", err)
	}

	select {
	case err := <-serverDone:
		if err != nil {
			t.Fatalf("Expected clean zero-exit shutdown, got: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Server did not shut down after exit notification within 2s")
	}

	_ = clientInW.Close()
	_ = serverInW.Close()
}

func TestCLI_WatchJSON(t *testing.T) {
	bin, err := detector.DetectBinary("")
	if err != nil {
		t.Skip("Neither tofu nor terraform found on system PATH; skipping watcher test")
	}

	tmpDir := t.TempDir()
	exFile := filepath.Join(tmpDir, "ex01.tf")
	_ = os.WriteFile(exFile, []byte("# I AM NOT DONE\nterraform {\n  required_version = \">= 1.6.0\"\n}\n"), 0644)

	exercises := []models.Exercise{
		{
			Name:        "ex01",
			ChapterName: "01_test",
			Path:        exFile,
			Mode:        models.ModeValidate,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var out lspSafeBuffer
	r := runner.NewRunner(bin)
	doneChan := make(chan error, 1)

	go func() {
		doneChan <- watcher.RunWatchJSON(ctx, r, exercises, nil, tmpDir, &out)
	}()

	// Wait for NDJSON output
	deadline := time.Now().Add(5 * time.Second)
	var eventsFound bool
	for time.Now().Before(deadline) {
		outStr := out.String()
		if strings.Contains(outStr, `"exercise_start"`) && strings.Contains(outStr, `"exercise_result"`) {
			eventsFound = true
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if !eventsFound {
		t.Fatalf("Expected NDJSON events in stdout, got:\n%s", out.String())
	}

	// Verify line-by-line JSON validity
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	eventTypes := make(map[string]bool)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(line), &parsed); err != nil {
			t.Fatalf("Output line is not valid JSON: %v. Line: %s", err, line)
		}
		if evt, ok := parsed["event"].(string); ok {
			eventTypes[evt] = true
		}
	}

	if !eventTypes["exercise_start"] {
		t.Error("Missing 'exercise_start' event in NDJSON stream")
	}
	if !eventTypes["exercise_result"] {
		t.Error("Missing 'exercise_result' event in NDJSON stream")
	}

	cancel()
	select {
	case <-doneChan:
	case <-time.After(1 * time.Second):
		t.Fatal("WatchJSON did not terminate cleanly on context cancel")
	}
}
