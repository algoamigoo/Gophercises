package cyoa

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"text/template"
)

type Book map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func JsonStory (r io.Reader) (Book,error){
	var book Book
	err:=json.NewDecoder(r).Decode(&book)
	if err!=nil{
		return nil,err
	}
	return book,nil
}

var defaultChapterTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{range .Story}}
        <p>{{.}}</p>
    {{end}}
    <ul>
        {{range .Options}}
            <li><a href="/{{.Arc}}">{{.Text}}</a></li>
        {{end}}
    </ul>
</body>
</html>`

var tmpl = template.Must(template.New("").Parse(defaultChapterTemplate))

func NewHandler(b Book) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSpace(r.URL.Path)
		if path == "" || path == "/" {
			path = "/intro"
		}
		path = path[1:] // Strip leading "/"

		if chapter, ok := b[path]; ok {
			// Now 'tmpl' is defined in the package scope and works here!
			err := tmpl.Execute(w, chapter)
			if err != nil {
				http.Error(w, "Something went wrong...", http.StatusInternalServerError)
			}
			return
		}

		http.Error(w, "Chapter not found.", http.StatusNotFound)
	})
}