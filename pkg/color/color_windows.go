package color

import (
	"os"
	"runtime"
	"syscall"
)

func init() {
	if runtime.GOOS == "windows" {
		// do we have working ANSI
		handle := syscall.Handle(os.Stdout.Fd())
		kernel32DLL := syscall.NewLazyDLL("kernel32.dll")
		setConsoleModeProc := kernel32DLL.NewProc("SetConsoleMode")

		// fallback to no colors if not
		if _, _, err := setConsoleModeProc.Call(uintptr(handle), 0x0001|0x0002|0x0004); err != nil && err.Error() != "The operation completed successfully." {
			Reset = ""
			Red = ""
			Yellow = ""
			Green = ""
			Blue = ""
		}
	}
}
