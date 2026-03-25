.PHONY: run tidy

run:
	GONOSUMDB="github.com/cjksplunk/*" go run . --config config.yaml

tidy:
	GONOSUMDB="github.com/cjksplunk/*" go mod tidy
