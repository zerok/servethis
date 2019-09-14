package server

import (
	"bytes"
	"context"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

var tpl = template.Must(template.New("main").Parse(htmlTemplate))

type frontmatter struct {
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
}

func New(ctx context.Context, rootDir string) *Server {
	handler := chi.NewRouter()
	handler.Use(requestLogger(ctx))
	srv := &Server{
		root:    rootDir,
		fs:      http.Dir(rootDir),
		handler: handler,
	}
	handler.Handle("/*", http.HandlerFunc(srv.dispatch))
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

func (s *Server) dispatch(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".md") {
		s.handleMarkdown(w, r)
		return
	}
	http.FileServer(s.fs).ServeHTTP(w, r)
}

func (s *Server) handleMarkdown(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())
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
	var fm frontmatter

	if bytes.HasPrefix(data, []byte("---")) {
		end := bytes.Index(data[3:], []byte("---"))
		if end != -1 {
			if err := yaml.Unmarshal(data[3:end+3], &fm); err != nil {
				http.Error(w, "Failed to decode frontmatter", http.StatusInternalServerError)
				logger.Error().Err(err).Msg("Failed to decode frontmatter.")
			}
		}
		data = data[end+3+3:]
	}

	extensions := parser.NewWithExtensions(parser.CommonExtensions)
	content := markdown.ToHTML(data, extensions, nil)

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	tpl.Execute(w, htmlTemplateContent{
		Content:     template.HTML(content),
		Frontmatter: fm,
	})
}
