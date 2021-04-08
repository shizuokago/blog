echo "build wasm"

GOOS=js GOARCH=wasm go build -o editor.wasm wasm.go

echo "gzip wasm"

gzip -9 -v -c editor.wasm > editor.wasm.gz

rm editor.wasm

echo "move wasm"
mv editor.wasm.gz ../../handler/internal/_assets/static/admin

echo "updated"
