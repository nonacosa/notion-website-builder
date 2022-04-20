package generator

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/fyne-io/examples/file"
	"github.com/fyne-io/examples/storage"
	"github.com/gohugoio/hugo/commands"
	"github.com/skratchdot/open-golang/open"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var counter = 0

func BuildLocal(website storage.Website) {
	storage.Scan(website.PageID, func(items []string, bdItem *badger.DB) {
		log.Printf("scan item size is : %b ", len(items))
		if items != nil {
			for i := 0; i < len(items); i++ {
				var item storage.WebsiteItem
				err := json.Unmarshal([]byte(items[i]), &item)
				if err != nil {
					fmt.Println(err.Error())
				}
				if item.FrontMatter.Position != nil {
					fileName := strings.ToLower(item.Name) + ".md"
					// todo not empty verify
					toMovePosition := item.FrontMatter.Position
					toMovePath := storage.HugoTheme(website.Theme, "exampleSite", "content", toMovePosition.(string))
					filePath := filepath.Join(storage.MdSavePath(website.Name), fileName)

					if err := file.Copy(filePath, filepath.Join(toMovePath, fileName)); err != nil {
						fmt.Errorf("move source file to theme path error : %w", err)
					}
				} else {
					file.CopyDirectory(storage.MdSavePath(website.Name), storage.HugoPost(website.Theme))
				}
				file.CopyDirectory(storage.MdImageSavePath(website.Name), storage.HugoStatic(website.Theme))

			}

		}
	})
	Build(website)
}

func Build(website storage.Website) {
	resp := commands.Execute([]string{
		"--source", storage.HugoSource(website.Theme),
		"--destination", storage.HugoDest(website.Theme),
		"--themesDir", storage.HugoTheme()},
	)
	fmt.Println(resp)
	counter++
	if counter == 1 {
		go func() {
			//https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve
			http.Handle("/", http.FileServer(http.Dir(storage.HugoDest(website.Theme))))
			open.Run("http://localhost:3000")
			http.ListenAndServe(":3000", nil)

		}()
	}
}
