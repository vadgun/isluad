package indexmodel

import (
	"context"
	"errors"
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
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Nombre        string             `json:"nombre" bson:"Nombre"`
	Correo        string             `json:"correo" bson:"Correo"`
	Invitados     string             `json:"invitados" bson:"Invitados"`
	Fecha         string             `json:"fecha" bson:"Fecha"`
	Telefono      string             `json:"telefono" bson:"Telefono"`
	Agregados     map[string]Extra   `json:"agregados" bson:"Agregados"`
	FechaReserva  time.Time          `json:"fechaReserva" bson:"FechaReserva"`
	FechaCreacion time.Time          `json:"fechaCreacion" bson:"FechaCreacion"`
	Activa        bool               `json:"activa" bson:"Activa"`
	Revisada      bool               `json:"revisada" bson:"Revisada"`
	CostoTotal    float64            `json:"costoTotal" bson:"CostoTotal"`
	ReservacionID string             `json:"reservaid" bson:"ReservacionID"`
}

type Reservacion struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Nombre          string             `json:"nombre" bson:"Nombre"`
	Correo          string             `json:"correo" bson:"Correo"`
	Invitados       string             `json:"invitados" bson:"Invitados"`
	Fecha           string             `json:"fecha" bson:"Fecha"`
	Telefono        string             `json:"telefono" bson:"Telefono"`
	Agregados       map[string]Extra   `json:"agregados" bson:"Agregados"`
	FechaReserva    time.Time          `json:"fechaReserva" bson:"FechaReserva"`
	FechaCreacion   time.Time          `json:"fechaCreacion" bson:"FechaCreacion"`
	Activa          bool               `json:"activa" bson:"Activa"`
	Revisada        bool               `json:"revisada" bson:"Revisada"`
	CostoTotal      float64            `json:"costoTotal" bson:"CostoTotal"`
	CotizacionStrID string             `json:"currentid" bson:"CotizacionID"`
}

type Reservaciones struct {
	FechaReserva  time.Time          `json:"fechaReserva" bson:"FechaReserva"`
	FechaCreacion time.Time          `json:"fechaCreacion" bson:"FechaCreacion"`
	CotizacionID  primitive.ObjectID `bson:"cotizacionID,omitempty"`
}

type Extra struct {
	Comentario      string  `json:"comentario" bson:"Comentario"`
	CostoIndividual float64 `json:"costoIndividual" bson:"CostoIndividual"`
	NombreCompleto  string  `json:"nombreCompleto" bson:"NombreCompleto"`
}

type Servicio struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Service     string             `json:"service" bson:"Service"`
	Titulo      string             `json:"titulo" bson:"Titulo"`
	Categoria   string             `json:"categoria" bson:"Categoria"`
	Descripcion string             `json:"descripcion" bson:"Descripcion"`
	Costo       float64            `json:"costo" bson:"Costo"`
	Activo      bool               `json:"activo" bson:"Activo"`
	CurrentID   string             `json:"currentid"`
}

type DiasApartados struct {
	Fecha      string `json:"date"`
	Celebrcion string `json:"label"`
}

type Dashboard struct {
	ProximoEvento      string `json:"proximoEvento"`
	ActiveServices     string `json:"activeServices"`
	CotizacionesSinVer string `json:"cotizacionesSinVer"`
	Reservaciones      string `json:"reservaciones"`
	Eventos            string `json:"eventos"`
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

func GetService(id string) (Servicio, error) {

	var servicio Servicio
	client, _ := db.ConectarMongoDB()
	serviciosCollection := serviciosCollection(client)
	idobj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return servicio, err
	}
	filter := bson.M{"_id": idobj}

	err = serviciosCollection.FindOne(context.TODO(), filter).Decode(&servicio)
	if err != nil {
		log.Println(err, "Error en obtener un solo servicio")
	}

	log.Println("Servicio encontrado")

	return servicio, nil
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

func GetReservaciones(s ...string) ([]Reservacion, error) {
	var reservaciones []Reservacion
	//obtener los servicios de la base de datos
	client, _ := db.ConectarMongoDB()
	reservacionesCollection := reservacionesCollection(client)

	filter := bson.M{}
	if len(s) > 0 && s[0] != "Todos las Reservaciones" {
		filter = bson.M{"_id": s[0]}
	}
	opts := options.Find().SetSort(bson.M{"Nombre": 1})

	cursor, err := reservacionesCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Println(err, "Error en buscar todas las reservaciones")
	}

	if err = cursor.All(context.TODO(), &reservaciones); err != nil {
		log.Println(err, "Error en retornar la consulta de reservaciones ")
	}

	return reservaciones, nil
}

