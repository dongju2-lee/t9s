# WebApp Configuration Test

Simple Terraform configuration for testing T9s - **NO installation, just text files!**

## What it Does

Creates simple text files:
- `app-{env}.txt` - App information
- `config-{env}.json` - JSON configuration
- `logging-{env}.conf` - Logging config (if enabled)
- `logs/setup.log` - Setup log

**Only uses:**
- ✅ `echo` commands
- ✅ `mkdir` for directories
- ✅ Text file creation
- ❌ NO installations
- ❌ NO downloads
- ❌ NO Docker/services

## Quick Test

```bash
terraform init
terraform apply -var-file=config/dev.tfvars
```

## Configurations

### dev.tfvars
- Version: 1.0.0
- Logging: ✅ Enabled

### prod.tfvars
- Version: 2.0.0
- Logging: ✅ Enabled

### staging.tfvars
- Version: 1.5.0
- Logging: ❌ Disabled (only 2 files created)

## Testing in T9s

1. Navigate to `webapp` folder
2. Press `i` - Init
3. Select a `.tfvars` file
4. Press `p` - Plan
5. Press `a` - Apply
6. Press `Tab` - Scroll output
7. Press `h` - View history

## View Results

```bash
ls -la output/
cat output/app-dev.txt
cat output/config-dev.json
```

## Clean Up

```bash
rm -rf output/ .terraform/ terraform.tfstate*
```

Super simple, safe, and perfect for testing! 🎯
