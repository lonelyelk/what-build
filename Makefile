binary=what-build
version=`git describe --tags`
flags=-ldflags="-X github.com/lonelyelk/what-build/what.Version=${version}"
cyclomax=15
cyclocmd=`gocyclo -over ${cyclomax} . | wc -c`

build:
	go build ${flags} -o ./bin/$(binary)

install:
	go install ${flags}

run:
	go run ${flags} ./main.go ${ARGS}

lint:
	if [[ ${cyclocmd} -ne 0 ]] ;then echo "Cyclomatic complexity over threshold:" && gocyclo -over $(cyclomax) . && exit 1; fi
