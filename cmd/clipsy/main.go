package main

import (
	windowsClip "github.com/DavidManga254/clipsy/lib/clipboard/windows"
	"github.com/DavidManga254/clipsy/lib/logger"
)

func main() {
	logger.Info("Starting up Clipsy oe.....")

	windowsClip.ReadWindowsClipboard()

}