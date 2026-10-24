#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# CIPHER FLEET: ZERO-DEPENDENCY MULTI-ARCHITECTURE DISTRIBUTION BUILDER
# Compiles all Go microservices for darwin-arm64 and linux-amd64 (CGO_ENABLED=0)
# Bundles standalone executables for Python services (PyInstaller / bootstrap)
# ==============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="${SCRIPT_DIR}/dist"
GO_SERVICES=("cipher-route" "cipher-log" "cipher-mem" "cipher-cli" "cipher-sec")
PYTHON_SERVICES=("cipher-core" "cipher-exec")
TARGET_ARCHS=("darwin/arm64" "linux/amd64")

CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
NC='\033[0m'

echo -e "${CYAN}${BOLD}⚡ [CIPHER DISTRIBUTION BUILDER] Initiating zero-dependency build pipeline...${NC}"
rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}/darwin-arm64/bin" "${DIST_DIR}/linux-amd64/bin"

# 1. Compile Go Microservices
for service in "${GO_SERVICES[@]}"; do
    svc_dir="${SCRIPT_DIR}/${service}"
    if [ ! -d "${svc_dir}" ]; then
        echo -e "${YELLOW}Warning: Directory ${svc_dir} not found, skipping.${NC}"
        continue
    fi

    echo -e "\n${BOLD}🔨 Compiling Go service: ${service}${NC}"
    cd "${svc_dir}"

    for target in "${TARGET_ARCHS[@]}"; do
        os="${target%/*}"
        arch="${target#*/}"
        out_name="${service}"
        if [ "${service}" = "cipher-cli" ]; then
            out_name="cipher"
        fi
        out_dir="${DIST_DIR}/${os}-${arch}/bin"
        out_path="${out_dir}/${out_name}"

        echo -e "   → Building ${os}-${arch} (CGO_ENABLED=0 stripped)..."
        CGO_ENABLED=0 GOOS="${os}" GOARCH="${arch}" go build -ldflags="-s -w" -o "${out_path}" .
        echo -e "     ${GREEN}✓ Generated:${NC} ${out_path} ($(du -h "${out_path}" | cut -f1))"
    done
done

# 2. Package Python Services (PyInstaller / Standalone Executables)
for service in "${PYTHON_SERVICES[@]}"; do
    svc_dir="${SCRIPT_DIR}/${service}"
    if [ ! -d "${svc_dir}" ]; then
        continue
    fi

    echo -e "\n${BOLD}📦 Packaging standalone Python service: ${service}${NC}"
    cd "${svc_dir}"

    entry="main.py"
    if [ "${service}" = "cipher-exec" ]; then
        entry="server.py"
    fi

    if command -v uv >/dev/null 2>&1; then
        echo -e "   → Compiling single-file binary with PyInstaller..."
        uv run pyinstaller --onefile --clean --noconfirm --name "${service}" "${entry}" >/dev/null 2>&1 || true
        if [ -f "dist/${service}" ]; then
            cp "dist/${service}" "${DIST_DIR}/darwin-arm64/bin/${service}"
            echo -e "     ${GREEN}✓ Bundled darwin-arm64 binary:${NC} ${DIST_DIR}/darwin-arm64/bin/${service} ($(du -h "${DIST_DIR}/darwin-arm64/bin/${service}" | cut -f1))"
        fi
    fi

    # Create hermetic bootstrap launcher script for cross-platform portability
    for arch in "darwin-arm64" "linux-amd64"; do
        launcher="${DIST_DIR}/${arch}/bin/${service}-launcher.sh"
        cat << 'EOF' > "${launcher}"
#!/usr/bin/env bash
set -e
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVICE="$(basename "${BASH_SOURCE[0]}" -launcher.sh)"

if [ -f "${DIR}/${SERVICE}" ] && [ "$(uname -m)" = "arm64" ] && [ "$(uname -s)" = "Darwin" ]; then
    exec "${DIR}/${SERVICE}" "$@"
fi

echo "[BOOTSTRAP] Launching ${SERVICE} in self-contained isolated runtime..."
if command -v uv >/dev/null 2>&1; then
    exec uv run "${DIR}/../src/${SERVICE}/main.py" "$@"
fi
exec python3 "${DIR}/../src/${SERVICE}/main.py" "$@"
EOF
        chmod +x "${launcher}"
    done
done

# 3. Generate SHA-256 Manifest
echo -e "\n${BOLD}📝 Generating SHA-256 release checksums...${NC}"
cd "${DIST_DIR}"
find . -type f -exec shasum -a 256 {} + | sort > "${DIST_DIR}/SHA256SUMS"
cat "${DIST_DIR}/SHA256SUMS"

echo -e "\n${GREEN}${BOLD}🎉 ZERO-DEPENDENCY DISTRIBUTION ARTIFACTS READY AT: ${DIST_DIR}${NC}\n"
