#!/bin/bash

set -e

echo "🧹 Cleaning up unnecessary Docker Compose files for k3d setup"
echo ""

# Create archive directory
mkdir -p archive/docker-compose
echo "✅ Created archive directory: archive/docker-compose"

# List files to archive
FILES_TO_ARCHIVE=(
    "docker-compose.prod.yml"
    "docker-compose.kafka.yml"
    "docker-compose.monitoring.yml"
    "docker-compose.secrets.yml"
)

echo ""
echo "📦 Files to be archived:"
for file in "${FILES_TO_ARCHIVE[@]}"; do
    if [ -f "$file" ]; then
        echo "   - $file"
    fi
done

echo ""
read -p "Continue with archiving? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Cancelled"
    exit 0
fi

# Archive files
for file in "${FILES_TO_ARCHIVE[@]}"; do
    if [ -f "$file" ]; then
        mv "$file" "archive/docker-compose/"
        echo "✅ Archived: $file"
    else
        echo "⏭️  Skipped (not found): $file"
    fi
done

echo ""
echo "✅ Cleanup complete!"
echo ""
echo "📁 Remaining Docker Compose files:"
ls -1 docker-compose*.yml 2>/dev/null || echo "   (none)"
echo ""
echo "📦 Archived files location: archive/docker-compose/"
echo ""
echo "💡 Next steps:"
echo "   1. Use k3d for Kubernetes-based development"
echo "   2. Refer to MACBOOK_SETUP.md for full setup guide"
echo "   3. Use remaining docker-compose.yml only for quick local testing"
