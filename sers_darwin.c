#include <sys/ioctl.h>

int ioctl1(int i, unsigned int r, void *d) {
    return ioctl(i, r, d);
}

int fcntl1(int i, unsigned int r, unsigned int d) {
    return fcntl(i, r, d);
}
