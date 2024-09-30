package main

import (
	"embed"
	"io"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/Younes-khadraoui/starter/handlers"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//go:embed all:views
var ViewFS embed.FS

//go:embed all:static
var StaticFS embed.FS

func main() {
    // Initialize the template renderer
    portfolioT := &Template{
        templates: template.Must(template.New("portfolio").Funcs(sprig.FuncMap()).ParseFS(ViewFS, "views/index.html","views/pages/*.html")),
    }

    server := echo.New()
    server.Logger.SetLevel(log.INFO)
    server.Use(middleware.Logger())
    server.Renderer = portfolioT

    server.GET("/static/*", func(c echo.Context) error {
        path := c.Param("*")
        data, err := StaticFS.ReadFile("static/" + path)
        if err != nil {
            return echo.NewHTTPError(404, "File not found")
        }
    
        // Detect content type
        contentType := "text/plain"
        if strings.HasSuffix(path, ".css") {
            contentType = "text/css"
        } else if strings.HasSuffix(path, ".js") {
            contentType = "application/javascript"
        } else if strings.HasSuffix(path, ".svg") {
            contentType = "image/svg+xml"
        } else if strings.HasSuffix(path, ".png") {
            contentType = "image/png"
        } else if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
            contentType = "image/jpeg"
        }
    
        return c.Blob(200, contentType, data)
    })
    

    server.GET("/", handlers.HandleHomePage)

    if err := server.Start(":8080"); err != nil {
        server.Logger.Fatal(err)
    }
}
