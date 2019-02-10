.PHONY: fmt lint golint test
PACKAGES=`go list ./... | grep -v vendor | grep -v mocks`

fmt:
	for pkg in ${PACKAGES}; do \
		go fmt $$pkg; \
	done;

lint:
	gometalinter --tests --disable-all --deadline=120s -E vet -E gofmt -E misspell -E ineffassign -E goimports -E deadcode ./...

golint:
	for pkg in ${PACKAGES}; do \
		golint -set_exit_status $$pkg || GOLINT_FAILED=1; \
	done; \
	[ -z "$$GOLINT_FAILED" ]

test:
	TEST_FAILED= ; \
	for pkg in ${PACKAGES}; do \
		go test $$pkg || TEST_FAILED=1; \
	done; \
	[ -z "$$TEST_FAILED" ]