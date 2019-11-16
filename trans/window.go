package trans

/*
 * The window file will ultimately be the glue between the rest
 * of krex and the user's terminal
 *
 * TODO @kris-nova Can we please make this beautiful and pretty and wonderful
 *
 */

import (
	. "github.com/rthornton128/goncurses"
)

const (
	DefaultHeight = 30
	//DefaultWidth  = 200
)

var Version = ""

type TransWindow struct {
	height int
	width  int
	my     int
	mx     int

	// The menu window
	window *Window

	// The full terminal
	stdscr *Window
}

func GetNewWindow() (*TransWindow, error) {
	stdscr, err := Init()
	if err != nil {
		return nil, err
	}
	my, mx := stdscr.MaxYX()

	// Calculate 20 percent of the terminal for margin around our window
	width := mx - int(float64(mx)*float64(.2))
	height := my - int(float64(my)*float64(.2))

	// Offset of the inner window
	y, x := 2, (mx/2)-(width/2)

	//fmt.Println(height, width, y, x)

	// 30 200 2 25
	win, _ := NewWindow(height, width, y, x)
	win.Keypad(true)

	Raw(true)
	Echo(false)
	Cursor(0)

	//stdscr.Clear()
	stdscr.Keypad(true)
	defer End()
	//stdscr.Print(msg)
	//stdscr.Refresh()

	return &TransWindow{
		width:  width,
		height: height,
		window: win,
		stdscr: stdscr,
		my:     my,
		mx:     mx,
	}, nil
}

func (tw *TransWindow) StartScreen(msg string) error {

	return nil
}

func (tw *TransWindow) Prompt(title string, items []string) string {

	// Init the prompt
	defer End()

	// Current cursor position in the list
	var active int

	// Virtual window definition
	var wLines int        // Nb lines displayed in the window
	var xPosition int = 0 // Left top corner of the virtual window

	// Clear the window
	tw.window.Clear()
	tw.window.Refresh()

	// Clear the main screen
	tw.stdscr.Clear()
	tw.stdscr.Refresh()

	tw.stdscr.Printf("Krex version [%s] -- Kubernetes Resource Explore by Kris Nova <kris@fabulous.af>\n", Version)
	tw.stdscr.Printf("Use navigation arrows, pgup, pgdwn, home, end, backspace, [q] to exit\n")

	wLines = min(20, len(items)) // Set the number of displayed lines tp 20 max.

	// Draw the inital window
	draw(tw.window, items, active, xPosition, wLines)

	// Event loop
	for {
		ch := tw.stdscr.GetChar()
		switch Key(ch) {
		case 'q':
			//tw.stdscr.Clear()
			return ""
		case KEY_HOME:
			xPosition = 0
			active = 0
		case KEY_END:
			xPosition = len(items) - wLines
			active = xPosition
		case KEY_PAGEUP:
			xPosition = max(xPosition-wLines, 0)
			active = xPosition
		case KEY_PAGEDOWN:
			xPosition = min(xPosition+(2*wLines)-1, len(items)) - wLines
			active = xPosition
		case KEY_UP:
			active = max(active-1, 0)
			if active < xPosition {
				xPosition = active
			}
		case KEY_DOWN:
			active = min(active+1, len(items)-1)
			if active >= xPosition+wLines {
				xPosition++
			}
		case KEY_RETURN, KEY_ENTER, Key('\r'), KEY_RIGHT:
			tw.stdscr.MovePrint(tw.my-2, 0, "Choice #%d: %s selected",
				active,
				items[active])
			tw.stdscr.Refresh()
			tw.stdscr.Clear()
			return items[active]
		case KEY_LEFT, KEY_BACKSPACE:
			return "[Action] Go back ../"
		default:
			// Todo
			tw.stdscr.MovePrint(tw.my-2, 0, "Character pressed = %3d/%c",
				ch, ch)
			tw.stdscr.ClearToEOL()
			tw.stdscr.Refresh()
		}
		draw(tw.window, items, active, xPosition, wLines)
	}
}

func (tw *TransWindow) End() {
	End()
}

func draw(w *Window, items []string, active int, xPosition int, wLines int) {
	var clear string = "                                                                               "
	y, x := 2, 2
	w.Box(0, 0)
	//w.Background() //??
	for i, s := range items {
		if i >= xPosition && i < (xPosition+wLines) {
			w.MovePrint(y+i-xPosition, x, clear)
			if i == active {
				w.AttrOn(A_REVERSE)
				w.MovePrint(y+i-xPosition, x, s)
				w.AttrOff(A_REVERSE)
			} else {
				w.MovePrint(y+i-xPosition, x, s)
			}
		}
	}
	w.Refresh()
}

// Return the min value from 2 integers
func min(x int, y int) int {
	if x <= y {
		return x
	}
	return y
}

// Return the max value from 2 integers
func max(x int, y int) int {
	if x >= y {
		return x
	}
	return y
}
