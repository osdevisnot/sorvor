EXAMPLES_MINIMAL := examples/minimal*
EXAMPLES_NPM := examples/*counter examples/server*
NPM_PACKAGES := npm/*

export version

.PHONY: clean upgrade build start test release verify check

clean:
	@git clean -fdX

upgrade:
	@rm go.sum
	@go get -t -u  ./...
	@for dir in $(EXAMPLES_NPM); do cd ${CURDIR}/$${dir}; ncu --upgrade --silent; done

build-local: main.go
	@go build -ldflags="-X 'main.version=local'" -o /usr/local/bin/sorvor

build: main.go
	GOOS=darwin  GOARCH=amd64 go build -ldflags="-X 'main.version=$(version)'" -o npm/sorvor-darwin-64/sorvor
	GOOS=darwin  GOARCH=arm64 go build -ldflags="-X 'main.version=$(version)'" -o npm/sorvor-darwin-arm64/sorvor
	GOOS=freebsd GOARCH=amd64 go build -ldflags="-X 'main.version=$(version)'" -o npm/sorvor-freebsd-64/sorvor
	GOOS=freebsd GOARCH=arm64 go build -ldflags="-X 'main.version=$(version)'" -o npm/sorvor-freebsd-arm64/sorvor
	GOOS=linux   GOARCH=386   go build -ldflags="-X 'main.version=$(version)'" -o npm/sorvor-linux-32/sorvor
	GOOS=linux   GOARCH=amd64 go build -ldflags="-X 'main.version=$(version)'" -o npm/sorvor-linux-64/sorvor
	GOOS=linux   GOARCH=arm   go build -ldflags="-X 'main.version=$(version)'" -o npm/sorvor-linux-arm/sorvor
	GOOS=linux   GOARCH=arm64 go build -ldflags="-X 'main.version=$(version)'" -o npm/sorvor-linux-arm64/sorvor
	GOOS=windows GOARCH=386   go build -ldflags="-X 'main.version=$(version)'" -o npm/sorvor-windows-32/sorvor.exe
	GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.version=$(version)'" -o npm/sorvor-windows-64/sorvor.exe

start:
	@make build && cd ${CURDIR}/examples/preact-counter && yarn install --silent --no-lockfile && yarn start

test:
	@make clean && make build-local && make check

release:
	@printf "Current Version: " && git describe --tags `git rev-list --tags --max-count=1`
	@read -p "Enter New Version: " version; make build; \
	node scripts/version.js $$version; \
	git commit -am "publish $$version"; \
	git tag -a v$$version -m "publish $$version"; \
	git push && git push --tags
	make verify

verify:
	@curl -sf https://gobinaries.com/osdevisnot/sorvor | sh && \
	make check


check:
	@for dir in $(EXAMPLES_MINIMAL); do cd ${CURDIR}/$${dir}; sorvor; done && \
    for dir in $(EXAMPLES_NPM); do cd ${CURDIR}/$${dir}; yarn install --silent --no-lockfile; yarn build; done
