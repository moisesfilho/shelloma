package sysinfo

import (
	"runtime"
	"testing"
)

func TestGetSystemContext(t *testing.T) {
	ctx := GetSystemContext()

	if ctx.OS != runtime.GOOS {
		t.Errorf("Esperava OS %s, obteve %s", runtime.GOOS, ctx.OS)
	}

	if ctx.Arch != runtime.GOARCH {
		t.Errorf("Esperava Arch %s, obteve %s", runtime.GOARCH, ctx.Arch)
	}

	if ctx.WorkingDir == "" {
		t.Error("WorkingDir não deveria estar vazio")
	}

	if ctx.Shell == "" {
		t.Error("Shell não deveria estar vazio")
	}

	if ctx.User == "" {
		t.Error("User não deveria estar vazio")
	}

	if ctx.DistroName == "" {
		t.Error("DistroName não deveria estar vazio")
	}
}
