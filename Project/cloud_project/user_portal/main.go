package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/microsoft/go-mssqldb"
	// "github.com/shopspring/decimal"
)

type Dish struct {
	DishName  string `json:"DishName,omitempty"`
	CountryOfOrigin string `json:"CountryOfOrigin,omitempty"`
	TypeOfDish string  `json:"TypeOfDish,omitempty"`
	Price float64 `json:"Price,omitempty"`
	IsVegetarian bool `json:"IsVegetarian"`
	IsNonVegetarian bool `json:"IsNonVegetarian"`
	IsVegan bool `json:"IsVegan"`
}

type FilterOptions struct {
	CountryOfOrigin []string `json:"CountryOfOrigin"`
	TypeOfDish      []string `json:"TypeOfDish"`
	MinPrice        float64  `json:"minPrice"`
	MaxPrice        float64  `json:"maxPrice"`
}

type RegisterUser struct {
    UserName     string `json:"userName"`
    Email        string `json:"email"`
    PasswordHash string `json:"passwordHash"`
}

type LoginUser struct {
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
	addProductQuery = "INSERT INTO DataTable (DishName, CountryOfOrigin, TypeOfDish, Price, IsVegetarian, IsNonVegetarian, IsVegan) VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7);"
	listProductsQuery = "select DishName, CountryOfOrigin, TypeOfDish, Price, IsVegetarian, IsNonVegetarian, IsVegan from DataTable WHERE Status='APPROVED';"
)


func addProductHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Function called")
	var dish Dish
	err := json.NewDecoder(r.Body).Decode(&dish)
	if err != nil {
		fmt.Println("Hello I am here")
		fmt.Println("The error is", err)
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(addProductQuery, dish.DishName, dish.CountryOfOrigin, dish.TypeOfDish, dish.Price, dish.IsVegetarian, dish.IsNonVegetarian, dish.IsVegan)
    if err != nil {
		fmt.Println("The error is", err)
        http.Error(w, "Error2 in executing query", http.StatusBadRequest)
		return
    }

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dish)
}

