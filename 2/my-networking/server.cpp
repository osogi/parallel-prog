// This server is based on Hallo Weeks's server
// https://github.com/halloweeks/networking/tree/main
// #define __cplusplus 201403L

#include "profile.h"
#include <arpa/inet.h>
#include <errno.h>
#include <pthread.h>
#include <shared_mutex>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <time.h>
#include <unistd.h>
#include <vector>

#define ADDRESS               "0.0.0.0"
#define PORT                  12345
#define CONCURRENT_CONNECTION 200
#define QUEUE_CONNECTION      100
#define BUFFER_SIZE           65356
#define THREAD_STACK_SIZE     524288

enum state_t { LOGIN, MAIN_MENU, PROFILE_EDIT, USERS_LIST, WATCHING_OTHER_PROFILE };

struct user_t {
    profile_t* profile;
    std::shared_mutex profile_edit_mutex;
    std::shared_mutex last_like_mutex;
    state_t state;
    int wathched_profile;

    user_t() {
        profile = (profile_t*)malloc(sizeof(profile));
        profile->liked_by_last_visiter = 0;
        profile->love_cats = 0;
        profile->love_dogs = 0;
        state = LOGIN;
        wathched_profile = 0;
    };

    int sprintf_profile(char* target_buffer) {
        profile_edit_mutex.lock_shared();
        last_like_mutex.lock_shared();
        int res = sprintf(target_buffer, "Love dogs: %d\nLove cats: %d\nLiked by last user: %d\n", profile->love_dogs,
                          profile->love_cats, profile->liked_by_last_visiter);
        profile_edit_mutex.unlock_shared();
        last_like_mutex.unlock_shared();
        return res;
    };

    void set_love(int to_dogs, int to_cats) {
        profile_edit_mutex.lock();
        profile->love_dogs = (to_dogs);
        profile->love_cats = (to_cats);
        profile_edit_mutex.unlock();
    }

    void set_last_liked(int liked) {
        last_like_mutex.lock();
        profile->liked_by_last_visiter = liked;
        last_like_mutex.unlock();
    }
};

int char2bool(char c) {
    c |= 0x20;
    if (c != 'y') {
        return 0;
    } else {
        return 1;
    }
}

std::shared_mutex users_mutex;

std::vector<user_t*> users;

pthread_mutex_t connection_mutex = PTHREAD_MUTEX_INITIALIZER;

int connection = 0;

void* connection_handler(void*);

int main(int argc, char* argv[]) {

    pthread_t thread_id;
    pthread_attr_t attr;

    if (pthread_attr_init(&attr) != 0) {
        printf("[ERROR][THREAD][INIT] %s\n", strerror(errno));
        return 1;
    }

    if (pthread_attr_setstacksize(&attr, THREAD_STACK_SIZE) != 0) {
        printf("[ERROR][THREAD][STACK] %s\n", strerror(errno));
        return 1;
    }

    if (pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED) != 0) {
        printf("[ERROR][THREAD][DETACH] %s\n", strerror(errno));
        return 1;
    }

    int master_socket, conn_id;
    struct sockaddr_in server, client;

    memset(&server, 0, sizeof(server));
    memset(&client, 0, sizeof(client));

    signal(SIGPIPE, SIG_IGN);

    if ((master_socket = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
        printf("[ERROR] CAN'T CREATE TO SOCKET\n");
        return 1;
    } else {
        printf("[NOTE] SOCKET CREATED DONE\n");
    }

    server.sin_family = AF_INET;
    server.sin_addr.s_addr = inet_addr(ADDRESS);
    server.sin_port = htons(PORT);

    socklen_t addrlen = sizeof(struct sockaddr_in);

    if (bind(master_socket, (struct sockaddr*)&server, sizeof(server)) == -1) {
        printf("[ERROR][BIND] %s\n", strerror(errno));
        return 1;
    } else {
        printf("[NOTE] BIND %s:%d\n", ADDRESS, PORT);
    }

    if (listen(master_socket, QUEUE_CONNECTION) == -1) {
        printf("[ERROR][LISTEN] %s\n", strerror(errno));
        return 1;
    } else {
        printf("[INFO] WAITING FOR INCOMING CONNECTIONS\n");
    }

    while (1) {
        conn_id = accept(master_socket, (struct sockaddr*)&client, (socklen_t*)&addrlen);

        if (conn_id == -1) {
            printf("[WARNING] CAN'T ACCEPT NEW CONNECTION\n");
        } else {
            pthread_mutex_lock(&connection_mutex);

            if (connection >= CONCURRENT_CONNECTION) {
                printf("[WARNING] CONNECTION LIMIT REACHED\n");
                close(conn_id);
            } else {
                int* conn_id_heap = (int*)malloc(sizeof(int));

                if (conn_id_heap == NULL) {
                    perror("[ERROR] Memory allocation failed");
                    close(conn_id);
                    continue;
                }

                *conn_id_heap = conn_id;

                printf("[INFO] NEW CONNECTION ACCEPTED FROM %s:%d\n", inet_ntoa(client.sin_addr),
                       ntohs(client.sin_port));

                if (pthread_create(&thread_id, &attr, connection_handler, (void*)conn_id_heap) == -1) {
                    printf("[WARNING] CAN'T CREATE NEW THREAD\n");
                    close(conn_id);
                    free(conn_id_heap);
                } else {
                    printf("[INFO] NEW THREAD CREATED\n");
                    connection++;
                }
            }
            pthread_mutex_unlock(&connection_mutex);
        }
    }

    return 0;
}

