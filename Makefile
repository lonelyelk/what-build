binary=what-build
version=`git describe --tags`
flags=-ldflags="-X github.com/lonelyelk/what-build/what.Version=${version}"

build:
	go build ${flags} -o ./bin/$(binary)

install:
	go install ${flags}

run:
	go run ${flags} ./main.go ${ARGS}
