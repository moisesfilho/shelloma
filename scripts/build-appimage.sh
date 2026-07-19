#!/usr/bin/env bash
set -e

export APPIMAGE_EXTRACT_AND_RUN=1

ARCH=${ARCH:-x86_64}
APP_DIR="build/AppDir"
OUTPUT_DIR="dist"

echo "🔨 Construindo AppImage para Shelloma (${ARCH})..."

mkdir -p "${APP_DIR}/usr/bin"
mkdir -p "${OUTPUT_DIR}"

if [ ! -f shelloma ]; then
    echo "🔨 Compilando binário do Shelloma..."
    go build -ldflags="-s -w" -o shelloma .
fi

cp shelloma "${APP_DIR}/usr/bin/shelloma"
cp scripts/shelloma.desktop "${APP_DIR}/shelloma.desktop"

# Criar diretórios de metainfo e ícones
mkdir -p "${APP_DIR}/usr/share/metainfo"
mkdir -p "${APP_DIR}/usr/share/icons/hicolor/64x64/apps"

cp scripts/org.shelloma.Shelloma.appdata.xml "${APP_DIR}/usr/share/metainfo/org.shelloma.Shelloma.appdata.xml"

# Criar um ícone genérico SVG/PNG caso não exista
echo '<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 64 64"><rect width="64" height="64" rx="12" fill="#2d3748"/><text x="12" y="42" font-family="monospace" font-size="32" fill="#48bb78" font-weight="bold">&gt;_</text></svg>' > "${APP_DIR}/org.shelloma.Shelloma.svg"
cp "${APP_DIR}/org.shelloma.Shelloma.svg" "${APP_DIR}/org.shelloma.Shelloma.png"
cp "${APP_DIR}/org.shelloma.Shelloma.png" "${APP_DIR}/usr/share/icons/hicolor/64x64/apps/org.shelloma.Shelloma.png"
cp "${APP_DIR}/org.shelloma.Shelloma.png" "${APP_DIR}/shelloma.png"

# Criar o script AppRun
cat <<'EOF' > "${APP_DIR}/AppRun"
#!/bin/sh
HERE="$(dirname "$(readlink -f "$0")")"
exec "$HERE/usr/bin/shelloma" "$@"
EOF
chmod +x "${APP_DIR}/AppRun"

# Baixar o appimagetool se não existir localmente nem no PATH
if ! command -v appimagetool >/dev/null 2>&1; then
    if [ ! -f ./appimagetool ]; then
        echo "⬇️ Baixando appimagetool..."
        curl -fsSL -o appimagetool "https://github.com/AppImage/appimagetool/releases/download/continuous/appimagetool-x86_64.AppImage"
        chmod +x ./appimagetool
    fi
    APPIMAGETOOL="./appimagetool"
else
    APPIMAGETOOL="appimagetool"
fi

echo "📦 Gerando arquivo AppImage..."
ARCH=${ARCH} $APPIMAGETOOL "${APP_DIR}" "${OUTPUT_DIR}/Shelloma-${ARCH}.AppImage"

echo "✅ AppImage gerado com sucesso em: ${OUTPUT_DIR}/Shelloma-${ARCH}.AppImage"
