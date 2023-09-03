package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	Name      string
	Completed bool
}

func queryDB(db *sql.DB) sql.Stmt {
	statement, err := db.Prepare("")
	if err != nil {
		log.Fatalf("Failed to create DB: %v", err)
	}

	return *statement
}

func main() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatalf("Failed to create DB: %v", err)
	}

	// create tables

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {

		tasks := []Task{
			{Name: "wash dishes", Completed: false},
			{Name: "take the dog out for a walk", Completed: true},
			{Name: "water the plants", Completed: false},
		}

		return c.Render("index", fiber.Map{
			"Tasks": tasks,
		}, "layouts/main")
	})

	log.Fatal(app.Listen(":3000"))
}

func database() {
	database, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer database.Close()

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Println(err)
		return
	}

	statement2, err := database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = statement2.Exec("shariq", "ali")
	if err != nil {
		fmt.Println(err)
		return
	}

	var id int
	var firstname string
	var lastname string

	rows, err := database.Query("SELECT id, firstname, lastname FROM people")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &firstname, &lastname)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(id, firstname, lastname)
	}
}
