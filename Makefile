default: build

test:
	go test ./...

build:
	go build -o bin/yatas-aws

update:
	go get -u 
	go mod tidy

install: build
	mkdir -p ~/.yatas.d/plugins/Thibaut/yatas-aws/local/
	mv ./bin/yatas-aws ~/.yatas.d/plugins/Thibaut/yatas-aws/local/

release: test
	npm run release
	git push --follow-tags origin main 