package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/term"
)

var PipeReader *io.PipeReader

func main() {
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
		// fmt.Println("listening...")
		buf := make([]byte, 16)
		for {
			_, err := port.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println(buf)
			rrr := strings.NewReader(string(buf))
			io.Copy(w, rrr)
			w.Close()
			PipeReader, w = io.Pipe()
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(w, PipeReader)
		if err != nil {
			log.Fatal(err)
		}
	})
	// fmt.Println("server started on port 8080")
	log.Fatal(http.ListenAndServeTLS(":8040", "localhost.crt", "localhost.key", nil))
}
