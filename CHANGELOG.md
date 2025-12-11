# Changelog

T9s 프로젝트의 모든 주요 변경 사항을 기록합니다.

## [0.3.0] - 2025-12-11

### 🚀 Major Features

#### Terraform 실행 개선
- **Execute/Auto Approve 버튼 분리**
  - Execute: Plan 결과를 보여준 후 Yes/No 다이얼로그로 확인
  - Auto Approve: `-auto-approve` 플래그로 즉시 실행
  - Cancel: 취소

- **실시간 스트리밍 출력**
  - Terraform 실행 결과를 실시간으로 화면에 표시
  - 자동 스크롤로 최신 로그 항상 표시
  - 사용자 입력과 분리되어 안전한 스크롤

- **Init/Destroy 템플릿 추가**
  - 설정에서 Init/Destroy 명령어 템플릿 커스터마이징
  - `init_template`, `destroy_template` 설정 필드

#### 히스토리 기능
- **SQLite 기반 히스토리 DB** (`~/.t9s/history.db`)
  - Apply/Destroy 실행 이력 영구 저장
  - 사용자, 브랜치, tfvars 내용 기록
  - 타임스탬프 정확히 저장 (RFC3339 포맷)

- **히스토리 뷰 개선**
  - 사용자/브랜치 정보 표시
  - `Shift+M`: 상세 내용(tfvars/config) 토글
  - `d`: 더 보기 (Load More)
  - `u`: 접기 (Load Less)

#### Git 통합
- **브랜치 전환** (`Shift+B`)
  - 로컬 브랜치 목록 표시
  - 현재 브랜치 표시 (헤더에 ● 또는 ✓)
  - Dirty 상태 처리: Stash/Commit/Force 옵션

#### UI/UX 개선
- **Help View** (`?` 또는 `Shift+H`)
  - 모든 단축키 카테고리별 표시
  - Resource, General, Git, History 섹션

- **Command Mode** (`/`)
  - 현재 디렉토리 표시
  - 커스텀 명령어 실행

- **Home 화면** (`Shift+C`)
  - Available Commands 표시
  - 디렉토리 이동 시 자동 표시

- **포커스 자동 전환**
  - Apply/Destroy 실행 시 Content View로 자동 전환
  - 스크롤 키가 Terraform 입력으로 들어가는 것 방지

### 📝 설정 변경
- `init_template`: Terraform Init 명령어 템플릿
- `destroy_template`: Terraform Destroy 명령어 템플릿
- `tfvars_file`: 기본 tfvars 파일 (기존 `var_file`에서 변경)
- `init_conf_file`: Init config 파일 (기존 `backend_config`에서 변경)

### 🔧 Internal
- `internal/db/history.go`: SQLite 히스토리 DB
- `internal/view/help_view.go`: 도움말 뷰
- `internal/view/history_view.go`: 히스토리 뷰
- `internal/view/command_view.go`: 커맨드 입력 뷰
- `internal/ui/dialog/branch.go`: 브랜치 선택 다이얼로그
- `internal/ui/dialog/commit.go`: 커밋 다이얼로그
- `internal/ui/dialog/dirty_branch.go`: 더티 브랜치 다이얼로그
- `internal/ui/dialog/file_selection.go`: 파일 선택 다이얼로그
- `internal/ui/dialog/terraform_confirm.go`: Terraform 확인 다이얼로그

---

## [0.2.0] - 2025-12-02

### 🎨 Architecture - k9s 스타일 적용

#### Added
- **Model Layer** (`internal/model/`)
  - `terraform.go`: TerraformDirectory, TerraformStatus, HelmRelease 모델
  - `git.go`: GitStatus 모델

- **DAO Layer** (`internal/dao/`)
  - `terraform.go`: TerraformDAO - 데이터 접근 계층
  - `git.go`: GitDAO - Git 데이터 접근

- **View Layer** (`internal/view/`)
  - `tree_view.go`: TreeView - 파일 트리 뷰 컴포넌트
  - `content_view.go`: ContentView - 컨텐츠 표시 뷰
  - `header_view.go`: HeaderView - 헤더 뷰
  - `status_bar.go`: StatusBar - 상태바 뷰

- **UI Components** (`internal/ui/components/`)
  - `executor.go`: CommandExecutor - 명령 실행 컴포넌트

- **Dialogs** (`internal/ui/dialog/`)
  - `confirm.go`: ConfirmDialog - 확인 다이얼로그
  - `settings.go`: SettingsDialog - 설정 다이얼로그

- **New App** (`internal/ui/`)
  - `app_new.go`: 새로운 아키텍처 기반 애플리케이션

#### Improved
- **코드 품질**
  - 관심사 분리 (Separation of Concerns)
  - 단일 파일 650+ 줄 → 여러 파일로 분산 (각 50-200줄)
  - 재사용 가능한 컴포넌트 8개 생성

---

## [0.1.0] - 2025-12-01

### Added
- 🎉 초기 릴리스
- 📁 Tree View 기반 파일 탐색
- 🚀 Terraform Plan/Apply 실행
- ⏰ Terraform History 조회
- ✏️ 파일 편집 기능 (`$EDITOR` 연동)
- ⚙️ 설정 관리 (Settings)
- 🎨 k9s 스타일 UI
- 🛡️ Apply 전 확인 모달

---

## 버전 관리 규칙

이 프로젝트는 [Semantic Versioning](https://semver.org/)을 따릅니다:

- **MAJOR**: 호환되지 않는 API 변경
- **MINOR**: 하위 호환되는 기능 추가
- **PATCH**: 하위 호환되는 버그 수정

### 변경 유형

- `Added`: 새로운 기능
- `Changed`: 기존 기능 변경
- `Deprecated`: 곧 제거될 기능
- `Removed`: 제거된 기능
- `Fixed`: 버그 수정
- `Security`: 보안 관련 수정