func serviciosCollection(client *mongo.Client) *mongo.Collection {
	servicios := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_S)
	return servicios
}

func cotizacionesCollection(client *mongo.Client) *mongo.Collection {
	cotizaciones := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_C)
	return cotizaciones
}

func reservacionesCollection(client *mongo.Client) *mongo.Collection {
	reservaciones := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_R)
	return reservaciones
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

func GetCotizaciones(s ...string) ([]Cotizacion, error) {
	var cotizaciones []Cotizacion
	//obtener los servicios de la base de datos
	client, _ := db.ConectarMongoDB()
	cotizacionesCollection := cotizacionesCollection(client)

	filter := bson.M{}
	if len(s) > 0 {
		idobj, err := primitive.ObjectIDFromHex(s[0])
		if err != nil {
			log.Println(err)
		}
		filter = bson.M{"_id": idobj}
	}
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

func VerificarFechaDeReserva(fecha string) bool {
	layout := "2006-01-02"
	location, _ := time.LoadLocation("America/Mexico_City")

	t, _ := time.ParseInLocation(layout, fecha, location)

	var cotizacion Cotizacion

	client, _ := db.ConectarMongoDB()
	// reservacionescol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_R)
	cotizacionescol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_C)

	// opts := options.FindOne().SetSort(bson.M{"Nombre": 1})
	opts := options.FindOne().SetSort(bson.M{"Nombre": 1})
	err := cotizacionescol.FindOne(context.TODO(), bson.M{"FechaReserva": t}, opts).Decode(&cotizacion)
	if err != nil {
		log.Println(err, "Error en Verificar fecha")
	}

	log.Println("Cotizacion encontrada")

	return true
}

func GuardarReservacion(reservacion Reservacion) bool {
	client, _ := db.ConectarMongoDB()
	reservacionesCollection := reservacionesCollection(client)

	res, err := reservacionesCollection.InsertOne(context.TODO(), reservacion)
	if err != nil {
		log.Println(err, "Error creando la reservacion")
		return false
	}
	log.Printf("Reservacion creada : %v, %v\n", reservacion.Nombre, res.InsertedID)
	return true
}

func ActualizaCotizacion(reservacion Reservacion, id string) bool {
	client, _ := db.ConectarMongoDB()
	cotizacionescol := client.Database(conexiones.MONGO_DB).Collection(conexiones.MONGO_DB_C)
	opts := options.FindOneAndUpdate().SetUpsert(true)
	var updatedDocument bson.M
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
	}
	filter := bson.M{"_id": idObj}
	update := bson.M{"$set": bson.M{"Revisada": true, "ReservacionID": reservacion.ID.Hex()}}

	log.Printf("Cotizacion actualizada  : %v\n", reservacion.Nombre)

	err = cotizacionescol.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		opts,
	).Decode(&updatedDocument)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println(err)
			return false
		}
	}

	return true
}

func EditarServicio(servicio Servicio) bool {
	client, _ := db.ConectarMongoDB()
	serviciosCollection := serviciosCollection(client)
	opts := options.FindOneAndUpdate().SetUpsert(true)
	idObj, err := primitive.ObjectIDFromHex(servicio.CurrentID)
	if err != nil {
		log.Println(err)
	}
	filter := bson.M{"_id": idObj}
	update := bson.M{"$set": bson.M{"Titulo": servicio.Titulo, "Descripcion": servicio.Descripcion, "Costo": servicio.Costo, "Activo": servicio.Activo, "Categoria": servicio.Categoria}}
	var updatedDocument bson.M
	err = serviciosCollection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		opts,
	).Decode(&updatedDocument)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false
		}
		log.Println(err)
	}
	log.Printf("Servicio %v actualizado correctamente\n", servicio.Titulo)
	return true
}
