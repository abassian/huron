BUILD_TAGS?=huron

# vendor uses Glide to install all the Go dependencies in vendor/
vendor:
	glide install

# install compiles and places the binary in GOPATH/bin
install:
	go install --ldflags '-extldflags "-static"' \
		--ldflags "-X github.com/abassian/huron/src/version.GitCommit=`git rev-parse HEAD`" \
		./cmd/huron

# build compiles and places the binary in /build
build:
	CGO_ENABLED=0 go build \
		--ldflags "-X github.com/abassian/huron/src/version.GitCommit=`git rev-parse HEAD`" \
		-o build/huron ./cmd/huron/

# dist builds binaries for all platforms and packages them for distribution
dist:
	@BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/dist.sh'"

tests:  test

test:
	glide novendor | xargs go test -count=1 -tags=unit

extratests:
	glide novendor | xargs go test -count=1 -run Extra

alltests:
	glide novendor | xargs go test -count=1 

.PHONY: vendor install build dist test extratests alltests tests
