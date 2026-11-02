default:
    @just --list

setup:
    go get -tool go.bytecodealliance.org/cmd/wit-bindgen-go@latest

generate:
    wkg build -o fluor-fun.wasm
    go generate ./...
    go mod tidy

build: setup generate
    mkdir -p out
    tinygo build -target=wasip2 -o out/{{project-name}}.wasm --wit-package ./fluor-fun.wasm --wit-world function .

inspect: build
    wasm-tools component wit out/{{project-name}}.wasm
