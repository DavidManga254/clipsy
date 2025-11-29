package clip

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/DavidManga254/clipsy/lib/logger"
	"golang.design/x/clipboard"
)

type Kind int

const (
    KindUnknown Kind = iota
    KindText
    KindImage
    KindFiles
)

//single clipboard item
type ClipboardItem struct {
    ID        string   
    Kind      Kind      
    Text      string    
    ImageData []byte    
    Files     []string 
    Source    string    // "local" or device ID in future
    CreatedAt time.Time
}

type Clipboard struct {
	mu sync.RWMutex

	initialized bool
	accessible bool

	current *ClipboardItem
	history []ClipboardItem
	maxHistory int

	events chan ClipboardItem

}

func InitClipboard() error {
	err := clipboard.Init()

	return err
}

func LatestClipboardtext() (string, error) {
	ClipboardInitErr := InitClipboard() 

	if ClipboardInitErr != nil {
		logger.Error("Error initializing clipboard package: %v", ClipboardInitErr)

		return "",ClipboardInitErr
	}

	clipboardText := clipboard.Read(clipboard.FmtText)
	clipboardTextToString := string(clipboardText)


	if clipboardText == nil || clipboardTextToString == "" {
		logger.Error("Current copy item is empty or not a text")

		return "",fmt.Errorf("current copy item is empty or not a text")
	}


	return clipboardTextToString,nil
}

func LatestClipboardImage() ([]byte, error) {
	ClipboardInitErr := InitClipboard() 

	if ClipboardInitErr != nil {
		logger.Error("Error initializing clipboard package: %v", ClipboardInitErr)

		return nil,ClipboardInitErr
	}

	clipboardImage := clipboard.Read(clipboard.FmtImage)

	if len(clipboardImage) == 0 {
		logger.Error("Current copy item is empty or not an image")

		return nil,fmt.Errorf("current copy item is empty or not an Image")
	}


	return clipboardImage,nil
}


func subscribeToClipboard(){
	ClipboardInitErr := InitClipboard() 

	if ClipboardInitErr != nil {
		logger.Error("Error initializing clipboard package: %v", ClipboardInitErr)

	}

	imageChannel := clipboard.Watch(context.TODO(),clipboard.FmtImage)
	textChannel := clipboard.Watch(context.TODO(),clipboard.FmtText)
	
}



/*
windows stores text and images directly in the clipboard when copied 
Other files on the other hand are stored as CF_HDROP format 
this is essentially a reference to the file like paths and other info and not the file itself exactly 

*/