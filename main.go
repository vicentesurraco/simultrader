package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"database/sql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	db "github.com/vicentesurraco/simutrader2/internal/database"
	"golang.org/x/crypto/bcrypt"
	"os"
)

var _, _ = fmt.Print("")
var queries *db.Queries

type User struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" validate:"required,min=3,max=30" db:"unique"`
	Password     string    `json:"password" validate:"required,min=8"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email" validate:"required,email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SubStonkReq struct {
	UserID int32  `json:"user_id"`
	Symbol string `json:"symbol"`
}

type Stonk struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type Position struct {
	Symbol string  `json:"symbol"`
	Shares int     `json:"shares"`
	Price  float64 `json:"price"`
}

type Trade struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Symbol    string    `json:"symbol"`
	Action    string    `json:"action"` // buy or sell
	Shares    int       `json:"shares"`
	Price     float64   `json:"price"` // average price
	Total     float64   `json:"total"`
	Timestamp time.Time `json:"timestamp"`
}

type Portfolio struct {
	UserID     int                 `json:"user_id"`
	Cash       float64             `json:"cash"`
	Positions  map[string]Position `json:"positions"` // key: symbol
	TotalValue float64             `json:"total_value"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("env did not load")
	}
	connStr := os.Getenv("DATABASE_URL")
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Errorf((err.Error()))
	}
	queries = db.New(conn)
	if err := openFinnhubWebsocket(); err != nil {
		panic("failed to open websocket")
	}
}

func openFinnhubWebsocket() error {
	return nil
}

// TODO: jwt token, rate limiting, input validation, logging
func loginUser(c echo.Context) error {
	var user = User{}
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	ctx := c.Request().Context()

	dbUser, err := queries.GetUser(ctx, user.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(user.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": dbUser.ID,
		"message": "Login successful",
	})
}

func createUser(c echo.Context) error {
	var user = User{}
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to create user early"})
	}
	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.PasswordHash = string(hashedPassword)
	user.Password = ""
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to encrypt password"})
	}

	ctx := c.Request().Context()
	_, err = queries.CreateUser(ctx, db.CreateUserParams{
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		Email:        user.Email,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to create user at db step: %s", err.Error()),
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
	// check if user is authenticated
	// delete from db
	return c.JSON(http.StatusOK, id)
}

func subStonk(c echo.Context) error {

	ctx := c.Request().Context()

	var subStonkReq SubStonkReq
	if err := c.Bind(&subStonkReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	stonk := db.SubStonkParams{UserID: subStonkReq.UserID, Symbol: subStonkReq.Symbol}
	if err := queries.SubStonk(ctx, stonk); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// add to db
	// open websocket

	return c.JSON(http.StatusOK, "")
}

func unsubStonk(c echo.Context) error {
	id := c.Param("id")

	// remove from db
	// close websocket

	return c.JSON(http.StatusOK, id)
}

func tradeStonk(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to trade stonk"})
	}

	data := make(map[string]interface{})

	if err = json.Unmarshal(body, &data); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to trade stonk"})
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

	// user routes
	e.POST("/login", loginUser)
	e.POST("/users", createUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	// sub/unsub routes
	e.POST("/stonks", subStonk)
	e.DELETE("/stonks/:id", unsubStonk)

	// buy/sell routes
	e.POST("/trade", tradeStonk)
	e.Logger.Fatal(e.Start(":1323"))

}

// account setup
// sub / unsub -> rabbitmq
// open websocket for data
// "trade" - buy / sell at cur price -> logged in db
// calculated p&l
