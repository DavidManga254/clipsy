package clip

import (
	"context"
	"fmt"

	// "sync"
	"time"

	db "github.com/DavidManga254/clipsy/lib/DB"
	clipboarditemcontroller "github.com/DavidManga254/clipsy/lib/DB/controllers/clipBoardItemController"
	windowsClip "github.com/DavidManga254/clipsy/lib/clipboard/windows"
	"github.com/DavidManga254/clipsy/lib/logger"
	"github.com/google/uuid"
	"golang.design/x/clipboard"
)

// type Kind int

// const (
//     KindUnknown Kind = iota
//     KindText
//     KindImage
//     KindFiles
// )

// //single clipboard item
// type ClipboardItem struct {
//     ID        string
//     Kind      Kind
//     Text      string
//     ImageData []byte
//     Files     []string
//     Source    string    // "local" or device ID in future
//     CreatedAt time.Time
// }

// type Clipboard struct {
// 	mu sync.RWMutex

// 	initialized bool
// 	accessible bool

// 	current *ClipboardItem
// 	history []ClipboardItem
// 	maxHistory int

// 	events chan ClipboardItem

// }

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


func SubscribeToClipboard(ctx context.Context){
	ClipboardInitErr := InitClipboard() 

	if ClipboardInitErr != nil {
		logger.Error("Error initializing clipboard package: %v", ClipboardInitErr)

	}

	imageChannel := clipboard.Watch(ctx,clipboard.FmtImage)
	textChannel := clipboard.Watch(ctx,clipboard.FmtText)

	fileChannel := make(chan []string,8)
	go watchwindowsClipboard(ctx,fileChannel)


	go func() {
		for {
			select {
				case <-ctx.Done() :
					return
				case imageValue, ok := <- imageChannel : {
					if !ok {
						return
					}

					go processTextAndImage(imageValue,"IMAGE",ctx)
				}
				case textValue, ok := <-textChannel : {
					if !ok {
						return
					}

					go processTextAndImage(textValue,"TEXT",ctx)
				}
				case filePaths,ok := <-fileChannel : {
					if !ok {
						return
					}

					go processFile(filePaths,ctx)

				}
			}
		}
		
	}()
	
}


func processTextAndImage(value []byte, valueType string, ctx context.Context) error{
	logger.Info("%v",value)

	itemId := uuid.New().String()

	switch valueType {
		case "IMAGE" : {
			imageItem := clipboarditemcontroller.ClipboardItem{
				Kind: clipboarditemcontroller.KindImage,
				ImageData: value,
				Text: "",
				Files: make([]string, 0),
				ID: itemId,
				Source: "",
				CreatedAt: time.Now(),
			}

			err := clipboarditemcontroller.New(db.Database,clipboarditemcontroller.ClipboardItem(imageItem),ctx)

			if err != nil {
				logger.Error("%v",err)
			}

			return err

		}
		case "TEXT" : {
			textItem := clipboarditemcontroller.ClipboardItem{
				Kind: clipboarditemcontroller.KindText,
				ImageData: make([]byte, 0),
				Text: string(value),
				Files: make([]string, 0),
				ID: itemId,
				Source: "",
				CreatedAt: time.Now(),
			}

			err := clipboarditemcontroller.New(db.Database,textItem,ctx)

			if err != nil {
				logger.Error("%v",err)
			}

			return err
		}
	}

	return nil
}

func processFile(value []string, ctx context.Context) error{
	logger.Info("%v",value)
	
	fileItem := clipboarditemcontroller.ClipboardItem{
		Kind: clipboarditemcontroller.KindFiles,
		ImageData: make([]byte, 0),
		Text: "",
		Files: value,
		ID: uuid.New().String(),
		Source: "",
		CreatedAt: time.Now(),
		
	}

	err := clipboarditemcontroller.New(db.Database,fileItem,ctx)

	return err
}

func watchwindowsClipboard(ctx context.Context, out chan<- []string){
	//poll every 300ms
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	//last polled file list
	var last []string

	for {
		select {
			case <-ctx.Done(): {
				return
			}
			case <-ticker.C : {
				copiedFilePaths, err := windowsClip.ReadWindowsClipboard()

				if err != nil {
					//TODO check how to work on specific error issues 
					continue
				}

				if len(copiedFilePaths) > 0 && !equalStringSlices(last,copiedFilePaths){
					last = append([]string(nil), copiedFilePaths...)        // copy for last
    				out <- append([]string(nil), copiedFilePaths...) 
				}
			}
		}
	}
}

func equalStringSlices(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}



/*
windows stores text and images directly in the clipboard when copied 
Other files on the other hand are stored as CF_HDROP format 
this is essentially a reference to the file like paths and other info and not the file itself exactly 

*/
