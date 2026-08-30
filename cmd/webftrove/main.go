package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	flag "github.com/spf13/pflag"
	_ "github.com/mattn/go-sqlite3"
)

const port = "9000"

func main() {
	dbPath := flag.StringP("db", "d", "", "Path to the FileTrove database file (required)")
	flag.Parse()

	if *dbPath == "" {
		log.Fatal("--db is required: provide the path to a filetrove.db file")
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?mode=ro", *dbPath))
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("connect to database %s: %v", *dbPath, err)
	}

	h := &handler{db: db}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.sessionsHandler)
	mux.HandleFunc("GET /session/{uuid}", h.sessionHandler)
	mux.HandleFunc("GET /session/{uuid}/files", h.filesPartialHandler)
	mux.HandleFunc("GET /session/{uuid}/file/{fileuuid}", h.fileHandler)
	mux.HandleFunc("GET /session/{uuid}/dirs", h.dirsHandler)
	mux.HandleFunc("GET /open", h.openHandler)

	addr := "http://localhost:" + port
	log.Printf("FileTrove Web running at %s", addr)

	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(addr)
	}()

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func openBrowser(url string) {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "windows":
		cmd = "start"
	default:
		cmd = "xdg-open"
	}
	if err := exec.Command(cmd, url).Start(); err != nil {
		log.Printf("could not open browser: %v", err)
	}
}
