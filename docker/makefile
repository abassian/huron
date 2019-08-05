HURON_VERSION=0.5.0
DUMMY_VERSION=0.5.0
WATCHER_VERSION=0.0.2
GLIDER_VERSION=0.0.3

all: huron-image dummy-image watcher-image glider-image

demo: huron-image dummy-image watcher-image

huron-image:
	CGO_ENABLED=0 go build \
		--ldflags "-X github.com/abassian/huron/src/version.GitCommit=`git rev-parse HEAD`" \
		-o huron/huron ../cmd/huron/
	docker build --no-cache=true -t abassian/huron:$(HURON_VERSION)  -t abassian/huron:latest huron/

dummy-image:
	CGO_ENABLED=0 go build \
		-o dummy/dummy ../cmd/dummy/
	docker build --no-cache=true -t abassian/dummy:$(DUMMY_VERSION) -t abassian/dummy:latest dummy/

watcher-image:
	docker build --no-cache=true -t abassian/watcher:$(WATCHER_VERSION) -t abassian/watcher:latest watcher/

glider-image:
	docker build --no-cache=true -t abassian/glider:$(GLIDER_VERSION) glider/

.PHONY: all demo huron-image dummy-image watcher-image glider-image
