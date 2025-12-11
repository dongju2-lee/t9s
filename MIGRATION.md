# T9s 마이그레이션 가이드

기존 구조에서 k9s 스타일의 새로운 구조로 마이그레이션하는 가이드입니다.

## 변경 사항

### 디렉토리 구조

**이전 (v0.1.0)**:
```
internal/
├── ui/
│   └── app.go           # 모든 UI 로직이 하나의 파일에
├── terraform/
│   └── manager.go       # Terraform 로직
├── git/
│   └── manager.go       # Git 로직
└── config/
    └── config.go        # 설정
```

**이후 (v0.2.0)**:
```
internal/
├── model/               # 새로 추가: 데이터 모델
│   ├── terraform.go
│   └── git.go
├── dao/                 # 새로 추가: 데이터 접근 계층
│   ├── terraform.go
│   └── git.go
├── view/                # 새로 추가: UI 뷰 컴포넌트
│   ├── tree_view.go
│   ├── content_view.go
│   ├── header_view.go
│   └── status_bar.go
├── ui/
│   ├── app.go           # 레거시 (호환성 유지)
│   ├── app_new.go       # 새로운 구조
│   ├── components/      # 새로 추가: UI 컴포넌트
│   │   └── executor.go
│   └── dialog/          # 새로 추가: 다이얼로그
│       ├── confirm.go
│       └── settings.go
├── terraform/           # 레거시 (추후 제거 예정)
│   └── manager.go
├── git/                 # 레거시 (추후 제거 예정)
│   └── manager.go
└── config/
    └── config.go
```

## 사용 방법

### 기존 버전 사용 (v0.1.0)
```go
// cmd/t9s/main.go
app := ui.NewApp()  // 기존 앱
app.Run()
```

### 새로운 버전 사용 (v0.2.0)
```go
// cmd/t9s/main_new.go
app := ui.NewAppNew()  // 새로운 구조의 앱
app.Run()
```

## 새로운 구조의 장점

### 1. 관심사 분리 (Separation of Concerns)
- **Model**: 데이터 구조만 정의
- **DAO**: 데이터 접근 및 외부 시스템 호출
- **View**: UI 표시 로직만
- **UI**: 컴포넌트 조합 및 이벤트 처리

### 2. 재사용성
```go
// 다른 곳에서도 TreeView를 재사용 가능
treeView := view.NewTreeView("/path/to/dir")
treeView.SetFileSelectHandler(myHandler)
```

### 3. 테스트 용이성
```go
// DAO를 독립적으로 테스트
dao := dao.NewTerraformDAO("/test/path")
dirs, err := dao.ListDirectories()
assert.NoError(t, err)
```

### 4. 확장성
```go
// 새로운 뷰를 쉽게 추가
type TerraformListView struct {
    *tview.Table
}

func NewTerraformListView() *TerraformListView {
    // 구현
}
```

## 코드 비교

### 파일 표시 (기존)
```go
// app.go - 650+ 줄의 단일 파일에서
func (a *App) displayFile(path string) {
    a.currentFile = path
    content, err := ioutil.ReadFile(path)
    if err != nil {
        a.contentView.Clear()
        fmt.Fprintf(a.contentView, "[red]Error reading file: %v[white]", err)
        return
    }
    a.contentView.Clear()
    a.contentView.SetTitle(fmt.Sprintf(" 📄 %s ", filepath.Base(path)))
    fmt.Fprintf(a.contentView, "[yellow]File:[white] %s\n", path)
    fmt.Fprintf(a.contentView, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))
    fmt.Fprintf(a.contentView, "%s", string(content))
}
```

### 파일 표시 (새로운)
```go
// view/content_view.go - 명확한 책임 분리
func (cv *ContentView) DisplayFile(path string) error {
    content, err := ioutil.ReadFile(path)
    if err != nil {
        cv.Clear()
        fmt.Fprintf(cv, "[red]Error reading file: %v[white]", err)
        return err
    }

    cv.Clear()
    cv.SetTitle(fmt.Sprintf(" 📄 %s ", filepath.Base(path)))
    fmt.Fprintf(cv, "[yellow]File:[white] %s\n", path)
    fmt.Fprintf(cv, "[cyan]%s[white]\n\n", strings.Repeat("─", 60))
    fmt.Fprintf(cv, "%s", string(content))
    return nil
}

// app_new.go - 간결한 호출
a.treeView.SetFileSelectHandler(func(path string) {
    a.currentFile = path
    a.contentView.DisplayFile(path)
})
```

## 마이그레이션 단계

### Phase 1: ✅ 완료
- [x] model 패키지 생성
- [x] dao 패키지 생성
- [x] view 패키지 생성
- [x] ui/components 패키지 생성
- [x] ui/dialog 패키지 생성
- [x] app_new.go 작성
- [x] 빌드 테스트

### Phase 2: 진행 중
- [ ] main.go를 main_new.go 사용하도록 변경
- [ ] 새로운 기능 추가:
  - [ ] Terraform 디렉토리 리스트 뷰
  - [ ] State 정보 테이블 뷰
  - [ ] Drift 감지 실시간 표시
  - [ ] Workspace 전환 기능

### Phase 3: 추후 계획
- [ ] 레거시 코드 제거
  - [ ] terraform/manager.go → dao/terraform.go로 완전 이전
  - [ ] git/manager.go → dao/git.go로 완전 이전
  - [ ] app.go 제거
- [ ] app_new.go → app.go로 리네임
- [ ] main_new.go → main.go로 병합

## 빌드 및 실행

```bash
# 빌드
go build -o t9s ./cmd/t9s

# 실행 (기존 버전)
./t9s

# 실행 (새로운 버전) - main.go 수정 후
./t9s
```

## 테스트

```bash
# 모든 패키지 빌드 테스트
go build ./...

# 특정 패키지 테스트
go test ./internal/model
go test ./internal/dao
go test ./internal/view
```

## 주의사항

1. **레거시 호환성**: 기존 app.go는 당분간 유지되어 롤백 가능
2. **점진적 마이그레이션**: 새로운 기능부터 새 구조 사용
3. **문서 업데이트**: ARCHITECTURE.md 참고

## 참고 자료

- [k9s 프로젝트](https://github.com/derailed/k9s)
- ARCHITECTURE.md - 새로운 아키텍처 상세 설명
- STRUCTURE.md - 기존 구조 문서



