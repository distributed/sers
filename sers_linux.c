#include <unistd.h>
#include <fcntl.h>
#include <linux/serial.h>
#include <linux/ioctl.h>
#include <asm-generic/ioctls.h>
#include <asm/termios.h>

int setbaudrate(int fd, int br) {
	struct termios2 tio;
	int ret;
	
	ret = ioctl(fd, TCGETS2, &tio);
	if (ret == -1) return ret;

	tio.c_cflag &= ~CBAUD;
	tio.c_cflag |= BOTHER;
	tio.c_ispeed = br;
	tio.c_ospeed = br;
	
	return ioctl(fd, TCSETS2, &tio);
}

int clearnonblocking(int fd) {
  return fcntl(fd, F_SETFL, 0);
}
