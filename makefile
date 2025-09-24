run: 
	go run *.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bbb-multicleaner main.go logging.go

build:
	go build

clean:
	rm -f bbb-multicleaner

# Pre-commit checks (same as the Git hook)
pre-commit:
	@echo "🔨 Running pre-commit checks..."
	@go fmt ./...
	@echo "📝 Code formatted"
	@go vet ./...
	@echo "🔍 go vet passed"
	@go build
	@echo "🏗️  Build successful"
	@if ls *_test.go >/dev/null 2>&1; then go test ./...; fi
	@echo "✅ All pre-commit checks passed!"

# Install the Git pre-commit hook
install-hooks:
	@echo "📦 Installing Git pre-commit hook..."
	@mkdir -p .git/hooks
	@cp scripts/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "✅ Pre-commit hook installed!"

# Test the pre-commit hook
test-hook:
	@echo "🧪 Testing pre-commit hook..."
	@./.git/hooks/pre-commit