#!/usr/bin/env bash
set -e

VERSION="1.0.2"
ARCH="amd64"
PKG_NAME="shelloma_${VERSION}_${ARCH}"
BUILD_DIR="build/deb/${PKG_NAME}"

echo "🔨 Construindo binário otimizado para produção..."
go build -ldflags="-s -w" -o shelloma .

echo "📦 Preparando estrutura do pacote Debian..."
rm -rf build/deb
mkdir -p "${BUILD_DIR}/DEBIAN"
mkdir -p "${BUILD_DIR}/usr/bin"
mkdir -p "${BUILD_DIR}/etc/shelloma"
mkdir -p "${BUILD_DIR}/usr/share/doc/shelloma"
mkdir -p "${BUILD_DIR}/usr/share/bash-completion/completions"

# Copiar executável
cp shelloma "${BUILD_DIR}/usr/bin/shelloma"
chmod 755 "${BUILD_DIR}/usr/bin/shelloma"

# Gerar arquivo DEBIAN/control
cat <<EOF > "${BUILD_DIR}/DEBIAN/control"
Package: shelloma
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: Moises <moises@localhost>
Depends: xclip | wl-clipboard | xsel
Description: Native natural language to Shell command translator powered by Ollama.
 Native Linux CLI application written in Go to convert natural language instructions
 into executable Bash/Zsh terminal commands for Debian and Ubuntu.
EOF

# Script DEBIAN/postinst para questionar o idioma no momento da instalação
cat <<'EOF' > "${BUILD_DIR}/DEBIAN/postinst"
#!/bin/sh
set -e

if [ "$1" = "configure" ]; then
    echo "======================================================="
    echo "  Shelloma Language Setup / Configuração de Idioma"
    echo "======================================================="
    
    LANG_CHOICE="en"

    if [ -t 0 ]; then
        echo "Select application language / Selecione o idioma padrão:"
        echo "  1) English (en) [Default]"
        echo "  2) Português Brasileiro (pt)"
        echo "  3) Español (es)"
        printf "Choice / Escolha [1-3] (default 1): "
        read input_choice || true

        case "$input_choice" in
            2|pt|PT)
                LANG_CHOICE="pt"
                ;;
            3|es|ES)
                LANG_CHOICE="es"
                ;;
            *)
                LANG_CHOICE="en"
                ;;
        esac
    fi

    echo "Language set to / Idioma selecionado: $LANG_CHOICE"

    # Salvar a configuração global do sistema em /etc/shelloma/config.json
    mkdir -p /etc/shelloma
    cat <<JSON > /etc/shelloma/config.json
{
  "ollama_url": "http://localhost:11434",
  "model": "",
  "temperature": 0.1,
  "auto_execute": false,
  "language": "$LANG_CHOICE"
}
JSON

    # Sincronizar o idioma escolhido na instalação com as configurações existentes dos usuários em /home e /root
    for user_home in /home/* /root; do
        if [ -f "$user_home/.config/shelloma/config.json" ]; then
            sed -i 's/"language": *"[^"]*"/"language": "'"$LANG_CHOICE"'"/' "$user_home/.config/shelloma/config.json" || true
        fi
    done

    # Instalação automática do utilitário de área de transferência (clipboard) se necessário
    if ! command -v xclip >/dev/null 2>&1 && ! command -v wl-copy >/dev/null 2>&1 && ! command -v xsel >/dev/null 2>&1; then
        echo "📦 Instalando utilitário de área de transferência (xclip)..."
        if command -v apt-get >/dev/null 2>&1; then
            apt-get update -qq && apt-get install -y -qq xclip || true
        elif command -v dnf >/dev/null 2>&1; then
            dnf install -y xclip || true
        elif command -v pacman >/dev/null 2>&1; then
            pacman -S --noconfirm xclip || true
        fi
    fi
fi
EOF

chmod 755 "${BUILD_DIR}/DEBIAN/postinst"

# Script de autocomplete para Bash
cat <<'EOF' > "${BUILD_DIR}/usr/share/bash-completion/completions/shelloma"
_shelloma_completions() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    if [ "$COMP_CWORD" -eq 1 ]; then
        COMPREPLY=( $(compgen -W "models config --help --version -m --model -l --lang -y --yes" -- "$cur") )
    fi
}
complete -F _shelloma_completions shelloma
EOF

chmod 644 "${BUILD_DIR}/DEBIAN/control"
chmod 644 "${BUILD_DIR}/usr/share/bash-completion/completions/shelloma"

echo "⚙️ Criando arquivo .deb com dpkg-deb..."
dpkg-deb --root-owner-group --build "${BUILD_DIR}" "${PKG_NAME}.deb"

echo "✅ Pacote gerado com sucesso: ${PKG_NAME}.deb"
ls -lh "${PKG_NAME}.deb"
