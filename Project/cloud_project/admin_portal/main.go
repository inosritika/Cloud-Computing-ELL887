package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
	_ "github.com/microsoft/go-mssqldb"
	// "github.com/rs/cors"
	// "github.com/shopspring/decimal"
)

type Dishes struct {
	ID int `json:"ID,omitempty"`
	DishName  string `json:"DishName,omitempty"`
	CountryOfOrigin string `json:"CountryOfOrigin,omitempty"`
	TypeOfDish string  `json:"TypeOfDish,omitempty"`
	Price float64 `json:"Price,omitempty"`
	IsVegetarian bool `json:"IsVegetarian"`
	IsNonVegetarian bool `json:"IsNonVegetarian"`
	IsVegan bool `json:"IsVegan"`
}

type ID struct{
	ID int `json:"ID,omitempty"`
}

type User struct {
    UserName     string `json:"userName"`
    Email        string `json:"email"`
    PasswordHash string `json:"passwordHash"`
}

var db *sql.DB
var server = "cloudcomputingproject.database.windows.net"
var port = 1433
var user = "Niketh"
var password = "Sai@1234"
var database = "displaytable"

const(
	// addProductQuery = "INSERT INTO DataTable (DishName, CountryOfOrigin, TypeOfDish, Price, IsVegetarian, IsNonVegetarian, IsVegan) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7);"
	listProductsQuery = "select ID, DishName, CountryOfOrigin, TypeOfDish, Price, IsVegetarian, IsNonVegetarian, IsVegan from DataTable WHERE Status='NOT_APPROVED';"
)

func listProductsHandler(w http.ResponseWriter, r *http.Request) {
	var dishes []Dishes
	stmt, err := db.Prepare(listProductsQuery)
	if err != nil {
		fmt.Println("The error is", err)
        http.Error(w, "Error2 in executing query", http.StatusInternalServerError)
		return
	}
	
	// 3. Execute the statement
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("The error is", err)
        http.Error(w, "Error2 in executing query", http.StatusInternalServerError)
		return
	}
	defer rows.Close() 
	
	for rows.Next() {
		var d Dishes
		err := rows.Scan(&d.ID, &d.DishName, &d.CountryOfOrigin, &d.TypeOfDish, &d.Price, &d.IsVegetarian, &d.IsNonVegetarian, &d.IsVegan) // Scan directly into struct fields
		if err != nil {
			fmt.Println("The error is", err)
			http.Error(w, "Error2 in executing query", http.StatusInternalServerError)
			return
		}
		dishes = append(dishes, d) // Append product to the slice
	}
	
	err = rows.Err()
	if err != nil {
		fmt.Println("The error is", err)
        http.Error(w, "Error2 in executing query", http.StatusInternalServerError)
		return
	}
  
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dishes)
}

func StatusApproveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Function called")
	var id ID
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		fmt.Println("Hello I am here")
		fmt.Println("The error is", err)
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE DataTable SET Status='APPROVED' WHERE ID = @p1;", id.ID)
	if err != nil {
		fmt.Println("The error is", err)
		http.Error(w, "Error2 in executing query", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id.ID)
}

func StatusRejectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Reject function called")
	var id ID
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		fmt.Println("Hello I am here")
		fmt.Println("The error is", err)
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE DataTable SET Status='REJECTED' WHERE ID = @p1;", id.ID)
	if err != nil {
		fmt.Println("The error is", err)
		http.Error(w, "Error2 in executing query", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id.ID)
}

// func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// 	var user User
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
// 		return
// 	}

// 	// Hash the password
// 	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
// 	if err != nil {
// 		http.Error(w, "Failed to hash the password", http.StatusInternalServerError)
// 		return
// 	}

// 	// Insert the user into the database
// 	_, err = db.Exec("INSERT INTO Admins (UserName, Email, PasswordHash) VALUES (@p1, @p2, @p3);", user.UserName, user.Email, passwordHash)
// 	if err != nil {
// 		http.Error(w, "Failed to insert user into the database", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(user)
// }

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	row := db.QueryRow("SELECT PasswordHash FROM Admins WHERE Email = @p1;", user.Email)

	var passwordHash string
	err = row.Scan(&passwordHash)
	if err != nil {
		http.Error(w, "Failed to get user from the database", http.StatusInternalServerError)
		return
	}
	// Compare the password with the hash
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.PasswordHash))
	if err != nil {
		fmt.Println("The error is", err)
		fmt.Print(passwordHash, user.PasswordHash)
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // Here you would clear any session or authentication token
    // For simplicity, we'll just send a success response
 
    // Respond to the client that logout was successful
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func main() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!")

	// http.HandleFunc("/api/addProduct", addProductHandler)
	http.HandleFunc("/api/listProducts", listProductsHandler)
	http.HandleFunc("/api/statusApprove", StatusApproveHandler)
	http.HandleFunc("/api/statusReject", StatusRejectHandler)
	// http.HandleFunc("/api/register", RegisterHandler)
	http.HandleFunc("/api/login", LoginHandler)
	http.HandleFunc("/api/logout", LogoutHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// log.Println("Server is running on http://localhost")
	log.Fatal(http.ListenAndServe("", nil))
}
