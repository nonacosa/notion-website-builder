package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func welcomeScreen(win fyne.Window, main *fyne.Container) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneScene)
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(171, 125))
	} else {
		logo.SetMinSize(fyne.NewSize(228, 167))
	}

	return container.NewCenter(container.NewVBox(

		widget.NewLabelWithStyle("Welcome to the notion meta app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		container.NewHBox(
			widget.NewHyperlink("notionwb.com", parseURL("https://notionwb.com/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("documentation", parseURL("https://doc.notionwb.com/")),
			widget.NewLabel("-"),
			widget.NewHyperlink("sponsor", parseURL("https://notionwb.com/sponsor/")),
		),
	))
}
