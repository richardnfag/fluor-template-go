//go:generate go tool wit-bindgen-go generate --world function --out internal ./fluor-fun.wasm

package main

import (
    "example.com/{{project-name}}/internal/fluor/fun/handler"
    "go.bytecodealliance.org/cm"
)

func init() {
    handler.Exports.Handle = func(request handler.Request) cm.Result[handler.ResponseShape, handler.Response, handler.Error] {
        return cm.OK[cm.Result[handler.ResponseShape, handler.Response, handler.Error]](handler.Response{
            Status: 200,
            Headers: cm.ToList([]handler.Header{
                {
                    Name:  "content-type",
                    Value: "application/json",
                },
            }),
            Body: `{"message":"Hello from {{project-name}}"}`,
        })
    }
}

func main() {}
