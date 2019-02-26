package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.GET("/check/:id/:sheet", check)

	log.Printf("Listening on port %s", port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

const (
	rowStart   = 3
	columnMail = 3 // D
)

var (
	mailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$")
)

func validateMail(v interface{}) bool {
	return mailRegexp.MatchString(fmt.Sprintf("%s", v))
}

func check(c echo.Context) error {

	sheet, err := getSheet(c.Param("id"), c.Param("sheet"))
	if err != nil {
		return err
	}

	for i, clms := range sheet {
		if i < rowStart-1 || len(clms) < columnMail+1 {
			continue
		}
		if !validateMail(clms[columnMail]) {
			return c.String(http.StatusOK, fmt.Sprintf("invalid email address: row=%d, value=%s", i+1, clms[columnMail]))
		}
	}

	return c.NoContent(http.StatusOK)
}

func getSheet(id, sheet string) ([][]interface{}, error) {

	config, err := google.ConfigFromJSON([]byte(os.Getenv("OAUTH_CREDENTIALS")), "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, fmt.Errorf("Unable to parse json to config: %v", err)
	}

	tok := &oauth2.Token{}
	if err := json.NewDecoder(strings.NewReader(os.Getenv("OAUTH_TOKEN"))).Decode(tok); err != nil {
		return nil, fmt.Errorf("Unable to parse json to token: %v", err)
	}

	srv, err := sheets.New(config.Client(context.Background(), tok))
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(id, sheet).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet: %v", err)
	}

	return resp.Values, nil
}
