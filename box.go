package tview

import (
	"github.com/gdamore/tcell"
)

// Box implements Primitive with a background and optional elements such as a
// border and a title. Most subclasses keep their content contained in the box
// but don't necessarily have to.
//
// See https://github.com/rivo/tview/wiki/Box for an example.
type Box struct {
	// The position of the rect.
	x, y, width, height int

	// Border padding.
	paddingTop, paddingBottom, paddingLeft, paddingRight int

	// The box's background color.
	backgroundColor tcell.Color

	// Whether or not a border is drawn, reducing the box's space for content by
	// two in width and height.
	border bool

	// The color of the border.
	borderColor tcell.Color

	// The title. Only visible if there is a border, too.
	title string

	// The color of the title.
	titleColor tcell.Color

	// The alignment of the title.
	titleAlign int

	// Provides a way to find out if this box has focus. We always go through
	// this interface because it may be overridden by implementing classes.
	focus Focusable

	// Whether or not this box has focus.
	hasFocus bool
}

// NewBox returns a Box without a border.
func NewBox() *Box {
	b := &Box{
		width:           15,
		height:          10,
		backgroundColor: Styles.PrimitiveBackgroundColor,
		borderColor:     Styles.BorderColor,
		titleColor:      Styles.TitleColor,
		titleAlign:      AlignCenter,
	}
	b.focus = b
	return b
}

// SetBorderPadding sets the size of the borders around the box content.
func (b *Box) SetBorderPadding(top, bottom, left, right int) *Box {
	b.paddingTop, b.paddingBottom, b.paddingLeft, b.paddingRight = top, bottom, left, right
	return b
}

// GetRect returns the current position of the rectangle, x, y, width, and
// height.
func (b *Box) GetRect() (int, int, int, int) {
	return b.x, b.y, b.width, b.height
}

// GetInnerRect returns the position of the inner rectangle, without the border
// and without any padding.
func (b *Box) GetInnerRect() (int, int, int, int) {
	x, y, width, height := b.GetRect()
	if b.border {
		x++
		y++
		width -= 2
		height -= 2
	}
	return x + b.paddingLeft,
		y + b.paddingTop,
		width - b.paddingLeft - b.paddingRight,
		height - b.paddingTop - b.paddingBottom
}

// SetRect sets a new position of the rectangle.
func (b *Box) SetRect(x, y, width, height int) {
	b.x = x
	b.y = y
	b.width = width
	b.height = height
}

// InputHandler returns nil.
func (b *Box) InputHandler() func(event *tcell.EventKey, setFocus func(p Primitive)) {
	return nil
}

// SetBackgroundColor sets the box's background color.
func (b *Box) SetBackgroundColor(color tcell.Color) *Box {
	b.backgroundColor = color
	return b
}

// SetBorder sets the flag indicating whether or not the box should have a
// border.
func (b *Box) SetBorder(show bool) *Box {
	b.border = show
	return b
}

// SetBorderColor sets the box's border color.
func (b *Box) SetBorderColor(color tcell.Color) *Box {
	b.borderColor = color
	return b
}

// SetTitle sets the box's title.
func (b *Box) SetTitle(title string) *Box {
	b.title = title
	return b
}

// SetTitleColor sets the box's title color.
func (b *Box) SetTitleColor(color tcell.Color) *Box {
	b.titleColor = color
	return b
}

// SetTitleAlign sets the alignment of the title, one of AlignLeft, AlignCenter,
// or AlignRight.
func (b *Box) SetTitleAlign(align int) *Box {
	b.titleAlign = align
	return b
}

// Draw draws this primitive onto the screen.
func (b *Box) Draw(screen tcell.Screen) {
	// Don't draw anything if there is no space.
	if b.width <= 0 || b.height <= 0 {
		return
	}

	def := tcell.StyleDefault

	// Fill background.
	background := def.Background(b.backgroundColor)
	for y := b.y; y < b.y+b.height; y++ {
		for x := b.x; x < b.x+b.width; x++ {
			screen.SetContent(x, y, ' ', nil, background)
		}
	}

	// Draw border.
	if b.border && b.width >= 2 && b.height >= 2 {
		border := background.Foreground(b.borderColor)
		var vertical, horizontal, topLeft, topRight, bottomLeft, bottomRight rune
		if b.focus.HasFocus() {
			vertical = GraphicsDbVertBar
			horizontal = GraphicsDbHorBar
			topLeft = GraphicsDbTopLeftCorner
			topRight = GraphicsDbTopRightCorner
			bottomLeft = GraphicsDbBottomLeftCorner
			bottomRight = GraphicsDbBottomRightCorner
		} else {
			vertical = GraphicsHoriBar
			horizontal = GraphicsVertBar
			topLeft = GraphicsTopLeftCorner
			topRight = GraphicsTopRightCorner
			bottomLeft = GraphicsBottomLeftCorner
			bottomRight = GraphicsBottomRightCorner
		}
		for x := b.x + 1; x < b.x+b.width-1; x++ {
			screen.SetContent(x, b.y, vertical, nil, border)
			screen.SetContent(x, b.y+b.height-1, vertical, nil, border)
		}
		for y := b.y + 1; y < b.y+b.height-1; y++ {
			screen.SetContent(b.x, y, horizontal, nil, border)
			screen.SetContent(b.x+b.width-1, y, horizontal, nil, border)
		}
		screen.SetContent(b.x, b.y, topLeft, nil, border)
		screen.SetContent(b.x+b.width-1, b.y, topRight, nil, border)
		screen.SetContent(b.x, b.y+b.height-1, bottomLeft, nil, border)
		screen.SetContent(b.x+b.width-1, b.y+b.height-1, bottomRight, nil, border)

		// Draw title.
		if b.title != "" && b.width >= 4 {
			width := b.width - 2
			title := []rune(b.title)
			if width < len(title) && width > 0 {
				title = append(title[:width-1], GraphicsEllipsis)
			}
			Print(screen, string(title), b.x+1, b.y, width, b.titleAlign, b.titleColor)
		}
	}
}

// Focus is called when this primitive receives focus.
func (b *Box) Focus(delegate func(p Primitive)) {
	b.hasFocus = true
}

// Blur is called when this primitive loses focus.
func (b *Box) Blur() {
	b.hasFocus = false
}

// HasFocus returns whether or not this primitive has focus.
func (b *Box) HasFocus() bool {
	return b.hasFocus
}

// GetFocusable returns the item's Focusable.
func (b *Box) GetFocusable() Focusable {
	return b.focus
}
