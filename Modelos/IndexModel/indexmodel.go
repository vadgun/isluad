package indexmodel

import (
	"context"
	"log"
	"strconv"
	"time"

	conexiones "github.com/vadgun/isluad/Conexiones"
	db "github.com/vadgun/isluad/Modelos/Db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Cotizacion struct {
	Nombre        string           `json:"nombre" bson:"Nombre"`
	Correo        string           `json:"correo" bson:"Correo"`
	Invitados     string           `json:"invitados" bson:"Invitados"`
	Fecha         string           `json:"fecha" bson:"Fecha"`
	Telefono      string           `json:"telefono" bson:"Telefono"`
	Agregados     map[string]Extra `json:"agregados" bson:"Agregados"`
	FechaReserva  time.Time        `json:"fechaReserva" bson:"FechaReserva"`
	FechaCreacion time.Time        `json:"fechaCreacion" bson:"FechaCreacion"`
	Activa        bool             `json:"activa" bson:"Activa"`
	Revisada      bool             `json:"revisada" bson:"Revisada"`
	CostoTotal    float64          `json:"costoTotal" bson:"CostoTotal"`
}

type Reservaciones struct {
	FechaReserva  time.Time          `json:"fechaReserva" bson:"FechaReserva"`
	FechaCreacion time.Time          `json:"fechaCreacion" bson:"FechaCreacion"`
	CotizacionID  primitive.ObjectID `bson:"CotizacionID,omitempty"`
}

type Extra struct {
	Comentario      string  `json:"Comentario" bson:"Comentario"`
	CostoIndividual float64 `json:"CostoIndividual" bson:"CostoIndividual"`
}

type Servicio struct {
	Service     string  `json:"Service" bson:"Service"`
	Titulo      string  `json:"Titulo" bson:"Titulo"`
	Categoria   string  `json:"Categoria" bson:"Categoria"`
	Descripcion string  `json:"Descripcion" bson:"Descripcion"`
	Costo       float64 `json:"costo" bson:"Costo"`
	Activo      bool    `json:"Activo" bson:"Activo"`
}

type Dashboard struct {
	ProximoEvento      string `json:"ProximoEvento"`
	ActiveServices     string `json:"ActiveServices"`
	CotizacionesSinVer string `json:"CotizacionesSinVer"`
	Reservaciones      string `json:"Reservaciones"`
	Eventos            string `json:"Eventos"`
}

func GetDashboard() Dashboard {

	//Obtener servicios
	serviciosChannel := make(chan []Servicio)
	defer close(serviciosChannel)
	go func() {
		servicios, err := GetServices()
		if err != nil {
			log.Println("Error al obtener servicios:", err)
			serviciosChannel <- nil
			return
		}
		serviciosChannel <- servicios
	}()
	servicios := <-serviciosChannel
	serviciosInt := strconv.Itoa(len(servicios))

	//Obtener cotizaciones
	cotizacionesChannel := make(chan []Cotizacion)
	defer close(cotizacionesChannel)
	go func() {
		cotizaciones, err := GetCotizacionesSinVer()
		if err != nil {
			log.Println("Error al obtener servicios:", err)
			cotizacionesChannel <- nil
			return
		}
		cotizacionesChannel <- cotizaciones
	}()
	cotizacionesSinVer := <-cotizacionesChannel
	cotizacionesSinVerInt := strconv.Itoa(len(cotizacionesSinVer))

	//Obtener reservaciones
	// reservacionesChannel := make(chan []Reservaciones)
	// defer close(reservacionesChannel)
	// go func() {
	// 	cotizaciones, err := GetCotizacionesSinVer()
	// 	if err != nil {
	// 		log.Println("Error al obtener servicios:", err)
	// 		cotizacionesChannel <- nil
	// 		return
	// 	}
	// 	cotizacionesChannel <- cotizaciones
	// }()
	// cotizacionesSinVer := <-cotizacionesChannel
	// cotizacionesSinVerInt := strconv.Itoa(len(cotizacionesSinVer))

	dashboard := Dashboard{
		ProximoEvento:      "12",
		ActiveServices:     serviciosInt,
		CotizacionesSinVer: cotizacionesSinVerInt,
		Reservaciones:      "4",
		Eventos:            "7",
	}
	return dashboard
}

func GetServices(s ...string) ([]Servicio, error) {
	var servicios []Servicio
	//obtener los servicios de la base de datos
	client, _ := db.ConectarMongoDB()
	serviciosCollection := serviciosCollection(client)

	filter := bson.M{}
	if len(s) > 0 && s[0] != "Todos los Servicios" {
		filter = bson.M{"Categoria": s[0]}
	}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := serviciosCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Println(err, "Error en Buscar todos los servicios")
	}

	if err = cursor.All(context.TODO(), &servicios); err != nil {
		log.Println(err, "Error en retornar la consulta de servicios ")
	}

	return servicios, nil
}

func serviciosCollection(client *mongo.Client) *mongo.Collection {
	servicios := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_S)
	return servicios
}

func cotizacionesCollection(client *mongo.Client) *mongo.Collection {
	servicios := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_C)
	return servicios
}

func GuardarCotizacion(cotizacion Cotizacion) bool {
	client, _ := db.ConectarMongoDB()
	cotizacionesCollection := cotizacionesCollection(client)

	res, err := cotizacionesCollection.InsertOne(context.TODO(), cotizacion)
	if err != nil {
		log.Println(err, "Error creando la cotizacion")
		return false
	}
	log.Printf("Cotizacion creada : %v, %v\n", cotizacion.Nombre, res.InsertedID)
	return true
}

func GetCotizacionesSinVer() ([]Cotizacion, error) {
	var cotizaciones []Cotizacion
	//obtener los servicios de la base de datos
	client, _ := db.ConectarMongoDB()
	cotizacionesCollection := cotizacionesCollection(client)

	filter := bson.M{"Activa": true}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := cotizacionesCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Println(err, "Error en buscar las cotizaciones sin ver")
	}

	if err = cursor.All(context.TODO(), &cotizaciones); err != nil {
		log.Println(err, "Error en retornar la consulta de cotizaciones sin ver")
	}

	return cotizaciones, nil
}

// cotizaciones, err := GetCotizacionesSinVer()
