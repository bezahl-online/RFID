package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
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
				str := string(buf[:14])
				decID, err := strconv.ParseInt(str, 16, 64)
				if err != nil {
					fmt.Printf("error while converting '%s' to decimal: %s\n", str, err)
				}
				resp := fmt.Sprintf("%d", decID)
				rrr := strings.NewReader(resp)
				io.Copy(w, rrr)
				w.Close()
				fmt.Printf("sent %d as response", decID)
				PipeReader, w = io.Pipe()
			}
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		waitingForID = true
		_, err := io.Copy(w, PipeReader)
		if err != nil {
			fmt.Printf("error while sending response: %s\n", err)
		}
		waitingForID = false
	})
	if *unsecure {
		fmt.Println("http server started on port 8040")
		log.Fatal(http.ListenAndServe(":8040", nil))
	}
	fmt.Println("https server started on port 8040")
	log.Fatal(http.ListenAndServeTLS(":8040", "localhost.crt", "localhost.key", nil))
}
