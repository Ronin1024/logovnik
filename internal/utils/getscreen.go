package utils

import "syscall"

var (
	user32           = syscall.NewLazyDLL("User32.dll")
	getSystemMetrics = user32.NewProc("GetSystemMetrics")
)

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

func GetScreenResolution() (int, int) {
	width, _, _ := getSystemMetrics.Call(SM_CXSCREEN)
	height, _, _ := getSystemMetrics.Call(SM_CYSCREEN)
	return int(width), int(height)
}
