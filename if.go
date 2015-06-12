package tuntap

import (
	"syscall"
	"time"

	"github.com/jaracil/poll"
)

// Interface is a TUN/TAP interface.
type Interface struct {
	isTAP bool
	file  *poll.File
	name  string
}

var ErrTimeout = poll.ErrTimeout
var ErrClosed = poll.ErrClosed

// Create a new TAP interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func NewTAP(ifName string) (ifce *Interface, err error) {
	fd, err := syscall.Open("/dev/net/tun", syscall.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	name, err := createInterface(uintptr(fd), ifName, cIFF_TAP|cIFF_NO_PI)
	if err != nil {
		return nil, err
	}
	file, err := poll.NewFile(uintptr(fd), "/dev/net/tun")
	if err != nil {
		return nil, err
	}
	ifce = &Interface{isTAP: true, file: file, name: name}
	return
}

// Create a new TUN interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func NewTUN(ifName string) (ifce *Interface, err error) {
	fd, err := syscall.Open("/dev/net/tun", syscall.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	name, err := createInterface(uintptr(fd), ifName, cIFF_TUN|cIFF_NO_PI)
	if err != nil {
		return nil, err
	}
	file, err := poll.NewFile(uintptr(fd), "/dev/net/tun")
	if err != nil {
		return nil, err
	}
	ifce = &Interface{isTAP: false, file: file, name: name}
	return
}

// Returns true if ifce is a TUN interface, otherwise returns false;
func (ifce *Interface) IsTUN() bool {
	return !ifce.isTAP
}

// Returns true if ifce is a TAP interface, otherwise returns false;
func (ifce *Interface) IsTAP() bool {
	return ifce.isTAP
}

// Returns the interface name of ifce, e.g. tun0, tap1, etc..
func (ifce *Interface) Name() string {
	return ifce.name
}

// Implement io.Writer interface.
func (ifce *Interface) Write(p []byte) (int, error) {
	return ifce.file.Write(p)
}

// Implement io.Reader interface.
func (ifce *Interface) Read(p []byte) (int, error) {
	return ifce.file.Read(p)

}

// Close closes interface. All goroutines blocked in read/write operations are awakened.
func (ifce *Interface) Close() error {
	return ifce.file.Close()
}

// SetDeadLine sets read and write deadline time.
func (ifce *Interface) SetDeadLine(t time.Time) error {
	return ifce.file.SetDeadline(t)
}

// SetReadDeadLine sets read deadline time.
func (ifce *Interface) SetReadDeadLine(t time.Time) error {
	return ifce.file.SetReadDeadline(t)
}

// SetWriteDeadLine sets write deadline time.
func (ifce *Interface) SetWriteLine(t time.Time) error {
	return ifce.file.SetWriteDeadline(t)
}
