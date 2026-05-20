#!/bin/bash
#
# Full Recon Toolkit Installer
# Installs dependencies for:
# subenum + alive pipeline
#

set -e

GREEN="\e[32m"
RED="\e[31m"
BLUE="\e[34m"
END="\e[0m"

log() {
    echo -e "${BLUE}[+]${END} $1"
}

warn() {
    echo -e "${RED}[!]${END} $1"
}

GOlang() {
    printf "                                \r"
    sudo apt update
    sudo apt install golang-go -y

    export GOPATH=$HOME/go
    export PATH=$PATH:$HOME/go/bin

    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc

    printf "[+] Golang Installed !.\n"
}

Findomain() {
    log "Installing Findomain..."
    sudo apt install findomain -y
    log "Findomain installed."
}

Subfinder() {
    log "Installing Subfinder..."
    go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
}

Amass() {
    log "Installing Amass..."
    go install github.com/owasp-amass/amass/v4/...@latest
}

Assetfinder() {
    log "Installing Assetfinder..."
    go install github.com/tomnomnom/assetfinder@latest
}

Dnsx() {
    log "Installing Dnsx..."
    go install github.com/projectdiscovery/dnsx/cmd/dnsx@latest
}

Httpx() {
    log "Installing Httpx..."
    go install github.com/projectdiscovery/httpx/cmd/httpx@latest
}

Naabu() {
    log "Installing Naabu..."
    go install github.com/projectdiscovery/naabu/v2/cmd/naabu@latest
}

Anew() {
    log "Installing Anew..."
    go install github.com/tomnomnom/anew@latest
}

Parallel() {
    log "Installing GNU Parallel..."
    sudo apt install -y parallel
}

SystemDeps() {
    log "Installing system dependencies..."
    sudo apt install -y wget unzip curl jq
}

# Install Go if missing
if ! command -v go >/dev/null 2>&1; then
    GOlang
else
    warn "Golang already installed."
    export PATH=$PATH:$HOME/go/bin
fi

SystemDeps

command -v findomain >/dev/null 2>&1 || Findomain
command -v subfinder >/dev/null 2>&1 || Subfinder
command -v amass >/dev/null 2>&1 || Amass
command -v assetfinder >/dev/null 2>&1 || Assetfinder
command -v dnsx >/dev/null 2>&1 || Dnsx
command -v httpx >/dev/null 2>&1 || Httpx
command -v naabu >/dev/null 2>&1 || Naabu
command -v anew >/dev/null 2>&1 || Anew
command -v parallel >/dev/null 2>&1 || Parallel

TOOLS=(
    go
    findomain
    subfinder
    amass
    assetfinder
    dnsx
    httpx
    naabu
    parallel
    anew
)

echo ""
echo "Verification:"
echo "-------------"

for tool in "${TOOLS[@]}"; do
    if command -v "$tool" >/dev/null 2>&1; then
        echo -e "[$tool] ${GREEN}Installed${END}"
    else
        echo -e "[$tool] ${RED}FAILED${END}"
    fi
done
