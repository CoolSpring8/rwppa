# Binary name
BINARY=rwppa
VERSION=1.0

# Builds
build:
		GO111MODULE=on go build -o ${BINARY} -ldflags "-X main.Version=${VERSION}"
		GO111MODULE=on go test -v

# Installs to $GOPATH/bin
install:
		GO111MODULE=on go install

release_mac:
		go clean
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -ldflags "-s -w -X main.Version=${VERSION}"
		mv ./${BINARY} ${BINARY}-mac64 

release_linux:
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags "-s -w -X main.Version=${VERSION}"
		mv ./${BINARY} ${BINARY}-linux64 

release_windows:
		go clean
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GO111MODULE=on go build -ldflags "-s -w -X main.Version=${VERSION}"
		mv ./${BINARY}.exe ${BINARY}-win64.exe 

# Release for different platforms
release: release_mac release_linux release_windows

clean:
		go clean
		rm ${BINARY}*
.PHONY: clean build release_mac release_linux release_windows
