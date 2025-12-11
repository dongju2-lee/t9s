# T9s 프로젝트 구조 (v0.3.0)

```
T9s/
├── cmd/
│   └── t9s/
│       └── main.go              # CLI 진입점
│
├── internal/
│   ├── config/
│   │   └── config.go            # 설정 관리 (YAML)
│   │
│   ├── db/
│   │   └── history.go           # SQLite 히스토리 DB
│   │
│   ├── git/
│   │   └── manager.go           # Git 통합
│   │                            # - 상태 확인, 브랜치 전환
│   │                            # - Stash, Commit, Diff
│   │
│   ├── terraform/
│   │   └── manager.go           # Terraform 작업 관리
│   │
│   ├── model/                   # 데이터 모델
│   │   ├── terraform.go         # Terraform 관련 모델
│   │   └── git.go               # Git 관련 모델
│   │
│   ├── dao/                     # Data Access Object
│   │   ├── terraform.go         # Terraform 데이터 접근
│   │   └── git.go               # Git 데이터 접근
│   │
│   ├── view/                    # UI View 컴포넌트
│   │   ├── tree_view.go         # 파일 트리 뷰
│   │   ├── content_view.go      # 컨텐츠 표시 뷰
│   │   ├── header_view.go       # 헤더 뷰 (브랜치 표시)
│   │   ├── status_bar.go        # 상태바 뷰
│   │   ├── help_view.go         # 도움말 뷰
│   │   ├── history_view.go      # 히스토리 뷰
│   │   └── command_view.go      # 커맨드 입력 뷰
│   │
│   └── ui/
│       ├── app.go               # 레거시 앱
│       ├── app_new.go           # 새로운 구조의 앱
│       ├── components/
│       │   ├── executor.go      # 명령 실행기
│       │   └── terraform_helper.go  # Terraform 헬퍼
│       └── dialog/
│           ├── confirm.go       # 확인 다이얼로그
│           ├── settings.go      # 설정 다이얼로그
│           ├── file_selection.go    # 파일 선택 다이얼로그
│           ├── terraform_confirm.go # Terraform 확인 다이얼로그
│           ├── branch.go        # 브랜치 선택 다이얼로그
│           ├── commit.go        # 커밋 다이얼로그
│           └── dirty_branch.go  # 더티 브랜치 다이얼로그
│
├── terra/                       # 테스트용 Terraform 코드
│   └── test/
│       ├── monitoring/
│       └── webapp/
│
├── README.md                    # 프로젝트 설명
├── QUICKSTART.md                # 빠른 시작 가이드
├── SETTINGS_GUIDE.md            # 설정 가이드
├── ARCHITECTURE.md              # 아키텍처 설명
├── CHANGELOG.md                 # 변경 이력
├── TODO.md                      # 로드맵
├── install.sh                   # 설치 스크립트
├── go.mod                       # Go 모듈 정의
└── go.sum                       # 의존성 체크섬
```

## 패키지 설명

### cmd/t9s/main.go
- 애플리케이션 진입점
- 버전 플래그 처리
- UI 앱 초기화 및 실행

### internal/config/config.go
- `~/.t9s/config.yaml` 관리
- Terraform 명령어 템플릿
- Init/Plan/Apply/Destroy 설정
- 기본 tfvars/conf 파일 경로

### internal/db/history.go
- SQLite 기반 히스토리 DB
- Apply/Destroy 실행 이력 저장
- 사용자, 브랜치, tfvars 내용 기록
- 타임스탬프 관리

### internal/git/manager.go
- Git 저장소 상태 확인
- 브랜치 목록 조회 및 전환
- Stash/Commit/Force Checkout
- Diff 조회

### internal/view/
- **tree_view.go**: 파일 트리 (좌측)
- **content_view.go**: 컨텐츠 표시 (우측)
- **header_view.go**: 헤더 (브랜치, 경로 표시)
- **status_bar.go**: 하단 상태바
- **help_view.go**: 도움말 화면
- **history_view.go**: 히스토리 화면
- **command_view.go**: 커맨드 입력

### internal/ui/dialog/
- **confirm.go**: 기본 확인 다이얼로그
- **settings.go**: 설정 다이얼로그
- **file_selection.go**: 파일 선택 (tfvars/conf)
- **terraform_confirm.go**: Terraform 실행 확인 (Execute/Auto Approve/Cancel)
- **branch.go**: 브랜치 선택
- **commit.go**: 커밋 메시지 입력
- **dirty_branch.go**: 더티 브랜치 처리 (Stash/Commit/Force)

## 의존성

```go
require (
    github.com/rivo/tview v0.42.0          // 터미널 UI 프레임워크
    github.com/gdamore/tcell/v2 v2.8.1     // 터미널 제어
    github.com/mattn/go-sqlite3 v1.14.24   // SQLite 드라이버
    gopkg.in/yaml.v3 v3.0.1                // YAML 파싱
)
```

## 데이터 저장

| 파일 | 경로 | 설명 |
|------|------|------|
| 설정 | `~/.t9s/config.yaml` | 앱 설정 |
| 히스토리 | `~/.t9s/history.db` | Apply/Destroy 이력 |

## 빌드

```bash
# 빌드
go build -o t9s ./cmd/t9s

# 설치
sudo mv t9s /usr/local/bin/

# 또는 스크립트 사용
./install.sh
```
