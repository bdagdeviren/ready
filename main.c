#ifdef _WIN32
#include <Windows.h>
#else
#include <unistd.h>
#endif

#include <log.h>
#include <stdlib.h>
#include "libs/ready/libready.h"

int main() {
    int waitTime;
    char *env = getenv("WAIT_TIME");
    if (env){
        waitTime = atoi(env);
    }else {
        waitTime = 5;
    }

    char *clusterName = getenv("CLUSTER_NAME");
    int workerCount = (int) getenv("WORKER_MACHINE_COUNT");
    int masterCount = (int) getenv("CONTROL_PLANE_MACHINE_COUNT");
    if ( clusterName && workerCount && masterCount) {
        while (1) {
            struct node_create_check_return node_create_check_return;
            node_create_check_return = node_create_check();
            if (node_create_check_return.r0 != 0) {
                log_error(node_create_check_return.r2);
#ifdef _WIN32
                Sleep(waitTime);
#else
                sleep(waitTime);
#endif
            } else {
                log_info(node_create_check_return.r1);
                break;
            }
        }
    }else {
        log_info("Set CLUSTER_NAME,CONTROL_PLANE_MACHINE_COUNT,WORKER_MACHINE_COUNT environment variables!");
    }

    return 0;
}
