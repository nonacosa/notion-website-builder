package gui

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"github.com/dgraph-io/badger/v3"
	"github.com/fyne-io/examples/gui/extend"
	"github.com/fyne-io/examples/notion/generator"
	"github.com/fyne-io/examples/storage"
	"log"
	"sort"
	"strings"
	"time"
)

func FetchNotion(id string, bdSite *badger.DB, website storage.Website, win fyne.Window, fetch bool, newWin bool) {
	log.Printf("click card is is > %s ", id)
	pageId := strings.Split(id, "_")[1]
	storage.Scan(pageId, func(items []string, bdItem *badger.DB) {
		log.Printf("scan item size is : %b ", len(items))
		if items == nil || fetch {
			if bdSite != nil {
				bdSite.Close()
			}

			bdItem.Close()
			var config generator.Config
			config.Notion.DatabaseID = website.PageID
			config.Notion.Secret = website.Secret
			preProgressBar(website, win, func(barData binding.ExternalFloat, msgData binding.ExternalString, window fyne.Window) {
				generator.Run(config, website, func(this float64, next float64, msg string) {
					time.Sleep(time.Duration(5) * time.Millisecond)
					msgData.Set(msg)
					barData.Set(this)
				})
				window.Close()
			})
			if newWin {
				storage.Scan(pageId, func(reItems []string, bdItem *badger.DB) {
					websiteItemWindow(website, win, reItems)
				})
			}
		} else {
			if newWin {
				websiteItemWindow(website, win, items)
			}
		}
	})
}

// mark card from local storage prefix [website.*]
func makeCardItem(win fyne.Window) *fyne.Container {
	var cards []fyne.CanvasObject

	storage.Scan("website", func(val []string, bdSite *badger.DB) {
		var websites []storage.Website
		for _, info := range val {
			var website storage.Website

			err := json.Unmarshal([]byte(info), &website)
			if err != nil {
				fmt.Println(err.Error())
			}
			websites = append(websites, website)
		}

		// sort by createTime
		sort.Slice(websites, func(i, j int) bool {
			return websites[i].CreateTime.UnixMilli() > websites[j].CreateTime.UnixMilli()
		})

		for _, website := range websites {
			card := extend.NewCard(website.Id, cardName(website.Name), website.Description, nil, func(id string) {
				FetchNotion(id, bdSite, website, win, false, true)
			})
			card.Image = canvas.NewImageFromResource(theme.FyneLogo())
			//card.Image.Resize()
			cards = append(cards, card)
		}
	})
	return container.NewGridWithColumns(3, cards...)
}

func cardName(origin string) string {
	var ret = origin
	cardLen := 12
	if len(origin) > cardLen {
		ret = origin[:12] + ".."
	} else {
		for i := 0; i < cardLen-len(origin); i++ {
			ret += " "
		}
	}
	return ret
}
