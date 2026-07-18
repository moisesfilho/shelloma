.PHONY: all build deb install-user test clean help

BINARY_NAME=shelloma
VERSION=1.0.0

all: build

build:
	@echo "🔨 Compilando $(BINARY_NAME)..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) .
	@echo "✔ Binário gerado em ./$(BINARY_NAME)"

test:
	@echo "🧪 Executando suíte de testes..."
	go test -v ./...

deb:
	@./scripts/build-deb.sh

install-user: build
	@mkdir -p $(HOME)/.local/bin
	@cp $(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME)
	@chmod +x $(HOME)/.local/bin/$(BINARY_NAME)
	@echo "✔ $(BINARY_NAME) instalado com sucesso em $(HOME)/.local/bin/$(BINARY_NAME)"

clean:
	@rm -rf $(BINARY_NAME) build *.deb
	@echo "🧹 Limpeza concluída."

help:
	@echo "Comandos disponíveis:"
	@echo "  make build        - Compila o binário nativo"
	@echo "  make test         - Executa os testes automatizados do projeto"
	@echo "  make deb          - Gera o pacote Debian (.deb) para instalação no Ubuntu/Debian"
	@echo "  make install-user - Instala o executável em ~/.local/bin/shelloma (sem sudo)"
	@echo "  make clean        - Remove os artefatos compilados"
