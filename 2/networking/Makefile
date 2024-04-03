all: server.elf client.elf sanit_server.elf
.PHONY: clean

CPP = clang++


server.elf: server.cpp Makefile
	$(CPP) server.cpp -o server.elf -pthread

client.elf: client.cpp Makefile
	$(CPP) client.cpp -o client.elf

sanit_server.elf: server.cpp Makefile
	$(CPP) server.cpp -fsanitize=thread -o sanit_server.elf -pthread 

clean:
	rm -f client.elf server.elf sanit_server.elf