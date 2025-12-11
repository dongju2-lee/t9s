# T9s 빠른 시작 가이드

## 1. 설치

```bash
cd /path/to/T9s
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
# Terraform 루트 디렉토리
terraform_root: /path/to/your/terraform

# 명령어 템플릿
commands:
  init_template: "terraform init -backend-config={initconf}"
  plan_template: "terraform plan -var-file={varfile}"
  apply_template: "terraform apply -var-file={varfile}"
  destroy_template: "terraform destroy -var-file={varfile}"
  tfvars_file: "config/env.tfvars"
  init_conf_file: "config/env.conf"
```

**중요**: `terraform_root`를 실제 Terraform 디렉토리 경로로 변경하세요!

## 3. 실행

```bash
t9s
```

## 4. 핵심 단축키

### 기본 네비게이션
| 키 | 설명 |
|---|---|
| `Tab` | 포커스 전환 (Tree ↔ Content) |
| `↑/↓` | 탐색 |
| `Enter` | 선택/확장 |
| `Esc` | 뒤로 가기 |
| `q` | 종료 |

### Terraform 작업
| 키 | 설명 |
|---|---|
| `i` | Terraform Init |
| `p` | Terraform Plan |
| `a` | Terraform Apply |
| `d` | Terraform Destroy |
| `h` | History (이력 조회) |

### 유틸리티
| 키 | 설명 |
|---|---|
| `e` | 파일 편집 ($EDITOR) |
| `s` | 설정 열기 |
| `Shift+B` | Git 브랜치 전환 |
| `?` 또는 `Shift+H` | 도움말 |
| `/` | 커맨드 모드 |
| `Shift+C` | 홈 화면 |

### Content View (스크롤)
| 키 | 설명 |
|---|---|
| `u` / `d` | 빠른 스크롤 (10줄) |
| `Shift+↑/↓` | 빠른 스크롤 (10줄) |

### History View
| 키 | 설명 |
|---|---|
| `d` | 더 보기 |
| `u` | 접기 |
| `Shift+M` | 상세 내용 토글 |

## 5. Terraform 실행 흐름

1. **디렉토리 선택**: Tree에서 Terraform 디렉토리로 이동
2. **명령 실행**: `a` (Apply) 또는 다른 단축키 입력
3. **파일 선택**: tfvars 파일 선택 다이얼로그에서 선택
4. **확인 다이얼로그**: 
   - **Execute**: Plan 결과를 보고 Yes/No 선택
   - **Auto Approve**: 즉시 실행 (-auto-approve)
   - **Cancel**: 취소
5. **결과 확인**: 실시간으로 로그 확인

## 6. 히스토리 기능

Apply/Destroy 실행 이력이 자동으로 저장됩니다:
- **저장 위치**: `~/.t9s/history.db` (SQLite)
- **저장 정보**: 디렉토리, 액션, 시간, 사용자, 브랜치, tfvars 내용
- **조회**: `h` 키로 현재 디렉토리의 이력 확인

## 7. Git 브랜치 전환

`Shift+B`를 누르면:
1. 로컬 브랜치 목록 표시
2. 브랜치 선택
3. 로컬 변경사항이 있으면:
   - **Stash**: 변경사항 임시 저장
   - **Commit**: 변경사항 커밋
   - **Force**: 변경사항 버리고 전환

## 8. 예상 디렉토리 구조

```
terraform/
├── s3-buckets/
│   ├── config/
│   │   ├── prod.tfvars
│   │   ├── dev.tfvars
│   │   └── env.conf
│   ├── main.tf
│   ├── variables.tf
│   └── backend.tf
├── eks-cluster/
│   ├── config/
│   │   └── prod.tfvars
│   ├── main.tf
│   └── ...
└── vpc-network/
    └── ...
```

## 9. 문제 해결

### Terraform이 인식되지 않는 경우
```bash
which terraform
# Terraform이 PATH에 있는지 확인
```

### 히스토리가 저장되지 않는 경우
```bash
# SQLite 파일 확인
ls -la ~/.t9s/history.db
```

### 설정이 적용되지 않는 경우
```bash
# 설정 파일 확인
cat ~/.t9s/config.yaml
```

## 10. 버전 확인

```bash
t9s --version
# T9s version 0.3.0
```

---

**도움이 필요하신가요?**
- 앱 내에서 `?` 또는 `Shift+H`를 눌러 도움말 확인
- README.md 참조
