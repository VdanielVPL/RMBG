package image

import (
	"bytes"
	"errors"
	"os/exec"
	"syscall"
)

func RemBG(model string, path string, img []byte) ([]byte, error) {
	if path != "" {
		cmd := exec.Command("rembg", "i", "-m", model, path, "-")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: 0x08000000,
			HideWindow:    true,
		}
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			cmd = nil
			if errors.Is(err, exec.ErrNotFound) {
				return nil, errors.New("REMBG_NOT_FOUND")
			}
			return nil, err
		}
		cmd = nil
		return out.Bytes(), nil
	} else if img != nil {
		cmd := exec.Command("rembg", "i", "-m", model, "-", "-")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: 0x08000000,
			HideWindow:    true,
		}
		cmd.Stdin = bytes.NewReader(img)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			cmd = nil
			if errors.Is(err, exec.ErrNotFound) {
				return nil, errors.New("REMBG_NOT_FOUND")
			}
			return nil, err
		}
		cmd = nil
		return out.Bytes(), nil
	}
	return nil, errors.New("NO_IMAGE")
}
