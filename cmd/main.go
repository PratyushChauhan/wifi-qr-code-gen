package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/skip2/go-qrcode"
)

type Templates struct {
	templates *template.Template
}

func newtemplate() *Templates {
	return &Templates{templates: template.Must(template.ParseGlob("views/*.html"))}
}

func generateQRCode(ssid string, password string) ([]byte, error) {
	encryption := "WPA2"

	// Generate the Wi-Fi QR code string
	wifiString := fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;;", encryption, ssid, password)

	// Generate QR Code
	qr, err := qrcode.Encode(wifiString, qrcode.Medium, 256)
	return qr, err

}

type FormData struct {
	data Data
	//	Errors map[string]string
}
type Data struct {
	ssid     string
	password string
	qr       []byte
}

type Page struct {
	Data Data
}

func newData(ssid string, password string, qr []byte) Data {
	return Data{ssid: "", password: "", qr: qr}
}

func newPage() Page {
	return Page{Data: newData("", "", nil)}
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
		ssid := c.FormValue("ssid")
		password := c.FormValue("password")
		qr, err := generateQRCode(ssid, password)
		if err != nil {
			return c.String(404, "Invalid id")
		}

		data := newData(ssid, password, qr)
		page.Data = data
		c.Render(200, "qrData", page.Data)
		return c.Render(200, "oob-contact", data)
	})

	e.Logger.Fatal(e.Start(":42069"))

}
