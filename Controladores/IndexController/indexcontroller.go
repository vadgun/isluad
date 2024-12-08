package indexcontroller

import (
	"fmt"
	"html/template"
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
	//Obtener todos los servicios
	ctx.ViewData("dashboard", indexmodel.GetDashboard())
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

func Categoria(ctx iris.Context) {
	cadena := ctx.PostValue("data")

	servicios, err := indexmodel.GetServices(cadena)
	if err != nil {
		log.Println(err)
	}

	htmlcode := crearTabla(servicios)
	ctx.HTML(htmlcode)
}

func Dashboard(ctx iris.Context) {
	cadena := ctx.PostValue("data")
	tmpContentHeader, err := template.ParseFiles("./Vistas/templates.html")
	if err != nil {
		log.Println(err)
	}

	data := struct {
		Cadena string
	}{
		Cadena: cadena,
	}

	err = tmpContentHeader.Execute(ctx.ResponseWriter(), data)
	if err != nil {
		log.Println(err)
	}

	htmlcode := crearHTML(cadena)

	ctx.HTML(htmlcode)
	//Responder al cliente
	// 	ctx.JSON(iris.Map{
	// 		"status":  cadena,
	// 		"message": 55,
	// 	})
}

func crearHTML(data string) string {

	var htmlcode string

	switch data {
	case "Panel":
		values := indexmodel.GetDashboard()
		htmlcode += fmt.Sprintf(`<div class="container">
            <div class="row">
                <div class="col-12 col-sm-4 mb-4 py-1 rectangle">
                    <h1>%v</h1>
                    <span class="icon">
                        <i class="bi bi-alarm"></i>
                    </span>
                    <span>Días próximo evento</span>
                </div>
                <div class="col-12 col-sm-4 mb-4 py-1 rectangle">
                    <h1>%v</h1>
                    <span class="icon">
                        <i class="bi bi-person-lines-fill"></i>
                    </span>
                    <span>Servicios activos</span>
                </div>
                <div class="col-12 col-sm-4 mb-4 py-1 rectangle">
                    <h1>%v</h1>
                    <span class="icon">
                        <i class="bi bi-exclamation-circle"></i>
                    </span>
                    <span>Cotizaciones sin revisar</span>
                </div>
                <div class="col-12 col-sm-4 mb-4 py-1 rectangle">
                    <h1>%v</h1>
                    <span class="icon">
                        <i class="bi bi-calendar-check-fill"></i>
                    </span>
                    <span>Reservaciones</span>
                </div>
                <div class="col-12 col-sm-4 mb-4 py-1 rectangle">
                    <h1>%v</h1>
                    <span class="icon">
                        <i class="bi bi-calendar2-check"></i>
                    </span>
                    <span>Eventos por realizar</span>
                </div>
            </div>
        </div>`, values.ProximoEvento, values.ActiveServices, values.CotizacionesSinVer, values.Reservaciones, values.Eventos)
	case "Servicios":
		servicios, _ := indexmodel.GetServices()
		categorias := make(map[string]int)

		for _, v := range servicios {
			categorias[v.Categoria]++
		}

		htmlcode += `<div class="container fluid"><div class="row text-white"><span>Categorias disponibles:</span><div class="d-flex py-2">`

		for k, v := range categorias {
			htmlcode += fmt.Sprintf(`<a href="Javascript:Categoria('%v')" class="btn btn-sm btn-outline-light px-2 mx-2" role="button">%v %v</a>`, k, k, v)

		}
		htmlcode += fmt.Sprintf(`<a href="Javascript:Categoria('Todos los Servicios')" class="btn btn-sm btn-outline-light px-2 mx-2" role="button">Todos los Servicios %v</a>`, len(servicios))
		htmlcode += `</div></div></div><div class="container fluid"><div id="tables" class="row text-white"></div></div>`
	case "Cotizaciones":
	case "Reservaciones":
	case "Calendario":
		htmlcode += fmt.Sprintf(`<div class="container fluid calendario">
            <div class="calendar-container">
                <div class="calendar-header">
                  <button id="prev-month">◀</button>
                  <span class="month-year" id="month-year">
                    <span id="month-label"></span> <span id="year-label"></span>
                  </span>
                  <button id="next-month">▶</button>
                </div>
                <div class="month-selector" id="month-selector">
                  <div class="month-options">
                    <span data-month="0">Enero</span><span data-month="1">Febrero</span>
                    <span data-month="2">Marzo</span><span data-month="3">Abril</span>
                    <span data-month="4">Mayo</span><span data-month="5">Junio</span>
                    <span data-month="6">Julio</span><span data-month="7">Agosto</span>
                    <span data-month="8">Septiembre</span><span data-month="9">Octubre</span>
                    <span data-month="10">Noviembre</span><span data-month="11">Diciembre</span>
                  </div>
                </div>
                <input type="date" id="date-input">
                <div class="calendar-body">
                  <div class="days-header">
                    <span>Dom</span><span>Lun</span><span>Mar</span>
                    <span>Mie</span><span>Jue</span><span>Vie</span><span>Sab</span>
                  </div>
                  <div class="days-grid" id="calendar-days"></div>
                </div>
              </div>
        </div>`)
	}

	return htmlcode
}

func crearTabla(servicios []indexmodel.Servicio) string {

	var htmlcode string

	htmlcode += `<table id="example" class="table table-hover table-striped" style="width:100%">`
	htmlcode += `<thead><tr><th>Nombre</th><th>Costo</th><th>Activo</th><th>Descripcion</th><th>Categoria</th><th>Acciones</th></tr></thead><tbody>`

	for _, v := range servicios {
		htmlcode += `<tr>`
		htmlcode += fmt.Sprintf(`<td>%v</td>`, v.Titulo)
		htmlcode += fmt.Sprintf(`<td class="text-center">%v</td>`, v.Costo)
		if v.Activo {
			htmlcode += `<td class="text-center"><i class="bi bi-check-circle"></i></td>`
		} else {
			htmlcode += `<td class="text-center"><i class="bi bi-dash-circle"></i></td>`
		}
		htmlcode += fmt.Sprintf(`<td>%v</td>`, v.Descripcion)
		htmlcode += fmt.Sprintf(`<td>%v</td>`, v.Categoria)
		htmlcode += `<td><a href="#" class="me-2 text-decoration-none"><i class="bi bi-pencil"></i><span>Editar</span></a></td>`
		htmlcode += `</tr>`

	}

	htmlcode += `</tbody></table>`
	return htmlcode
}
