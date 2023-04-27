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
	"text/template"
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
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

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
	line, lastLine, err := ReadLine(file, lineNum+1)
	if err != nil {
		c.l.Println("Error reading line ", line, lastLine, err)
	}

	// select random image
	imgNum := randgenerator.Intn(18)

	const tpl = `
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>{{.Title}}</title>
			
			<style>
			body {
			  background-color: #121212;
			}
			
			div {
				position: absolute;
				top: 50%;
				left: 50%;
				transform: translate(-50%, -50%);
				text-align: center;
			}
			h1 {
				color: white;
			}

			img {
				width: 60%;
				margin: auto;
				display: block;
			}

			</style>
		</head>
		<body>
			<div>
				<img src="assets/cats/cat-{{.Image}}.jpg" alt="Cute Cat Photo">
				<h1>{{.Fact}}</h1>
			</div>
		</body>
	</html>`

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Title string
		Fact  string
		Image int
	}{
		Title: "Cat Facts",
		Fact:  line,
		Image: imgNum,
	}

	err = t.Execute(rw, data)
	if err != nil {
		log.Fatal(err)
	}
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

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/favicon-32x32.png")
}
