package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/term"
)

var PipeReader *io.PipeReader
var PipeWriter *io.PipeWriter
var waitingForID = false
var unsecure *bool

const BAUD = 9600
const WAIT_FOR_RFID_MODULE_TIME = 3 * time.Second

func main() {
	unsecure = flag.Bool("u", false, "unsecure")
	flag.Parse()
	rfid_port := "/dev/serial/by-id/usb-1a86_USB2.0-Ser_-if00-port0"
	if len(os.Getenv("RFID_PORT_NAME")) > 3 {
		rfid_port = os.Getenv("RFID_PORT_NAME")
	}
	var err error
	PipeReader, PipeWriter = io.Pipe()
	go func() {
		var port *term.Term
		for {
			port, err = term.Open(rfid_port)
			if err != nil {
				fmt.Printf("error connecting to RFID modul on port %s\nwaiting for module to come up..", rfid_port)
				if waitingForID {
					writeToPipe("lost connection to RFID module")
				}
				time.Sleep(WAIT_FOR_RFID_MODULE_TIME)
				continue
			}
			port.SetRaw()
			port.SetSpeed(BAUD)
			buf := make([]byte, 16)
			for {
				fmt.Println("listening to RFID module ..")
				n, err := port.Read(buf)
				if err != nil {
					fmt.Printf("error while reading from RFID module: %s\n", err)
					break
				}
				fmt.Printf("got '%s' from RFID module\n", string(buf[:14]))
				if waitingForID && n > 7 {
					waitingForID = false
					str := string(buf[:14])
					decID, err := strconv.ParseInt(str, 16, 64)
					if err != nil {
						fmt.Printf("error while converting '%s' to decimal: %s\n", str, err)
					}
					resp := fmt.Sprintf("%d", decID)
					writeToPipe(resp)
				}
			}
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var data string
		buf := bytes.NewBufferString(data)
		waitingForID = true
		_, err := io.Copy(buf, PipeReader)
		fmt.Printf("sent %s as response\n", buf.String())
		if strings.Contains(buf.String(), "lost") {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write(buf.Bytes())
		if err != nil {
			fmt.Printf("error while sending response: %s\n", err)
		}
	})
	if *unsecure {
		fmt.Println("http server started on port 8040")
		log.Fatal(http.ListenAndServe(":8040", nil))
	}
	fmt.Println("https server started on port 8040")
	log.Fatal(http.ListenAndServeTLS(":8040", "localhost.crt", "localhost.key", nil))
}

func writeToPipe(resp string) {
	rrr := strings.NewReader(resp)
	io.Copy(PipeWriter, rrr)
	PipeWriter.Close()
	PipeReader, PipeWriter = io.Pipe()
}
