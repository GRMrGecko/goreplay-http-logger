package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Separator between http requests.
var payloadSeparator = "\n🐵🙈🙉\n"

// Generate random hex.
func randByte(len int) []byte {
	b := make([]byte, len/2)
	rand.Read(b)

	h := make([]byte, len)
	hex.Encode(h, b)

	return h
}

// Generate a uuid for a request.
func uuid() []byte {
	return randByte(24)
}

// Generate a header with type 1 for a request.
func payloadHeader(uuid []byte, timing int64) (header []byte) {
	return []byte(fmt.Sprintf("1 %s %d 0\n", uuid, timing))
}

// The log file structure for writing logs.
type LogFile struct {
	sync.Mutex
	file        *os.File
	currentName string
	nextUpdate  time.Time
}

// Write to log file.
func (l *LogFile) write(data []byte) (n int, err error) {
	// Lock to prevent multiple writes at the same time.
	l.Lock()
	defer l.Unlock()

	// Get current time for file name generator.
	now := time.Now()
	// If no file defined or we are after the next update time, update the tile name.
	if l.file == nil || now.After(l.nextUpdate) {
		// Set next update to a second later, truncating the nanoseconds.
		l.nextUpdate = now.Truncate(time.Second).Add(time.Second)

		// Generate the log file name based on config.
		name := config.LogFile
		// Year
		name = strings.ReplaceAll(name, "%Y", now.Format("2006"))
		// Month
		name = strings.ReplaceAll(name, "%m", now.Format("01"))
		// Day
		name = strings.ReplaceAll(name, "%d", now.Format("02"))
		// Hour
		name = strings.ReplaceAll(name, "%H", now.Format("15"))
		// Minute
		name = strings.ReplaceAll(name, "%M", now.Format("04"))
		// Second
		name = strings.ReplaceAll(name, "%S", now.Format("05"))
		l.currentName = filepath.Clean(name)

		// If new name generated is different from existing open file or no file is opened, open it.
		if l.file == nil || l.currentName != l.file.Name() {
			l.file, err = os.OpenFile(l.currentName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)

			if err != nil {
				log.Fatalf("Cannot open file %q. Error: %s", l.currentName, err)
			}
		}
	}

	// Write data.
	n, err = l.file.Write(data)
	return
}

// Global log file definition.
var logFile *LogFile

// Log the http request.
func logRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request %s %s", r.Method, r.URL)

	// Generate log entry.
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "%s", payloadHeader(uuid(), time.Now().UnixNano()))
	fmt.Fprintf(&buff, "%s %s %s\n", r.Method, r.URL, r.Proto)
	r.Header.Write(&buff)
	fmt.Fprint(&buff, "\r\n")
	body, _ := io.ReadAll(r.Body)
	buff.Write(body)
	fmt.Fprint(&buff, payloadSeparator)

	// Write to log file.
	logFile.write(buff.Bytes())
}

// The config for this run, set in flags.
type Config struct {
	HTTPBind string
	HTTPPort int
	LogFile  string
}

// Global config variable.
var config Config

// Generate help information.
func usage() {
	fmt.Println("http log server")
	flag.PrintDefaults()
	os.Exit(2)
}

// The main program.
func main() {
	// Parse flags.
	flag.Usage = usage
	flag.StringVar(&config.HTTPBind, "bind", "", "HTTP bind address")
	flag.IntVar(&config.HTTPPort, "port", 8080, "HTTP port")
	flag.StringVar(&config.LogFile, "log-file", "http-%Y%m%d.log", "Log file name with date")
	flag.Parse()

	// Setup the log file.
	logFile = new(LogFile)

	// Handle all http requests with the logRequest handler.
	http.HandleFunc("/", logRequest)

	// Start the HTTP server/
	fmt.Printf("Starting server at port %d\n", config.HTTPPort)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.HTTPBind, config.HTTPPort), nil); err != nil {
		log.Fatal(err)
	}
}