void* connection_handler(void* sock_fd) {
    int conn_id = *((int*)sock_fd);
    free(sock_fd);

    clock_t start, end;
    start = clock();

    char buffer[BUFFER_SIZE] = {0};
    char* response = (char*)malloc(1024);
    char* start_response = response;

    user_t user = user_t();
    users_mutex.lock();
    int id = users.size();
    users.push_back(&user);
    users_mutex.unlock();

    while (recv(conn_id, buffer, BUFFER_SIZE, 0) > 0) {
        printf("[RECEIVED] %s\n", buffer);

        int buffer_int = atoi(buffer);

        response = start_response;
        char is_response_ready = 0;

        switch (user.state) {
        default:
            user.state = MAIN_MENU;
            break;
        case MAIN_MENU:
            switch (buffer_int) {
            case 1:
                user.state = PROFILE_EDIT;
                break;
            case 2:
                user.state = USERS_LIST;
                break;
            default:
                break;
            }
            break;
        case PROFILE_EDIT:
            user.set_love(char2bool(buffer[0]), char2bool(buffer[1]));
            user.state = MAIN_MENU;
            break;
        case USERS_LIST:
            users_mutex.lock_shared();
            buffer_int--;
            if ((users.size() > buffer_int) && (buffer_int >= 0)) {
                response += users[buffer_int]->sprintf_profile(response);
                user.wathched_profile = buffer_int;
                user.state = WATCHING_OTHER_PROFILE;
            } else {
                user.state = MAIN_MENU;
            }
            users_mutex.unlock_shared();
            break;
        case WATCHING_OTHER_PROFILE:
            users_mutex.lock_shared();
            if ((users.size() > user.wathched_profile) && (user.wathched_profile >= 0)) {
                users[user.wathched_profile]->set_last_liked(char2bool(buffer[0]));
            }

            user.state = MAIN_MENU;
            users_mutex.unlock_shared();
            break;
        }

        if (!is_response_ready) {
            switch (user.state) {
            default: // MAIN MENU
                sprintf(response, "Choose action: \n[1] Edit profile\n[2] Look other profiles\n");
                break;
            case PROFILE_EDIT:
                sprintf(response, "Do you love dogs? Do you love cats?\n[YY]: Yes, Yes;\n[YN]: Yes, No;\n[NY]: No, "
                                  "Yes;\n[NN]: No, No;\n\n");
                break;
            case USERS_LIST:
                users_mutex.lock_shared();
                sprintf(response, "There is current %zu profiles. Write number of profile which you want watch\n",
                        users.size());
                users_mutex.unlock_shared();
                break;
            case WATCHING_OTHER_PROFILE:
                sprintf(response, "Like this profile? [Y/N]\n");
                break;
            }
            is_response_ready = 1;
        }

        if (send(conn_id, start_response, strlen(start_response), 0) > 0) {
            printf("[SEND] %s\n", start_response);
        } else {
            printf("[WARNING][SEND] %s\n", strerror(errno));
        }

        memset(buffer, 0, BUFFER_SIZE);
    }

    close(conn_id);

    printf("[INFO] CONNECTION CLOSED\n");

    pthread_mutex_lock(&connection_mutex);
    connection--;
    printf("[DEBUG] CONNECTION LIVE: %d\n", connection);
    pthread_mutex_unlock(&connection_mutex);

    printf("[INFO] THREAD TERMINATED\n");

    end = clock();
    double time_taken = (double)(end - start) / CLOCKS_PER_SEC;
    printf("[TIME] PROCESS COMPLETE IN %.5f SEC\n", time_taken);
    printf("------------------------\n");

    pthread_exit(NULL);
}
