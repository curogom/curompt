# Homebrew 배포 가이드

## 개요

`curompt`를 Homebrew를 통해 배포하는 방법입니다.

## 배포 방법 선택

### 방법 1: Homebrew Core (공식 레포) - 권장 ❌
- **장점**: `brew install curompt`로 바로 설치 가능
- **단점**: 
  - 리뷰 프로세스가 까다로움
  - 유명도 필요 (GitHub Stars 100+ 등)
  - 승인까지 시간 소요

### 방법 2: Homebrew Tap (개인 레포) - 추천 ✅
- **장점**: 
  - 빠른 배포
  - 완전한 제어권
  - 커스터마이징 가능
- **단점**: 사용자는 `brew install {username}/curompt/curompt` 형식으로 설치

## 구현 단계 (Homebrew Tap 방식)

### 1단계: GitHub Release 준비

```bash
# 1. 버전 태그 생성
git tag -a v0.1.0 -m "Initial release"
git push origin v0.1.0

# 2. GitHub에서 Release 생성
# - Assets에 소스 tarball 첨부
# - 또는 자동 생성된 소스 코드 다운로드 사용
```

### 2단계: Formula 파일 작성

Formula 파일: `Formula/c/curompt.rb`

```ruby
class CuroPrompt < Formula
  desc "CLI tool for analyzing, evaluating, and optimizing LLM prompts"
  homepage "https://github.com/curogom/curompt"
  url "https://github.com/curogom/curompt/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "CHECKSUM_HERE"  # shasum -a 256 명령으로 계산
  license "Apache-2.0"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"curompt", "./cmd/curompt"
  end

  test do
    system "#{bin}/curompt", "--version"
  end
end
```

### 3단계: Homebrew Tap 레포지토리 생성

```bash
# 새로운 레포지토리 생성 (예: homebrew-curompt)
# Formula/c/curompt.rb 파일 추가
git add Formula/c/curompt.rb
git commit -m "Add curompt formula"
git push
```

### 4단계: SHA256 계산

```bash
# Release tarball 다운로드 후
shasum -a 256 curompt-0.1.0.tar.gz
# 결과를 Formula의 sha256에 입력
```

## 자동화 스크립트

### GitHub Actions로 자동 배포

`.github/workflows/release.yml` 예시:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Build
        run: |
          make build
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            bin/curompt
          generate_release_notes: true
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 자동 Formula 업데이트

`scripts/update-homebrew.sh`:

```bash
#!/bin/bash
VERSION=$1
SHA256=$(curl -sL "https://github.com/curogom/curompt/archive/refs/tags/v${VERSION}.tar.gz" | shasum -a 256 | awk '{print $1}')

# Formula 업데이트
sed -i '' "s/url \".*\"/url \"https:\/\/github.com\/curogom\/curompt\/archive\/refs\/tags\/v${VERSION}.tar.gz\"/" Formula/c/curompt.rb
sed -i '' "s/sha256 \".*\"/sha256 \"${SHA256}\"/" Formula/c/curompt.rb
```

## 사용자 설치 방법

### Tap 레포지토리 방식

```bash
# Tap 추가
brew tap curogom/curompt

# 설치
brew install curompt

# 업데이트
brew upgrade curompt
```

## 업데이트 프로세스

1. 새 버전 태그 생성
2. GitHub Release 생성
3. Formula의 url과 sha256 업데이트
4. Tap 레포지토리에 PR 또는 직접 push

## 체크리스트

- [ ] GitHub Release 생성 (tarball 포함)
- [ ] SHA256 해시 계산
- [ ] Formula 파일 작성 및 테스트
- [ ] Tap 레포지토리 생성 또는 기존 레포 업데이트
- [ ] 설치 테스트 (`brew install --build-from-source`)
- [ ] 문서 업데이트 (README에 설치 방법 추가)

