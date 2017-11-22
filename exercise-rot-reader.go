package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (s rot13Reader) Read(buffer []byte) (int, error) {
	n, err := s.r.Read(buffer)

	for i := 0; i < n; i++ {
		if buffer[i] >= 'a' && buffer[i] <= 'z' {
			buffer[i] += 13
			if buffer[i] > 'z' {
				buffer[i] -= (13 * 2)
			}
		} else if buffer[i] >= 'A' && buffer[i] <= 'Z' {
			buffer[i] += 13
			if buffer[i] > 'Z' {
				buffer[i] -= (13 * 2)
			}
		}
	}
	
	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
