package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"soumik-serverless/constants"
	"strconv"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	model "soumik-serverless/postgres/public/model"
	table "soumik-serverless/postgres/public/table"
)

// User represents a user in our system
type User struct {
	ID   int64  `json:"id" sql:"primary_key"`
	Name string `json:"name" sql:"not null"`
	Age  int    `json:"age" sql:"not null"`
}

var db *sql.DB

func main() {
	// Connect to the database
	host := os.Getenv(constants.POSTGRES_HOST)
	port := os.Getenv(constants.POSTGRES_PORT)
	user := os.Getenv(constants.POSTGRES_USER)
	password := os.Getenv(constants.POSTGRES_PASSWORD)
	database := os.Getenv(constants.POSTGRES_DATABASE)

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	e := echo.New()

	// Routes
	e.POST("/users", createUser)
	e.GET("/users", getUsers)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// createUser handles POST /users
func createUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	var err error
	stmt := SELECT(table.Users.ID, table.Users.Age, table.Users.Name).FROM(table.Users)
	var dest []struct {
		model.Users
	}

	err = stmt.Query(db, &dest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dest)
}

func getUsers(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	var err error
	stmt := SELECT(table.Users.ID, table.Users.Age, table.Users.Name).FROM(table.Users)
	var dest []struct {
		model.Users
	}

	err = stmt.Query(db, &dest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, dest)
}

// getUser handles GET /users/:id
func getUser(c echo.Context) error {
	strconv.ParseInt(c.Param("id"), 10, 64)

	var user User
	var err error
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// updateUser handles PUT /users/:id
func updateUser(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	var err error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	u.ID = id
	return c.JSON(http.StatusOK, u)
}

// deleteUser handles DELETE /users/:id
func deleteUser(c echo.Context) error {
	strconv.ParseInt(c.Param("id"), 10, 64)

	var err error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
