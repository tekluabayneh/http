package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(out)

		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				out <- str
				str = ""
			}
			str = string(data)

		}
		if len(str) != 0 {
			fmt.Printf("read: %s\n", str)
		}

	}()

	return out
}

func main() {

	listner, err := net.Listen("tcp", ":42069")

	if err != nil {
		log.Fatal("error", "error", err)
	}
	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		go func(con net.Conn) {
			for line := range getLinesChannel(con) {
				fmt.Printf("read: %s\n", line)
			}

		}(conn)
	}

}
