package main

import (
	"github.com/kataras/iris/v12"
	indexcontroller "github.com/vadgun/isluad/Controladores/IndexController"
)

func main() {
	app := iris.New()
	// Archivos staticos
	app.HandleDir("/Recursos", "./Recursos")
	app.Favicon("./Recursos/Imagenes/favicon.ico")
	app.RegisterView(iris.HTML("./Vistas", ".html").Reload(true))
	app.Use(iris.Compression)

	// Rutas
	app.Get("/", indexcontroller.Index)
	app.Get("/servicios", indexcontroller.Servicios)
	app.Get("/paquetes", indexcontroller.Paquetes)
	app.Get("/reservaciones", indexcontroller.Reservaciones)
	app.Get("/cotizaciones", indexcontroller.Cotizaciones)

	// Server
	app.Listen(":8080")
}
