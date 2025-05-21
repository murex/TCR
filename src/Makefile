.PHONY: default
default: build ;

.PHONY: deps
deps:
	@go get -u -t tool ./...

.PHONY: lint
lint:
	@golangci-lint run --build-tags=test_helper -v

.PHONY: modernize
modernize:
	@go run -tags test_helper golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -test -fix ./...

.PHONY: vet
vet:
	@go vet --tags=test_helper ./...

define RUN_TESTS
@mkdir -p _test_results
@go tool gotestsum --format pkgname --junitfile _test_results/output.xml --hide-summary=skipped -- $(1) -covermode atomic -coverprofile=_test_results/coverage-all.out -tags=test_helper -v ./...
@grep -v "_test_" _test_results/coverage-all.out >| _test_results/coverage.out
@sed -e 's:\(/tcr/\):\1src/:g' _test_results/coverage.out >| _test_results/coverage-sonar.out
endef

.PHONY: test
test:
	$(call RUN_TESTS,)

.PHONY: test-short
test-short:
	$(call RUN_TESTS,-short)

.PHONY: slow-tests
slow-tests:
	@go test -json -tags=test_helper ./... | go tool gotestsum tool slowest --threshold 500ms

.PHONY: cov
cov:
	@go tool cover -html=_test_results/coverage.out

.PHONY: tidy
tidy:
	@go mod tidy

CONFIG_PKG="github.com/murex/tcr/settings"
export CGO_ENABLED=0
.PHONY: build
build:
	@go build -ldflags "-s -w \
    		-X ${CONFIG_PKG}.BuildVersion="`git describe --tags`" \
            -X ${CONFIG_PKG}.BuildOs="`go env GOOS`" \
            -X ${CONFIG_PKG}.BuildArch="`go env GOARCH`" \
            -X ${CONFIG_PKG}.BuildCommit="`git rev-list --max-count=1 --tags`" \
            -X ${CONFIG_PKG}.BuildDate="`date -u +%FT%TZ`" \
            -X ${CONFIG_PKG}.BuildAuthor="`id -un`" \
            "

.PHONY: run
run:
	@tcr-local

.PHONY: download
download:
	@go mod download
