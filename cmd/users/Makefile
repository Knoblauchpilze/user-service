
APPLICATION ?= users

setup:
	mkdir -p build/bin
	cp configs/users-template-dev.yml configs/users-dev.yml

release:
	go build -o build/bin/${APPLICATION} main.go

run: release
	./build/bin/${APPLICATION} users-dev

clean:
	rm -rf build configs/users-dev.yml
