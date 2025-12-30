// Package uiutil provides utility functions for UI message handling.
package uiutil

import (
	"context"
	"errors"
	"log/slog"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"mvdan.cc/sh/v3/shell"
)

// Cursor interface - not used in prepf's bubbletea v1.3.10
// Kept for compatibility with cloned TUI components
type Cursor interface {
	// Cursor() method - implementation depends on bubbletea version
	// Stub for compatibility
}

func CmdHandler(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func ReportError(err error) tea.Cmd {
	slog.Error("Error reported", "error", err)
	return CmdHandler(InfoMsg{
		Type: InfoTypeError,
		Msg:  err.Error(),
	})
}

type InfoType int

const (
	InfoTypeInfo InfoType = iota
	InfoTypeSuccess
	InfoTypeWarn
	InfoTypeError
	InfoTypeUpdate
)

func ReportInfo(info string) tea.Cmd {
	return CmdHandler(InfoMsg{
		Type: InfoTypeInfo,
		Msg:  info,
	})
}

func ReportWarn(warn string) tea.Cmd {
	return CmdHandler(InfoMsg{
		Type: InfoTypeWarn,
		Msg:  warn,
	})
}

type (
	InfoMsg struct {
		Type InfoType
		Msg  string
		TTL  time.Duration
	}
	ClearStatusMsg struct{}
)

// ExecShell parses a shell command string and executes it with exec.Command.
// Uses shell.Fields for proper handling of shell syntax like quotes and
// arguments while preserving TTY handling for terminal editors.
func ExecShell(ctx context.Context, cmdStr string, callback tea.ExecCallback) tea.Cmd {
	fields, err := shell.Fields(cmdStr, nil)
	if err != nil {
		return ReportError(err)
	}
	if len(fields) == 0 {
		return ReportError(errors.New("empty command"))
	}

	cmd := exec.CommandContext(ctx, fields[0], fields[1:]...)
	return tea.ExecProcess(cmd, callback)
}

