sers
====

Overview
--------

Package sers offers serial port access for the Go programming language.
It is a stated goal of this project to allow for configuration of
non-traditional bit rates as the may be useful in a wide range of
embedded projects.

Functionality
-------------

Primarily the library provides an implementation of the `SerialPort` interface
for the supported platforms. For every serial port, you get the functionality
of an `io.ReadWriteCloser` as well as methods to set and query various
parameters. Secondarily it helps you communicate with the users of your
programs through the use of mode strings that represent the data than can be
written and read with `{Get,Set}Mode`.

The `SetMode()` and `GetMode()` methods allow setting and retrieving key
parameters such as baud rate, number of data bits, number of stop bits and the
used parity scheme. As far as the underlying platforms support it, `SetMode`
allows you to set arbitrary, non-traditional baud rate.  The generation of
break conditions is supported through the `SetBreak()` method. The behavior
of the `Read()` method can be fine tuned with the `SetReadParams()` method.

The library offers the concept of mode strings like `9600,7e1` or
`57600,8n1,rtscts`. These can be parsed by `ParseModestring` into an internal
representation that can be set through `SetModeStruct`. The mode string is quite
flexible and allows omission of certain parts. Mode strings are easily read and
written by humans. This makes the user interface of programs that use serial
ports easier to create and use: programs do not need to provide a plethora
of command line switches, but can accept one string and delegate handling
to `sers`.

Due to backwards compatibility there is a difference in data representation
between `SetMode` and `GetMode`.


Platforms 
---------

`sers` has been successfully used on Mac OS X, Linux and Windows.

The following known restrictions apply.

### Linux

- `cgo` is needed to compile `sers`
- Using non-traditional baud rates may not work. This depends on whether the
  system headers have correct definitions for `struct termios2`. Traditional
  baud rates use traditional termios baud rate setting methods via defines such
  as `B1200` etc. The author is successfully using non-traditional baud rates
  on Linux, though.

### OS X

- Calling `GetMode` before having called `SetMode` will result in an error.
- `cgo` is needed to compile `sers`

### Windows

- Only `NO_HANDSHAKE` is supported.


Feature ideas
-------------

### Not relying on `cgo`

Compilation for Linux and OS X involved use of `cgo` and thus a C compiler.
This makes cross compilation harder for every user of `sers`.

If all necessary definitions for serial port handling can be expressed in Go,
it is conceivable to make a pure-Go version for both these operation systems.

### Deadlines

Adding deadline methods, `SetDeadline` as well as `Set{Read,Write}Deadline`.

On termios platforms this relatively easy as we can let the runtime netpoller
do the heavy lifting. On Windows I do not yet understand the most promising
path to implementation.

### `net.Conn` support

Building on deadlines, `{Local,Remote}Addr` methods would enable a `SerialPort`
to implement `net.Conn`. This makes it easy to plug a `SerialPort` into various
network protocol implementations that build on byte streams.

### Regularize interface

With a working deadline implementation, the wart of `SetReadParams` could be
dropped or, at least, relegated to platform specific support. `{Get,Set}Mode`
can be made to both take a `Mode` struct. Representation of parity, stop bit
and handshake settings will be move to their own type to prevent mixups.

These changes would mean a new major version as they imply modifications that are
not backwards compatible.


Release History
---------------

### untagged

- add `GetMode()` to `SerialPort`
- add modestrings, `ParseModestring`
- add `Mode` struct

### v1.0.2

- termios platofrms: make `Close()` unblock readers

### v1.0.1

- use our own definition of `struct termios2` as the headers on a number of
  linux distributions are wrong.
- add `verification/setbaudrate`

### v1.0.0

`sers` appears in the primordial soup.
