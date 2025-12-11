# T9s - TODO & Roadmap

## ✅ 완료된 기능 (v0.1.0)

- [x] 기본 프로젝트 구조 설정
- [x] tview 기반 TUI 인터페이스
- [x] 디렉토리 목록 표시 (Tree View)
- [x] 헤더/푸터 UI
- [x] 키보드 단축키 (q: 종료, Ctrl+C: 강제종료)
- [x] Terraform 디렉토리 스캔 로직
- [x] Git 상태 확인 매니저
- [x] 설정 파일 관리 (YAML)
- [x] 빌드 및 설치 스크립트

## ✅ 완료된 기능 (v0.2.0)

- [x] k9s 스타일 아키텍처 적용 (Model/DAO/View)
- [x] UI 컴포넌트 분리 (Header, Tree, Content, StatusBar)
- [x] 다이얼로그 시스템 (Confirm, Settings)
- [x] Terraform Plan/Apply 실행
- [x] 파일 편집 기능 ($EDITOR 연동)

## ✅ 완료된 기능 (v0.3.0)

### Terraform 통합
- [x] Terraform Init 실행 (설정 파일 선택)
- [x] Terraform Plan 실행 (tfvars 선택)
- [x] Terraform Apply 실행 (tfvars 선택)
- [x] Terraform Destroy 실행 (tfvars 선택)
- [x] 명령어 템플릿 커스터마이징 (Init/Plan/Apply/Destroy)
- [x] 파일 선택 다이얼로그 (미리보기 지원)
- [x] Execute/Auto Approve 버튼 분리
- [x] 실시간 스트리밍 출력
- [x] Yes/No 확인 다이얼로그 자동 팝업

### 히스토리
- [x] SQLite 기반 히스토리 DB
- [x] Apply/Destroy 실행 이력 저장
- [x] 사용자, 브랜치, tfvars 내용 기록
- [x] 히스토리 뷰 (h 키)
- [x] 상세 내용 토글 (Shift+M)
- [x] 페이지네이션 (d/u)

### Git 통합
- [x] 현재 브랜치 헤더에 표시
- [x] 브랜치 전환 (Shift+B)
- [x] Dirty 상태 처리 (Stash/Commit/Force)
- [x] 브랜치 상태 아이콘 (● 더티, ✓ 클린)

### UI/UX
- [x] Help View (? 또는 Shift+H)
- [x] Command Mode (/)
- [x] Home 화면 (Shift+C)
- [x] ANSI 컬러 지원
- [x] 빠른 스크롤 (u/d, Shift+방향키)
- [x] Apply/Destroy 시 자동 포커스 전환

## 🚧 진행 예정 (v0.4.0)

### Terraform
- [ ] Workspace 전환 UI
- [ ] State 정보 테이블 뷰
- [ ] Drift 감지 (terraform plan -detailed-exitcode)
- [ ] Plan 결과 저장 및 비교

### Git
- [ ] Git Diff 뷰어
- [ ] 커밋 이력 보기
- [ ] 수정된 파일 목록

### UI/UX
- [ ] 필터 및 검색 기능
- [ ] 즐겨찾기 디렉토리
- [ ] 작업 로그 저장

## 📋 계획된 기능 (v0.5.0+)

### 고급 기능
- [ ] 여러 환경 (dev/stage/prod) 지원
- [ ] 환경 간 비교 기능
- [ ] Terraform 모듈 의존성 그래프
- [ ] 리소스 비용 추정

### UI/UX 개선
- [ ] 테마 지원 (다크/라이트 모드)
- [ ] 커스텀 색상 설정
- [ ] 레이아웃 커스터마이징

### 통합
- [ ] AWS CLI 통합
- [ ] kubectl 통합 (EKS 관련)
- [ ] Slack/Discord 알림

## 🐛 알려진 이슈

- [ ] 터미널 환경이 없으면 실행 불가 (정상 동작)

## 💡 아이디어 목록

### 편의성
- [ ] 자주 사용하는 명령어 매크로
- [ ] 디렉토리별 커스텀 스크립트 실행
- [ ] Terraform 문법 검사

### 보안
- [ ] tfvars 파일 암호화 지원
- [ ] 민감한 변수 마스킹
- [ ] 작업 승인 워크플로우

### 모니터링
- [ ] 리소스 변경 이력
- [ ] 비용 추세 분석
- [ ] 성능 메트릭

## 🎯 버전 목표

| 버전 | 목표 | 상태 |
|------|------|------|
| v0.1.0 | 기본 UI 및 구조 | ✅ 완료 |
| v0.2.0 | k9s 스타일 아키텍처 | ✅ 완료 |
| v0.3.0 | 실용적인 Terraform 관리 | ✅ 완료 |
| v0.4.0 | Workspace 및 State 관리 | 📝 진행 예정 |
| v0.5.0 | 고급 기능 및 통합 | 📝 계획 |
| v1.0.0 | 프로덕션 준비 | 🔮 미래 |

---

**최종 업데이트**: 2025-12-11
