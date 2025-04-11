package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// User struct represents a user in the system
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var db *sql.DB

// Initialize the database connection
func initDB() {
	var err error
	dsn := "user=sukhbat password=new_password dbname=user_auth sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected!")
}

// Helper function to hash passwords
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// Helper function to compare hashed passwords
func checkPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Handler for user registration
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password before storing it
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Printf("Error hashing password: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO users(name, email, password) VALUES($1, $2, $3)",
		user.Name, user.Email, hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error inserting user: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("User %s registered successfully", user.Email)
}

// Handler for user login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var storedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE email = $1", user.Email).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		log.Printf("Login failed for user %s: %v", user.Email, err)
		return
	}

	// Check the password against the stored hashed password
	err = checkPassword(storedPassword, user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		log.Printf("Password mismatch for user %s", user.Email)
		return
	}

	// Return a dummy token (replace with a real token in production)
	token := "dummy-token"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
	log.Printf("User %s logged in successfully", user.Email)
}

// Handler for getting all registered users
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT name, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error fetching users: %v", err)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Name, &user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error scanning user: %v", err)
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error during rows iteration: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Handler to serve HTML files (example for sign.html)
func serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(".", "sign.html"))
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	// Enable CORS for all origins (adjust as needed for production)
	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	r.Use(handlers.CORS(originsOk, headersOk, methodsOk))

	// Allow OPTIONS for both endpoints to handle preflight
	r.HandleFunc("/api/register", registerHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/login", loginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/users", getUsersHandler).Methods("GET")
	r.HandleFunc("/sign", serveHTML).Methods("GET")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
