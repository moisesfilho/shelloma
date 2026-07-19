#!/usr/bin/env bash
set -e

ARCH=${ARCH:-x86_64}
OUTPUT_DIR="dist"
BUILD_DIR="build/flatpak-build"
REPO_DIR="build/flatpak-repo"

echo "🔨 Construindo Flatpak para Shelloma (${ARCH})..."

mkdir -p "${OUTPUT_DIR}"
rm -rf "${BUILD_DIR}" "${REPO_DIR}"

if [ ! -f shelloma ]; then
    echo "🔨 Compilando binário do Shelloma..."
    go build -ldflags="-s -w" -o shelloma .
fi

if ! command -v flatpak-builder >/dev/null 2>&1; then
    echo "⚠️ flatpak-builder não instalado. Pulando geração do Flatpak..."
    exit 0
fi

# Adicionar repositório flathub se necessário
flatpak remote-add --user --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo || true

echo "📦 Compilando repositório Flatpak com flatpak-builder..."
flatpak-builder --force-clean --disable-rofiles-fuse --repo="${REPO_DIR}" --arch="${ARCH}" "${BUILD_DIR}" scripts/org.shelloma.Shelloma.yml

echo "📦 Gerando pacote .flatpak..."
flatpak build-bundle "${REPO_DIR}" "${OUTPUT_DIR}/Shelloma-${ARCH}.flatpak" org.shelloma.Shelloma

echo "✅ Flatpak gerado com sucesso em: ${OUTPUT_DIR}/Shelloma-${ARCH}.flatpak"
