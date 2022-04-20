package gui

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dgraph-io/badger/v3"
	"github.com/fyne-io/examples/notion/generator"
	"github.com/fyne-io/examples/storage"
	"github.com/skratchdot/open-golang/open"
	"log"
	"strconv"
)

func makeWebsiteItemList(win fyne.Window, items []string, website storage.Website) fyne.CanvasObject {

	data := make([]string, len(items))
	var websiteItems []storage.WebsiteItem
	webItemList := binding.BindStringList(&data)
	for i := range data {
		data[i] = "Test Item " + strconv.Itoa(i)
	}
	for k, item := range items {
		var websiteItem storage.WebsiteItem
		err := json.Unmarshal([]byte(item), &websiteItem)
		if err != nil {
			fmt.Println(err.Error())
		}
		data[k] = websiteItem.Name
		websiteItems = append(websiteItems, websiteItem)

	}

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	hbox := container.NewHBox(icon, label)

	list := widget.NewListWithData(
		webItemList,
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(text binding.DataItem, item fyne.CanvasObject) {
			t, _ := text.(binding.String).Get()
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(t)
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		log.Printf("id is > %d", id)
		label.SetText(fmt.Sprintf("%s \n \n 发布状态：%s ", websiteItems[id].Name, websiteItems[id].Status))
		icon.SetResource(theme.DocumentIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	list.Select(125)

	tabs := container.NewAppTabs(
		container.NewTabItem("Tab 1", widget.NewLabel("Content of tab 1")),
		container.NewTabItem("Tab 2 bigger", widget.NewLabel("Content of tab 2")),
		container.NewTabItem("Tab 3", widget.NewLabel("Content of tab 3")),
	)
	for i := 4; i <= 12; i++ {
		tabs.Append(container.NewTabItem(fmt.Sprintf("Tab %d", i), widget.NewLabel(fmt.Sprintf("Content of tab %d", i))))
	}
	locations := makeTabLocationSelect(tabs.SetTabLocation)

	buttons := container.NewGridWithColumns(3,
		container.NewVBox(widget.NewButton("fresh website", func() {
			FetchNotion(website.Id, nil, website, win, true, false)
			storage.Scan(website.PageID, func(reItems []string, bdItem *badger.DB) {
				data := make([]string, len(items)+1)
				var newWebsiteItems []storage.WebsiteItem

				for i := range data {
					data[i] = "Test Item " + strconv.Itoa(i)
				}
				for k, item := range items {
					var websiteItem storage.WebsiteItem
					err := json.Unmarshal([]byte(item), &websiteItem)
					if err != nil {
						fmt.Println(err.Error())
					}
					data[k] = websiteItem.Name
					websiteItems = newWebsiteItems
				}
				data[4] = "test"
				webItemList.Set(data)

			})
		})),
		container.NewVBox(widget.NewButton("build website", func() {
			generator.BuildLocal(website)
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "notion-hugo",
				Content: "build succeed !",
			})

		})),
		container.NewVBox(widget.NewButton("open in finder", func() {
			log.Printf(" open in finder button ... ")
			//open.Run("https://google.com/")
			open.Run(storage.GetSavePath(website.Name))

		})),
	)
	//  widget.NewSeparator()
	buttonBox := container.NewGridWithColumns(2, locations, buttons)

	return container.NewBorder(buttonBox, nil, nil, nil, container.NewHSplit(list, container.NewCenter(hbox)))
	//return container.NewHSplit(list, container.NewCenter(hbox))
}

func makeTabLocationSelect(callback func(container.TabLocation)) *widget.Select {
	locations := widget.NewSelect([]string{"Top", "Bottom", "Leading", "Trailing"}, func(s string) {
		callback(map[string]container.TabLocation{
			"Top":      container.TabLocationTop,
			"Bottom":   container.TabLocationBottom,
			"Leading":  container.TabLocationLeading,
			"Trailing": container.TabLocationTrailing,
		}[s])
	})
	locations.SetSelected("Top")
	return locations
}
