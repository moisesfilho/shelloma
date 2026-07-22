.PHONY: all build build-all build-linux build-darwin build-windows deb tar appimage flatpak install-user test clean help release-snapshot

BINARY_NAME=shelloma
VERSION=1.2.0

all: build

build: lint test
	@echo "🔨 Compilando $(BINARY_NAME) (nativo local)..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) .
	@echo "✔ Binário gerado em ./$(BINARY_NAME)"

build-all: lint test build-linux build-darwin build-windows

build-linux:
	@echo "🔨 Compilando binários nativos Linux (amd64 e arm64)..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/$(BINARY_NAME)_linux_amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/$(BINARY_NAME)_linux_arm64 .
	@echo "✔ Binários Linux gerados em ./dist/"

build-darwin:
	@echo "🔨 Compilando binários nativos macOS (amd64 e arm64)..."
	@mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/$(BINARY_NAME)_darwin_amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/$(BINARY_NAME)_darwin_arm64 .
	@echo "✔ Binários macOS gerados em ./dist/"

build-windows:
	@echo "🔨 Compilando binários nativos Windows (amd64 e arm64)..."
	@mkdir -p dist
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/$(BINARY_NAME)_windows_amd64.exe .
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o dist/$(BINARY_NAME)_windows_arm64.exe .
	@echo "✔ Binários Windows gerados em ./dist/"

test:
	@echo "🧪 Executando suíte de testes..."
	go test -v ./...

lint:
	@echo "🔍 Executando análise estática de código..."
	@go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	elif [ -f $(shell go env GOPATH)/bin/golangci-lint ]; then \
		$(shell go env GOPATH)/bin/golangci-lint run; \
	elif command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	elif [ -f $(shell go env GOPATH)/bin/staticcheck ]; then \
		$(shell go env GOPATH)/bin/staticcheck ./...; \
	else \
		echo "💡 Dica: instale o 'golangci-lint' ou 'staticcheck' para análises estáticas ainda mais profundas."; \
	fi

deb:
	@./scripts/build-deb.sh

tar: build-all
	@echo "📦 Gerando pacotes comprimidos (.tar.gz e .zip)..."
	@mkdir -p dist
	tar -czvf dist/$(BINARY_NAME)_$(VERSION)_linux_amd64.tar.gz -C dist $(BINARY_NAME)_linux_amd64
	tar -czvf dist/$(BINARY_NAME)_$(VERSION)_linux_arm64.tar.gz -C dist $(BINARY_NAME)_linux_arm64
	tar -czvf dist/$(BINARY_NAME)_$(VERSION)_darwin_amd64.tar.gz -C dist $(BINARY_NAME)_darwin_amd64
	tar -czvf dist/$(BINARY_NAME)_$(VERSION)_darwin_arm64.tar.gz -C dist $(BINARY_NAME)_darwin_arm64
	@echo "✔ Pacotes gerados em ./dist/"

appimage: build
	@./scripts/build-appimage.sh

flatpak: build
	@./scripts/build-flatpak.sh

release-snapshot:
	@echo "🚀 Testando build de release com GoReleaser (modo snapshot local)..."
	goreleaser release --snapshot --clean

install-user: build
	@mkdir -p $(HOME)/.local/bin
	@cp $(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME)
	@chmod +x $(HOME)/.local/bin/$(BINARY_NAME)
	@echo "✔ $(BINARY_NAME) instalado com sucesso em $(HOME)/.local/bin/$(BINARY_NAME)"

clean:
	@rm -rf $(BINARY_NAME) build dist *.deb *.rpm *.AppImage *.flatpak appimagetool
	@echo "🧹 Limpeza concluída."

help:
	@echo "Comandos disponíveis:"
	@echo "  make build            - Compila o binário nativo para a arquitetura local"
	@echo "  make build-all        - Compila binários para Linux, macOS e Windows (amd64/arm64)"
	@echo "  make build-linux      - Compila binários para Linux (amd64/arm64)"
	@echo "  make build-darwin     - Compila binários para macOS (amd64/arm64)"
	@echo "  make build-windows    - Compila binários para Windows (amd64/arm64)"
	@echo "  make test             - Executa os testes automatizados do projeto"
	@echo "  make deb              - Gera o pacote Debian (.deb)"
	@echo "  make tar              - Gera os pacotes comprimidos em ./dist/"
	@echo "  make appimage         - Gera o pacote AppImage (.AppImage)"
	@echo "  make release-snapshot - Simula a release completa localmente usando GoReleaser"
	@echo "  make install-user     - Instala o executável em ~/.local/bin/shelloma"
	@echo "  make clean            - Remove os artefatos compilados"
