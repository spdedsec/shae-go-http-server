.PHONY: build test release clean install help
GO ?= go
BIN = shae-http-server

build:
	@$(GO) build -o $(BIN) src/main.go

test: build
	@echo "Running smoke test..."
	@./$(BIN) & sleep 1
	@curl -s http://localhost:8080/ | grep -q '"status":"ok"'
	@kill %1 2>/dev/null || true
	@echo "PASS"

release: build
	@tar czf dist/$(BIN)-v0.1.0-src.tar.gz src/ scripts/ Makefile README.md DESIGN.md LICENSE docs/

install: build
	@cp $(BIN) ~/.local/bin/

clean:
	@rm -f $(BIN)
	@rm -rf dist/

help:
	@echo "Targets: build test release install clean help"
