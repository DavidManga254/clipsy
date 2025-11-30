package windowsClip

import (
	"fmt"
	"unsafe"

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

const (
	cfHDrop      = 15         
	maxUint32Idx = 0xFFFFFFFF  
)

// ReadWindowsClipboard reads file paths from the Windows clipboard when the
// current clipboard content is CF_HDROP (file drop). It returns:
//
// - (nil, nil) if the clipboard does not currently contain files
// - ([]string, nil) with file paths if it does
// - (nil, error) on a real error (open clipboard, lock failure, etc.)
func ReadWindowsClipboard() ([]string, error) {
	// Open the clipboard for this process (hwnd = 0)
	r1, _, err := procOpenClipboard.Call(0)
	if r1 == 0 {
		// Failed to open clipboard at all
		if err != nil {
			return nil, fmt.Errorf("OpenClipboard failed: %w", err)
		}
		return nil, fmt.Errorf("OpenClipboard failed")
	}
	defer procCloseClipboard.Call()

	// Get handle to CF_HDROP data
	hDrop, _, _ := procGetClipboardData.Call(cfHDrop)
	if hDrop == 0 {
		// Clipboard doesn't contain CF_HDROP normal case
		return nil, nil
	}

	// Lock the global memory handle to get a pointer to DROPFILES + file list
	ptr, _, err := procGlobalLock.Call(hDrop)
	if ptr == 0 {
		if err != nil {
			return nil, fmt.Errorf("GlobalLock failed: %w", err)
		}
		return nil, fmt.Errorf("GlobalLock failed")
	}
	defer procGlobalUnlock.Call(hDrop)

	// First call to DragQueryFileW with index = 0xFFFFFFFF gets the file count
	fileCount, _, _ := procDragQueryFileW.Call(
		ptr,
		maxUint32Idx,
		0,
		0,
	)
	if fileCount == 0 {
		return nil, nil
	}

	paths := make([]string, 0, int(fileCount))

	// For each file, query its length then its path
	for i := 0; i < int(fileCount); i++ {
		nameLen, _, _ := procDragQueryFileW.Call(
			ptr,
			uintptr(i),
			0,
			0,
		)
		if nameLen == 0 {
			continue
		}

		buf := make([]uint16, nameLen+1)

		_, _, _ = procDragQueryFileW.Call(
			ptr,
			uintptr(i),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(nameLen+1),
		)

		paths = append(paths, windows.UTF16ToString(buf))
	}

	return paths, nil
}
