#!/bin/sh

echo "going to listen to :${PORT}"
goexec 'http.ListenAndServe(`:` + os.Getenv(`PORT`), http.FileServer(http.Dir(`./static`)))'
