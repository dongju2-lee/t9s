# T9s 디렉토리 구조 (v0.3.0)

```
T9s/
│
├── 📄 README.md                    # 프로젝트 개요 및 사용법
├── 📄 QUICKSTART.md                # 빠른 시작 가이드
├── 📄 SETTINGS_GUIDE.md            # 설정 가이드
├── 📄 TODO.md                      # 로드맵 및 TODO
├── 📄 CHANGELOG.md                 # 변경 이력
│
├── 📄 STRUCTURE.md                 # 프로젝트 구조
├── 📄 ARCHITECTURE.md              # 아키텍처 설명
├── 📄 MIGRATION.md                 # 마이그레이션 가이드
├── 📄 REFACTORING_SUMMARY.md       # 리팩토링 요약
├── 📄 DIRECTORY_TREE.md            # 이 파일
│
├── 📄 go.mod                       # Go 모듈 정의
├── 📄 go.sum                       # 의존성 체크섬
├── 🔧 install.sh                   # 설치 스크립트
│
├── 📁 cmd/                         # 명령줄 애플리케이션
│   └── 📁 t9s/
│       └── 📄 main.go              # CLI 진입점
│
└── 📁 internal/                    # 내부 패키지
    │
    ├── 📁 config/                  # 설정 관리
    │   └── 📄 config.go            # YAML 설정 로드/저장
    │
    ├── 📁 db/                      # 데이터베이스 ⭐
    │   └── 📄 history.go           # SQLite 히스토리 DB
    │       └── HistoryDB, HistoryEntry
    │
    ├── 📁 git/                     # Git 통합
    │   └── 📄 manager.go           # Git 명령 실행
    │       └── Manager
    │           ├── GetStatus()
    │           ├── ListBranches()
    │           ├── CheckoutBranch()
    │           ├── StashChanges()
    │           ├── CommitChanges()
    │           └── ForceCheckout()
    │
    ├── 📁 terraform/               # Terraform 통합
    │   └── 📄 manager.go           # Terraform 명령 실행
    │
    ├── 📁 model/                   # 데이터 모델
    │   ├── 📄 terraform.go         # Terraform 관련 모델
    │   └── 📄 git.go               # Git 관련 모델
    │
    ├── 📁 dao/                     # Data Access Object
    │   ├── 📄 terraform.go         # Terraform 데이터 접근
    │   └── 📄 git.go               # Git 데이터 접근
    │
    ├── 📁 view/                    # UI View 컴포넌트
    │   ├── 📄 tree_view.go         # 파일 트리 뷰
    │   │   └── TreeView
    │   ├── 📄 content_view.go      # 컨텐츠 표시 뷰
    │   │   └── ContentView
    │   ├── 📄 header_view.go       # 헤더 뷰 (브랜치 표시) ⭐
    │   │   └── HeaderView
    │   │       └── SetGitBranch()
    │   ├── 📄 status_bar.go        # 상태바 뷰
    │   │   └── StatusBar
    │   ├── 📄 help_view.go         # 도움말 뷰 ⭐
    │   │   └── HelpView
    │   ├── 📄 history_view.go      # 히스토리 뷰 ⭐
    │   │   └── HistoryView
    │   └── 📄 command_view.go      # 커맨드 입력 뷰 ⭐
    │       └── CommandView
    │
    └── 📁 ui/                      # UI 관련
        ├── 📄 app.go               # 레거시 앱
        ├── 📄 app_new.go           # 새로운 구조의 앱 ⭐
        │   └── AppNew
        │       ├── executeTerraformCommand()
        │       ├── showApplyConfirmDialog()
        │       ├── showBranchSelection()
        │       └── showHistory()
        │
        ├── 📁 components/          # 재사용 가능한 UI 컴포넌트
        │   ├── 📄 executor.go      # 명령 실행기
        │   │   └── CommandExecutor
        │   └── 📄 terraform_helper.go  # Terraform 헬퍼 ⭐
        │       └── GetTerraformCommandInfo()
        │
        └── 📁 dialog/              # 다이얼로그 컴포넌트
            ├── 📄 confirm.go       # 기본 확인 다이얼로그
            ├── 📄 settings.go      # 설정 다이얼로그
            ├── 📄 file_selection.go    # 파일 선택 다이얼로그 ⭐
            │   └── FileSelectionDialog
            ├── 📄 terraform_confirm.go # Terraform 확인 ⭐
            │   └── TerraformConfirmDialog
            │       └── (Execute/Auto Approve/Cancel)
            ├── 📄 branch.go        # 브랜치 선택 ⭐
            │   └── BranchDialog
            ├── 📄 commit.go        # 커밋 다이얼로그 ⭐
            │   └── CommitDialog
            └── 📄 dirty_branch.go  # 더티 브랜치 처리 ⭐
                └── DirtyBranchDialog
                    └── (Stash/Commit/Force/Cancel)
```

## 패키지별 역할

### 🆕 v0.3.0에서 추가된 패키지/파일

#### `internal/db/` - 데이터베이스
- **목적**: 영구 데이터 저장
- **파일**:
  - `history.go`: SQLite 히스토리 DB
    - `HistoryDB`: DB 연결 및 쿼리
    - `HistoryEntry`: 히스토리 엔트리 모델

#### `internal/view/` - 추가된 뷰들
- **help_view.go**: 도움말 화면
- **history_view.go**: 히스토리 화면 (페이지네이션, 상세보기)
- **command_view.go**: 커맨드 입력 모드

#### `internal/ui/dialog/` - 추가된 다이얼로그들
- **file_selection.go**: 파일 선택 (미리보기 지원)
- **terraform_confirm.go**: Terraform 실행 확인 (3버튼)
- **branch.go**: Git 브랜치 선택
- **commit.go**: 커밋 메시지 입력
- **dirty_branch.go**: 더티 브랜치 처리

---

## 데이터 흐름

```
사용자 입력 (키보드)
    ↓
AppNew (internal/ui/app_new.go)
    ↓
Dialog (file_selection, terraform_confirm 등)
    ↓
Components (executor.go, terraform_helper.go)
    ↓
Git Manager / Terraform CLI
    ↓
HistoryDB (internal/db/history.go)
    ↓
View (content_view, history_view 등)
    ↓
화면 표시
```

---

## 파일 크기

| 파일 | 라인 수 | 역할 |
|------|---------|------|
| `app_new.go` | ~1100 | 메인 앱 로직 |
| `history.go` (db) | ~240 | 히스토리 DB |
| `content_view.go` | ~200 | 컨텐츠 뷰 |
| `header_view.go` | ~130 | 헤더 뷰 |
| `tree_view.go` | ~150 | 트리 뷰 |
| `help_view.go` | ~150 | 도움말 뷰 |
| `history_view.go` | ~200 | 히스토리 뷰 |
| `terraform_confirm.go` | ~100 | 확인 다이얼로그 |
| `file_selection.go` | ~150 | 파일 선택 다이얼로그 |

---

## 빌드 정보

```bash
# 빌드
go build -o t9s ./cmd/t9s

# 결과
-rwxr-xr-x  t9s  ~5MB

# 버전
./t9s --version
# T9s version 0.3.0
```

---

**주의**: ⭐ 표시된 항목은 v0.3.0에서 새로 추가/수정된 부분입니다.
