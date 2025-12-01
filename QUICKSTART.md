# T9s 빠른 시작 가이드

## 1. 설치

```bash
cd /Users/idongju/dev/T9s
./install.sh
```

또는 수동으로:

```bash
go build -o t9s ./cmd/t9s
sudo mv t9s /usr/local/bin/
```

## 2. 설정

첫 실행 시 `~/.t9s/config.yaml` 파일이 자동으로 생성됩니다.

```yaml
# Terraform 루트 디렉토리 (서브 디렉토리에 s3, eks 등이 있는 경로)
terraform_root: /path/to/your/terraform

# S3 Backend 설정
backend:
  bucket: terraform-state
  region: ap-northeast-2
  
# 기본 설정
defaults:
  auto_refresh: true
  refresh_interval: 60
```

**중요**: `terraform_root`를 실제 Terraform 디렉토리 경로로 변경하세요!

## 3. 실행

```bash
t9s
```

## 4. 키보드 단축키

### 네비게이션
- `↑/↓` - 디렉토리 목록 탐색
- `Enter` - 디렉토리 선택
- `Tab` - 패널 간 이동

### Terraform 작업
- `p` - Terraform Plan 실행
- `a` - Terraform Apply 실행 (확인 필요)
- `i` - Terraform Init 실행

### 정보 조회
- `s` - State 정보 보기
- `g` - Git diff 보기
- `h` - Helm 차트 목록 보기
- `l` - 로그 보기

### 편집
- `e` - tfvars 파일 편집

### 기타
- `r` - 새로고침
- `?` - 도움말
- `q` - 종료
- `Ctrl+C` - 강제 종료

## 5. 예상 디렉토리 구조

T9s는 다음과 같은 구조를 권장합니다:

```
terraform/
├── s3-buckets/
│   ├── config/
│   │   ├── prod.tfvars
│   │   └── dev.tfvars
│   ├── main.tf
│   ├── variables.tf
│   └── backend.tf
├── eks-cluster/
│   ├── config/
│   │   └── prod.tfvars
│   ├── main.tf
│   └── ...
├── vpc-network/
│   └── ...
└── rds-databases/
    └── ...
```

## 6. 문제 해결

### Terraform이 인식되지 않는 경우
```bash
which terraform
# Terraform이 PATH에 있는지 확인
```

### Git 상태가 표시되지 않는 경우
해당 디렉토리가 Git 저장소인지 확인:
```bash
cd /path/to/terraform/directory
git status
```

### State를 읽을 수 없는 경우
Terraform init 실행:
```bash
cd /path/to/terraform/directory
terraform init
```

## 7. 개발 모드

소스에서 직접 실행하려면:

```bash
cd /Users/idongju/dev/T9s
go run ./cmd/t9s
```

## 8. 다음 단계

- [ ] 실제 Terraform 디렉토리 경로로 설정 업데이트
- [ ] 각 디렉토리에서 `terraform init` 실행
- [ ] T9s 실행 및 테스트
- [ ] 키보드 단축키 익히기

## 9. 버전 확인

```bash
t9s --version
```

---

**도움이 필요하신가요?**
- GitHub Issues: https://github.com/idongju/t9s/issues
- 문서: README.md 참조
