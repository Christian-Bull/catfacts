package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	fact := "Test Fact"

	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<p>{{.Fact}}</p>
	</body>
</html>`

	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Title string
		Fact  string
	}{
		Title: "Cat Facts",
		Fact:  fact,
	}

	err = t.Execute(os.Stdout, data)
	if err != nil {
		log.Fatal(err)
	}
}
