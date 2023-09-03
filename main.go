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
	Id        int
	Name      string
	Completed bool
}

func main() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Println("Failed to create DB: ", err)
		return
	}
	defer db.Close()

	// intialize the database
	initStatement, err := db.Prepare("create table if not exists task (id integer primary key AUTOINCREMENT, name text, completed boolean)")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = initStatement.Exec()
	if err != nil {
		fmt.Println(err)
		return
	}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./public")

	app.Get("/get-tasks", func(c *fiber.Ctx) error {
		// template for get tasks statement
		var tasks []Task
		var id int
		var name string
		var completed bool
		rows, err := db.Query("select id, name, completed from task")
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&id, &name, &completed)
			if err != nil {
				fmt.Println(err)
				return err
			}
			newTask := Task{
				Id:        id,
				Name:      name,
				Completed: completed,
			}
			tasks = append(tasks, newTask)
		}
		return c.Render("partials/task", fiber.Map{
			"Tasks": tasks,
		})
	})

	app.Post("/add-todo", func(c *fiber.Ctx) error {
		taskName := c.FormValue("name")
		insertTaskStatement, err := db.Prepare("insert into task (name, completed) values (?, false)")
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = insertTaskStatement.Exec(taskName)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return c.Render("partials/task", fiber.Map{
			"Name": taskName,
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {

		var tasks []Task
		var id int
		var name string
		var completed bool
		rows, err := db.Query("select id, name, completed from task")
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&id, &name, &completed)
			if err != nil {
				fmt.Println(err)
				return err
			}
			newTask := Task{
				Id:        id,
				Name:      name,
				Completed: completed,
			}
			tasks = append(tasks, newTask)
		}

		return c.Render("index", fiber.Map{
			"Tasks": tasks,
		}, "layouts/main")
	})

	log.Fatal(app.Listen(":3000"))
}
