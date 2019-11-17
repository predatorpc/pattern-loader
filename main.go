package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"


	//"github.com/jimlawless/whereami"
	"io"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"os"
	"path/filepath"
	"fileloader/models"
	"fileloader/utils"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var Config []models.GeneralConfigFile

func init() {
	fds, _ := filepath.Glob("conf.d/*.json")
	if len(fds) > 0 {
		for _, item := range fds {
			//fdsReplace := filepath.Base(item)

			// прочитываем весь файл в буфер по 32кб
			file, _ := os.Open(item)
			w := bytes.NewBuffer(nil)
			_, _ = io.Copy(w, file)
			_ = file.Close()

			var tmp models.GeneralConfigFile
			_ = json.Unmarshal(w.Bytes(), &tmp)

			Config = append(Config, tmp)
		}
	}

	spew.Dump(Config)
}

func main() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/:integration", func(c echo.Context) error {

		var integration = c.Param("integration")

		if (len(Config) > 0) {
			for _, v:= range Config {
				if v.IntegrationName == integration {

					keyMap := buildKeyMap(v.Params)
					fmt.Println("KEYMAP IS = ", keyMap)

					resultMap, _ := utils.URIByMap(c, keyMap)
					fmt.Println("RESULT MAP IS = ", resultMap)

					return c.String(http.StatusOK, spew.Sdump(v) + spew.Sdump(resultMap))
				}
			}
		}

		return c.String(http.StatusOK, "Nothing found!")
	})

	// Start server
	e.Logger.Fatal(e.Start(":80"))
}

func buildKeyMap(entry interface{}) []string {
	var keyMap []string
	fields := reflect.TypeOf(entry)
	values := reflect.ValueOf(entry)

	num := fields.NumField()

	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		if  value.Bool() == true {
			keyMap = append(keyMap, field.Name)
		}
		//fmt.Print("Type:", field.Type, ",", field.Name, "=", value, "\n")
	}
	return keyMap
}