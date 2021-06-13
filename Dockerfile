FROM golang:1.16-alpine

WORKDIR /go/src/github.com/saimanwong/go-game-of-life
COPY . .
RUN GOOS=js GOARCH=wasm go build -o main.wasm main.go world.go
RUN cp $(go env GOROOT)/misc/wasm/wasm_exec.js main.wasm ./static
EXPOSE 8080
RUN go get -u github.com/shurcooL/goexec
ENTRYPOINT goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`./static`)))'
