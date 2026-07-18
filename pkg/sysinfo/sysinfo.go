package sysinfo

import (
	"bufio"
	"os"
	"os/user"
	"runtime"
	"strings"
)

// SystemContext guarda informações do ambiente para enviar ao prompt da LLM
type SystemContext struct {
	OS          string
	DistroName  string
	DistroVer   string
	Shell       string
	User        string
	WorkingDir  string
	IsRoot      bool
	Arch        string
}

// GetSystemContext lê o ambiente atual do sistema
func GetSystemContext() SystemContext {
	ctx := SystemContext{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	// Diretório atual
	if pwd, err := os.Getwd(); err == nil {
		ctx.WorkingDir = pwd
	}

	// Shell em uso ($SHELL)
	shellEnv := os.Getenv("SHELL")
	if shellEnv != "" {
		parts := strings.Split(shellEnv, "/")
		ctx.Shell = parts[len(parts)-1]
	} else {
		ctx.Shell = "bash"
	}

	// Usuário atual
	if usr, err := user.Current(); err == nil {
		ctx.User = usr.Username
		ctx.IsRoot = (usr.Uid == "0")
	}

	// Ler informações da distribuição Linux (/etc/os-release)
	ctx.DistroName, ctx.DistroVer = getDistroInfo()

	return ctx
}

func getDistroInfo() (string, string) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "Linux", ""
	}
	defer file.Close()

	name := "Linux"
	version := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "NAME=") {
			name = strings.Trim(strings.TrimPrefix(line, "NAME="), `"`)
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), `"`)
		} else if strings.HasPrefix(line, "PRETTY_NAME=") && name == "Linux" {
			name = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
		}
	}

	return name, version
}
