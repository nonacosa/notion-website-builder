package gui

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/fyne-io/examples/storage"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func preProgressBar(website storage.Website, win fyne.Window, bind func(binding.ExternalFloat, binding.ExternalString, fyne.Window)) {
	f := 0.0
	text := "begin fetch from notion ..."
	data := binding.BindFloat(&f)
	msg := binding.BindString(&text)
	label := widget.NewLabelWithData(msg)
	//entry := widget.NewEntryWithData(binding.FloatToString(data))
	floats := container.NewGridWithColumns(2, label)

	//slide := widget.NewSliderWithData(0, 1, data)
	//slide.Step = 0.01
	bar := widget.NewProgressBarWithData(data)

	item := container.NewVBox(floats, bar, widget.NewSeparator())

	w := fyne.CurrentApp().NewWindow(website.Name)
	w.SetContent(item)
	w.Resize(fyne.NewSize(480, 120))
	w.SetFixedSize(false)
	w.CenterOnScreen()
	w.Show()

	bind(data, msg, w)
}

func websiteItemWindow(website storage.Website, win fyne.Window, items []string) {
	w := fyne.CurrentApp().NewWindow(website.Name)
	w.SetContent(makeWebsiteItemList(win, items, website))
	w.Resize(fyne.NewSize(720, 540))
	w.SetFixedSize(false)
	w.CenterOnScreen()
	w.Show()
}

func windowScreen(_ fyne.Window, main *fyne.Container) fyne.CanvasObject {
	windowGroup := container.NewVBox(
		widget.NewButton("New window", func() {
			w := fyne.CurrentApp().NewWindow("Hello")
			w.SetContent(widget.NewLabel("Hello World!"))
			w.Show()
		}),
		widget.NewButton("Fixed size window", func() {
			w := fyne.CurrentApp().NewWindow("Fixed")
			w.SetContent(fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewLabel("Hello World!")))

			w.Resize(fyne.NewSize(240, 180))
			w.SetFixedSize(true)
			w.Show()
		}),
		widget.NewButton("Toggle between fixed/not fixed window size", func() {
			w := fyne.CurrentApp().NewWindow("Toggle fixed size")
			w.SetContent(fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewCheck("Fixed size", func(toggle bool) {
				if toggle {
					w.Resize(fyne.NewSize(240, 180))
				}
				w.SetFixedSize(toggle)
			})))
			w.Show()
		}),
		widget.NewButton("Centered window", func() {
			w := fyne.CurrentApp().NewWindow("Central")
			w.SetContent(fyne.NewContainerWithLayout(layout.NewCenterLayout(), widget.NewLabel("Hello World!")))
			w.CenterOnScreen()
			w.Show()
		}))

	drv := fyne.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		windowGroup.Objects = append(windowGroup.Objects,
			widget.NewButton("Splash Window (only use on start)", func() {
				w := drv.CreateSplashWindow()
				w.SetContent(widget.NewLabelWithStyle("Hello World!\n\nMake a splash!",
					fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
				w.Show()

				go func() {
					time.Sleep(time.Second * 3)
					w.Close()
				}()
			}))
	}

	otherGroup := widget.NewCard("Other", "",
		widget.NewButton("Notification", func() {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Fyne Demo",
				Content: "Testing notifications...",
			})
		}))

	return container.NewVBox(widget.NewCard("Windows", "", windowGroup), otherGroup)
}
