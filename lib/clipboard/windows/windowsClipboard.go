package windowsClip

import (
	"fmt"
	"unsafe"

	"github.com/DavidManga254/clipsy/lib/logger"
	"golang.org/x/sys/windows"
)

var (
    user32   = windows.NewLazySystemDLL("user32.dll")
    kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	shell32  = windows.NewLazySystemDLL("shell32.dll")
)

var (
    procOpenClipboard    = user32.NewProc("OpenClipboard")
    procCloseClipboard   = user32.NewProc("CloseClipboard")
    procGetClipboardData = user32.NewProc("GetClipboardData")
    procGlobalLock       = kernel32.NewProc("GlobalLock")
    procGlobalUnlock     = kernel32.NewProc("GlobalUnlock")
    procDragQueryFileW   = shell32.NewProc("DragQueryFileW")
)


func ReadWindowsClipboard() ([]string, error) {
	r1, _, _ := procOpenClipboard.Call(0)
	logger.Info("%v",r1)

	if r1 == 0 {
		logger.Error("Failed to open clipboard")

		return nil, fmt.Errorf("failed to open clipboard")
	}


	
	// uintptr is an integer type that is large enough to hold the bit pattern of
	// any pointer.
	//r2 is  A handle to a global memory block containing the DROPFILES structure + file list.
	fileGlobalMemLock, _, _ := procGetClipboardData.Call(15)

	if fileGlobalMemLock == 0 {
		logger.Error("Current copied item in the clipboard is not a file(s)")

		return nil,fmt.Errorf("current copied item in the clipboard is not a file(s)")
	}

	logger.Info("%v",fileGlobalMemLock)


	filesAddress,_,_ := procGlobalLock.Call(fileGlobalMemLock)

	if filesAddress == 0 {
		logger.Error("Failed to get memory addres of copied files ")

		return nil,fmt.Errorf("failed to get memory addres of copied files ")

	}


	filecount,_,_ := procDragQueryFileW.Call(
		filesAddress,
		uintptr(0xFFFFFFFF),
		0,
		0,
	)


	filepaths := make([]string,0)

	for i := 0; i < int(filecount); i++ {
		fileNameLength,_,_ := procDragQueryFileW.Call(
			filesAddress,
			uintptr(i),
			0,
			0,
		)

		fileNameBuffer := make([]uint16, fileNameLength+1)

		_,_,_ = procDragQueryFileW.Call(
			filesAddress,
			uintptr(i),
			uintptr(unsafe.Pointer(&fileNameBuffer[0])),
			uintptr(fileNameLength + 1),
		)

		filepaths = append(filepaths,windows.UTF16ToString(fileNameBuffer))
	}
	

	defer func() {
		//close clipboard
		procCloseClipboard.Call()

		//unlock global
    	procGlobalUnlock.Call(fileGlobalMemLock)
	}()


	return filepaths,nil
}