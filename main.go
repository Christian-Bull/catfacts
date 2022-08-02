package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	l := log.New(os.Stdout, "catfacts-api", log.LstdFlags)

	bindAddr := fmt.Sprintf(":%s", port)

	ch := newCatfact(l)
	http.Handle("/", ch)

	l.Printf("Starting server on port %s", port)
	l.Fatal(http.ListenAndServe(bindAddr, nil))
}

type Catfact struct {
	l *log.Logger
}

func newCatfact(l *log.Logger) *Catfact {
	return &Catfact{l}
}

func (c *Catfact) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c.l.Println("Serving catfact")

	// count max lines
	fileLines, err := os.Open("src/facts.txt")
	if err != nil {
		c.l.Println(err)
	}
	defer fileLines.Close()

	maxLines, err := lineCounter(fileLines)
	if err != nil {
		c.l.Println("Error counting lines")
		maxLines = 101
	}

	// select random line
	randsource := rand.NewSource(time.Now().UnixNano())
	randgenerator := rand.New(randsource)
	lineNum := randgenerator.Intn(maxLines)

	// read file and serve cat fact
	file, err := os.Open("src/facts.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer fileLines.Close()
	line, lastLine, err := ReadLine(file, lineNum)
	if err != nil {
		c.l.Println("Error reading line ", line, lastLine, err)
	}

	fmt.Fprintf(rw, line)
}

// Read a specific line in a text file
func ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

// Count max number of lines in a text file
func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
