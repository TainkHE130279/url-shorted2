# 🔧 Pre-commit Hooks Setup Guide

## 📋 Tổng quan

Dự án này đã được setup với **pre-commit hooks** để đảm bảo code quality trước khi commit. Hooks sẽ tự động chạy các checks sau:

- ✅ **Code formatting** (gofmt)
- ✅ **Code linting** (golangci-lint)
- ✅ **Unit tests** (go test)
- ✅ **Build check** (go build)

## 🚀 Setup cho Developer mới

### **Bước 1: Clone repository**
```bash
git clone <repository-url>
cd url_shorted2
```

### **Bước 2: Setup pre-commit hooks**
```bash
# Chạy script setup tự động
./scripts/setup-hooks.sh
```

Script này sẽ:
- Kiểm tra Go đã được cài đặt
- Cài đặt golangci-lint nếu chưa có
- Setup pre-commit hook

### **Bước 3: Verify setup**
```bash
# Test pre-commit hook
./scripts/pre-commit.sh
```

## 🔍 Cách hoạt động

### **Khi commit:**
```bash
git add .
git commit -m "Your commit message"
```

Pre-commit hook sẽ tự động chạy và:
- ✅ **Pass**: Commit thành công
- ❌ **Fail**: Commit bị block, cần fix lỗi trước

### **Khi hook fail:**
```bash
🔍 Running pre-commit checks...
❌ Code is not properly formatted!
Files that need formatting:
internal/usecases/url_usecase.go

Run 'gofmt -w .' to fix formatting issues
```

**Cách fix:**
1. Sửa lỗi theo hướng dẫn
2. Chạy lại commit

## 🛠️ Manual Commands

### **Format code:**
```bash
gofmt -w .
```

### **Run linting:**
```bash
golangci-lint run
```

### **Run tests:**
```bash
go test ./... -v
```

### **Run build check:**
```bash
go build ./cmd/main.go
```

### **Run all checks manually:**
```bash
./scripts/pre-commit.sh
```

## ⚙️ Configuration

### **golangci-lint config** (`.golangci.yml`):
```yaml
linters:
  enable:
    - gofmt          # Code formatting
    - goimports      # Import organization
    - govet          # Go vet
    - errcheck       # Error handling
    - staticcheck    # Static analysis
    - unused         # Unused code
    - gosimple       # Code simplification
    - ineffassign    # Inefficient assignments
    - gosec          # Security issues
    - goconst        # Constants
    - misspell       # Spelling
    - revive         # Code style
```

### **Excluded linters:**
- `typecheck` - Excluded for test files và utils
- `errcheck` - Excluded for test files
- `gosec` - Excluded for test files

## 🚨 Troubleshooting

### **Lỗi: "golangci-lint not found"**
```bash
# Cài đặt golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
```

### **Lỗi: "Go not installed"**
- Cài đặt Go từ: https://golang.org/dl/
- Đảm bảo `go` command có trong PATH

### **Lỗi: "Permission denied"**
```bash
chmod +x scripts/*.sh
```

### **Skip pre-commit hook (không khuyến khích):**
```bash
git commit --no-verify -m "Skip pre-commit checks"
```

## 📁 Files Structure

```
scripts/
├── pre-commit.sh      # Pre-commit hook script
├── setup-hooks.sh     # Setup script for new developers
├── golang.sh          # Go development commands
└── docker.sh          # Docker commands

.git/hooks/
└── pre-commit         # Git pre-commit hook (auto-generated)

.golangci.yml          # golangci-lint configuration
```

## 🎯 Benefits

### **Code Quality:**
- ✅ Consistent code formatting
- ✅ No linting errors
- ✅ All tests pass
- ✅ Code builds successfully

### **Team Collaboration:**
- ✅ Consistent code style across team
- ✅ Catch issues early
- ✅ Reduce code review time
- ✅ Prevent broken builds

### **CI/CD Integration:**
- ✅ Same checks run in GitHub Actions
- ✅ Local checks match CI checks
- ✅ Faster CI pipeline (fewer failures)

## 🔄 Updates

### **Update golangci-lint:**
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
```

### **Update pre-commit hook:**
```bash
cp scripts/pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

## 📞 Support

Nếu gặp vấn đề với pre-commit hooks:
1. Chạy `./scripts/pre-commit.sh` để debug
2. Kiểm tra logs trong terminal
3. Liên hệ team lead để được hỗ trợ

---

**Lưu ý:** Pre-commit hooks là bắt buộc cho tất cả commits. Không được skip hooks trừ trường hợp khẩn cấp và phải có approval từ team lead.
