all: server.elf client.elf
.PHONY: clean

server.elf: server.cpp Makefile
	g++ server.cpp -o server.elf -pthread

client.elf: client.cpp Makefile
	g++ client.cpp -o client.elf

clean:
	rm -f client.elf server.elf