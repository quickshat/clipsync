package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"
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

	e.Static("/", "public")

	// Public
	e.POST("/send", postData)
	e.POST("/register", postRegister)
	e.GET("/metrics", getMetrics)

	// Local
	local := e.Group("local")
	local.Use(onlyLocal)
	local.POST("/settings", postSettings)
	local.GET("/discoveredDevices", getDiscoveredDevices)
	local.GET("/settings", getSettings)
	local.GET("/stop", getStop)
	local.GET("/interfaces", getInterfaces)
	local.GET("/logs", getLog)

	e.Logger.Fatal(e.Start(":8081"))
}

func postData(c echo.Context) error {
	group := c.QueryParam("group")
	if group != settings.Group {
		return c.NoContent(http.StatusUnauthorized)
	}
	emitLog("REST", "Clipboard data recieved!")
	data, _ := ioutil.ReadAll(c.Request().Body)
	recievedBoard = data
	clip.WriteAll(string(data))
	return c.NoContent(http.StatusOK)
}

func postRegister(c echo.Context) error {
	port, _ := strconv.Atoi(c.FormValue("port"))
	ip := strings.Split(c.Request().RemoteAddr, ":")[0]

	val, found := activeDevices[ip]
	if found {
		val.LastPing = time.Now()
		val.Port = int64(port)
	} else {
		a := &activeDevice{ip, int64(port), time.Now(), metrics{}}
		err := a.loadMetrics()
		if err == nil {
			emitLog("REST", "register", ip, port)
			activeDevices[ip] = a
		}
	}
	return c.NoContent(http.StatusOK)
}

func getMetrics(c echo.Context) error {
	hn, _ := os.Hostname()
	u, _ := user.Current()

	metrics := metrics{hn, os.Getuid(), os.Getgid(), u.Username, u.Name, runtime.GOOS}

	return c.JSON(http.StatusOK, metrics)
}

func postSettings(c echo.Context) error {
	i := c.FormValue("Interface")
	webPort := c.FormValue("webPort")
	p, err := strconv.Atoi(webPort)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	settings.Interface = i
	settings.WebPort = p
	return c.NoContent(http.StatusOK)
}

func getDiscoveredDevices(c echo.Context) error {
	return c.JSON(http.StatusOK, activeDevices)
}

func getSettings(c echo.Context) error {
	return c.JSON(http.StatusOK, settings)
}

func getStop(c echo.Context) error {
	saveSettings()
	go func() {
		time.Sleep(time.Second * 2)
		os.Exit(0)
	}()
	return c.NoContent(http.StatusOK)
}

func getInterfaces(c echo.Context) error {
	return c.JSON(http.StatusOK, availableInterfaces)
}

func getLog(c echo.Context) error {
	return c.JSON(http.StatusOK, logger.Logs)
}
