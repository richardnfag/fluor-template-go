# {{project-name}}

Starter template for a Fluor function in Go using TinyGo and `wit-bindgen-go`.

## Requirements

- Go 1.24+
- TinyGo 0.34+
- `just`: `cargo install just`
- `wit-bindgen-go`: `go get -tool go.bytecodealliance.org/cmd/wit-bindgen-go`
- `wkg`: `cargo install wkg`
- `wasm-tools`: `cargo install --locked wasm-tools`
- `wasm-opt` on `PATH` (for example via Binaryen)

## Build

```bash
just build
```

This runs `go generate` and produces `out/{{project-name}}.wasm`.

## Notes

The Go Component Model toolchain is still evolving. If `wit-bindgen-go` changes its generated package layout, adjust the import in `main.go` accordingly.

Fluor now targets the `v2` contract with an exported `handler` interface and typed `request`/`response` records. Regenerate bindings after cloning before adjusting the sample implementation if `wit-bindgen-go` changes its emitted result helpers.

TinyGo currently expects `wasm-opt` to be available during the `wasip2` build. If `just build` fails with `could not find wasm-opt`, install Binaryen or point `WASMOPT` at your local `wasm-opt` binary.

This scaffold uses `wkg build` to resolve the WASI dependencies required by TinyGo's `wasip2` target before generating bindings and compiling the component.

Validated locally on November 2, 2026 with Go 1.25.4, TinyGo 0.44.0, `wit-bindgen-go` 0.7.0, `wkg` 0.16.1, `wasm-tools` 1.263.0, and Binaryen `wasm-opt` 133.
