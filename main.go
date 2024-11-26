package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
	"context"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Generate a random string of length n
func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}


func main() {
    // Get login from .env 
    if err := godotenv.Load(); err != nil {
	log.Fatal(err)
    }

    d := os.Getenv("DB")
    u := os.Getenv("DB_U")
    p := os.Getenv("DB_P")

    // DB conn
    pg := "postgresql://" + u + ":" + p +
	    "@localhost:5432/" + d + "?sslmode=disable"

    db, err := pgx.Connect(context.Background(), pg)
    
    // Check errors
    if err != nil  {
	log.Fatal(err)
    } else {
	fmt.Println("Connected successfully.")
	fmt.Println()
    }

    // Close DB connection when we're done
    defer db.Close(context.Background())
    
    // Try to ping DB
    if err := db.Ping(context.Background()); err != nil {
	log.Fatal(err)
    }

    /*
     * Create the table
     */
    { 
	q := `
	CREATE TABLE IF NOT EXISTS users (
	    id SERIAL PRIMARY KEY,
	    username TEXT NOT NULL,
	    password TEXT NOT NULL,
	    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(context.Background(), q); err != nil {
	    log.Fatal(err) 
	}
    }

    /*
     * Insert into table
     */
    {
	uname 	:= randSeq(8) 
	pwd   	:= randSeq(8)
	time  	:= time.Unix(time.Now().Unix(), 0).UTC()
	
	var id int

	q 	:= `
	INSERT INTO users 
	    (username, password, created_at) VALUES ($1, $2, $3)
	    RETURNING id`

	if err := db.QueryRow(
	context.Background(), q, uname, pwd, time).Scan(&id); err != nil {
	    log.Fatal(err)
	} else {
	    fmt.Println("Created user with id ", id)
	    fmt.Println()
	}
     }

    /*
     * Query DB
     */
    {
	var (
	    id		int
	    uname	string
	    pwd		string
	    createdAt	time.Time
	)
	
	q := `
	SELECT id, username, password, created_at
	FROM users WHERE id = $1
	`

	if err := db.QueryRow(context.Background(), q, 1).Scan(
	&id, &uname, &pwd, &createdAt); err != nil {
	      
	    log.Fatal(err)
	}
	fmt.Println("Getting first user...")
	fmt.Printf("%3d %s %s %v\n\n", id, uname, pwd, createdAt)
    }	
    
    /*
     * Query db for all users and create user struct
     */
    {
	type user struct {
	    id		int
	    uname	string
	    pwd		string
	    createdAt	time.Time
	}

	q 	:= `
	SELECT id, username, password, created_at FROM users
	`
	rows, err := db.Query(context.Background(), q)
	defer rows.Close()

	if err != nil {
	    log.Fatal(err)
	}

	var users []user
	    
	fmt.Println("Getting all users...")

	for rows.Next() {
	    var u user
	    if err := rows.Scan(
	    &u.id, &u.uname, &u.pwd, &u.createdAt); err != nil {
		log.Fatal(err)
	    }
	    users = append(users, u)
	}

	for _, u := range users {
	    fmt.Printf("%3d %s %s %v\n", u.id, u.uname, u.pwd, u.createdAt)
	}
    }

    /*
     * To delete we just run a SQL-query, for example:
     * _, err := db.Exec(`DELETE FROM users WHERE id = ?`, 1) 
     */
}