func listProductsHandler(w http.ResponseWriter, r *http.Request) {
	var dishes []Dish
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
		var d Dish
		err := rows.Scan(&d.DishName, &d.CountryOfOrigin, &d.TypeOfDish, &d.Price, &d.IsVegetarian, &d.IsNonVegetarian, &d.IsVegan) // Scan directly into struct fields
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

func filterOptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Execute SQL query to retrieve unique values of CountryOfOrigin column
	rows, err := db.Query("SELECT DISTINCT CountryOfOrigin FROM DataTable WHERE Status='APPROVED'")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to execute SQL query for CountryOfOrigin: %v", err)
		return
	}
	defer rows.Close()

	// Extract unique values from the query result for CountryOfOrigin
	var countryOfOriginValues []string
	for rows.Next() {
		var countryOfOrigin string
		if err := rows.Scan(&countryOfOrigin); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to scan query result for CountryOfOrigin: %v", err)
			return
		}
		countryOfOriginValues = append(countryOfOriginValues, countryOfOrigin)
	}
	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error iterating over query result for CountryOfOrigin: %v", err)
		return
	}

	// Execute SQL query to retrieve unique values of TypeOfDish column
	rows, err = db.Query("SELECT DISTINCT TypeOfDish FROM DataTable WHERE Status='APPROVED'")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to execute SQL query for TypeOfDish: %v", err)
		return
	}
	defer rows.Close()

	// Extract unique values from the query result for TypeOfDish
	var typeOfDishValues []string
	for rows.Next() {
		var typeOfDish string
		if err := rows.Scan(&typeOfDish); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to scan query result for TypeOfDish: %v", err)
			return
		}
		typeOfDishValues = append(typeOfDishValues, typeOfDish)
	}
	if err := rows.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error iterating over query result for TypeOfDish: %v", err)
		return
	}

	// Execute SQL queries to retrieve minimum and maximum values in Price column
	var minPrice, maxPrice float64
	err = db.QueryRow("SELECT MIN(Price), MAX(Price) FROM DataTable WHERE STATUS = 'APPROVED'").Scan(&minPrice, &maxPrice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to execute SQL query for min and max Price: %v", err)
		return
	}

	// Create FilterOptions struct with retrieved values
	filterOptions := FilterOptions{
		CountryOfOrigin: countryOfOriginValues,
		TypeOfDish:      typeOfDishValues,
		MinPrice:        minPrice,
		MaxPrice:        maxPrice,
	}

	// Marshal filter options data to JSON
	jsonData, err := json.Marshal(filterOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to marshal filter options data: %v", err)
		return
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func filterProductHandler(w http.ResponseWriter, r *http.Request) {
	var dishes []Dish
	var myMap map[string][]string

	err := json.NewDecoder(r.Body).Decode(&myMap)
	if err != nil {
		fmt.Println("Hello I am here")
		fmt.Println("The error is", err)
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}
	fmt.Println(myMap)

	// Initialize an empty slice to store SQL query conditions
	var conditions []string

	// Loop through each column in userInputMap
	for column, values := range myMap {
		// Check if values are provided for filtering on this column
		if len(values) > 0 {
			var condition []string
			var temp string
			// Generate SQL condition for IN operator based on user's input values
			if(column=="Price"){
				placeholders := make([]string, len(values))
				for i,value := range values {
					placeholders[i] = fmt.Sprintf("%s",value)
					fmt.Println(value)
					k, err := strconv.ParseFloat(placeholders[i],8)
					if err != nil {
						http.Error(w, "Error in conversion of string to int", http.StatusInternalServerError)
					}
					con := fmt.Sprintf("(%s <= %g)", column, k)
					fmt.Println(con)
					condition = append(condition,con)
				}
				temp = strings.Join(condition," OR ")
				final_condition := fmt.Sprintf("(%s)",temp)
				fmt.Println(final_condition)
				// condition := fmt.Sprintf("(%s = '%s')", column, strings.Join(placeholders, ", "))
				conditions = append(conditions, final_condition)
			} else{
				placeholders := make([]string, len(values))
				for i,value := range values {
					placeholders[i] = fmt.Sprintf("%s",value)
					fmt.Println(value)
					con := fmt.Sprintf("(%s = '%s')", column, placeholders[i])
					condition = append(condition,con)
				}
				temp = strings.Join(condition," OR ")
				final_condition := fmt.Sprintf("(%s)",temp)
				fmt.Println(final_condition)
				// condition := fmt.Sprintf("(%s = '%s')", column, strings.Join(placeholders, ", "))
				conditions = append(conditions, final_condition)
			}
		}
	}

	// Combine all conditions using AND operator
	sqlCondition := strings.Join(conditions, " AND ")
	// Construct the final SQL query
	sqlQuery := fmt.Sprintf("SELECT DishName, CountryOfOrigin, TypeOfDish, Price, IsVegetarian, IsNonVegetarian, IsVegan FROM DataTable WHERE Status='APPROVED' AND %s;", sqlCondition)

	fmt.Println("SQL Query is %s",sqlQuery)

	// Prepare the SQL statement
	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
    	http.Error(w, "Error preparing SQL statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the SQL query with user input values
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("Error executing SQL query:", err)
		http.Error(w, "Error executing SQL query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var d Dish
		err := rows.Scan(&d.DishName, &d.CountryOfOrigin, &d.TypeOfDish, &d.Price, &d.IsVegetarian, &d.IsNonVegetarian, &d.IsVegan) // Scan directly into struct fields
		if err != nil {
			fmt.Println("The error is", err)
			http.Error(w, "Issue here : Error2 in executing query", http.StatusInternalServerError)
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user RegisterUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash the password", http.StatusInternalServerError)
		return
	}

	// Insert the user into the database
	_, err = db.Exec("INSERT INTO Users (UserName, Email, PasswordHash) VALUES (@p1, @p2, @p3);", user.UserName, user.Email, passwordHash)
	if err != nil {
		http.Error(w, "Failed to insert user into the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	row := db.QueryRow("SELECT PasswordHash FROM Users WHERE Email = @p1;", user.Email)

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

	http.HandleFunc("/api/addProduct", addProductHandler)
	http.HandleFunc("/api/listProducts", listProductsHandler)
	http.HandleFunc("/api/filterProducts",filterProductHandler)
	http.HandleFunc("/api/filterOptions", filterOptionsHandler)
	http.HandleFunc("/api/register", RegisterHandler)
	http.HandleFunc("/api/login", LoginHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// log.Println("Server is running on http://localhost")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


