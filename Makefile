upgrade:
	@go get -u  ./...

build: main.go
	@go build -o /usr/local/bin/sorvor

verify:
	@make clean && make build && \
	cd ${CURDIR}/examples/minimal && sorvor && \
	cd ${CURDIR}/examples/minimal-css && sorvor && \
	cd ${CURDIR}/examples/minimal-typescript && sorvor && \
	cd ${CURDIR}/examples/preact-counter && npm run build

clean:
	@git clean -fdX

release:
	@make verify && \
	printf "current version: " && git describe --tags --abbrev=0
	@read -p "enter new version: " version; git tag v$$version
	@git push --tags
