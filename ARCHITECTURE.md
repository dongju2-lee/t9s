# T9s Architecture

k9s를 참고하여 개선된 T9s 아키텍처입니다.

## 디렉토리 구조

```
T9s/
├── cmd/
│   └── t9s/
│       └── main.go                 # CLI 진입점
│
├── internal/
│   ├── model/                      # 데이터 모델 (k9s의 model 패키지 참고)
│   │   ├── terraform.go           # Terraform 관련 모델
│   │   └── git.go                 # Git 관련 모델
│   │
│   ├── dao/                        # Data Access Object (k9s의 dao 패키지 참고)
│   │   ├── terraform.go           # Terraform 데이터 접근 계층
│   │   └── git.go                 # Git 데이터 접근 계층
│   │
│   ├── view/                       # UI View 컴포넌트 (k9s의 view 패키지 참고)
│   │   ├── tree_view.go           # 파일 트리 뷰
│   │   ├── content_view.go        # 컨텐츠 표시 뷰
│   │   ├── header_view.go         # 헤더 뷰
│   │   └── status_bar.go          # 상태바 뷰
│   │
│   ├── ui/                         # UI 관련 (기존 호환성 유지)
│   │   ├── app.go                 # 기존 앱 (레거시)
│   │   ├── app_new.go             # 새로운 구조의 앱
│   │   ├── components/            # 재사용 가능한 UI 컴포넌트
│   │   │   └── executor.go        # 명령 실행기
│   │   └── dialog/                # 다이얼로그 컴포넌트
│   │       ├── confirm.go         # 확인 다이얼로그
│   │       └── settings.go        # 설정 다이얼로그
│   │
│   ├── config/                     # 설정 관리
│   │   └── config.go
│   │
│   ├── terraform/                  # Terraform 매니저 (레거시)
│   │   └── manager.go
│   │
│   └── git/                        # Git 매니저 (레거시)
│       └── manager.go
│
├── README.md
├── QUICKSTART.md
├── STRUCTURE.md                    # 기존 구조 문서
├── ARCHITECTURE.md                 # 이 파일 - 새로운 아키텍처 문서
├── TODO.md
├── go.mod
├── go.sum
└── install.sh
```

## 계층 구조 (k9s 스타일)

### 1. Model Layer (internal/model/)
- **목적**: 순수한 데이터 구조 정의
- **특징**: 
  - 비즈니스 로직 없음
  - 다른 패키지에 대한 의존성 최소화
  - 재사용 가능한 타입 정의

**파일들**:
- `terraform.go`: TerraformDirectory, TerraformStatus, HelmRelease 등
- `git.go`: GitStatus 등

### 2. DAO Layer (internal/dao/)
- **목적**: 데이터 접근 및 외부 시스템과의 상호작용
- **특징**:
  - Terraform CLI 실행
  - Git 명령 실행
  - 파일 시스템 접근
  - 실제 데이터 조작

**파일들**:
- `terraform.go`: TerraformDAO - Terraform 관련 모든 데이터 작업
- `git.go`: GitDAO - Git 관련 모든 데이터 작업

### 3. View Layer (internal/view/)
- **목적**: UI 컴포넌트 정의
- **특징**:
  - tview 위젯 래핑
  - 재사용 가능한 UI 컴포넌트
  - 각 뷰는 독립적으로 동작
  - 표시 로직만 포함

**파일들**:
- `tree_view.go`: TreeView - 파일 트리
- `content_view.go`: ContentView - 메인 컨텐츠 영역
- `header_view.go`: HeaderView - 애플리케이션 헤더
- `status_bar.go`: StatusBar - 하단 상태바

### 4. UI Layer (internal/ui/)
- **목적**: 애플리케이션 조합 및 이벤트 핸들링
- **특징**:
  - View 컴포넌트들을 조합
  - 키보드 이벤트 처리
  - 애플리케이션 상태 관리

**구성**:
- `app_new.go`: 새로운 구조의 메인 애플리케이션
- `components/`: 복잡한 UI 로직을 가진 컴포넌트
  - `executor.go`: 명령 실행 로직
- `dialog/`: 다이얼로그 컴포넌트
  - `confirm.go`: 확인 다이얼로그
  - `settings.go`: 설정 다이얼로그

### 5. Config Layer (internal/config/)
- **목적**: 설정 파일 관리
- **특징**:
  - YAML 설정 로드/저장
  - 기본 설정 생성

## 데이터 흐름

```
사용자 입력
    ↓
UI Layer (app_new.go)
    ↓
Components (executor.go)
    ↓
DAO Layer (terraform.go, git.go)
    ↓
외부 시스템 (terraform CLI, git CLI)
    ↓
Model Layer (terraform.go, git.go)
    ↓
View Layer (content_view.go 등)
    ↓
화면 표시
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

## 마이그레이션 계획

1. ✅ **Phase 1**: 새로운 구조 생성
   - model, dao, view 패키지 생성
   - app_new.go 작성

2. **Phase 2**: 기능 추가
   - Terraform 디렉토리 리스트 뷰
   - State 정보 뷰
   - Drift 감지 뷰

3. **Phase 3**: 레거시 제거
   - app.go → app_new.go로 완전 전환
   - terraform/manager.go, git/manager.go 제거
   - app_new.go → app.go로 리네임

## 장점

1. **관심사 분리**: 각 레이어가 명확한 책임을 가짐
2. **테스트 용이성**: 각 레이어를 독립적으로 테스트 가능
3. **재사용성**: View 컴포넌트를 다른 곳에서도 재사용 가능
4. **확장성**: 새로운 뷰나 기능을 쉽게 추가 가능
5. **유지보수성**: 코드 위치를 예측 가능하고 찾기 쉬움

## 사용 예시

### 새로운 View 추가
```go
// internal/view/terraform_list_view.go
package view

type TerraformListView struct {
    *tview.Table
}

func NewTerraformListView() *TerraformListView {
    // 구현
}
```

### 새로운 DAO 기능 추가
```go
// internal/dao/terraform.go
func (d *TerraformDAO) GetWorkspace(dir string) (string, error) {
    // 구현
}
```

### 새로운 Model 타입 추가
```go
// internal/model/terraform.go
type Workspace struct {
    Name    string
    Current bool
}
```

