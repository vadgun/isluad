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
	app.Get("/login", indexcontroller.Login)

	app.Post("/services", indexcontroller.Dashboard)
	app.Post("/categoria", indexcontroller.Categoria)

	app.Post("/administracion", indexcontroller.Administracion)
	app.Post("/guardarcotizacion", indexcontroller.GuardarCotizacion)

	// Server
	app.Listen(":8080")
}
