#!/bin/bash

# Alive Recon Pipeline
# Usage: alive -f subdomains.txt

THREADS_DNS=800
THREADS_HTTP=400
THREADS_NAABU=400

show_help() {
    echo "Usage: alive -f <subdomains_file>"
    echo ""
    echo "Example:"
    echo "  alive -f all_subdomains.txt"
    exit 1
}

if [[ $# -eq 0 ]]; then
    show_help
fi

while getopts "f:h" opt; do
    case $opt in
        f) INPUT="$OPTARG" ;;
        h) show_help ;;
        *) show_help ;;
    esac
done

if [[ -z "$INPUT" ]]; then
    echo "[!] Input file required"
    show_help
fi

if [[ ! -f "$INPUT" ]]; then
    echo "[!] File not found: $INPUT"
    exit 1
fi

for tool in dnsx naabu httpx; do
    if ! command -v $tool >/dev/null 2>&1; then
        echo "[!] Missing dependency: $tool"
        exit 1
    fi
done

echo "[+] Starting DNS resolution..."
dnsx -l "$INPUT" -silent -resp -threads $THREADS_DNS > resolved.txt

DNS_COUNT=$(wc -l < resolved.txt)
echo "[+] DNS resolved hosts: $DNS_COUNT"

echo "[+] Starting port scan..."
naabu -l resolved.txt -top-ports 1000 -silent -c $THREADS_NAABU > ports.txt

PORT_COUNT=$(wc -l < ports.txt)
echo "[+] Hosts with open ports: $PORT_COUNT"

echo "[+] Starting HTTP probing..."
httpx -l resolved.txt \
-silent \
-title \
-tech-detect \
-status-code \
-follow-redirects \
-threads $THREADS_HTTP \
-timeout 8 \
-retries 1 \
-o alive.txt

HTTP_COUNT=$(wc -l < alive.txt)
echo "[+] Web alive hosts: $HTTP_COUNT"

echo ""
echo "[+] Recon completed."
echo "    resolved.txt"
echo "    ports.txt"
echo "    alive.txt"
