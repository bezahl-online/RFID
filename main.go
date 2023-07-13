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
var nodevice *bool

const BAUD = 9600
const WAIT_FOR_RFID_MODULE_TIME = 3 * time.Second

func main() {
	unsecure = flag.Bool("u", false, "unsecure")
	nodevice = flag.Bool("d", false, "no device")
	flag.Parse()
	rfid_port := "/dev/serial/by-id/usb-1a86_USB2.0-Ser_-if00-port0"
	if len(os.Getenv("RFID_PORT_NAME")) > 3 {
		rfid_port = os.Getenv("RFID_PORT_NAME")
	}
	var err error
	PipeReader, PipeWriter = io.Pipe()
	go func() {
		var port *term.Term
		if !*nodevice {
			for {
				port, err = term.Open(rfid_port)
				if err != nil {
					log.Printf("error connecting to RFID modul on port %s\nwaiting for module to come up..", rfid_port)
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
					log.Println("listening to RFID module ..")
					n, err := port.Read(buf)
					if err != nil {
						log.Printf("error while reading from RFID module: %s\n", err)
						break
					}
					log.Printf("got '%s' from RFID module(%d bytes)\n", string(buf[:14]), n)
					if n != 16 {
						log.Printf("error: data length is not 16 (it is %d)", n)
						port.Flush()
						time.Sleep(5 * time.Second)
						// FIXME: reset buf?
						continue
					}
					if waitingForID {
						waitingForID = false
						str := string(buf[:14])
						decID, err := strconv.ParseInt(str, 16, 64)
						if err != nil {
							log.Printf("error while converting '%s' to decimal: %s\n", str, err)
						}
						resp := fmt.Sprintf("%d", decID)
						writeToPipe(resp)
					}
				}
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// FIXME: if data older then discard...
		w.Header().Set("Access-Control-Allow-Origin", "*")
		var data string
		buf := bytes.NewBufferString(data)
		waitingForID = true
		_, err := io.Copy(buf, PipeReader)
		log.Printf("sent %s as response\n", buf.String())
		if strings.Contains(buf.String(), "lost") {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write(buf.Bytes())
		if err != nil {
			log.Printf("error while sending response: %s\n", err)
		}
	})
	http.HandleFunc("/mock", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		code := r.URL.Query().Get("code")
		fmt.Printf("mock: '%s'", code)
		writeToPipe(code)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	log.Println("RFID server v0.2.0")
	if *nodevice {
		log.Println("*** starting without device")
	}
	if *unsecure {
		log.Println("*** unsecure http server started on port 8040")
		log.Fatal(http.ListenAndServe(":8040", nil))
	} else {
		log.Println("https server started on port 8040")
		log.Fatal(http.ListenAndServeTLS(":8040", "localhost.crt", "localhost.key", nil))
	}
}

func writeToPipe(resp string) {
	rrr := strings.NewReader(resp)
	io.Copy(PipeWriter, rrr)
	PipeWriter.Close()
	PipeReader, PipeWriter = io.Pipe()
}
