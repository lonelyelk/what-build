binary=what-build
version=`git describe --tags`

build:
	go build -ldflags="-X github.com/lonelyelk/what-build/what.Version=${version}" -o ./bin/$(binary)

install:
	go install -ldflags="-X github.com/lonelyelk/what-build/what.Version=${version}"
