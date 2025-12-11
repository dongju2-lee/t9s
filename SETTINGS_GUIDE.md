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

---

### 2. Terraform Init Template
```
Terraform Init Template: terraform init -backend-config={initconf}
```

**설명**: `i` 키를 눌렀을 때 실행될 Terraform Init 명령어 템플릿입니다.

**템플릿 변수**:
- `{initconf}`: 선택된 .conf 파일 경로로 자동 치환됩니다

**예시**:
```bash
# 기본
terraform init -backend-config={initconf}

# 재초기화
terraform init -reconfigure -backend-config={initconf}

# 마이그레이션
terraform init -migrate-state -backend-config={initconf}
```

---

### 3. Terraform Plan Template
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

---

### 4. Terraform Apply Template
```
Terraform Apply Template: terraform apply -var-file={varfile}
```

**설명**: `a` 키를 눌렀을 때 실행될 Terraform Apply 명령어 템플릿입니다.

**템플릿 변수**:
- `{varfile}`: 선택된 .tfvars 파일 경로로 자동 치환됩니다

**예시**:
```bash
# 기본 (Execute 버튼 - Yes/No 확인)
terraform apply -var-file={varfile}

# Auto Approve 버튼 선택 시 자동으로 -auto-approve 추가
terraform apply -var-file={varfile} -auto-approve
```

---

### 5. Terraform Destroy Template
```
Terraform Destroy Template: terraform destroy -var-file={varfile}
```

**설명**: `d` 키를 눌렀을 때 실행될 Terraform Destroy 명령어 템플릿입니다.

**템플릿 변수**:
- `{varfile}`: 선택된 .tfvars 파일 경로로 자동 치환됩니다

---

### 6. Default tfvars File
```
Default tfvars File: config/env.tfvars
```

**설명**: 파일 선택 다이얼로그에서 기본으로 선택될 .tfvars 파일의 상대 경로입니다.

**예시**:
```bash
config/prod.tfvars
config/dev.tfvars
vars/production.tfvars
```

---

### 7. Init Config File
```
Init Config File: config/env.conf
```

**설명**: Terraform Init 시 기본으로 선택될 .conf 파일의 상대 경로입니다.

**예시**:
```bash
config/backend.conf
config/env.conf
```

---

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
  init_template: terraform init -backend-config={initconf}
  plan_template: terraform plan -var-file={varfile}
  apply_template: terraform apply -var-file={varfile}
  destroy_template: terraform destroy -var-file={varfile}
  tfvars_file: config/env.tfvars
  init_conf_file: config/env.conf
```

---

## 사용 예시

### 예시 1: 개인 프로젝트 설정
```yaml
terraform_root: /Users/dongju/dev/my-terraform
commands:
  init_template: terraform init -backend-config={initconf}
  plan_template: terraform plan -var-file={varfile}
  apply_template: terraform apply -var-file={varfile}
  destroy_template: terraform destroy -var-file={varfile}
  tfvars_file: environments/production.tfvars
  init_conf_file: config/backend.conf
```

### 예시 2: 회사 프로젝트 설정
```yaml
terraform_root: /home/user/company/infrastructure
commands:
  init_template: terraform init -reconfigure -backend-config={initconf}
  plan_template: terraform plan -var-file={varfile} -out=tfplan
  apply_template: terraform apply -var-file={varfile}
  destroy_template: terraform destroy -var-file={varfile}
  tfvars_file: config/prod.tfvars
  init_conf_file: config/prod.conf
```

---

## 팁

### 1. 절대 경로 사용
```
✅ 권장: /Users/dongju/dev/terraform
❌ 비권장: ../terraform (상대 경로는 작동하지 않을 수 있음)
```

### 2. 환경별 설정 변경
개발/프로덕션 환경을 전환할 때:
1. `s` 키로 설정 열기
2. Default tfvars File만 변경 (`config/dev.tfvars` ↔ `config/prod.tfvars`)
3. Save

### 3. Execute vs Auto Approve
- **Execute**: Terraform이 Plan 결과를 보여주고 Yes/No 다이얼로그로 확인
- **Auto Approve**: `-auto-approve` 플래그로 즉시 실행

---

## 트러블슈팅

### 설정이 저장되지 않음
```bash
# 설정 디렉토리 권한 확인
ls -la ~/.t9s/

# 권한 수정 (필요시)
chmod 755 ~/.t9s
chmod 644 ~/.t9s/config.yaml
```

### 히스토리가 저장되지 않음
```bash
# SQLite 파일 확인
ls -la ~/.t9s/history.db
```

---

## 단축키

설정 화면에서:
- `Tab`: 다음 필드로 이동
- `Shift+Tab`: 이전 필드로 이동
- `Enter`: 버튼 클릭
- `Esc` 또는 Cancel: 설정 창 닫기

---

**참고**: 설정 변경 후 반드시 **Save** 버튼을 눌러야 저장됩니다!
