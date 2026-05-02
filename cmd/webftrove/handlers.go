package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	ft "github.com/steffenfritz/FileTrove"
)

//go:embed templates
var templateFS embed.FS

const pageSize = 100

type handler struct {
	db *sql.DB
}

type sessionsPageData struct {
	Summaries []ft.SessionSummary
}

type sessionPageData struct {
	Session   ft.SessionMD
	Stats     ft.SessionInfoMD
	Files     []ft.WebFileMD
	Total     int
	Page      int
	Pages     int
	Filters   ft.FileFilters
	Mimes     []string
	Exts      []string
	IsPartial bool
}

type filePageData struct {
	Session ft.SessionMD
	Detail  ft.FileDetail
}

type dirsPageData struct {
	Session   ft.SessionMD
	Dirs      []ft.WebDirMD
	Total     int
	Query     string
	IsPartial bool
}

var funcMap = template.FuncMap{
	"contains": func(slice []string, s string) bool {
		for _, v := range slice {
			if v == s {
				return true
			}
		}
		return false
	},
	"formatBytes": func(n int64) string {
		switch {
		case n < 1024:
			return fmt.Sprintf("%d B", n)
		case n < 1024*1024:
			return fmt.Sprintf("%.1f KB", float64(n)/1024)
		case n < 1024*1024*1024:
			return fmt.Sprintf("%.1f MB", float64(n)/1024/1024)
		default:
			return fmt.Sprintf("%.2f GB", float64(n)/1024/1024/1024)
		}
	},
	"formatEntropy": func(e float64) string {
		return fmt.Sprintf("%.4f", e)
	},
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
	"seq": func(n int) []int {
		s := make([]int, n)
		for i := range s {
			s[i] = i + 1
		}
		return s
	},
}

func parseTemplates(names ...string) *template.Template {
	paths := make([]string, len(names))
	for i, n := range names {
		paths[i] = "templates/" + n
	}
	return template.Must(template.New("").Funcs(funcMap).ParseFS(templateFS, paths...))
}

var (
	tmplSessions      = parseTemplates("layout.html", "sessions.html")
	tmplSession       = parseTemplates("layout.html", "session.html", "files_section.html")
	tmplFilesPartial  = parseTemplates("files_section.html")
	tmplFile          = parseTemplates("layout.html", "file.html")
	tmplDirs          = parseTemplates("layout.html", "dirs.html", "dirs_table.html")
	tmplDirsTable     = parseTemplates("dirs_table.html")
)

func (h *handler) sessionsHandler(w http.ResponseWriter, r *http.Request) {
	summaries, err := ft.GetSessionSummaries(h.db)
	if err != nil {
		httpError(w, "loading sessions", err)
		return
	}
	if err := tmplSessions.ExecuteTemplate(w, "layout", sessionsPageData{Summaries: summaries}); err != nil {
		log.Printf("render sessions: %v", err)
	}
}

func (h *handler) sessionHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	session, err := ft.GetSessionByUUID(h.db, uuid)
	if err != nil {
		httpError(w, "loading session", err)
		return
	}

	filters, page := parseFilters(r)
	files, total, err := ft.QueryFiles(h.db, uuid, filters)
	if err != nil {
		httpError(w, "querying files", err)
		return
	}

	mimes, _ := ft.GetDistinctMimes(h.db, uuid)
	exts, _ := ft.GetDistinctExtensions(h.db, uuid)

	stats, _ := ft.ListSession(h.db, uuid)

	pages := (total + pageSize - 1) / pageSize

	data := sessionPageData{
		Session: session,
		Stats:   stats,
		Files:   files,
		Total:   total,
		Page:    page,
		Pages:   pages,
		Filters: filters,
		Mimes:   mimes,
		Exts:    exts,
	}

	if err := tmplSession.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("render session: %v", err)
	}
}

func (h *handler) filesPartialHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	filters, page := parseFilters(r)
	files, total, err := ft.QueryFiles(h.db, uuid, filters)
	if err != nil {
		httpError(w, "querying files", err)
		return
	}

	pages := (total + pageSize - 1) / pageSize

	data := sessionPageData{
		Session:   ft.SessionMD{UUID: uuid},
		Files:     files,
		Total:     total,
		Page:      page,
		Pages:     pages,
		Filters:   filters,
		IsPartial: true,
	}

	if err := tmplFilesPartial.ExecuteTemplate(w, "files_section", data); err != nil {
		log.Printf("render files partial: %v", err)
	}
}

func (h *handler) fileHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	fileuuid := r.PathValue("fileuuid")

	session, err := ft.GetSessionByUUID(h.db, uuid)
	if err != nil {
		httpError(w, "loading session", err)
		return
	}

	detail, err := ft.GetFileDetail(h.db, fileuuid)
	if err != nil {
		httpError(w, "loading file detail", err)
		return
	}

	if err := tmplFile.ExecuteTemplate(w, "layout", filePageData{Session: session, Detail: detail}); err != nil {
		log.Printf("render file: %v", err)
	}
}

func (h *handler) dirsHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	query := r.URL.Query().Get("q")

	session, err := ft.GetSessionByUUID(h.db, uuid)
	if err != nil {
		httpError(w, "loading session", err)
		return
	}

	dirs, err := ft.GetSessionDirs(h.db, uuid, query)
	if err != nil {
		httpError(w, "loading directories", err)
		return
	}

	data := dirsPageData{Session: session, Dirs: dirs, Total: len(dirs), Query: query}

	if r.Header.Get("HX-Request") == "true" {
		data.IsPartial = true
		if err := tmplDirsTable.ExecuteTemplate(w, "dirs_table", data); err != nil {
			log.Printf("render dirs table: %v", err)
		}
		return
	}

	if err := tmplDirs.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("render dirs: %v", err)
	}
}

func (h *handler) openHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path required", http.StatusBadRequest)
		return
	}
	dir := filepath.Dir(path)
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "explorer"
	default:
		cmd = "xdg-open"
	}
	if err := exec.Command(cmd, dir).Start(); err != nil {
		log.Printf("open dir %s: %v", dir, err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseFilters(r *http.Request) (ft.FileFilters, int) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	return ft.FileFilters{
		Query:       q.Get("q"),
		QueryNegate: q.Get("qneg") == "1",
		Ext:      q.Get("ext"),
		Mimes:    q["mime"],
		NSRL:     q.Get("nsrl"),
		YaraOnly: q.Get("yara") == "1",
		SortBy:   q.Get("sort"),
		Order:    q.Get("order"),
		Limit:    pageSize,
		Offset:   offset,
	}, page
}

func httpError(w http.ResponseWriter, op string, err error) {
	log.Printf("%s: %v", op, err)
	http.Error(w, fmt.Sprintf("Error %s: %v", op, err), http.StatusInternalServerError)
}
