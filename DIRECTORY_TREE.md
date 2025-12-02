# T9s 디렉토리 구조 (v0.2.0)

```
T9s/
│
├── 📄 README.md                    # 프로젝트 개요 및 사용법
├── 📄 QUICKSTART.md                # 빠른 시작 가이드
├── 📄 TODO.md                      # 로드맵 및 TODO
│
├── 📄 STRUCTURE.md                 # 기존 구조 문서 (v0.1.0)
├── 📄 ARCHITECTURE.md              # 새로운 아키텍처 설명 (v0.2.0) ⭐
├── 📄 MIGRATION.md                 # 마이그레이션 가이드 ⭐
├── 📄 REFACTORING_SUMMARY.md       # 리팩토링 요약 ⭐
├── 📄 DIRECTORY_TREE.md            # 이 파일
│
├── 📄 go.mod                       # Go 모듈 정의
├── 📄 go.sum                       # 의존성 체크섬
├── 📄 .gitignore                   # Git 제외 파일
├── 🔧 install.sh                   # 설치 스크립트
│
├── 📁 cmd/                         # 명령줄 애플리케이션
│   └── 📁 t9s/
│       └── 📄 main.go              # CLI 진입점 (v0.2.0 - NewAppNew 사용)
│
└── 📁 internal/                    # 내부 패키지
    │
    ├── 📁 model/                   # ⭐ 데이터 모델 (k9s 스타일)
    │   ├── 📄 terraform.go         # Terraform 관련 모델
    │   │   └── TerraformDirectory, TerraformStatus, HelmRelease
    │   └── 📄 git.go               # Git 관련 모델
    │       └── GitStatus
    │
    ├── 📁 dao/                     # ⭐ Data Access Object (k9s 스타일)
    │   ├── 📄 terraform.go         # Terraform 데이터 접근
    │   │   └── TerraformDAO
    │   │       ├── ListDirectories()
    │   │       ├── CheckDrift()
    │   │       ├── Plan()
    │   │       ├── Apply()
    │   │       └── GetHelmReleases()
    │   └── 📄 git.go               # Git 데이터 접근
    │       └── GitDAO
    │           ├── GetStatus()
    │           └── GetDiff()
    │
    ├── 📁 view/                    # ⭐ UI View 컴포넌트 (k9s 스타일)
    │   ├── 📄 tree_view.go         # 파일 트리 뷰
    │   │   └── TreeView
    │   ├── 📄 content_view.go      # 컨텐츠 표시 뷰
    │   │   └── ContentView
    │   ├── 📄 header_view.go       # 헤더 뷰
    │   │   └── HeaderView
    │   └── 📄 status_bar.go        # 상태바 뷰
    │       └── StatusBar
    │
    ├── 📁 ui/                      # UI 관련
    │   ├── 📄 app.go               # 레거시 앱 (호환성 유지)
    │   │   └── App (v0.1.0)
    │   ├── 📄 app_new.go           # ⭐ 새로운 구조의 앱
    │   │   └── AppNew (v0.2.0)
    │   │
    │   ├── 📁 components/          # ⭐ 재사용 가능한 UI 컴포넌트
    │   │   └── 📄 executor.go      # 명령 실행기
    │   │       └── CommandExecutor
    │   │           ├── ExecutePlan()
    │   │           ├── ExecuteApply()
    │   │           ├── ShowHistory()
    │   │           ├── ShowHelm()
    │   │           └── EditFile()
    │   │
    │   └── 📁 dialog/              # ⭐ 다이얼로그 컴포넌트
    │       ├── 📄 confirm.go       # 확인 다이얼로그
    │       │   └── ConfirmDialog
    │       └── 📄 settings.go      # 설정 다이얼로그
    │           └── SettingsDialog
    │
    ├── 📁 config/                  # 설정 관리
    │   └── 📄 config.go            # 설정 파일 로드/저장
    │       └── Config
    │
    ├── 📁 terraform/               # 레거시 (추후 제거 예정)
    │   └── 📄 manager.go           # Terraform 매니저 (v0.1.0)
    │
    └── 📁 git/                     # 레거시 (추후 제거 예정)
        └── 📄 manager.go           # Git 매니저 (v0.1.0)
```

## 패키지별 역할

### 🆕 새로 추가된 패키지 (v0.2.0)

