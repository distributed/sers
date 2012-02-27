// +build linux
package sers

func (bp *baseport) SetBaudRate(br int) error {
	panic("not implemented yet on linux")
}

func (bp *baseport) ClearNonBlocking() error {
	panic("not implemented yet on linux")
}
