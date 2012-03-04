#include <termios.h>
#include <unistd.h>
#include <fcntl.h>
#include <linux/serial.h>
#include <linux/ioctl.h>
#include <asm-generic/ioctls.h>

int setbaudrate(int fd, int br) {
  struct serial_struct ss;

  int ret = ioctl(fd, TIOCGSERIAL, &ss);
  if (ret < 0) return ret;

  ss.flags = (ss.flags & (~ASYNC_SPD_MASK)) | ASYNC_SPD_CUST;
  ss.custom_divisor = (ss.baud_base + (br/2)) / br;
  
  return ioctl(fd, TIOCSSERIAL, &ss);
}

int clearnonblocking(int fd) {
  return fcntl(fd, F_SETFL, 0);
}
