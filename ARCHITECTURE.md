# T9s Architecture (v0.3.0)

k9s를 참고하여 설계된 T9s 아키텍처입니다.

## 디렉토리 구조

```
T9s/
├── cmd/
│   └── t9s/
│       └── main.go                 # CLI 진입점
│
└── internal/
    ├── config/                     # 설정 관리
    │   └── config.go               # YAML 설정 로드/저장
    │
    ├── db/                         # 데이터베이스
    │   └── history.go              # SQLite 히스토리 DB
    │
    ├── git/                        # Git 통합
    │   └── manager.go              # Git 명령 실행
    │
    ├── terraform/                  # Terraform 통합
    │   └── manager.go              # Terraform 명령 실행
    │
    ├── model/                      # 데이터 모델 (k9s 스타일)
    │   ├── terraform.go            # Terraform 관련 모델
    │   └── git.go                  # Git 관련 모델
    │
    ├── dao/                        # Data Access Object (k9s 스타일)
    │   ├── terraform.go            # Terraform 데이터 접근
    │   └── git.go                  # Git 데이터 접근
    │
    ├── view/                       # UI View 컴포넌트 (k9s 스타일)
    │   ├── tree_view.go            # 파일 트리 뷰
    │   ├── content_view.go         # 컨텐츠 표시 뷰
    │   ├── header_view.go          # 헤더 뷰 (브랜치 표시)
    │   ├── status_bar.go           # 상태바 뷰
    │   ├── help_view.go            # 도움말 뷰
    │   ├── history_view.go         # 히스토리 뷰
    │   └── command_view.go         # 커맨드 입력 뷰
    │
    └── ui/                         # UI 관련
        ├── app.go                  # 레거시 앱
        ├── app_new.go              # 새로운 구조의 앱
        │
        ├── components/             # 재사용 가능한 UI 컴포넌트
        │   ├── executor.go         # 명령 실행기
        │   └── terraform_helper.go # Terraform 헬퍼
        │
        └── dialog/                 # 다이얼로그 컴포넌트
            ├── confirm.go          # 기본 확인 다이얼로그
            ├── settings.go         # 설정 다이얼로그
            ├── file_selection.go   # 파일 선택 다이얼로그
            ├── terraform_confirm.go # Terraform 확인 (Execute/Auto/Cancel)
            ├── branch.go           # 브랜치 선택 다이얼로그
            ├── commit.go           # 커밋 다이얼로그
            └── dirty_branch.go     # 더티 브랜치 처리 다이얼로그
```

## 계층 구조 (k9s 스타일)

### 1. Config Layer (internal/config/)
- **목적**: 설정 파일 관리
- **특징**:
  - YAML 설정 로드/저장
  - 기본 설정 생성
  - Terraform 명령어 템플릿 관리

### 2. Database Layer (internal/db/)
- **목적**: 영구 데이터 저장
- **특징**:
  - SQLite 기반 히스토리 DB
  - Apply/Destroy 실행 이력 저장
  - 사용자, 브랜치, tfvars 내용 기록

### 3. Model Layer (internal/model/)
- **목적**: 순수한 데이터 구조 정의
- **특징**:
  - 비즈니스 로직 없음
  - 다른 패키지에 대한 의존성 최소화
  - 재사용 가능한 타입 정의

### 4. DAO Layer (internal/dao/)
- **목적**: 데이터 접근 및 외부 시스템과의 상호작용
- **특징**:
  - Terraform CLI 실행
  - Git 명령 실행
  - 파일 시스템 접근

### 5. View Layer (internal/view/)
- **목적**: UI 컴포넌트 정의
- **특징**:
  - tview 위젯 래핑
  - 재사용 가능한 UI 컴포넌트
  - 각 뷰는 독립적으로 동작

**파일들**:
| 파일 | 컴포넌트 | 설명 |
|------|----------|------|
| `tree_view.go` | TreeView | 파일 트리 |
| `content_view.go` | ContentView | 메인 컨텐츠 영역 |
| `header_view.go` | HeaderView | 헤더 (브랜치, 경로) |
| `status_bar.go` | StatusBar | 하단 상태바 |
| `help_view.go` | HelpView | 도움말 |
| `history_view.go` | HistoryView | 히스토리 |
| `command_view.go` | CommandView | 커맨드 입력 |

### 6. UI Layer (internal/ui/)
- **목적**: 애플리케이션 조합 및 이벤트 핸들링
- **특징**:
  - View 컴포넌트들을 조합
  - 키보드 이벤트 처리
  - 애플리케이션 상태 관리

**구성**:
| 파일/폴더 | 설명 |
|-----------|------|
| `app_new.go` | 새로운 구조의 메인 앱 |
| `components/executor.go` | 명령 실행 로직 |
| `components/terraform_helper.go` | Terraform 헬퍼 함수 |
| `dialog/` | 다이얼로그 컴포넌트들 |

## 데이터 흐름

```
사용자 입력 (키보드)
    ↓
UI Layer (app_new.go)
    ↓
Dialog (terraform_confirm.go 등)
    ↓
Components (executor.go)
    ↓
Git Manager / Terraform Manager
    ↓
외부 시스템 (terraform CLI, git CLI)
    ↓
DB Layer (history.go) - 이력 저장
    ↓
View Layer (content_view.go 등)
    ↓
화면 표시
```

## Terraform 실행 흐름

```
사용자가 'a' (Apply) 입력
    ↓
showApplyConfirmation() - 파일 선택 다이얼로그
    ↓
showApplyConfirmationWithFile() - 확인 다이얼로그
    ↓
┌─────────────────────────────────────────┐
│  [Execute] [Auto Approve] [Cancel]      │
└─────────────────────────────────────────┘
    ↓
executeTerraformCommand()
    - 실시간 스트리밍 출력
    - (Execute) "Enter a value:" 감지 시 Yes/No 다이얼로그
    - 히스토리 DB에 저장
    ↓
결과 표시
```

## k9s와의 비교

| k9s | T9s | 설명 |
|-----|-----|------|
| internal/model/ | internal/model/ | 데이터 모델 |
| internal/dao/ | internal/dao/ | 데이터 접근 계층 |
| internal/view/ | internal/view/ | UI 뷰 컴포넌트 |
| internal/render/ | (미구현) | 렌더링 로직 |
| internal/ui/ | internal/ui/ | 애플리케이션 조합 |
| internal/config/ | internal/config/ | 설정 관리 |
| - | internal/db/ | 히스토리 DB |

## 장점

1. **관심사 분리**: 각 레이어가 명확한 책임을 가짐
2. **테스트 용이성**: 각 레이어를 독립적으로 테스트 가능
3. **재사용성**: View 컴포넌트를 다른 곳에서도 재사용 가능
4. **확장성**: 새로운 뷰나 기능을 쉽게 추가 가능
5. **유지보수성**: 코드 위치를 예측 가능하고 찾기 쉬움

## 새로운 기능 추가 예시

### 새로운 View 추가
```go
// internal/view/workspace_view.go
package view

type WorkspaceView struct {
    *tview.List
}

func NewWorkspaceView() *WorkspaceView {
    // 구현
}
```

### 새로운 Dialog 추가
```go
// internal/ui/dialog/workspace_dialog.go
package dialog

type WorkspaceDialog struct {
    *tview.Modal
}

func NewWorkspaceDialog(workspaces []string, onSelect func(string)) *WorkspaceDialog {
    // 구현
}
```

### 새로운 DB 기능 추가
```go
// internal/db/history.go
func (h *HistoryDB) GetByUser(user string, limit int) ([]*HistoryEntry, error) {
    // 구현
}
```
