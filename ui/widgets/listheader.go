package widgets

import (
	"supersonic/ui/layouts"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ListColumn struct {
	Text          string
	AlignTrailing bool
}

type ListHeader struct {
	widget.BaseWidget

	columns       []ListColumn
	columnsLayout *layouts.ColumnsLayout

	columnsContainer *fyne.Container
	container        *fyne.Container
}

func NewListHeader(cols []ListColumn, layout *layouts.ColumnsLayout) *ListHeader {
	l := &ListHeader{
		columns:          cols,
		columnsLayout:    layout,
		columnsContainer: container.New(layout),
	}
	l.container = container.NewMax(canvas.NewRectangle(theme.BackgroundColor()), l.columnsContainer)
	l.ExtendBaseWidget(l)
	l.buildColumns()
	return l
}

func (l *ListHeader) SetColumnVisible(colNum int, visible bool) {
	if visible {
		l.columnsContainer.Objects[colNum].Show()
	} else {
		l.columnsContainer.Objects[colNum].Hide()
	}
}

func (l *ListHeader) buildColumns() {
	for _, c := range l.columns {
		t := widget.NewRichTextWithText(c.Text)
		t.Segments[0].(*widget.TextSegment).Style.TextStyle.Bold = true
		al := fyne.TextAlignLeading
		if c.AlignTrailing {
			al = fyne.TextAlignTrailing
		}
		t.Segments[0].(*widget.TextSegment).Style.Alignment = al
		l.columnsContainer.Add(t)
	}
}

func (l *ListHeader) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(l.container)
}
