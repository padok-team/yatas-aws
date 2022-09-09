default: build

test:
	go test ./...

build:
	go build -o bin/yatas-aws

install: build
	mkdir -p ~/.yatas.d/plugins
	mv ./bin/yatas-aws ~/.yatas.d/plugins

release: test
	npm run release
	git push --follow-tags origin main 