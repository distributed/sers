package sers

import "testing"

func TestParseModestring(t *testing.T) {
	cases := []struct {
		Str    string
		Parses bool
		Mode   Mode
	}{
		// all none default
		{"5000,5e2,rtscts", true, Mode{5000, 5, E, 2, RTSCTS_HANDSHAKE}},
		// all none default
		{"115200,7n2,rtscts", true, Mode{115200, 7, N, 2, RTSCTS_HANDSHAKE}},
		// only decimal baudrates
		{"0x5,8n1,rtscts", false, Mode{}},
		// a standard case
		{"57600,8n1", true, Mode{57600, 8, N, 1, NO_HANDSHAKE}},
		// 8 bit default
		{"20,n1,rtscts", true, Mode{20, 8, N, 1, RTSCTS_HANDSHAKE}},
		// 8 bit, N default
		{"25,2,rtscts", true, Mode{25, 8, N, 2, RTSCTS_HANDSHAKE}},
		// no handshake default
		{"112,8o2", true, Mode{112, 8, O, 2, NO_HANDSHAKE}},
		// parity can be left at default
		{"115,82,rtscts", true, Mode{115, 8, N, 2, RTSCTS_HANDSHAKE}},
		// stop bits can be left at default
		{"20,7e,rtscts", true, Mode{20, 7, E, 1, RTSCTS_HANDSHAKE}},
		// alternate default handshake
		{"20,7e,", true, Mode{20, 7, E, 1, NO_HANDSHAKE}},
		// unusual, but possible: just set the stopbits
		{"37,2", true, Mode{37, 8, N, 2, NO_HANDSHAKE}},
		// 0 instead of O
		{"21,801", false, Mode{}},
		// 6 databits
		{"52,6e1", true, Mode{52, 6, E, 1, NO_HANDSHAKE}},
		// wrong handshake format
		{"54,8n1,invalid", false, Mode{}},
		// all default framing format
		{"66,,rtscts", true, Mode{66, 8, N, 1, RTSCTS_HANDSHAKE}},
	}

	for i, c := range cases {
		mode, err := ParseModestring(c.Str)
		if c.Parses && err != nil {
			t.Errorf("case %d: expected to parse, but errors with %q", i, err)
			continue
		} else if !c.Parses && err == nil {
			t.Errorf("case %d: expected to error, but parses as %v", i, mode)
			continue
		}

		if !c.Parses && err != nil {
			t.Logf("str %q error %s", c.Str, err)
			continue
		}

		if mode != c.Mode {
			t.Errorf("case %d: got %v, want %v", i, mode, c.Mode)
			continue
		}
	}
}

func TestModestringStringMethod(t *testing.T) {
	cases := []struct {
		Mode Mode
		Str  string
	}{
		{Mode{2400, 5, E, 2, RTSCTS_HANDSHAKE}, "2400,5e2,rtscts"},
		{Mode{1200, 6, 20, 1, NO_HANDSHAKE}, "invalid_mode(1200,6,20,1,0)"},
		{Mode{4800, 6, N, 1, NO_HANDSHAKE}, "4800,6n1,none"},
		{Mode{9600, 7, O, 2, RTSCTS_HANDSHAKE}, "9600,7o2,rtscts"},
		{Mode{19200, 8, N, 1, NO_HANDSHAKE}, "19200,8n1,none"},
	}

	for i, c := range cases {
		act := c.Mode.String()
		exp := c.Str
		if act != exp {
			t.Errorf("case %d: got %q, expected %q", i, act, exp)
			continue
		}
	}
}
