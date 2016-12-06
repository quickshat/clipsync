package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"time"

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
	e.POST("/register", postRegister)

	e.GET("/metrics", getMetrics)
	e.GET("/discoveredDevices", getDiscoveredDevices)

	e.Logger.Fatal(e.Start(":8081"))
}

func postData(c echo.Context) error {
	data, _ := ioutil.ReadAll(c.Request().Body)
	recievedBoard = data
	clip.WriteAll(string(data))
	return c.String(http.StatusOK, "")
}

func postRegister(c echo.Context) error {
	port, _ := strconv.Atoi(c.FormValue("port"))
	ip := c.Request().RemoteAddr

	val, found := activeDevices[ip]
	if found {
		val.LastPing = time.Now()
		val.Port = int64(port)
	} else {
		a := &activeDevice{ip, int64(port), time.Now(), metrics{}}
		err := a.loadMetrics()
		if err != nil {
			log.Println("[REST] register", ip, port)
			activeDevices[ip] = a
		}
	}
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
