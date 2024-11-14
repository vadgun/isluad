package indexcontroller

import (
	"log"

	"github.com/kataras/iris/v12"
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
