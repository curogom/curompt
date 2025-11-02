# 설치 가이드

## 개요

`curo-prompt`를 시스템에 설치하여 어디서든 `curo-prompt` 명령으로 실행할 수 있도록 설정합니다.

## 설치 방법

### 방법 1: Homebrew (macOS 권장) ⭐

macOS에서 `curo-prompt`를 설치하는 가장 쉬운 방법입니다. Tap 레포지토리가 준비되어 바로 사용 가능합니다!

#### 1단계: Tap 추가

```bash
brew tap curogom/curo-prompt
```

#### 2단계: 설치

```bash
brew install curo-prompt
```

#### 3단계: 설치 확인

```bash
# 버전 확인
curo-prompt --version

# 도움말 확인
curo-prompt --help
```

#### 업데이트

```bash
brew upgrade curo-prompt
```

#### 제거

```bash
brew uninstall curo-prompt
brew untap curogom/curo-prompt  # 선택사항: Tap 제거
```

> **✅ 사용 가능**: Homebrew tap이 준비되었습니다. 위 명령어로 바로 설치할 수 있습니다!

---

### 방법 2: Go install (Linux 권장) ⭐

가장 간단하고 표준적인 방법입니다.

#### 1단계: 설치

```bash
# 방법 A: Makefile 사용
git clone https://github.com/curogom/curo-prompt.git
cd curo-prompt
make install

# 방법 B: go install 직접 사용
go install github.com/curogom/curo-prompt/cmd/curo-prompt@latest
```

설치 위치:
- `$GOPATH/bin/curo-prompt` (기본값: `~/go/bin/curo-prompt`)

#### 2단계: PATH 확인

```bash
# Go bin 경로 확인
go env GOPATH
# 출력 예: /Users/username/go

# PATH에 이미 포함되어 있는지 확인
echo $PATH | grep -q "$(go env GOPATH)/bin" && echo "✅ 이미 PATH에 포함됨" || echo "❌ PATH에 없음"
```

#### 3단계: PATH 추가 (필요한 경우)

**macOS (zsh)**:
```bash
# ~/.zshrc에 추가
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

**Linux (bash)**:
```bash
# ~/.bashrc 또는 ~/.bash_profile에 추가
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

**fish**:
```fish
# ~/.config/fish/config.fish에 추가
set -gx PATH $PATH (go env GOPATH)/bin
```

#### 4단계: 설치 확인

```bash
# 명령어 확인
which curo-prompt
# 출력: /Users/username/go/bin/curo-prompt

# 버전 확인
curo-prompt --version

# 도움말 확인
curo-prompt --help
```

---

### 방법 2: 로컬 빌드 후 수동 설치

개발 중이거나 특정 버전을 사용하고 싶은 경우:

#### 1단계: 빌드

```bash
git clone https://github.com/curogom/curo-prompt.git
cd curo-prompt
make build
```

바이너리 위치: `./bin/curo-prompt`

#### 2단계: 시스템 경로에 복사

**옵션 A: 시스템 전역 설치 (sudo 필요)**

```bash
# macOS/Linux
sudo cp ./bin/curo-prompt /usr/local/bin/
sudo chmod +x /usr/local/bin/curo-prompt

# 확인
/usr/local/bin/curo-prompt --version
```

**옵션 B: 사용자 로컬 설치**

```bash
# ~/bin 디렉토리 생성
mkdir -p ~/bin

# 복사
cp ./bin/curo-prompt ~/bin/

# PATH에 추가 (셸 설정 파일에)
echo 'export PATH="$PATH:$HOME/bin"' >> ~/.zshrc  # 또는 ~/.bashrc
source ~/.zshrc

# 확인
which curo-prompt
```

---

### 방법 3: GitHub Releases (향후)

공식 릴리스가 있을 경우:

```bash
# 다운로드 (예시)
wget https://github.com/curogom/curo-prompt/releases/download/v1.0.0/curo-prompt-linux-amd64

# 실행 권한 부여
chmod +x curo-prompt-linux-amd64

# 시스템 경로로 이동
sudo mv curo-prompt-linux-amd64 /usr/local/bin/curo-prompt

# 확인
curo-prompt --version
```

---

## PATH 문제 해결

### 문제: `command not found: curo-prompt`

#### 해결 1: PATH 확인

```bash
# 현재 PATH 확인
echo $PATH

# Go bin 경로 확인
go env GOPATH
# 출력 예: /Users/username/go

# 수동으로 경로 추가하여 테스트
export PATH="$PATH:$(go env GOPATH)/bin"
curo-prompt --help  # 이제 작동해야 함
```

#### 해결 2: 셸 설정 파일 수정

**zsh 사용자**:
```bash
# ~/.zshrc 확인
cat ~/.zshrc | grep GOPATH

# 없으면 추가
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc
source ~/.zshrc
```

**bash 사용자**:
```bash
# ~/.bashrc 또는 ~/.bash_profile 확인
cat ~/.bashrc | grep GOPATH

# 없으면 추가
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

#### 해결 3: 직접 경로 지정

```bash
# 직접 경로로 실행
$(go env GOPATH)/bin/curo-prompt --help

# 또는 alias 설정
echo 'alias curo-prompt="$(go env GOPATH)/bin/curo-prompt"' >> ~/.zshrc
source ~/.zshrc
```

---

## 설치 위치 확인

```bash
# Go install로 설치한 경우
which curo-prompt
# 출력: /Users/username/go/bin/curo-prompt

# 또는
go env GOPATH
# 출력: /Users/username/go
# 실제 위치: /Users/username/go/bin/curo-prompt
```

---

## 업데이트

### Go install 사용한 경우

```bash
# 최신 버전으로 업데이트
go install github.com/curogom/curo-prompt/cmd/curo-prompt@latest

# 또는 저장소를 업데이트 후
git pull
make install
```

### 로컬 빌드 사용한 경우

```bash
# 저장소 업데이트
git pull

# 다시 빌드 및 설치
make build
sudo cp ./bin/curo-prompt /usr/local/bin/  # 또는 ~/bin/
```

---

## 제거 (Uninstall)

### Go install로 설치한 경우

```bash
rm $(go env GOPATH)/bin/curo-prompt
```

### 수동 설치한 경우

```bash
# 시스템 경로에서 제거
sudo rm /usr/local/bin/curo-prompt

# 또는 사용자 로컬에서 제거
rm ~/bin/curo-prompt
```

---

## 설치 확인 체크리스트

- [ ] `curo-prompt --version` 명령이 작동함
- [ ] `which curo-prompt`로 경로 확인 가능
- [ ] PATH에 올바른 경로가 포함됨
- [ ] 새 터미널에서도 작동함 (셸 재시작 후 확인)

---

## 다음 단계

설치가 완료되었다면:

1. [시작하기 가이드](./GETTING_STARTED.md)로 실제 테스트 진행
2. [설정 가이드](./CONFIG.md)에서 설정 파일 구성
3. [아키텍처](./ARCHITECTURE.md) 문서로 이해도 향상

---

**문제가 발생하면**: GitHub Issues에 리포트하거나 문서를 확인하세요.

