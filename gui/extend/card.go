package extend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"log"
)

type MyCard struct {
	widget.Card
	Id       string
	OnTapped func(string) `json:"-"`
}

func (m *MyCard) TappedSecondary(*fyne.PointEvent) {
	log.Println("Right Click")
}

func (b *MyCard) Tapped(*fyne.PointEvent) {
	b.Refresh()

	if b.OnTapped != nil {
		b.OnTapped(b.Id)
	}
}

func NewCard(id, title, subtitle string, content fyne.CanvasObject, tapped func(string)) *MyCard {
	ret := &MyCard{}
	ret.ExtendBaseWidget(ret)
	ret.OnTapped = tapped
	ret.Title = title
	ret.Subtitle = subtitle
	ret.Content = content
	ret.Id = id
	return ret
}

// Layout the components of the card container.
func Layout(size fyne.Size) {

}
