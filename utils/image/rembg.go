package image

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"syscall"
)

// RemBG runs rembg, canceling ctx kills it (the pip launcher takes its Python process down with it)
func RemBG(ctx context.Context, model string, path string, img []byte) ([]byte, error) {
	if path == "" && img == nil {
		return nil, errors.New("NO_IMAGE")
	}

	input := path
	if path == "" {
		input = "-"
	}
	cmd := exec.CommandContext(ctx, "rembg", "i", "-m", model, input, "-")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x08000000,
		HideWindow:    true,
	}
	if path == "" {
		cmd.Stdin = bytes.NewReader(img)
	}
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, errors.New("REMBG_NOT_FOUND")
		}
		if msg := errorLine(stderr.String()); msg != "" {
			return nil, errors.New(msg)
		}
		return nil, err
	}
	return out.Bytes(), nil
}

// errorLine picks the useful line from rembg's stderr: the exception for Python tracebacks,
// otherwise the last line mentioning an error (CLI usage text ends with hints, not the error)
func errorLine(stderr string) string {
	lines := strings.Split(strings.TrimSpace(stderr), "\n")
	if strings.Contains(stderr, "Traceback (most recent call last):") {
		return strings.TrimSpace(lines[len(lines)-1])
	}
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(strings.ToLower(lines[i]), "error") {
			return strings.TrimSpace(lines[i])
		}
	}
	return strings.TrimSpace(lines[len(lines)-1])
}
