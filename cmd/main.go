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
}

type Page struct {
	Data Data
}

func newData() Data {
	return Data{ssid: "", password: ""}
}

func newPage() Page {
	return Page{Data: newData()}
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
		generateQRCode(ssid, password)
		if page.Data.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "This email already exists"
			return c.Render(422, "form", formData)
		}
		contact := newContact(name, email)
		page.Data.Contacts = append(page.Data.Contacts, contact)
		c.Render(200, "form", page.Form)
		return c.Render(200, "oob-contact", contact)
	})

	e.Logger.Fatal(e.Start(":42069"))

}
