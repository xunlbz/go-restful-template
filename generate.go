//go:generate cp -r ./webui/dist .
//go:generate go-bindata -o=bindata/bindata.go -pkg=bindata dist/...
//go:generate rm -rf ./dist

package main
