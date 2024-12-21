package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

// setupDB initializes the database connection and ensures it stays open
func setupDB() {
	var err error
	db, err = sql.Open("sqlite3", "./hlts.db")
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	// Create the table if it doesn't exist
	createTable()
}

func main() {
	setupDB()
	defer db.Close()

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Route handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/careers", careersHandler)
	http.HandleFunc("/submit", submitApplicationHandler)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Could not load the index template", http.StatusInternalServerError)
		log.Println("Index template execution error:", err)
	}
}

func careersHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "careers.html", nil)
	if err != nil {
		http.Error(w, "Could not load the careers template", http.StatusInternalServerError)
		log.Println("Careers template execution error:", err)
	}
}

func createTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS job_applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		phone TEXT,
		email TEXT,
		status TEXT,
		experience INTEGER,
		details TEXT,
		resume_path TEXT
	);`)
	if err != nil {
		log.Fatal("Error creating job_applications table:", err)
	}
}

func submitApplicationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	status := r.FormValue("status")
	experience := r.FormValue("experience")
	details := r.FormValue("details")

	_, err := db.Exec("INSERT INTO job_applications (name, phone, email, status, experience, details) VALUES (?, ?, ?, ?, ?, ?)", name, phone, email, status, experience, details)
	if err != nil {
		http.Error(w, "Unable to save application", http.StatusInternalServerError)
		log.Println("Database insert error:", err)
		return
	}

	http.Redirect(w, r, "/careers", http.StatusSeeOther)
}
