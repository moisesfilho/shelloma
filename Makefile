.PHONY: all build build-all deb tar appimage flatpak install-user test clean help release-snapshot

BINARY_NAME=shelloma
VERSION=1.0.2

all: build

build:
	@echo "🔨 Compilando $(BINARY_NAME) (nativo local)..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) .
	@echo "✔ Binário gerado em ./$(BINARY_NAME)"

build-all:
	@echo "🔨 Compilando binários nativos para amd64 e arm64..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/$(BINARY_NAME)_linux_amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/$(BINARY_NAME)_linux_arm64 .
	@echo "✔ Binários gerados em ./dist/"

test:
	@echo "🧪 Executando suíte de testes..."
	go test -v ./...

deb:
	@./scripts/build-deb.sh

tar: build-all
	@echo "📦 Gerando tarballs (.tar.gz) para amd64 e arm64..."
	@mkdir -p dist
	tar -czvf dist/$(BINARY_NAME)_$(VERSION)_linux_amd64.tar.gz $(BINARY_NAME) README.md README_pt.md LICENSE scripts/postinstall.sh
	tar -czvf dist/$(BINARY_NAME)_$(VERSION)_linux_arm64.tar.gz dist/$(BINARY_NAME)_linux_arm64 README.md README_pt.md LICENSE scripts/postinstall.sh
	@echo "✔ Tarballs gerados em ./dist/"

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
	@echo "  make build-all        - Compila binários nativos para linux/amd64 e linux/arm64 em ./dist/"
	@echo "  make test             - Executa os testes automatizados do projeto"
	@echo "  make deb              - Gera o pacote Debian (.deb)"
	@echo "  make tar              - Gera os pacotes comprimidos (.tar.gz) para amd64 e arm64 em ./dist/"
	@echo "  make appimage         - Gera o pacote AppImage (.AppImage)"
	@echo "  make release-snapshot - Simula a release completa localmente usando GoReleaser"
	@echo "  make install-user     - Instala o executável em ~/.local/bin/shelloma"
	@echo "  make clean            - Remove os artefatos compilados"
