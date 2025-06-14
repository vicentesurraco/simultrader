package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"database/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	db "github.com/vicentesurraco/simutrader2/internal/database"
	"golang.org/x/crypto/bcrypt"
	"os"
)

var _, _ = fmt.Print("")
var queries *db.Queries

type User struct {
	Name         string    `json:"name" validate:"required,min=3,max=30"`
	Password     string    `json:"password" validate:"required,min=8"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email" validate:"required"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Stonk struct {
	Name string `json:"stonk"`
}

type Trade struct {
	Name   string `json:"stonk"`
	Action string `json:"action"`
}

func init() {
	connStr := os.Getenv("DATABASE_URL")
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Errorf((err.Error()))
	}
	queries = db.New(conn)
}

func loginUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, id)
}

func createUser(c echo.Context) error {
	var user = User{}
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Failed to create user")
	}
	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.PasswordHash = string(hashedPassword)
	user.Password = ""
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Failed to encrypt password")
	}

	ctx := c.Request().Context()
	_, err = queries.CreateUser(ctx, db.CreateUserParams{
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		Email:        user.Email,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}
	return c.JSON(http.StatusCreated, map[string]string{
		"message":  "User created successfully",
		"username": user.Name,
	})
}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, id)
}

func deleteUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, id)
}

func saveStonk(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})

	if err = json.Unmarshal(body, &data); err != nil {
		return c.JSON(http.StatusBadRequest, "Failed to save stonk")
	}

	// add to db
	// open websocket

	return c.JSON(http.StatusOK, "")
}

func deleteStonk(c echo.Context) error {
	id := c.Param("id")

	// remove from db
	// close websocket

	return c.JSON(http.StatusOK, id)
}

func tradeStonk(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "failed to trade stonk")
	}

	data := make(map[string]interface{})

	if err = json.Unmarshal(body, &data); err != nil {
		return c.JSON(http.StatusBadRequest, "failed to trade stonk")
	}

	// remove from db
	// setup cloud run postgres instance
	// setup websocket
	// close websocket

	return c.JSON(http.StatusOK, "")
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))

	// user routes
	e.POST("/users", createUser)
	e.GET("/users/:id", loginUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	// sub/unsub routes
	e.POST("/stonks", saveStonk)
	e.DELETE("/stonks/:id", deleteStonk)

	// buy/sell routes
	e.POST("/trade", tradeStonk)

}

// account setup
// sub / unsub -> rabbitmq
// open websocket for data
// "trade" - buy / sell at cur price -> logged in db
// calculated p&l
