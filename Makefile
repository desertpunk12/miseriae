# Makefile for Cloudflare Worker with TinyGo

# Set environment variables for TinyGo WASM build
export GOOS=js
export GOARCH=wasm

generate:
	@echo "Generating templ templates"
	templ generate

# Default target
.PHONY: build
build: generate
	@echo "Building WASM module with TinyGo..."
	tinygo build -o main.wasm -target wasm .
	@powershell -Command "if (!(Test-Path wasm_exec.js)) { Copy-Item \"$$(tinygo env TINYGOROOT)/targets/wasm_exec.js\" -Destination . }"
	@echo "Build complete."
deploy: build
	@echo Deploying to Cloudflare
	wrangler deploy


.PHONY: clean
clean:
	rm -f main.wasm wasm_exec.js
