# T9s 마이그레이션 가이드

## v0.2.x → v0.3.0 마이그레이션

### 설정 파일 변경

**이전 (v0.2.x)**:
```yaml
commands:
  plan_template: terraform plan -var-file={varfile}
  apply_template: terraform apply -var-file={varfile}
  var_file: config/prod.tfvars
  backend_config: config/env.conf
```

**이후 (v0.3.0)**:
```yaml
commands:
  init_template: terraform init -backend-config={initconf}
  plan_template: terraform plan -var-file={varfile}
  apply_template: terraform apply -var-file={varfile}
  destroy_template: terraform destroy -var-file={varfile}
  tfvars_file: config/env.tfvars
  init_conf_file: config/env.conf
```

**변경사항**:
- `var_file` → `tfvars_file`
- `backend_config` → `init_conf_file`
- `init_template` 추가
- `destroy_template` 추가

### 새로운 데이터 파일

**v0.3.0에서 추가**:
```
~/.t9s/
├── config.yaml      # 설정 (기존)
└── history.db       # 히스토리 DB (신규)
```

---

## v0.1.x → v0.2.0 마이그레이션

### 디렉토리 구조

**이전 (v0.1.0)**:
```
internal/
├── ui/
│   └── app.go           # 모든 UI 로직이 하나의 파일에
├── terraform/
│   └── manager.go
├── git/
│   └── manager.go
└── config/
    └── config.go
```

**이후 (v0.2.0+)**:
```
internal/
├── model/               # 새로 추가: 데이터 모델
├── dao/                 # 새로 추가: 데이터 접근 계층
├── view/                # 새로 추가: UI 뷰 컴포넌트
├── db/                  # 새로 추가: 데이터베이스 (v0.3.0)
├── ui/
│   ├── app.go           # 레거시 (호환성 유지)
│   ├── app_new.go       # 새로운 구조
│   ├── components/      # UI 컴포넌트
│   └── dialog/          # 다이얼로그
├── terraform/
├── git/
└── config/
```

---

## 새로운 구조의 장점

### 1. 관심사 분리 (Separation of Concerns)
- **Model**: 데이터 구조만 정의
- **DAO**: 데이터 접근 및 외부 시스템 호출
- **View**: UI 표시 로직만
- **UI**: 컴포넌트 조합 및 이벤트 처리
- **DB**: 영구 데이터 저장

### 2. 재사용성
```go
// 다른 곳에서도 TreeView를 재사용 가능
treeView := view.NewTreeView("/path/to/dir")
treeView.SetFileSelectHandler(myHandler)
```

### 3. 테스트 용이성
```go
// DB를 독립적으로 테스트
db, _ := db.NewHistoryDB()
entry := &db.HistoryEntry{...}
db.AddEntry(entry)
```

### 4. 확장성
```go
// 새로운 뷰를 쉽게 추가
type WorkspaceView struct {
    *tview.List
}
```

---

## 마이그레이션 체크리스트

### v0.3.0 업그레이드
- [ ] 설정 파일 필드명 변경 확인
  - `var_file` → `tfvars_file`
  - `backend_config` → `init_conf_file`
- [ ] `init_template`, `destroy_template` 추가
- [ ] 히스토리 DB 자동 생성 확인 (`~/.t9s/history.db`)

### 일반
- [ ] `go mod download` 실행 (SQLite 드라이버 추가)
- [ ] `go build -o t9s ./cmd/t9s` 빌드
- [ ] `t9s --version`으로 버전 확인

---

## 빌드 및 실행

```bash
# 의존성 다운로드
go mod download

# 빌드
go build -o t9s ./cmd/t9s

# 실행
./t9s

# 또는 설치 스크립트 사용
./install.sh
```

---

## 참고 자료

| 문서 | 설명 |
|------|------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | 아키텍처 상세 설명 |
| [CHANGELOG.md](CHANGELOG.md) | 변경 이력 |
| [SETTINGS_GUIDE.md](SETTINGS_GUIDE.md) | 설정 가이드 |
| [README.md](README.md) | 프로젝트 개요 및 사용법 |
