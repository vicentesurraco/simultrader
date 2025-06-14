package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type User struct {
	Name string `json:"username"`
}

type Stonk struct {
	Name string `json:"stonk"`
}

type Trade struct {
	Name   string `json:"stonk"`
	Action string `json:"action"`
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, id)
}

func saveUser(c echo.Context) error {
	return c.JSON(http.StatusOK, "User saved.")
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
		return c.JSON(http.StatusBadRequest, "Invalid JSON")
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
		return err
	}

	data := make(map[string]interface{})

	if err = json.Unmarshal(body, &data); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid JSON")
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
	e.POST("/users", saveUser)
	e.GET("/users/:id", getUser)
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
