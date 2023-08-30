package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Transaction struct {
	ID              int
	ItemDescription string
	ItemCost        int
	AddBalance      int
	BalanceAfter    int
}

var tmpl = template.Must(template.ParseGlob("templates/*"))

func main() {
	// ... (other parts of your main function)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHandler)
	// Serve static assets from the "assets" directory
	fs := http.FileServer(http.Dir("assets"))

	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.ListenAndServe("localhost:9990", nil)
}

// ... (indexHandler, updateHandler, and deleteHandler functions here)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/balance")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT transactionid, item_description, item_cost, add_balance, balance_after FROM accounting")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.ItemDescription, &t.ItemCost, &t.AddBalance, &t.BalanceAfter)
		if err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	tmpl.ExecuteTemplate(w, "index.html", transactions)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/balance")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}

		r.ParseForm()
		itemDescription := r.FormValue("item_description")
		itemCost := r.FormValue("item_cost")
		addBalance := r.FormValue("add_balance")
		balanceAfter := r.FormValue("balance_after")

		// Convert form values to appropriate data types
		itemCostInt, err := strconv.Atoi(itemCost)
		if err != nil {
			log.Fatal(err)
		}
		addBalanceInt, err := strconv.Atoi(addBalance)
		if err != nil {
			log.Fatal(err)
		}
		balanceAfterInt, err := strconv.Atoi(balanceAfter)
		if err != nil {
			log.Fatal(err)
		}

		// Perform INSERT operation using the converted values
		insertQuery := "INSERT INTO accounting (item_description, item_cost, add_balance, balance_after) VALUES (?, ?, ?, ?)"
		_, err = db.Exec(insertQuery, itemDescription, itemCostInt, addBalanceInt, balanceAfterInt)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl.ExecuteTemplate(w, "create.html", nil)
}
func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Get the transaction ID from the query parameter
		transactionID := r.URL.Query().Get("transaction_id")

		// Convert the transaction ID to an integer
		transactionIDInt, err := strconv.Atoi(transactionID)
		if err != nil {
			log.Fatal(err)
		}

		// Fetch the transaction data from the database
		db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/balance")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		query := "SELECT transactionid, item_description, item_cost, add_balance, balance_after FROM accounting WHERE transactionid = ?"
		row := db.QueryRow(query, transactionIDInt)

		var t Transaction
		err = row.Scan(&t.ID, &t.ItemDescription, &t.ItemCost, &t.AddBalance, &t.BalanceAfter)
		if err != nil {
			log.Fatal(err)
		}

		tmpl.ExecuteTemplate(w, "update.html", t)
		return
	}

	if r.Method == http.MethodPost {
		// Parse form values
		r.ParseForm()
		transactionID := r.FormValue("transaction_id")
		itemDescription := r.FormValue("item_description")
		itemCost := r.FormValue("item_cost")
		addBalance := r.FormValue("add_balance")
		balanceAfter := r.FormValue("balance_after")

		// Convert form values to appropriate data types
		transactionIDInt, err := strconv.Atoi(transactionID)
		if err != nil {
			log.Fatal(err)
		}
		itemCostInt, err := strconv.Atoi(itemCost)
		if err != nil {
			log.Fatal(err)
		}
		addBalanceInt, err := strconv.Atoi(addBalance)
		if err != nil {
			log.Fatal(err)
		}
		balanceAfterInt, err := strconv.Atoi(balanceAfter)
		if err != nil {
			log.Fatal(err)
		}

		// Perform the update operation using the extracted values
		db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/balance")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		updateQuery := "UPDATE accounting SET item_description = ?, item_cost = ?, add_balance = ?, balance_after = ? WHERE transactionid = ?"
		_, err = db.Exec(updateQuery, itemDescription, itemCostInt, addBalanceInt, balanceAfterInt, transactionIDInt)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Updating transaction:", transactionID)
		// Redirect or display a success message
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl.ExecuteTemplate(w, "update.html", nil)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form values
		r.ParseForm()
		transactionID := r.FormValue("transaction_id")

		// Convert form values to appropriate data types
		transactionIDInt, err := strconv.Atoi(transactionID)
		if err != nil {
			log.Fatal(err)
		}

		// Perform DELETE operation using the converted value
		db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/balance")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		deleteQuery := "DELETE FROM accounting WHERE transactionid = ?"
		_, err = db.Exec(deleteQuery, transactionIDInt)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl.ExecuteTemplate(w, "delete.html", nil)
}
