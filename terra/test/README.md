# T9s Test Environment

This directory contains test Terraform configurations for testing T9s features without actual cloud infrastructure provisioning.

## Test Directories

### 1. 📊 Monitoring (`monitoring/`)
Simulates a monitoring system setup:
- Creates monitoring configuration files (JSON)
- Generates HTML dashboards
- Creates alert rules (YAML)
- Logs setup activities

**What it does:**
- ✅ Local file creation only
- ✅ Echo commands for logging
- ❌ NO cloud resources
- ❌ NO installations

### 2. 🚀 WebApp (`webapp/`)
Simulates web application configuration:
- Creates app info files (TXT)
- Generates JSON configurations
- Creates logging configs (conditional)
- Simple echo commands

**What it does:**
- ✅ Super simple text files
- ✅ Echo and mkdir only
- ❌ NO Docker/services
- ❌ NO installations

## Quick Start

### Option 1: Use T9s (Recommended)
```bash
# From project root
t9s
```

Then:
1. Navigate to `monitoring` or `webapp` folder
2. Press `i` - Init
3. Press `p` - Plan
4. Press `a` - Apply
5. Press `Tab` - Switch to content view to scroll
6. Press `h` - View history

### Option 2: Command Line
```bash
# Monitoring
cd monitoring
terraform init
terraform apply -var-file=config/env.tfvars

# WebApp
cd webapp
terraform init
terraform apply -var-file=config/dev.tfvars
```

## Testing Scenarios

### Scenario 1: Multiple Environments
Test with different configurations:
- `monitoring/config/env.tfvars` - Development
- `monitoring/config/prod.tfvars` - Production
- `monitoring/config/staging.tfvars` - Staging (no alerts)

### Scenario 2: Tab Focus Switching
1. Apply to generate output
2. Press `Tab` to switch to Content View
3. Use arrow keys to scroll through output
4. Press `Tab` again to return to File Tree

### Scenario 3: History Tracking
1. Apply with different tfvars files
2. Press `h` to view history
3. Compare different configurations

### Scenario 4: Conditional Resources
- `monitoring` with staging.tfvars - 3 resources (no alerts)
- `monitoring` with env.tfvars - 4 resources (with alerts)
- `webapp` with staging.tfvars - 2 resources (no logging)
- `webapp` with dev.tfvars - 3 resources (with logging)

## What Gets Created

### Monitoring Output
```
monitoring/output/
├── monitoring-{env}.json
├── dashboard-{env}.html
├── alerts-{env}.yaml (if enabled)
└── logs/setup.log
```

### WebApp Output
```
webapp/output/
├── app-{env}.txt
├── config-{env}.json
├── logging-{env}.conf (if enabled)
└── logs/setup.log
```

## Safety Guarantees

✅ **No Cloud Access** - Everything runs locally  
✅ **No Installations** - Only creates text files  
✅ **No Network Calls** - No external dependencies  
✅ **No Credentials Needed** - No AWS/GCP/Azure required  
✅ **Fast Execution** - Instant apply/destroy  
✅ **Safe to Delete** - Just remove output/ directories  

## Clean Up

```bash
# Clean monitoring
cd monitoring
rm -rf output/ .terraform/ terraform.tfstate*

# Clean webapp
cd webapp
rm -rf output/ .terraform/ terraform.tfstate*
```

## Why This is Perfect for Testing

1. **Realistic Workflow** - Same terraform commands as real infrastructure
2. **Visible Results** - Created files you can inspect
3. **Fast Feedback** - No waiting for cloud resources
4. **History Tracking** - T9s can track all applies/destroys
5. **Multiple Scenarios** - Different environments and configurations
6. **Tab Navigation** - Test focus switching with real output
7. **Completely Safe** - No risk of accidental cloud costs

Enjoy testing T9s! 🎉

