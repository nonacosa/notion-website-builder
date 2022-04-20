package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"log"
	"strings"
)

func makeFormItem(_ fyne.Window) []*widget.FormItem {
	name := widget.NewEntry()
	name.SetPlaceHolder("your website name")

	notion := widget.NewEntry()
	notion.SetPlaceHolder("your notion secret")

	pageId := widget.NewEntry()
	pageId.SetPlaceHolder("your notion meta database page")

	//disabled := widget.NewRadioGroup([]string{"Option 1", "Option 2"}, func(string) {})
	//disabled.Horizontal = true
	//disabled.Disable()

	var selects []string
	// todo test

	files, err := ioutil.ReadDir("/Users/nonacosa/.notion-wb/theme")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// filter system file or dir
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		selects = append(selects, f.Name())
	}

	selectEntry := widget.NewSelect(selects, func(s string) { fmt.Println("selected theme", s) })

	largeText := widget.NewMultiLineEntry()

	items := []*widget.FormItem{
		{Text: "Name", Widget: name, HintText: "Your full name"},
		{Text: "Secret", Widget: notion, HintText: "A valid email address"},
		{Text: "PageId", Widget: pageId},
		//{Text: "Disabled", Widget: disabled},
		{Text: "Theme", Widget: selectEntry},
		{Text: "Description", Widget: largeText},
	}

	return items
}
