package mouse

import (
	"bytes"
	"errors"
	"io"
)

func Listen(channel io.ReadWriter) (*Event, error) {
	if _, err := channel.Write([]byte("\033[?1000h\033[?25h\033[?25l")); err != nil {
		return nil, err
	}

	for {
		buf := make([]byte, 6)
		if _, err := channel.Read(buf); err != nil {
			return nil, err
		}

		if !bytes.HasPrefix(buf, []byte{27, 91, 77}) {
			continue
		}

		var alertType Alert = 0
		switch buf[3] {
		case 32:
			alertType = LeftClick
		case 33:
			alertType = ScrollClick
		case 34:
			alertType = RightClick
		case 96:
			alertType = ScrollUp
		case 97:
			alertType = ScrollDown
		case 35:
			return nil, errors.New("invalid")
		}

		return &Event{
			Click:  alertType,
			X:      ConvertPosition(buf[4]),
			Y:      ConvertPosition(buf[5]),
			String: alertType.String(),
		}, nil
	}
}
