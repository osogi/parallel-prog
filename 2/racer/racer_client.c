#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <unistd.h>

#define CLIENT_PATH "./client.elf"
#define MAX_SLEEP   50000

int main(int argc, char* argv[]) {
    static char* newargv[] = {NULL};
    newargv[0] = CLIENT_PATH;

    fork();
    fork();
    fork();
    fork();

    while (1) {
        pid_t chpid = fork();
        if (chpid == 0) {
            execv(newargv[0], newargv);
            perror("execve"); /* execve() returns only on error */
            exit(EXIT_FAILURE);
        } else {
            usleep(rand() % MAX_SLEEP);
            kill(chpid, 9);
        }
    }
}