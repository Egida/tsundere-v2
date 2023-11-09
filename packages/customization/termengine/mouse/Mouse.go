package mouse

type Alert int

type Event struct {
	Click  Alert
	X      int
	Y      int
	String string
}

const (
	LeftClick   Alert = 1
	ScrollClick Alert = 2
	RightClick  Alert = 3
	ScrollUp    Alert = 5
	ScrollDown  Alert = 6
)

func (a *Alert) String() string {
	switch *a {
	case LeftClick:
		return "LEFT_CLICK"
	case ScrollClick:
		return "SCROLL_CLICK"
	case RightClick:
		return "RIGHT_CLICK"
	case ScrollUp:
		return "SCROLL_UP"
	case ScrollDown:
		return "SCROLL_DOWN"
	default:
		return "EOF"
	}
}

func ConvertPosition(buf byte) int {
	return int(buf) - 32
}
