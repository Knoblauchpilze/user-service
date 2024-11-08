
APPLICATION ?= users

release:
	go build -o build/bin/${APPLICATION} main.go

run: release
	./build/bin/${APPLICATION}

clean:
	rm -rf build
