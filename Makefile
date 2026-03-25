.PHONY: run tidy

run:
	GOPROXY="off,direct" GONOSUMDB="github.com/cjksplunk/*" TMPDIR=/tmp go run . --config config.yaml

tidy:
	GOPROXY="off,direct" GONOSUMDB="github.com/cjksplunk/*" TMPDIR=/tmp go mod tidy
