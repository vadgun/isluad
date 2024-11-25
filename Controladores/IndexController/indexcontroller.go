package indexcontroller

import (
	"log"
	"time"

	"github.com/kataras/iris/v12"
	indexmodel "github.com/vadgun/isluad/Modelos/IndexModel"
)

func Index(ctx iris.Context) {
	log.Println("Index view served")
	if err := ctx.View("index.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

func Servicios(ctx iris.Context) {
	log.Println("Servicios view served")
	if err := ctx.View("servicios.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

func Login(ctx iris.Context) {
	log.Println("Login view served")
	if err := ctx.View("login.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

func Administracion(ctx iris.Context) {
	log.Println("Administracion view served")
	if err := ctx.View("administracion.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

func Paquetes(ctx iris.Context) {
	log.Println("Paquetes view served")
	if err := ctx.View("paquetes.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

func Reservaciones(ctx iris.Context) {
	log.Println("Reservaciones view served")
	if err := ctx.View("reservaciones.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

func Cotizaciones(ctx iris.Context) {
	log.Println("Cotizaciones view served")
	if err := ctx.View("cotizaciones.html"); err != nil {
		ctx.Application().Logger().Infof(err.Error())
	}
}

func GuardarCotizacion(ctx iris.Context) {
	//Obtener los servicios disponibles en la base de datos para compararlos, guardarlos, y enviarlos por correo.
	serviciosChannel := make(chan []indexmodel.Servicio)
	go func() {
		servicios, err := indexmodel.GetServices()
		if err != nil {
			log.Println("Error al obtener servicios:", err)
			serviciosChannel <- nil
			return
		}
		serviciosChannel <- servicios
	}()
	servicios := <-serviciosChannel

	var cotizacion indexmodel.Cotizacion
	if err := ctx.ReadJSON(&cotizacion); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{
			"error":  "Error al leer los datos",
			"detail": err.Error(),
		})
		return
	}
	// Procesar los datos de la cotización
	cotizacion.Activa = true

	//Comparar los servicios que eligio el usuario
	for _, servicio := range servicios {
		if _, ok := cotizacion.Agregados[servicio.Service]; ok {
			cotizacion.CostoTotal = servicio.Costo + cotizacion.CostoTotal

			extras := indexmodel.Extra{
				CostoIndividual: servicio.Costo,
				Comentario:      cotizacion.Agregados[servicio.Service].Comentario,
			}

			cotizacion.Agregados[servicio.Service] = extras
		}
	}

	//Fecha -> Reserva, Creación
	layout := "2006-01-02"
	location, _ := time.LoadLocation("America/Mexico_City")
	cotizacion.FechaReserva, _ = time.ParseInLocation(layout, cotizacion.Fecha, location)
	cotizacion.FechaCreacion = time.Now()

	//Guardar cotizacion
	log.Printf("Cotización recibida: %v\n", string(cotizacion.FechaCreacion.AppendFormat([]byte("Hora: "), time.Kitchen)))

	if !indexmodel.GuardarCotizacion(cotizacion) {
		log.Printf("Cotizacion no guardada")
	}

	//Responder al cliente
	ctx.JSON(iris.Map{
		"status":  "success",
		"message": "Gracias " + cotizacion.Nombre,
	})
}
