# T9s 프로젝트 구조

```
T9s/
├── cmd/
│   └── t9s/
│       └── main.go              # CLI 진입점
│
├── internal/
│   ├── ui/
│   │   └── app.go               # tview 기반 UI 컴포넌트
│   │
│   ├── terraform/
│   │   └── manager.go           # Terraform 작업 관리
│   │                            # - 디렉토리 스캔
│   │                            # - Plan/Apply 실행
│   │                            # - State 조회
│   │                            # - Helm 릴리스 추출
│   │
│   ├── git/
│   │   └── manager.go           # Git 통합
│   │                            # - 상태 확인
│   │                            # - Diff 조회
│   │                            # - 커밋 정보
│   │
│   └── config/
│       └── config.go            # 설정 관리
│                                # - YAML 로드/저장
│                                # - 기본 설정 생성
│
├── README.md                    # 프로젝트 설명
├── QUICKSTART.md                # 빠른 시작 가이드
├── TODO.md                      # 로드맵 및 TODO
├── install.sh                   # 설치 스크립트
├── go.mod                       # Go 모듈 정의
├── go.sum                       # 의존성 체크섬
├── .gitignore                   # Git 제외 파일
└── t9s                          # 빌드된 바이너리 (gitignore됨)
```

## 파일 설명

### cmd/t9s/main.go
- 애플리케이션 진입점
- 버전 플래그 처리
- UI 앱 초기화 및 실행

### internal/ui/app.go
- tview 기반 터미널 UI
- 헤더, 디렉토리 목록, 상세 패널, 푸터
- 키보드 단축키 처리
- k9s 스타일 인터페이스

### internal/terraform/manager.go
- Terraform 디렉토리 스캔 및 관리
- Plan/Apply 실행
- State 정보 조회
- Backend 설정 파싱
- Helm 릴리스 추출 (state에서)
- Drift 감지

### internal/git/manager.go
- Git 저장소 상태 확인
- 브랜치, 수정 파일 정보
- Diff 조회
- 커밋 이력

### internal/config/config.go
- ~/.t9s/config.yaml 관리
- Terraform root 경로
- S3 backend 설정
- 기본 옵션 (자동 새로고침 등)

## 의존성

```go
require (
    github.com/rivo/tview v0.42.0          // 터미널 UI 프레임워크
    github.com/gdamore/tcell/v2 v2.8.1     // 터미널 제어
    gopkg.in/yaml.v3 v3.0.1                // YAML 파싱
)
```

## 빌드 결과

- **바이너리 크기**: ~3.8MB
- **플랫폼**: macOS (현재)
- **설치 위치**: `/usr/local/bin/t9s`

## 설정 파일 위치

- **설정**: `~/.t9s/config.yaml`
- **로그** (예정): `~/.t9s/logs/`
- **캐시** (예정): `~/.t9s/cache/`

## 다음 단계

v0.2.0에서는 UI와 실제 Terraform/Git 매니저를 연결하여:
1. 설정 파일에서 terraform_root 읽기
2. 실제 디렉토리 스캔
3. Git 상태 표시
4. State 정보 조회
5. 실시간 drift 감지

구현 예정입니다!
