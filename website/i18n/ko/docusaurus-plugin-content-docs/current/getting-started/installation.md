---
sidebar_position: 1
title: 설치
---

# 설치

curompt는 Homebrew 또는 소스 빌드를 통해 설치할 수 있습니다.

## Homebrew

```bash
brew tap curogom/curompt
brew install curompt
curompt --version
```

업데이트:

```bash
brew update
brew upgrade curompt
```

## 소스에서 빌드

```bash
git clone https://github.com/curogom/curompt.git
cd curompt
make build
./bin/curompt --help
```

또는 `go install ./cmd/curompt`를 실행해 `$GOBIN`에 설치할 수 있습니다.

## 릴리스 바이너리

릴리스 페이지에서 OS에 맞는 파일을 내려받아 실행 권한을 부여한 뒤 `$PATH`에 배치합니다.

```bash
curl -L -o curompt https://github.com/curogom/curompt/releases/download/v0.2.0/curompt-darwin-amd64
chmod +x curompt
mv curompt /usr/local/bin/
```

## 설치 확인

```bash
curompt --version
curompt scan --help
```
