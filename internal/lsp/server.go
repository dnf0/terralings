package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/dnf0/terralings/internal/diagnostics"
	"github.com/dnf0/terralings/internal/models"
	"github.com/dnf0/terralings/internal/runner"
	"github.com/dnf0/terralings/internal/state"
	"github.com/dnf0/terralings/internal/version"
)

// Server is the Terralings JSON-RPC 2.0 Language Server.
type Server struct {
	runner     *runner.Runner
	manifest   *models.Manifest
	store      *state.Store
	writerMu   sync.Mutex
	writer     io.Writer
	mu         sync.Mutex
	isShutdown bool
}

// NewServer constructs a new LSP Server instance.
func NewServer(r *runner.Runner, m *models.Manifest, s *state.Store) *Server {
	return &Server{
		runner:   r,
		manifest: m,
		store:    s,
	}
}

// Run starts the LSP server on stdio or provided in/out streams.
func (s *Server) Run(in io.Reader, out io.Writer) error {
	return s.RunWithContext(context.Background(), in, out)
}

// RunWithContext starts the LSP server with context cancellation.
func (s *Server) RunWithContext(ctx context.Context, in io.Reader, out io.Writer) error {
	s.writerMu.Lock()
	s.writer = out
	s.writerMu.Unlock()

	reader := bufio.NewReader(in)

	type readResult struct {
		msg []byte
		err error
	}

	msgChan := make(chan readResult, 1)

	go func() {
		for {
			body, err := s.readFramedMessage(reader)
			if err != nil {
				msgChan <- readResult{err: err}
				return
			}
			msgChan <- readResult{msg: body}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil

		case res := <-msgChan:
			if res.err != nil {
				if res.err == io.EOF || strings.Contains(res.err.Error(), "closed") {
					return nil
				}
				return res.err
			}

			shouldExit, err := s.handleRawMessage(res.msg)
			if err != nil {
				// Continue serving even if single message had an error
				continue
			}
			if shouldExit {
				return nil
			}
		}
	}
}

func (s *Server) readFramedMessage(reader *bufio.Reader) ([]byte, error) {
	contentLength := 0

	// 1. Read header fields until empty line
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headerName := strings.TrimSpace(parts[0])
			headerVal := strings.TrimSpace(parts[1])
			if strings.EqualFold(headerName, "Content-Length") {
				cl, err := strconv.Atoi(headerVal)
				if err == nil {
					contentLength = cl
				}
			}
		}
	}

	if contentLength <= 0 {
		return nil, fmt.Errorf("invalid or missing Content-Length header")
	}

	// 2. Read exact payload body
	body := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, body); err != nil {
		return nil, err
	}

	return body, nil
}

func (s *Server) sendResponse(id *json.RawMessage, result interface{}, rpcErr *RPCError) error {
	resp := RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
		Error:   rpcErr,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))

	s.writerMu.Lock()
	defer s.writerMu.Unlock()

	if _, err := s.writer.Write([]byte(header)); err != nil {
		return err
	}
	_, err = s.writer.Write(data)
	return err
}

func (s *Server) sendNotification(method string, params interface{}) error {
	notif := RPCNotification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))

	s.writerMu.Lock()
	defer s.writerMu.Unlock()

	if _, err := s.writer.Write([]byte(header)); err != nil {
		return err
	}
	_, err = s.writer.Write(data)
	return err
}

func (s *Server) handleRawMessage(raw []byte) (bool, error) {
	var msg RPCMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		_ = s.sendResponse(nil, nil, &RPCError{
			Code:    ParseError,
			Message: fmt.Sprintf("Failed to parse JSON-RPC: %v", err),
		})
		return false, nil
	}

	switch msg.Method {
	case "initialize":
		res := InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   1, // Full document sync
				HoverProvider:      true,
				CodeActionProvider: true,
			},
			ServerInfo: &ServerInfo{
				Name:    "terralings-lsp",
				Version: version.Number,
			},
		}
		_ = s.sendResponse(msg.ID, res, nil)
		return false, nil

	case "initialized":
		// Client confirmation notification, no response required
		return false, nil

	case "textDocument/didOpen":
		var params DidOpenTextDocumentParams
		if err := json.Unmarshal(msg.Params, &params); err == nil {
			s.publishDiagnosticsForURI(params.TextDocument.URI)
		}
		return false, nil

	case "textDocument/didSave":
		var params DidSaveTextDocumentParams
		if err := json.Unmarshal(msg.Params, &params); err == nil {
			s.publishDiagnosticsForURI(params.TextDocument.URI)
		}
		return false, nil

	case "textDocument/hover":
		var params HoverParams
		if err := json.Unmarshal(msg.Params, &params); err != nil {
			_ = s.sendResponse(msg.ID, nil, &RPCError{
				Code:    InvalidParams,
				Message: "Invalid HoverParams",
			})
			return false, nil
		}
		hover := s.handleHover(params)
		_ = s.sendResponse(msg.ID, hover, nil)
		return false, nil

	case "textDocument/codeAction":
		var params CodeActionParams
		if err := json.Unmarshal(msg.Params, &params); err != nil {
			_ = s.sendResponse(msg.ID, nil, &RPCError{
				Code:    InvalidParams,
				Message: "Invalid CodeActionParams",
			})
			return false, nil
		}
		actions := s.handleCodeAction(params)
		_ = s.sendResponse(msg.ID, actions, nil)
		return false, nil

	case "shutdown":
		s.mu.Lock()
		s.isShutdown = true
		s.mu.Unlock()
		_ = s.sendResponse(msg.ID, nil, nil)
		return false, nil

	case "exit":
		return true, nil

	default:
		if msg.ID != nil {
			_ = s.sendResponse(msg.ID, nil, &RPCError{
				Code:    MethodNotFound,
				Message: fmt.Sprintf("Method '%s' not found", msg.Method),
			})
		}
		return false, nil
	}
}

