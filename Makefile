all: clean mediaproxy

mediaproxy:
	go build -o mediaproxy main.go

install:
	mv mediaproxy /usr/bin

clean:
	rm -f mediaproxy
