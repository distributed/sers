package sers

import (
	"fmt"
	"strconv"
	"strings"
)

// Parses a modestring like "115200,8n1,rtscts" into a struct Mode. The format
// is baudrate,framestring,handshake. Either the handshake part or both the
// framestring and handshake parts can be omitted. For the omitted parts,
// defaults of 8 data bits, no parity, 1 stop bit and no handshaking will be
// filled in. The framestring consists of a sequence databits, parity, stopbits.
// Any or all of the three components can be left out. Non-specified parts will
// take the default values mentioned before.
//
// Valid choices for databits are [5, 6, 7, 8], for parity it is [n, o, e] and
// for stopbits it's [1, 2]. Valid choices for the handshake parts are "", "none"
// and "rtscts". The function is not case sensitive.
//
// A couple of examples:
//
// 		60000	         - 60000 baud, 8 data bits, no parity, 1 stopbit, no handshake
//		115200,8e1       - 115200 baud, 8 data bits, even parity, 1 stopbit, no handshake
// 		57600,8o1,rtscts - 57600 baud, 8 data bits, odd parity, 1 stopbit, rts/cts handshake
//		19200,72		 - 19200 baud, 7 data bits, no parity, 2 stop bits, no handshake
//		9600,2,rtscts    - 9600 baud, 8 data bits, no parity, 2 stop bits, rts/cts handshake
func ParseModestring(s string) (Mode, error) {
	var mode Mode

	commaparts := strings.Split(strings.ToUpper(s), ",")

	brpart := commaparts[0]
	br64, err := strconv.ParseUint(brpart, 10, 32)
	if err != nil {
		return mode, fmt.Errorf("modestring %q cannot parse baudrate: %v", s, err)
	}

	mode.Baudrate = int(br64)

	mode.DataBits = 8
	mode.Parity = N
	mode.Stopbits = 1
	if len(commaparts) >= 2 {
		framepart := commaparts[1]
		idx := 0
		lastidx := len(framepart) - 1

		if idx <= lastidx {
			switch framepart[idx] {
			case '5':
				mode.DataBits = 5
				idx++
			case '6':
				mode.DataBits = 6
				idx++
			case '7':
				mode.DataBits = 7
				idx++
			case '8':
				mode.DataBits = 8
				idx++
			}
		}

		if idx <= lastidx {
			switch framepart[idx] {
			case 'N':
				mode.Parity = N
				idx++
			case 'O':
				mode.Parity = O
				idx++
			case 'E':
				mode.Parity = E
				idx++
			}
		}

		if idx <= lastidx {
			switch framepart[idx] {
			case '1':
				mode.Stopbits = 1
				idx++
			case '2':
				mode.Stopbits = 2
				idx++
			}
		}

		if idx <= lastidx {
			return mode, fmt.Errorf("cannot parse serial framing format %q in %q: unknown sequence %q", framepart, s, framepart[idx:])
		}
	}

	mode.Handshake = NO_HANDSHAKE
	if len(commaparts) >= 3 {
		hspart := commaparts[2]
		switch hspart {
		case "NONE", "":
			mode.Handshake = NO_HANDSHAKE
		case "RTSCTS":
			mode.Handshake = RTSCTS_HANDSHAKE
		default:
			return mode, fmt.Errorf("cannot parse serial modestring %q: unknown handshake format %q", s, hspart)
		}
	}

	return mode, nil
}
