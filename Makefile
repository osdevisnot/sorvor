EXAMPLES_MINIMAL := examples/minimal*
EXAMPLES_NPM := examples/*counter examples/server*

.PHONY: clean build upgrade start test release verify

clean:
	@git clean -fdX

upgrade:
	@go get -t -u  ./...

build: main.go
	@go build -o /usr/local/bin/sorvor

start:
	@make build && cd ${CURDIR}/examples/preact-counter && yarn install --silent --no-lockfile && yarn start

test:
	@make clean && make build && \
	for dir in $(EXAMPLES_MINIMAL); do cd $$dir; sorvor; cd ${CURDIR}; done
	for dir in $(EXAMPLES_NPM); do cd $${dir}; yarn install --silent --no-lockfile; yarn build; cd ${CURDIR}; done

release:
	@make test && \
	printf "current version: " && git describe --tags `git rev-list --tags --max-count=1`
	@read -p "enter new version: " version; git tag v$$version -m "publish $$version"
	@git push
	@git push --tags
	@make verify

verify:
	@curl -sf https://gobinaries.com/osdevisnot/sorvor | sh
	@cd ${CURDIR}/examples/minimal && sorvor && \
	cd ${CURDIR}/examples/minimal-css && sorvor && \
	cd ${CURDIR}/examples/minimal-typescript && sorvor && \
	cd ${CURDIR}/examples/preact-counter && npm install && npm run build && \
	cd ${CURDIR}/examples/react-counter && npm install && npm run build