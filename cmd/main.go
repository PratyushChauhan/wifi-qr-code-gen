package main

import (
	"fmt"
	"html/template"
	"image/color"
	"io"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/skip2/go-qrcode"
)

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Templates struct {
	templates *template.Template
}

func newtemplate() *Templates {
	return &Templates{templates: template.Must(template.ParseGlob("static/*.html"))}
}

func generateQRCode(ssid string, password string) error {
	encryption := "WPA2"

	// Generate the Wi-Fi QR code string
	wifiString := fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", encryption, ssid, password)

	// Generate QR Code
	err := qrcode.WriteColorFile(wifiString, qrcode.Medium, 256, color.Black, color.White, "./static/qr.png")
	return err

}

type Data struct {
	SSID     string
	Password string
	Qr       string
}

type Page struct {
	Data Data
}

func newData(ssid string, password string, qr string) Data {
	return Data{SSID: ssid, Password: password, Qr: qr}
}

func newPage() Page {
	return Page{Data: newData("", "", "")}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = newtemplate()
	e.Static("/images", "images")
	e.Static("/css", "css")
	page := newPage()
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})
	e.POST("/qr", func(c echo.Context) error {
		path := "./qr.png"
		ssid := c.FormValue("ssid")
		password := c.FormValue("password")
		err := generateQRCode(ssid, password)
		if err != nil {
			return c.String(404, "Invalid id")
		}

		data := newData(ssid, password, path)
		page.Data = data
		return c.Render(200, "data", page.Data)
	})

	e.Logger.Fatal(e.Start(":42069"))

}
