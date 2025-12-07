# T9s 설정 가이드

T9s 설정을 통해 작업 환경을 커스터마이징할 수 있습니다.

## 설정 열기

앱 실행 중 `s` 키를 누르면 설정 화면이 열립니다.

## 설정 항목

### 1. Terraform Root Directory
```
Terraform Root Directory: /Users/user/terraform
```

**설명**: Terraform 코드가 있는 루트 디렉토리 경로입니다.

**예시**:
- `/Users/dongju/dev/terraform`
- `/home/user/infrastructure/terraform`
- `~/projects/terraform-aws`

**동작**:
- 앱 시작 시 이 디렉토리를 기준으로 파일 트리를 표시합니다
- 설정 저장 시 디렉토리가 변경되면 UI가 자동으로 리빌드됩니다
- 각 사용자의 프로젝트 위치에 맞게 설정 가능

### 2. Terraform Plan Template
```
Terraform Plan Template: terraform plan -var-file={varfile}
```

**설명**: `p` 키를 눌렀을 때 실행될 Terraform Plan 명령어 템플릿입니다.

**템플릿 변수**:
- `{varfile}`: 선택된 .tfvars 파일 경로로 자동 치환됩니다

**예시**:
```bash
# 기본
terraform plan -var-file={varfile}

# Out 파일 지정
terraform plan -var-file={varfile} -out=tfplan

# 추가 옵션
terraform plan -var-file={varfile} -parallelism=10
```

### 3. Terraform Apply Template
```
Terraform Apply Template: terraform apply -var-file={varfile}
```

**설명**: `a` 키를 눌렀을 때 실행될 Terraform Apply 명령어 템플릿입니다.

**템플릿 변수**:
- `{varfile}`: 선택된 .tfvars 파일 경로로 자동 치환됩니다

**예시**:
```bash
# 기본
terraform apply -var-file={varfile}

# 자동 승인
terraform apply -var-file={varfile} -auto-approve

# Plan 파일 사용
terraform apply tfplan
```

### 4. Default Var File
```
Default Var File: config/prod.tfvars
```

**설명**: 디렉토리 선택 시 기본으로 사용할 .tfvars 파일의 상대 경로입니다.

**예시**:
```bash
# 프로덕션 환경
config/prod.tfvars

# 개발 환경
config/dev.tfvars

# 커스텀 경로
vars/production.tfvars
```

## 설정 파일 위치

```
~/.t9s/config.yaml
```

설정은 YAML 형식으로 저장됩니다:

```yaml
terraform_root: /Users/dongju/dev/terraform
backend:
  bucket: terraform-state
  region: ap-northeast-2
  prefix: ""
defaults:
  auto_refresh: true
  refresh_interval: 60
commands:
  plan_template: terraform plan -var-file={varfile}
  apply_template: terraform apply -var-file={varfile}
  var_file: config/prod.tfvars
```

## 사용 예시

### 예시 1: 개인 프로젝트 설정
```
Terraform Root Directory: /Users/dongju/dev/my-terraform
Terraform Plan Template: terraform plan -var-file={varfile}
Terraform Apply Template: terraform apply -var-file={varfile} -auto-approve
Default Var File: environments/production.tfvars
```

### 예시 2: 회사 프로젝트 설정
```
Terraform Root Directory: /home/user/company/infrastructure
Terraform Plan Template: terraform plan -var-file={varfile} -out=tfplan
Terraform Apply Template: terraform apply tfplan
Default Var File: config/prod.tfvars
```

### 예시 3: 멀티 환경 설정
```
Terraform Root Directory: ~/projects/multi-env-terraform
Terraform Plan Template: terraform plan -var-file={varfile} -parallelism=20
Terraform Apply Template: terraform apply -var-file={varfile}
Default Var File: vars/staging.tfvars
```

## 설정 UI 화면

```
╔═══════════════════════════════════════════════════════════════════╗
║  T9s - Settings                                                    ║
╚═══════════════════════════════════════════════════════════════════╝

┌─ ⚙️  Settings ─────────────────────────────────────────────────────┐
│                                                                     │
│  Terraform Root Directory: /Users/dongju/dev/terraform____________ │
│                                                                     │
│  Terraform Plan Template: terraform plan -var-file={varfile}______ │
│                                                                     │
│  Terraform Apply Template: terraform apply -var-file={varfile}____ │
│                                                                     │
│  Default Var File: config/prod.tfvars____________________________ │
│                                                                     │
│                    [ Save ]  [ Cancel ]                             │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

Terraform Root Directory:
  Directory where your Terraform code is located (e.g., /home/user/terraform)

Template Variables:
  {varfile} - Will be replaced with the var file path

Examples:
  terraform plan -var-file={varfile}
  terraform apply -var-file={varfile} -auto-approve
```

## 팁

### 1. 절대 경로 사용
```
✅ 권장: /Users/dongju/dev/terraform
❌ 비권장: ../terraform (상대 경로는 작동하지 않을 수 있음)
```

### 2. 경로에 공백이 있는 경우
```
✅ 가능: /Users/dong ju/my terraform
💡 하지만 공백 없는 경로를 권장
```

### 3. 홈 디렉토리 축약
```
✅ 가능: ~/terraform
✅ 권장: /Users/dongju/terraform (전체 경로)
```

### 4. 환경별 설정 변경
개발/프로덕션 환경을 전환할 때:
1. `s` 키로 설정 열기
2. Default Var File만 변경 (`config/dev.tfvars` ↔ `config/prod.tfvars`)
3. Save

### 5. 프로젝트별 설정
여러 Terraform 프로젝트 작업 시:
1. 프로젝트 A 작업: Terraform Root Directory를 `/path/to/projectA`로 설정
2. 프로젝트 B 작업: Terraform Root Directory를 `/path/to/projectB`로 설정
3. 필요할 때마다 설정에서 변경

## 트러블슈팅

### 설정이 저장되지 않음
```bash
# 설정 디렉토리 권한 확인
ls -la ~/.t9s/

# 권한 수정 (필요시)
chmod 755 ~/.t9s
chmod 644 ~/.t9s/config.yaml
```

### 경로를 변경했는데 트리가 업데이트되지 않음
1. Save 버튼을 눌렀는지 확인
2. T9s를 재시작
3. 경로가 올바른지 확인 (`ls /path/to/directory`)

### 디렉토리가 존재하지 않음
```bash
# 디렉토리 생성
mkdir -p /path/to/terraform

# 또는 기존 디렉토리 경로로 설정 변경
```

## 단축키

설정 화면에서:
- `Tab`: 다음 필드로 이동
- `Shift+Tab`: 이전 필드로 이동
- `Enter`: 버튼 클릭
- `Esc` 또는 Cancel: 설정 창 닫기

---

**참고**: 설정 변경 후 반드시 **Save** 버튼을 눌러야 저장됩니다!


