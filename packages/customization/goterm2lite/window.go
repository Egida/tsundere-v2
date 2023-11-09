package goterm2lite

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"strconv"
	"strings"
)

type Position struct {
	Horizontal, Vertical int
}

var ErrPositionOutput = fmt.Errorf("invalid position output")

func GetWindowPosition(channel ssh.Channel) (*Position, error) {
	// Send the request to get the window position
	_, err := channel.Write([]byte("\x1b[13t"))
	if err != nil {
		return nil, err
	}

	// Read the response
	buf := make([]byte, 64)
	n, err := io.ReadFull(channel, buf)
	if err != nil {
		return nil, err
	}

	// Convert the response to a string
	response := string(buf[:n])

	// Check if the response is in the expected format
	if !strings.HasSuffix(response, "t") {
		return nil, ErrPositionOutput
	}

	// Remove the trailing "t"
	response = strings.TrimSuffix(response, "t")

	// Split the response by semicolon to extract horizontal and vertical positions
	positions := strings.Split(response, ";")
	if len(positions) != 2 {
		return nil, ErrPositionOutput
	}

	// Convert the positions to integers
	horizontal, err := strconv.Atoi(positions[0])
	if err != nil {
		return nil, err
	}

	vertical, err := strconv.Atoi(positions[1])
	if err != nil {
		return nil, err
	}

	return &Position{Horizontal: horizontal, Vertical: vertical}, nil
}

// loops until the byte value is equal to 0
// this will ensure its only has proper values
func LoopTillZero(src []byte, dst []byte) []byte {
	//ranges through src
	for p := range src {

		//if byte value is 0
		//we ignore them properly
		if src[p] == 0 { //checks
			return dst //ends loop
		}

		//appends into the dst properly
		dst = append(dst, src[p])
	}

	//returns the values
	return dst
}
