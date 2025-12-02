# Changelog

T9s 프로젝트의 모든 주요 변경 사항을 기록합니다.

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

#### Changed
- **Main Entry Point**
  - `cmd/t9s/main.go`: NewApp() → NewAppNew() 사용
  - 버전: 0.1.0 → 0.2.0

#### Documentation
- 📄 `ARCHITECTURE.md`: k9s 스타일 아키텍처 상세 설명
- 📄 `MIGRATION.md`: v0.1.0에서 v0.2.0으로 마이그레이션 가이드
- 📄 `REFACTORING_SUMMARY.md`: 리팩토링 요약 및 개선 효과
- 📄 `DIRECTORY_TREE.md`: 디렉토리 구조 시각화
- 📄 `CHANGELOG.md`: 변경 이력 (이 파일)
- ✏️ `README.md`: v0.2.0 로드맵 업데이트

#### Improved
- **코드 품질**
  - 관심사 분리 (Separation of Concerns)
  - 단일 파일 650+ 줄 → 여러 파일로 분산 (각 50-200줄)
  - 재사용 가능한 컴포넌트 8개 생성

- **유지보수성**
  - 명확한 패키지 구조
  - 예측 가능한 코드 위치
  - 독립적인 테스트 가능

- **확장성**
  - 새로운 뷰 추가 용이
  - DAO 패턴으로 데이터 접근 표준화
  - 컴포넌트 재사용 가능

### 📊 Statistics
- **새로 생성된 파일**: 13개
- **새로운 패키지**: 4개 (model, dao, view, components/dialog)
- **총 Go 파일**: 5개 → 18개
- **빌드 크기**: 약 4.9MB (변화 없음)

### 🔄 Migration Path
- Phase 1: ✅ 새 구조 생성 (완료)
- Phase 2: 📝 새 기능 추가 (진행 예정)
- Phase 3: 🗑️ 레거시 제거 (추후)

### 🎯 Breaking Changes
없음 - 레거시 코드(`app.go`) 유지로 하위 호환성 보장

---

## [0.1.0] - 2025-12-01

### Added
- 🎉 초기 릴리스
- 📁 Tree View 기반 파일 탐색
- 🚀 Terraform Plan/Apply 실행
- ⏰ Terraform History 조회
- ⎈ Helm List 통합 (`helm list -A`)
- ✏️ 파일 편집 기능 (`$EDITOR` 연동)
- ⚙️ 설정 관리 (Settings)
- 🎨 k9s 스타일 UI
- 🛡️ Apply 전 확인 모달

### Features
- `internal/ui/app.go`: 단일 파일 기반 UI 애플리케이션
- `internal/terraform/manager.go`: Terraform 작업 관리
- `internal/git/manager.go`: Git 통합
- `internal/config/config.go`: 설정 파일 관리
- `cmd/t9s/main.go`: CLI 진입점

### Documentation
- 📄 `README.md`: 프로젝트 소개 및 사용법
- 📄 `QUICKSTART.md`: 빠른 시작 가이드
- 📄 `STRUCTURE.md`: 프로젝트 구조 설명
- 📄 `TODO.md`: 로드맵

### Infrastructure
- 🔧 `install.sh`: 설치 스크립트
- 📦 `go.mod`, `go.sum`: Go 모듈 의존성

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

