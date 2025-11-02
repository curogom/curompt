#!/bin/bash
# Homebrew Formula 업데이트 스크립트
# Usage: ./scripts/update-homebrew-formula.sh v0.1.0

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v0.1.0"
    exit 1
fi

# Remove 'v' prefix if present
VERSION_NUM=${VERSION#v}

echo "Updating Homebrew formula for version: $VERSION"
echo ""

# Download tarball and calculate SHA256
echo "Calculating SHA256..."
TARBALL_URL="https://github.com/curogom/curo-prompt/archive/refs/tags/${VERSION}.tar.gz"
SHA256=$(curl -sL "$TARBALL_URL" | shasum -a 256 | awk '{print $1}')

if [ -z "$SHA256" ]; then
    echo "Error: Could not calculate SHA256. Make sure the release exists."
    exit 1
fi

echo "SHA256: $SHA256"
echo ""

# Update formula file
FORMULA_FILE="Formula/curo-prompt.rb"
if [ ! -f "$FORMULA_FILE" ]; then
    echo "Error: Formula file not found: $FORMULA_FILE"
    exit 1
fi

# Update version in URL
sed -i '' "s|url \".*\"|url \"${TARBALL_URL}\"|" "$FORMULA_FILE"

# Update SHA256
sed -i '' "s/sha256 \".*\"/sha256 \"${SHA256}\"/" "$FORMULA_FILE"

echo "✅ Formula updated successfully!"
echo ""
echo "Next steps:"
echo "1. Review the updated formula: cat $FORMULA_FILE"
echo "2. Commit and push to your Homebrew tap repository"
echo "3. Users can install with: brew install curogom/curo-prompt/curo-prompt"

