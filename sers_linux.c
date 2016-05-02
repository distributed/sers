#include <asm/termios.h>
#include <unistd.h>
#include <fcntl.h>
#include <linux/serial.h>
#include <linux/ioctl.h>
#include <asm-generic/ioctls.h>

int setbaudrate(int fd, int br) {

  struct termios2 tio;
  ioctl(fd, TCGETS2, &tio);
  tio.c_cflag &= ~CBAUD;
  tio.c_cflag |= BOTHER;
  tio.c_ispeed = br;
  tio.c_ospeed = br;
  /* do other miscellaneous setup options with the flags here */
  return ioctl(fd, TCSETS2, &tio);

}

int clearnonblocking(int fd) {
  return fcntl(fd, F_SETFL, 0);
}
