package image

import (
	"bytes"
	"os/exec"
	"syscall"
)

func RemBG(model string, path string, img []byte) ([]byte, error) {
	if path != "" {
		cmd := exec.Command("rembg", "i", "-m", model, path, "-")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		cmd = nil
		return out.Bytes(), nil
	}
	return nil, nil
}
