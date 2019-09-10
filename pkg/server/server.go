package server

import (
	"context"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

var tpl = template.Must(template.New("main").Parse(htmlTemplate))

func New(ctx context.Context, rootDir string) *Server {
	handler := chi.NewRouter()
	handler.Use(requestLogger(ctx))
	srv := &Server{
		root:    rootDir,
		fs:      http.Dir(rootDir),
		handler: handler,
	}
	handler.Handle("/*", http.FileServer(srv.fs))
	handler.Get("/{path}.md", srv.handleMarkdown)
	return srv
}

type Server struct {
	root    string
	fs      http.FileSystem
	handler http.Handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) handleMarkdown(w http.ResponseWriter, r *http.Request) {
	fp, err := s.fs.Open(r.URL.Path)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	extensions := parser.NewWithExtensions(parser.CommonExtensions)
	content := markdown.ToHTML(data, extensions, nil)

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	tpl.Execute(w, htmlTemplateContent{
		Content: template.HTML(content),
	})
}
