all: racer.elf
.PHONY: clean

CC = clang


racer.elf: racer_client.c Makefile
	$(CC) racer_client.c -o racer.elf

clean:
	rm -f racer.elf