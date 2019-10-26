package server

import (
	"bytes"
	"context"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	path := filepath.Join(s.root, r.URL.Path)
	path, err := filepath.Abs(path)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	rel, err := filepath.Rel(s.root, path)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if strings.Contains(rel, "../") {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	st, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	if st.IsDir() {
		s.handleDir(w, r, path)
		return
	}
	if strings.HasSuffix(r.URL.Path, ".md") {
		s.handleMarkdown(w, r)
		return
	}
	if strings.HasSuffix(r.URL.Path, ".html") {
		fp, err := os.Open(path)
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer fp.Close()
		http.ServeContent(w, r, st.Name(), st.ModTime(), fp)
		return
	}
	http.FileServer(s.fs).ServeHTTP(w, r)
}

func (s *Server) handleDir(w http.ResponseWriter, r *http.Request, path string) {
	if !strings.HasSuffix(r.URL.Path, "/") {
		http.Redirect(w, r, r.URL.Path+"/", http.StatusTemporaryRedirect)
		return
	}
	var listingTpl = template.Must(template.New("listing").Funcs(map[string]interface{}{
		"datetime": func(v time.Time) string {
			return v.Format("2006-01-02")
		},
	}).Parse(listingTemplate))
	files, err := ioutil.ReadDir(path)
	if err != nil {
		http.Error(w, "Failed to load directory", http.StatusNotFound)
		return
	}
	var out bytes.Buffer
	if err := listingTpl.Execute(&out, files); err != nil {
		http.Error(w, "Failed to load directory", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	tpl.Execute(w, htmlTemplateContent{
		Content: template.HTML(out.String()),
	})
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
