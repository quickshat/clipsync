package main

import "github.com/labstack/echo"
import "net/http"
import "io/ioutil"
import "os"
import "os/user"
import "runtime"
import clip "github.com/atotto/clipboard"

func initWeb() {
	e := echo.New()

	e.POST("/send", postData)

	e.Logger.Fatal(e.Start(":6666"))
}

func postData(c echo.Context) error {
	data, _ := ioutil.ReadAll(c.Request().Body)
	clip.WriteAll(string(data))
	return c.String(http.StatusOK, "")
}

func getMetrics(c echo.Context) error {
	hn, _ := os.Hostname()
	u, _ := user.Current()

	metrics := struct {
		Hostname string
		UID      int
		ID       int
		Username string
		Name     string
		OS       string
	}{hn, os.Getuid(), os.Getgid(), u.Username, u.Name, runtime.GOOS}
	return c.JSON(http.StatusOK, metrics)
}
