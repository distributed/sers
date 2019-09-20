#include <stddef.h>

#include <termios.h>
#include <sys/ioctl.h>

#include <fcntl.h>

/*

So... custom baud rate support is a very sad topic under linux. Sometimes
there's no asm/termios.h, sometimes is clashes with termios.h, sometimes
termios.h has the wrong definition for NCCS, messing up the layout of a
struct termios2 that is defined here but dependent on NCCS.


The picocom has a whole rant/article about the situation:
https://github.com/Rosonix/picocom/blob/master/termios2.txt

Their copy-pasted implementation lives here:
https://github.com/Rosonix/picocom/blob/238547d7174571b2a463738f03f3dea128e5c676/termbits2.h

And is probably wrong for x86_64. But who knows?!


The pyserial approach can be found at:
https://github.com/pyserial/pyserial/blob/master/serial/serialposix.py

Look at the function _set_special_baudrate. I guess the code is wrong for
pretty much every platform that is not x86_64. Maybe nobody ever noticed?

*/

#if NCCS == 32
// i have not seen a platform with NCCS == 32 in the kernel header files. i saw
// NCCS == 32 On My Machine where it should have been 19, so we're using that.
// Getting this wrong will most likely result in a ioctl that returns an error.
#define NCCS_K 19
#else
#define NCCS_K NCCS
#endif
/* extended termios struct for custom baud rate */
struct termios3 {
	tcflag_t c_iflag;		/* input mode flags */
	tcflag_t c_oflag;		/* output mode flags */
	tcflag_t c_cflag;		/* control mode flags */
	tcflag_t c_lflag;		/* local mode flags */
	cc_t c_line;			/* line discipline */
	cc_t c_cc[NCCS_K];		/* control characters */
	speed_t c_ispeed;		/* input speed */
	speed_t c_ospeed;		/* output speed */
};

#define TCGETS3 _IOR('T', 0x2A, struct termios3)
#define TCSETS3 _IOW('T', 0x2B, struct termios3)

#include <stdio.h>
#include <errno.h>
#include <string.h>


int ioctl1(int i, unsigned int r, void *d) {
    return ioctl(i, r, d);
}

speed_t lookupbaudrate(int br) {
	switch (br) {
		case 50     : return B50      ;
		case 75     : return B75      ;
		case 110    : return B110     ;
		case 134    : return B134     ;
		case 150    : return B150     ;
		case 200    : return B200     ;
		case 300    : return B300     ;
		case 600    : return B600     ;
		case 1200   : return B1200    ;
		case 1800   : return B1800    ;
		case 2400   : return B2400    ;
		case 4800   : return B4800    ;
		case 9600   : return B9600    ;
		case 19200  : return B19200   ;
		case 38400  : return B38400   ;
		case 57600  : return B57600   ;
		case 115200 : return B115200  ;
		case 230400 : return B230400  ;
		case 460800 : return B460800  ;
		case 500000 : return B500000  ;
		case 576000 : return B576000  ;
		case 921600 : return B921600  ;
		case 1000000: return B1000000 ;
		case 1152000: return B1152000 ;
		case 1500000: return B1500000 ;
		case 2000000: return B2000000 ;
		case 2500000: return B2500000 ;
		case 3000000: return B3000000 ;
		case 3500000: return B3500000 ;
		case 4000000: return B4000000 ;
	}
	return 0;
}

int speedtobaudrate(speed_t speed) {
    printf("enum baud rate reverse matching\n");
	switch (speed) {
		case B50     : return 50      ;
		case B75     : return 75      ;
		case B110    : return 110     ;
		case B134    : return 134     ;
		case B150    : return 150     ;
		case B200    : return 200     ;
		case B300    : return 300     ;
		case B600    : return 600     ;
		case B1200   : return 1200    ;
		case B1800   : return 1800    ;
		case B2400   : return 2400    ;
		case B4800   : return 4800    ;
		case B9600   : return 9600    ;
		case B19200  : return 19200   ;
		case B38400  : return 38400   ;
		case B57600  : return 57600   ;
		case B115200 : return 115200  ;
		case B230400 : return 230400  ;
		case B460800 : return 460800  ;
		case B500000 : return 500000  ;
		case B576000 : return 576000  ;
		case B921600 : return 921600  ;
		case B1000000: return 1000000 ;
		case B1152000: return 1152000 ;
		case B1500000: return 1500000 ;
		case B2000000: return 2000000 ;
		case B2500000: return 2500000 ;
		case B3000000: return 3000000 ;
		case B3500000: return 3500000 ;
		case B4000000: return 4000000 ;
	}
	return 0;
}

int setbaudrate(int fd, int br) {
	// we try to set the baud rate via the old school interface first. this
	// should work more reliably, if only for certain baud rates.
	speed_t baudrateconstant = lookupbaudrate(br);
	if (baudrateconstant) {
		struct termios tio;
		int ret = tcgetattr(fd, &tio);
		if (ret == -1) return ret;

		ret = cfsetispeed(&tio, baudrateconstant);
		if (ret == -1) return ret;

		ret = cfsetospeed(&tio, baudrateconstant);
		if (ret == -1) return ret;

		ret = tcsetattr(fd, TCSANOW, &tio);
		return ret;
	}

	// if we have a custom baud rate, we use our own cobbled together approach.

	struct termios3 tio;
	int ret;

	ret = ioctl(fd, TCGETS3, &tio);
	if (ret == -1) return ret;

	tio.c_cflag &= ~CBAUD;
	tio.c_cflag |= CBAUDEX;
	tio.c_ispeed = br;
	tio.c_ospeed = br;
	
	return ioctl(fd, TCSETS3, &tio);
}

int getbaudrate(int fd, int *br) {
    struct termios3 tio;
	int ret = ioctl(fd, TCGETS3, &tio);
	if (ret == -1) return ret;

	if (tio.c_cflag & CBAUDEX) {
	    *br = tio.c_ispeed;
	    return 0;
	}

	*br = speedtobaudrate(tio.c_cflag & CBAUD);

	return 0;
}
