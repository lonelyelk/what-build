binary=what-build
version=`git describe --tags`
flags=-ldflags="-X github.com/lonelyelk/what-build/what.Version=${version}"
cyclomax=15
cyclocmd=`gocyclo -over ${cyclomax} . | wc -c`

build:
	go build ${flags} -o ./bin/$(binary)

build_all:
	env GOOS=darwin GOARCH=amd64 go build ${flags} -o ./bin/lgtm-darwin-amd64-$(version)
	env GOOS=linux GOARCH=amd64 go build ${flags} -o ./bin/lgtm-linux-amd64-$(version)

install:
	go install ${flags}

run:
	go run ${flags} ./main.go ${ARGS}

lint:
	if [[ ${cyclocmd} -ne 0 ]] ;then echo "Cyclomatic complexity over threshold:" && gocyclo -over $(cyclomax) . && exit 1; fi

test:
	go test ./... -cover

dep:
	go get -v -t -d ./...