#### 1. `internal/model/` - 데이터 모델
- **목적**: 순수한 데이터 구조 정의
- **특징**: 비즈니스 로직 없음, 재사용 가능
- **파일**:
  - `terraform.go`: TerraformDirectory, TerraformStatus, HelmRelease
  - `git.go`: GitStatus

#### 2. `internal/dao/` - 데이터 접근 계층
- **목적**: 외부 시스템과의 상호작용
- **특징**: CLI 실행, 파일 시스템 접근
- **파일**:
  - `terraform.go`: TerraformDAO - Terraform 작업
  - `git.go`: GitDAO - Git 작업

#### 3. `internal/view/` - UI 뷰 컴포넌트
- **목적**: 재사용 가능한 UI 컴포넌트
- **특징**: tview 위젯 래핑, 독립적 동작
- **파일**:
  - `tree_view.go`: TreeView - 파일 트리
  - `content_view.go`: ContentView - 메인 컨텐츠
  - `header_view.go`: HeaderView - 헤더
  - `status_bar.go`: StatusBar - 상태바

#### 4. `internal/ui/components/` - UI 컴포넌트
- **목적**: 복잡한 UI 로직
- **파일**:
  - `executor.go`: CommandExecutor - 명령 실행

#### 5. `internal/ui/dialog/` - 다이얼로그
- **목적**: 모달 다이얼로그
- **파일**:
  - `confirm.go`: ConfirmDialog - 확인 창
  - `settings.go`: SettingsDialog - 설정 창

### 🔄 기존 패키지

#### `cmd/t9s/` - CLI 진입점
- `main.go`: v0.2.0에서 NewAppNew() 사용

#### `internal/ui/` - UI 애플리케이션
- `app.go`: 레거시 (v0.1.0, 호환성 유지)
- `app_new.go`: 새 구조 (v0.2.0) ⭐

#### `internal/config/` - 설정 관리
- `config.go`: YAML 설정 로드/저장

#### `internal/terraform/` - 레거시 (추후 제거)
- `manager.go`: v0.1.0 Terraform 매니저

#### `internal/git/` - 레거시 (추후 제거)
- `manager.go`: v0.1.0 Git 매니저

## 데이터 흐름

```
사용자 입력 (키보드)
    ↓
AppNew (internal/ui/app_new.go)
    ↓
CommandExecutor (internal/ui/components/executor.go)
    ↓
TerraformDAO / GitDAO (internal/dao/)
    ↓
Terraform CLI / Git CLI
    ↓
Model (internal/model/)
    ↓
View (internal/view/)
    ↓
화면 표시
```

## 파일 크기 비교

### v0.1.0
- `internal/ui/app.go`: ~650 줄 (단일 파일)

### v0.2.0
- `internal/ui/app_new.go`: ~200 줄
- `internal/view/tree_view.go`: ~100 줄
- `internal/view/content_view.go`: ~80 줄
- `internal/view/header_view.go`: ~90 줄
- `internal/view/status_bar.go`: ~50 줄
- `internal/ui/components/executor.go`: ~180 줄
- `internal/ui/dialog/confirm.go`: ~30 줄
- `internal/ui/dialog/settings.go`: ~80 줄

**총 라인 수**: 비슷하지만 **관심사 분리**로 **유지보수성 대폭 향상** ⬆️

## k9s 스타일 비교

| k9s | T9s | 상태 |
|-----|-----|------|
| internal/model/ | internal/model/ | ✅ 구현 |
| internal/dao/ | internal/dao/ | ✅ 구현 |
| internal/view/ | internal/view/ | ✅ 구현 |
| internal/render/ | (미구현) | 📝 추후 |
| internal/ui/ | internal/ui/ | ✅ 구현 |
| internal/config/ | internal/config/ | ✅ 구현 |

## 빌드 정보

```bash
# 빌드
go build -o t9s ./cmd/t9s

# 결과
-rwxr-xr-x  t9s  4.9MB

# 버전
./t9s --version
# T9s version 0.2.0
```

## 다음 단계

1. ✅ k9s 스타일 아키텍처 적용
2. 📝 새로운 기능 추가 (새 구조 활용)
3. 🗑️ 레거시 코드 제거
4. 🧪 테스트 작성

---

**주의**: ⭐ 표시된 항목은 v0.2.0에서 새로 추가된 부분입니다.

