package pkg

import (
	"io"
	"strings"
)

func NewBashOutputReader(src io.Reader) *stdReader {
	return &stdReader{
		src:    src,
		behind: 0,
	}
}

type stdReader struct {
	src    io.Reader
	behind int
}

func (sr *stdReader) ReadString() (length int, text string, err error) {
	tmp := make([]byte, 1024)
	length, err = sr.src.Read(tmp)
	if err != nil {
		return 0, "", err
	}
	if length <= sr.behind {
		return 0, "", err
	}

	realReadLength := length - sr.behind

	text = strings.TrimSpace(string(tmp[sr.behind:length]))

	sr.behind += length

	return realReadLength, text, nil
}
