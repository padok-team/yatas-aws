default: build

test:
	go test ./...

build:
	go build

install: build
	mkdir -p ~/.yatas.d/plugins
	mv ./yatas-aws ~/.yatas.d/plugins

release: test
	npm run release
	git push --follow-tags origin main 