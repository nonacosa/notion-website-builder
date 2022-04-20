package gui

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/dgraph-io/badger/v3"
	"github.com/fyne-io/examples/notion/generator"
	"github.com/fyne-io/examples/storage"
	"log"
	"net/url"
	"time"
)

const PageId = "f339b27be68f4f30850d75f7675e6beb"

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func websiteScreen(win fyne.Window, main *fyne.Container) fyne.CanvasObject {
	cards := makeCardItem(win)

	buttons := container.NewAdaptiveGrid(4,
		container.NewHBox(widget.NewButton("Add WebSite", func() {
			items := makeFormItem(win)
			dialog.ShowForm("Website Setting", "Confirm", "Cancel", items, func(b bool) {

				if !b {
					return
				}
				name := items[0].Widget.(*widget.Entry).Text
				//secret := items[1].Widget.(*widget.Entry).Text
				pageId := items[2].Widget.(*widget.Entry).Text
				theme := items[3].Widget.(*widget.Select).Selected
				description := items[4].Widget.(*widget.Entry).Text

				website := storage.Website{
					//Id:          storage.TestKey,
					Id:     "website_" + pageId,
					Secret: "secret_f4YXEox1cwoP4qjB2gBMQN4PZLpCeylBSoFVTxDFs0J",
					PageID: pageId,
					//PageID:      storage.PageId,
					Name:        name,
					Theme:       theme,
					Description: description,
					CreateTime:  time.Now(),
					// todo image
				}

				json, _ := json.Marshal(website)
				// todo: create notion setting page
				log.Printf("create notion setting page ... %s\n", string(json))
				storage.Save("website_"+pageId, string(json))
				//storage.Save(storage.TestKey, string(json))
				//websiteScreen(win)
				// refresh card view

				main.Objects = []fyne.CanvasObject{websiteScreen(win, main)}
				main.Refresh()
			}, win)
		})),

		container.NewVBox(widget.NewButton("Website Test", func() {
			items := makeFormItem(win)
			dialog.ShowForm("Website Test", "Confirm", "Cancel", items, func(b bool) {
				if !b {
					return
				}
				storage.Scan(PageId, func(items []string, badger *badger.DB) {
					badger.Close()
					log.Printf("scan item size is : %b ", len(items))
					//if items == nil {
					var config generator.Config
					config.Notion.Secret = "secret_f4YXEox1cwoP4qjB2gBMQN4PZLpCeylBSoFVTxDFs0J"
					config.Notion.DatabaseID = PageId
					generator.Run(config, storage.Website{Name: "Test", Secret: "secret_f4YXEox1cwoP4qjB2gBMQN4PZLpCeylBSoFVTxDFs0J"}, func(this float64, next float64, msg string) {

					})

				})
			}, win)
		})),
		container.NewVBox(widget.NewButton("Test Bar", func() {
			preProgressBar(storage.Website{}, win, func(barData binding.ExternalFloat, msgData binding.ExternalString, window fyne.Window) {
				for i := 0; i < 100; i++ {
					time.Sleep(time.Duration(50) * time.Millisecond)
					num, _ := barData.Get()
					barData.Set(num + 0.01)
				}
				window.Close()
			})
		})),
	)

	buttonBox := container.NewVBox(buttons, widget.NewSeparator())

	return container.NewVBox(buttonBox, container.NewCenter(cards))
}
