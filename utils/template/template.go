package template

import (
	"bytes"
	"fmt"
	htmltpl "html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/aymerick/douceur/inliner"
	"gopkg.in/russross/blackfriday.v2"
)

type Message struct {
	Body   string
	Name   string
	FromID string
	ToID   string
}

type Thread struct {
	Subject  string
	Messages []Message
	Preview  string
}

type message struct {
	Body   htmltpl.HTML
	Name   string
	FromID string
	ToID   string
}

type thread struct {
	Subject  string
	Messages []message
	Preview  string
}

var templates map[string]*htmltpl.Template

func init() {
	if templates == nil {
		templates = make(map[string]*htmltpl.Template)
	}

	var basePath string
	if strings.HasSuffix(os.Args[0], ".test") {
		basePath = "../templates"
	} else {
		basePath = "./templates"
	}

	layouts, err := filepath.Glob(basePath + "/layouts/*.html")
	if err != nil {
		panic(err)
	}

	includes, err := filepath.Glob(basePath + "/includes/*.html")
	if err != nil {
		panic(err)
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, layout := range layouts {
		files := append(includes, layout)
		templates[filepath.Base(layout)] = htmltpl.Must(htmltpl.ParseFiles(files...))
	}
}

func renderTemplate(name string, data Thread) (string, error) {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return "", fmt.Errorf("the template %s does not exist", name)
	}

	t := thread{
		Subject:  data.Subject,
		Messages: make([]message, len(data.Messages)),
		Preview:  data.Preview,
	}

	// Render markdown to HTML
	for i := range data.Messages {
		t.Messages[i].Body = htmltpl.HTML(blackfriday.Run([]byte(data.Messages[i].Body)))
		t.Messages[i].Name = data.Messages[i].Name
		t.Messages[i].FromID = data.Messages[i].FromID
		t.Messages[i].ToID = data.Messages[i].ToID
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base.html", t); err != nil {
		return "", err
	}

	html, err := inliner.Inline(buf.String())
	if err != nil {
		return html, err
	}

	return html, nil
}

// Render returns the email message rendered to a string.
func Render(data Thread) (string, error) {
	return renderTemplate("thread.html", data)
}
