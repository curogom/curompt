#!/bin/bash
# Homebrew Tap 레포지토리 초기 설정 스크립트

set -e

TAP_NAME="homebrew-curompt"
GITHUB_USER="curogom"

echo "=== Homebrew Tap 설정 ==="
echo ""

# 1. Tap 레포지토리 생성 안내
echo "1. GitHub에서 새 레포지토리 생성:"
echo "   - 레포지토리 이름: $TAP_NAME"
echo "   - Public 레포지토리로 생성"
echo "   - README는 선택 사항"
echo ""
read -p "레포지토리를 생성했나요? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "먼저 GitHub에서 레포지토리를 생성하세요."
    exit 1
fi

# 2. 로컬에 Tap 레포지토리 클론 또는 생성
echo ""
echo "2. Tap 레포지토리 초기화..."
if [ -d "$TAP_NAME" ]; then
    echo "   기존 디렉토리 발견: $TAP_NAME"
    cd "$TAP_NAME"
    git pull || true
else
    echo "   레포지토리 클론 중..."
    git clone "https://github.com/$GITHUB_USER/$TAP_NAME.git" || {
        echo "   클론 실패. 새로 생성합니다..."
        mkdir -p "$TAP_NAME"
        cd "$TAP_NAME"
        git init
        git remote add origin "https://github.com/$GITHUB_USER/$TAP_NAME.git"
    }
fi

# 3. Formula 디렉토리 생성
cd "$TAP_NAME"
mkdir -p Formula/c

# 4. Formula 파일 복사
echo ""
echo "3. Formula 파일 준비..."
if [ -f "../Formula/curompt.rb" ]; then
    cp "../Formula/curompt.rb" "Formula/c/curompt.rb"
    echo "   ✅ Formula 파일 복사 완료"
else
    echo "   ⚠️  Formula 파일을 찾을 수 없습니다."
    echo "   ../Formula/curompt.rb 파일을 먼저 생성하세요."
fi

# 5. README 생성
cat > README.md << 'EOF'
# Homebrew Tap for curompt

Install curompt using Homebrew:

```bash
brew tap curogom/curompt
brew install curompt
```

## Update

```bash
brew upgrade curompt
```
EOF

echo "   ✅ README 생성 완료"

# 6. Git 초기 커밋
if [ -z "$(git status --porcelain)" ]; then
    echo ""
    echo "   이미 커밋된 상태입니다."
else
    echo ""
    echo "4. Git 커밋 준비..."
    git add .
    echo ""
    echo "다음 명령어로 커밋하고 푸시하세요:"
    echo "  cd $TAP_NAME"
    echo "  git commit -m 'Add curompt formula'"
    echo "  git push -u origin main"
fi

echo ""
echo "=== 완료 ==="
echo ""
echo "다음 단계:"
echo "1. Formula 파일의 SHA256을 실제 Release 값으로 업데이트"
echo "2. cd $TAP_NAME && git add . && git commit -m 'Add curompt formula'"
echo "3. git push -u origin main"
echo "4. 설치 테스트: brew tap $GITHUB_USER/curompt && brew install curompt"

