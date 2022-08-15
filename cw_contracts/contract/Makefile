build:
	RUSTFLAGS='-C link-arg=-s' cargo wasm
	cp target/wasm32-unknown-unknown/release/watcher.wasm .
.PHONY: build