func (s *Server) normalizeURItoPath(rawURI string) string {
	rawURI = strings.TrimPrefix(rawURI, "file://")
	decoded, err := url.PathUnescape(rawURI)
	if err != nil {
		decoded = rawURI
	}
	return filepath.Clean(filepath.ToSlash(decoded))
}

func (s *Server) findExerciseByURI(rawURI string) *models.Exercise {
	if s.manifest == nil {
		return nil
	}

	norm := s.normalizeURItoPath(rawURI)

	for _, ex := range s.manifest.AllExercises() {
		exNorm := filepath.Clean(filepath.ToSlash(ex.Path))
		if norm == exNorm || strings.HasSuffix(norm, exNorm) || strings.HasSuffix(exNorm, norm) {
			target := ex
			return &target
		}

		// Also check by exercise name match in path
		base := filepath.Base(norm)
		ext := filepath.Ext(base)
		nameWithoutExt := strings.TrimSuffix(base, ext)
		if nameWithoutExt == ex.Name {
			target := ex
			return &target
		}
	}

	return nil
}

func (s *Server) publishDiagnosticsForURI(uri string) {
	ex := s.findExerciseByURI(uri)
	if ex == nil {
		_ = s.sendNotification("textDocument/publishDiagnostics", PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: []Diagnostic{},
		})
		return
	}

	var res runner.RunResult
	if s.runner != nil {
		res = s.runner.Run(*ex)
		if s.store != nil {
			_ = s.store.RecordAttempt(ex.Name, ex.ChapterName, res.Passed)
		}
	}

	rawOutput := res.Output
	if res.Error != "" {
		rawOutput = rawOutput + "\n" + res.Error
	}

	parsedDiags := diagnostics.Parse(rawOutput, *ex)
	lspDiags := make([]Diagnostic, 0, len(parsedDiags))

	for _, d := range parsedDiags {
		startLine := d.Line - 1
		if startLine < 0 {
			startLine = 0
		}
		startCol := d.Column - 1
		if startCol < 0 {
			startCol = 0
		}

		endLine := d.EndLine - 1
		if endLine < startLine {
			endLine = startLine
		}
		endCol := d.EndColumn - 1
		if endCol <= startCol {
			endCol = startCol + 80
		}

		sev := SeverityError
		switch d.Severity {
		case diagnostics.SeverityWarning:
			sev = SeverityWarning
		case diagnostics.SeverityInfo:
			sev = SeverityInformation
		}

		msg := d.Summary
		if d.Detail != "" && d.Detail != d.Summary {
			msg = fmt.Sprintf("%s\n\n%s", d.Summary, d.Detail)
		}

		lspDiags = append(lspDiags, Diagnostic{
			Range: Range{
				Start: Position{Line: startLine, Character: startCol},
				End:   Position{Line: endLine, Character: endCol},
			},
			Severity: sev,
			Source:   "terralings",
			Message:  msg,
		})
	}

	_ = s.sendNotification("textDocument/publishDiagnostics", PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: lspDiags,
	})
}

func (s *Server) handleHover(params HoverParams) Hover {
	ex := s.findExerciseByURI(params.TextDocument.URI)
	if ex == nil {
		return Hover{
			Contents: MarkupContent{
				Kind:  "markdown",
				Value: "",
			},
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### %s (`%s`)\n\n", ex.Title, ex.Name))
	sb.WriteString(fmt.Sprintf("**Chapter**: `%s` | **Mode**: `%s`\n\n", ex.ChapterName, ex.Mode))
	sb.WriteString("---\n\n")

	if len(ex.Hints) > 0 {
		sb.WriteString("#### 💡 Hints\n")
		for i, hint := range ex.Hints {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, hint))
		}
	} else {
		sb.WriteString("*No hints available for this exercise.*\n")
	}

	return Hover{
		Contents: MarkupContent{
			Kind:  "markdown",
			Value: sb.String(),
		},
	}
}

func (s *Server) handleCodeAction(params CodeActionParams) []CodeAction {
	ex := s.findExerciseByURI(params.TextDocument.URI)
	if ex == nil {
		return []CodeAction{}
	}

	var actions []CodeAction

	if len(ex.Hints) > 0 {
		actions = append(actions, CodeAction{
			Title: "Show next hint",
			Kind:  "source",
			Command: &Command{
				Title:     "Show hint",
				Command:   "terralings.showHint",
				Arguments: []interface{}{ex.Name},
			},
		})
	}

	return actions
}
