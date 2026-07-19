#!/bin/sh
set -e

# Script de pós-instalação para pacotes deb/rpm do Shelloma
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
