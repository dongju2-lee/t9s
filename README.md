# T9s - Terraform Infrastructure Manager

![T9s Logo](https://img.shields.io/badge/T9s-Terraform%20TUI-blue)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.20-00ADD8)
![Version](https://img.shields.io/badge/version-v0.3.0-green)

**T9s**는 k9s에서 영감을 받은 Terraform 인프라 관리를 위한 터미널 UI 도구입니다. 복잡한 Terraform 작업을 직관적인 TUI 환경에서 쉽고 안전하게 수행할 수 있습니다.

## ✨ 주요 기능

- 📁 **Tree View 탐색**: 현재 디렉토리의 모든 파일과 폴더를 트리 구조로 직관적으로 탐색
- 🚀 **Terraform 실행**: Init, Plan, Apply, Destroy 명령어를 TUI에서 직접 실행
- 📝 **파일 선택 다이얼로그**: 실행 시 `.conf` 또는 `.tfvars` 파일을 목록에서 선택 (미리보기 지원)
- 🔄 **실시간 스트리밍 출력**: Terraform 실행 결과를 실시간으로 확인
- ✅ **Execute/Auto Approve 분리**: 일반 실행(Yes/No 확인) vs 자동 승인 선택 가능
- 🎨 **ANSI 컬러 지원**: Terraform의 컬러풀한 출력을 그대로 TUI에서 확인
- ⚡ **빠른 스크롤**: `u`/`d`, `Shift+방향키` 등을 이용한 대용량 로그의 빠른 탐색
- ⏰ **History 추적**: 각 디렉토리의 실행 이력(사용자, 브랜치, tfvars 내용) SQLite에 자동 저장
- 🔀 **Git Branch 전환**: `Shift+B`로 브랜치 전환 (Stash/Commit/Force 옵션)
- ✏️ **파일 편집**: 내장된 편집 기능(`$EDITOR` 연동)으로 tfvars 및 설정 파일 수정
- ⚙️ **유연한 설정**: Init/Plan/Apply/Destroy 명령어 템플릿 커스터마이징
- 🛡️ **안전한 작업**: 실행 전 확인 모달 및 tfvars/config 내용 미리보기
- 🎨 **k9s 스타일 UI**: 검정 배경의 깔끔하고 전문적인 인터페이스

## 📦 설치

### 방법 1: 소스에서 빌드

```bash
git clone https://github.com/idongju/t9s.git
cd t9s
go build -o t9s ./cmd/t9s
sudo mv t9s /usr/local/bin/
```

### 방법 2: 간편 설치 스크립트

```bash
./install.sh
```

## 🚀 사용법

### 기본 실행

```bash
t9s
```

### 버전 확인

```bash
t9s --version
```

## ⌨️ 키보드 단축키

### 전역 / 네비게이션
| 키 | 설명 |
|---|---|
| `Tab` | 포커스 전환 (File Tree ↔ Content View) |
| `q` | 종료 (Quit) |
| `?` / `Shift+H` | 도움말 (Help) |
| `/` | 커맨드 모드 |
| `Shift+C` | 홈 화면 (Available Commands) |
| `Esc` | 뒤로 가기 / 다이얼로그 닫기 |

### File Tree 포커스 시
| 키 | 설명 |
|---|---|
| `↑/↓` | 파일/폴더 탐색 |
| `Enter` | 디렉토리 확장/축소 또는 파일 선택 |
| `i` | **Init**: Terraform Init (설정 파일 선택) |
| `p` | **Plan**: Terraform Plan (tfvars 선택) |
| `a` | **Apply**: Terraform Apply (tfvars 선택) |
| `d` | **Destroy**: Terraform Destroy (tfvars 선택) |
| `h` | **History**: Terraform 실행 이력 확인 |
| `e` | **Edit**: 선택된 파일 편집 (`$EDITOR`) |
| `s` | **Settings**: 설정 창 열기 |
| `Shift+B` | **Branch**: Git 브랜치 전환 |

### Content View 포커스 시 (로그/출력 확인)
| 키 | 설명 |
|---|---|
| `↑/↓` | 스크롤 (1줄) |
| `u` / `d` | 빠른 스크롤 (10줄) |
| `Shift + ↑/↓` | 빠른 스크롤 (10줄) |
| `PageUp` / `PageDown` | 페이지 스크롤 |
| `Home` / `End` | 맨 위/아래로 이동 |

### History View (이력 확인)
| 키 | 설명 |
|---|---|
| `d` | 이력 더보기 (Load More) |
| `u` | 이력 접기 (Load Less) |
| `Shift+M` | 상세 내용(tfvars/config) 토글 |
| `Esc` | 뒤로 가기 |

### Confirmation Dialog (확인 창)
| 버튼 | 설명 |
|---|---|
| **Execute** | 일반 실행 (Terraform이 Yes/No 물어봄 → 자동 Yes) |
| **Auto Approve** | `-auto-approve` 플래그로 즉시 실행 |
| **Cancel** | 취소 |

## 🔧 설정

T9s는 `~/.t9s/config.yaml` 파일을 통해 설정을 관리합니다. 앱 내에서 `s` 키를 눌러 쉽게 수정할 수 있습니다.

```yaml
# Terraform 루트 디렉토리
terraform_root: /path/to/your/terraform

# Terraform 명령어 템플릿
commands:
  # {initconf}은 init 시 선택된 conf 파일 경로로 치환됩니다.
  # {varfile}은 선택된 tfvars 파일 경로로 자동 치환됩니다.
  init_template: "terraform init -backend-config={initconf}"
  plan_template: "terraform plan -var-file={varfile}"
  apply_template: "terraform apply -var-file={varfile}"
  destroy_template: "terraform destroy -var-file={varfile}"
  
  # 기본 파일 설정
  tfvars_file: "config/env.tfvars"
  init_conf_file: "config/env.conf"

# 기본 설정
defaults:
  auto_refresh: true
  refresh_interval: 60
```

## 📁 데이터 저장 위치

| 파일 | 경로 | 설명 |
|---|---|---|
| 설정 파일 | `~/.t9s/config.yaml` | 앱 설정 |
| 히스토리 DB | `~/.t9s/history.db` | Apply/Destroy 실행 이력 (SQLite) |

## 🛠️ 개발 로드맵

### v0.1.0 ✅
- [x] Tree View 기반 파일 탐색
- [x] Terraform Plan/Apply 실행
- [x] 설정 관리 (Settings)
- [x] 외부 에디터 연동
- [x] k9s 스타일 UI

### v0.2.0 ✅ - 아키텍처 개선
- [x] **k9s 스타일 아키텍처 적용** (Model/DAO/View)
- [x] UI 컴포넌트 분리 (Header, Tree, Content, StatusBar)
- [x] 다이얼로그 시스템 (Confirm, Settings, FileSelection)

### v0.3.0 ✅ (Current) - UX 및 기능 강화
- [x] **파일 선택 다이얼로그**: Init/Plan/Apply/Destroy 시 설정 파일 선택
- [x] **Terraform Init/Destroy 템플릿**: 설정에서 커스터마이징
- [x] **히스토리 기능**: Apply/Destroy 실행 이력 SQLite 저장
- [x] **히스토리 상세보기**: 사용자, 브랜치, tfvars 내용 표시
- [x] **Git Branch 전환**: Shift+B로 브랜치 전환 (Stash/Commit/Force)
- [x] **실시간 스트리밍 출력**: Terraform 로그 실시간 표시
- [x] **Execute/Auto Approve 분리**: Yes/No 확인 다이얼로그
- [x] **Help View**: 단축키 도움말 화면
- [x] **Command Mode**: `/`로 커맨드 입력
- [x] **Home 화면**: Shift+C로 Available Commands 표시
- [x] **ANSI 컬러 지원**: Terraform 출력 컬러 유지
- [x] **빠른 스크롤**: `u`/`d` 키 및 `Shift+방향키` 지원

### v0.4.0 (Next)
- [ ] Terraform Workspace 전환 UI
- [ ] State 정보 테이블 뷰
- [ ] Terraform Drift 감지
- [ ] 리소스 변경 이력 추적

## 📚 문서

| 문서 | 설명 |
|------|------|
| [QUICKSTART.md](QUICKSTART.md) | 빠른 시작 가이드 |
| [SETTINGS_GUIDE.md](SETTINGS_GUIDE.md) | 설정 가이드 |
| [ARCHITECTURE.md](ARCHITECTURE.md) | 아키텍처 설명 |
| [CHANGELOG.md](CHANGELOG.md) | 변경 이력 |

## 🤝 기여하기

이슈 및 PR은 언제나 환영합니다!

```bash
# 개발 환경 설정
git clone https://github.com/idongju/t9s.git
cd t9s
go mod download
go build ./...
go run ./cmd/t9s
```

## 📄 라이선스

MIT License

## 🙏 크레딧

- [tview](https://github.com/rivo/tview) - Terminal UI framework
- [k9s](https://k9scli.io/) - 영감을 준 Kubernetes CLI
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver
