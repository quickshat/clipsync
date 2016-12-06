package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"runtime"

	"github.com/labstack/echo"

	clip "github.com/atotto/clipboard"
)

type metrics struct {
	Hostname string
	UID      int
	ID       int
	Username string
	Name     string
	OS       string
}

func initWeb() {
	e := echo.New()

	e.POST("/send", postData)

	e.GET("/metrics", getMetrics)
	e.GET("/discoveredDevices", getMetrics)

	e.Logger.Fatal(e.Start(":8081"))
}

func postData(c echo.Context) error {
	data, _ := ioutil.ReadAll(c.Request().Body)
	clip.WriteAll(string(data))
	return c.String(http.StatusOK, "")
}

func getMetrics(c echo.Context) error {
	hn, _ := os.Hostname()
	u, _ := user.Current()

	metrics := metrics{hn, os.Getuid(), os.Getgid(), u.Username, u.Name, runtime.GOOS}

	return c.JSON(http.StatusOK, metrics)
}

func getDiscoveredDevices(c echo.Context) error {
	return c.JSON(http.StatusOK, activeDevices)
}
