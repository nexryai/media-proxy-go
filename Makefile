all: clean mediaproxy

mediaproxy:
	go build -ldflags="-s -w" -trimpath -o mediaproxy main.go

install:
	mv mediaproxy /usr/bin

clean:
	rm -f mediaproxy
