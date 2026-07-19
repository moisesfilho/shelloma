package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"shelloma/pkg/i18n"
)

func CopyToClipboard(text string, t i18n.Translations) error {
	osc52 := fmt.Sprintf("\033]52;c;%s\a", encodeBase64(text))
	fmt.Print(osc52)

	if runtime.GOOS == "darwin" {
		if isCommandAvailable("pbcopy") {
			cmd := exec.Command("pbcopy")
			cmd.Stdin = strings.NewReader(text)
			return cmd.Run()
		}
		return fmt.Errorf("pbcopy command not found")
	}

	if runtime.GOOS == "windows" {
		if isCommandAvailable("clip.exe") {
			cmd := exec.Command("clip.exe")
			cmd.Stdin = strings.NewReader(text)
			return cmd.Run()
		}
		if isCommandAvailable("clip") {
			cmd := exec.Command("clip")
			cmd.Stdin = strings.NewReader(text)
			return cmd.Run()
		}
		if isCommandAvailable("powershell.exe") {
			cmd := exec.Command("powershell.exe", "-NoProfile", "-Command", "Set-Clipboard -Value "+quoteCmdArg(text))
			return cmd.Run()
		}
		return fmt.Errorf("clip command not found")
	}

	// Linux / POSIX
	var lastErr error

	if isCommandAvailable("wl-copy") {
		cmd := exec.Command("wl-copy")
		cmd.Stdin = strings.NewReader(text)
		err := cmd.Run()
		if err == nil {
			return nil
		}
		lastErr = err
	}

	if isCommandAvailable("xclip") {
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		err := cmd.Run()
		if err == nil {
			return nil
		}
		lastErr = err
	}

	if isCommandAvailable("xsel") {
		cmd := exec.Command("xsel", "--clipboard", "--input")
		cmd.Stdin = strings.NewReader(text)
		err := cmd.Run()
		if err == nil {
			return nil
		}
		lastErr = err
	}

	if isCommandAvailable("clip.exe") {
		cmd := exec.Command("clip.exe")
		cmd.Stdin = strings.NewReader(text)
		err := cmd.Run()
		if err == nil {
			return nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return lastErr
	}

	if t.NoClipboardUtility != "" {
		return fmt.Errorf("%s", t.NoClipboardUtility)
	}
	return fmt.Errorf("no clipboard utility found (install 'xclip' or 'wl-clipboard')")
}

func quoteCmdArg(arg string) string {
	return "'" + strings.ReplaceAll(arg, "'", "''") + "'"
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func encodeBase64(s string) string {
	const b64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result []byte
	n := len(s)
	for i := 0; i < n; i += 3 {
		var b uint32
		b = uint32(s[i]) << 16
		if i+1 < n {
			b |= uint32(s[i+1]) << 8
		}
		if i+2 < n {
			b |= uint32(s[i+2])
		}
		result = append(result, b64[(b>>18)&63])
		result = append(result, b64[(b>>12)&63])
		if i+1 < n {
			result = append(result, b64[(b>>6)&63])
		} else {
			result = append(result, '=')
		}
		if i+2 < n {
			result = append(result, b64[b&63])
		} else {
			result = append(result, '=')
		}
	}
	return string(result)
}
