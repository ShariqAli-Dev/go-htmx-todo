package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// create a new engine
	engine := html.New("./views", ".html")

	// pass the engine to views
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	// this is for things like assets or css or js files
	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello World",
		}, "layouts/main")
	})
	// app.Get("/layout", func(c *fiber.Ctx) error {
	// 	return c.Render("index", fiber.Map{
	// 		"Title": "Hello, World",
	// 	}, "layouts/main")
	// })

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
