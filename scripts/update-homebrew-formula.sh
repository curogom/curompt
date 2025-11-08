#!/bin/bash
# Homebrew Formula 업데이트 스크립트
# Usage: ./scripts/update-homebrew-formula.sh v0.1.0 [tap-repo-path]

set -e

VERSION=$1
TAP_REPO_PATH=$2

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version> [tap-repo-path]"
    echo "Example: $0 v0.1.0"
    echo "Example: $0 v0.1.0 ../homebrew-curompt"
    exit 1
fi

# Remove 'v' prefix if present
VERSION_NUM=${VERSION#v}

echo "=== Homebrew Formula 업데이트 ==="
echo "Version: $VERSION"
echo ""

# Download tarball and calculate SHA256
echo "📥 Release tarball 다운로드 및 SHA256 계산 중..."
TARBALL_URL="https://github.com/curogom/curompt/archive/refs/tags/${VERSION}.tar.gz"
TEMP_TARBALL=$(mktemp)
curl -sL "$TARBALL_URL" -o "$TEMP_TARBALL" || {
    echo "❌ Error: Release tarball을 다운로드할 수 없습니다."
    echo "   Release가 생성되었는지 확인하세요: $TARBALL_URL"
    exit 1
}

SHA256=$(shasum -a 256 "$TEMP_TARBALL" | awk '{print $1}')
rm -f "$TEMP_TARBALL"

if [ -z "$SHA256" ]; then
    echo "❌ Error: SHA256 계산 실패"
    exit 1
fi

echo "✅ SHA256: $SHA256"
echo ""

# Determine formula file location
if [ -n "$TAP_REPO_PATH" ]; then
    FORMULA_FILE="$TAP_REPO_PATH/Formula/c/curompt.rb"
    if [ ! -f "$FORMULA_FILE" ]; then
        FORMULA_FILE="$TAP_REPO_PATH/Formula/curompt.rb"
    fi
else
    FORMULA_FILE="Formula/curompt.rb"
fi

if [ ! -f "$FORMULA_FILE" ]; then
    echo "❌ Error: Formula 파일을 찾을 수 없습니다: $FORMULA_FILE"
    exit 1
fi

echo "📝 Formula 파일 업데이트 중: $FORMULA_FILE"

# Backup
cp "$FORMULA_FILE" "${FORMULA_FILE}.bak"

# Update version in URL (macOS와 Linux 모두 지원)
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s|url \".*\"|url \"${TARBALL_URL}\"|" "$FORMULA_FILE"
    sed -i '' "s/sha256 \".*\"/sha256 \"${SHA256}\"/" "$FORMULA_FILE"
else
    sed -i "s|url \".*\"|url \"${TARBALL_URL}\"|" "$FORMULA_FILE"
    sed -i "s/sha256 \".*\"/sha256 \"${SHA256}\"/" "$FORMULA_FILE"
fi

echo "✅ Formula 업데이트 완료!"
echo ""
echo "📋 변경 사항:"
echo "   URL: $TARBALL_URL"
echo "   SHA256: $SHA256"
echo ""
echo "📝 다음 단계:"
echo "1. Formula 검토: cat $FORMULA_FILE"
echo "2. Tap 레포지토리에 커밋 및 푸시"
if [ -z "$TAP_REPO_PATH" ]; then
    echo "3. Tap 레포지토리로 복사:"
    echo "   cp $FORMULA_FILE <tap-repo>/Formula/c/curompt.rb"
fi
echo ""
echo "👥 사용자 설치 방법:"
echo "   brew tap curogom/curompt"
echo "   brew install curompt"

