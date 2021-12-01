package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/term"
)

var PipeReader *io.PipeReader
var waitingForID = false
var unsecure *bool

func main() {
	unsecure = flag.Bool("u", false, "unsecure")
	flag.Parse()
	var err error
	var w *io.PipeWriter
	PipeReader, w = io.Pipe()
	go func() {
		var port *term.Term
		port, err = term.Open("/dev/ttyUSB0")
		if err != nil {
			log.Fatal(err)
		}
		port.SetRaw()
		buf := make([]byte, 16)
		for {
			_, err := port.Read(buf)
			if err != nil {
				fmt.Printf("error while reading from RFID module: %s\n", err)
			}
			fmt.Printf("got '%s' from RFID module\n", string(buf[:14]))
			if waitingForID {
				buf[14] = 10
				rrr := strings.NewReader(string(buf[:15]))
				io.Copy(w, rrr)
				w.Close()
				PipeReader, w = io.Pipe()
			}
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		waitingForID = true
		_, err := io.Copy(w, PipeReader)
		if err != nil {
			fmt.Printf("error while sending response: %s\n", err)
		}
		fmt.Println("sent as response")
		waitingForID = false
	})
	if *unsecure {
		fmt.Println("http server started on port 8040")
		log.Fatal(http.ListenAndServe(":8040", nil))
	}
	fmt.Println("https server started on port 8040")
	log.Fatal(http.ListenAndServeTLS(":8040", "localhost.crt", "localhost.key", nil))
}
