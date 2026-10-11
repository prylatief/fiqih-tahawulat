#!/usr/bin/env bash
# update.sh — Download PDF dari Google Drive, extract, import ke DB
set -euo pipefail

DRIVE_URL="${1:-}"
PDF_DIR="/home/ubuntu/tmp/pdf-habib"

if [ -z "$DRIVE_URL" ]; then
    echo "Usage: $0 <google-drive-folder-url>"
    echo "   or: $0 (re-extract existing PDFs)"
    exit 1
fi

echo "=== Download from Google Drive ==="
~/.local/bin/gdown --folder "$DRIVE_URL" -O "$PDF_DIR"

echo "=== Extract PDF text ==="
cd "$PDF_DIR"
python3 extract.py

echo "=== Import to PostgreSQL ==="
python3 import_db5.py

echo "=== Done ==="
