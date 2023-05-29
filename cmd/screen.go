package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

func DrawText(s tcell.Screen, StartX, StartY, EndX, EndY int, style tcell.Style, text string) error {
	x1, y1, x2, y2 := StartX, StartY, EndX, EndY
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
	return nil
}

func NewTerminal() (*Terminal, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	StyleDefault := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(StyleDefault)
	return &Terminal{
		S:            s,
		StyleDefault: StyleDefault,
	}, nil
}

type Terminal struct {
	S            tcell.Screen
	StyleDefault tcell.Style
}

func (t *Terminal) Exec() error {
	return nil
}

// func (t *Terminal) Draw(i IDraw) error {
// return nil
// }

func (t *Terminal) DrawLineText(StartX, StartY, MaxLineW int, style tcell.Style, text string) error {
	text = runewidth.FillRight(text, MaxLineW)
	return t.DrawText(StartX, StartY, StartX+MaxLineW, StartY, style, text)
}

func (t *Terminal) DrawText(StartX, StartY, EndX, EndY int, style tcell.Style, text string) error {
	x1, y1, x2, y2 := StartX, StartY, EndX, EndY
	row := y1
	col := x1
	for _, r := range []rune(text) {
		t.S.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
	return nil
}

func (t *Terminal) DrawBox(b Box) error {
	x1, y1, x2, y2, style := b.StartX, b.StartY, b.EndX, b.EndY, b.Style
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			t.S.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		t.S.SetContent(col, y1, tcell.RuneHLine, nil, style)
		t.S.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		t.S.SetContent(x1, row, tcell.RuneVLine, nil, style)
		t.S.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		t.S.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		t.S.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		t.S.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		t.S.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	// return t.DrawText(x1+1, y1+1, x2-1, y2-1, style, text)
	return nil
}

func (t *Terminal) Quit() {
	maybePanic := recover()
	t.S.Fini()
	if maybePanic != nil {
		panic(maybePanic)
	}
}

// type IDraw interface {
// GetPosition() (StartX, StartY, EndX, EndY int)
// GetStyle() tcell.Style
// }

func NewBox(StartX, StartY, EndX, EndY int, Style tcell.Style) *Box {
	return &Box{
		StartX: StartX,
		StartY: StartY,
		EndX:   EndX,
		EndY:   EndY,
		Style:  Style,
	}
}

type Box struct {
	StartX int
	StartY int
	EndX   int
	EndY   int
	Style  tcell.Style
}

func (b Box) DrawRange() (StartX, StartY, EndX, EndY int) {
	return b.StartX + 1, b.StartY + 1, b.EndX - 1, b.EndY - 1
}

func (b *Box) DrawText(s tcell.Screen, style tcell.Style, text string) {
	sx, sy, ex, ey := b.DrawRange()
	DrawText(s, sx, sy, ex, ey, style, text)
}

type SelectItem struct {
	IsSelect bool
	Info     interface{}
}

type Select struct {
	StartX      int
	StartY      int
	MaxWidth    int
	MaxHeight   int
	Items       []*SelectItem
	SelectIndex int
	StyleSelect tcell.Style
}

func (s Select) GetSeleteItem() *SelectItem {
	return s.Items[s.SelectIndex]
}

func (s *Select) GetDrawItems() []*SelectItem {
	if len(s.Items) > s.MaxHeight {

	}
	return s.Items
}

func (s *Select) MoveDownSelect(step int) (isChange bool) {
	var minH = s.MaxHeight
	if len(s.Items) < minH {
		minH = len(s.Items)
	}
	if s.SelectIndex+1 < minH {
		s.SelectIndex += 1
		isChange = true
	}
	Log.Debugf("MoveDownSelect step: %d index %d", step, s.SelectIndex)
	return
}

func (s *Select) MoveUpSelect(step int) (isChange bool) {
	if s.SelectIndex != 0 {
		s.SelectIndex -= step
		if s.SelectIndex < 0 {
			s.SelectIndex = 0
		}
		isChange = true
	}
	Log.Debugf("MoveUpSelect step: %d index %d", step, s.SelectIndex)
	return
}
