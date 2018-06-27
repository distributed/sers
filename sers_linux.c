#include <termios.h>
#include <sys/ioctl.h>
#include <fcntl.h>

/* extended termios struct for custom baud rate */
struct termios2 {
	tcflag_t c_iflag;		/* input mode flags */
	tcflag_t c_oflag;		/* output mode flags */
	tcflag_t c_cflag;		/* control mode flags */
	tcflag_t c_lflag;		/* local mode flags */
	cc_t c_line;			/* line discipline */
	cc_t c_cc[NCCS];		/* control characters */
	speed_t c_ispeed;		/* input speed */
	speed_t c_ospeed;		/* output speed */
};

#include <stdio.h>
#include <errno.h>
#include <string.h>


int ioctl1(int i, unsigned int r, void *d) {
    return ioctl(i, r, d);
}

int setbaudrate(int fd, int br) {
	struct termios2 tio;
	int ret;
	
	printf("doing the tcgets2 ioctl\n");
	ret = ioctl(fd, TCGETS2, &tio);
	printf("ret %d errno %d strerror %s\n", ret, errno, strerror(errno));
	if (ret == -1) return ret;

	tio.c_cflag &= ~CBAUD;
	tio.c_cflag |= CBAUDEX;
	tio.c_ispeed = br;
	tio.c_ospeed = br;
	
	return ioctl(fd, TCSETS2, &tio);
}

int clearnonblocking(int fd) {
  return fcntl(fd, F_SETFL, 0);
}
