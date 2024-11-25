package indexmodel

import (
	"context"
	"log"
	"time"

	conexiones "github.com/vadgun/isluad/Conexiones"
	db "github.com/vadgun/isluad/Modelos/Db"
	"go.mongodb.org/mongo-driver/bson"

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
	Activa        bool             `json:"cctiva" bson:"Activa"`
	CostoTotal    float64          `json:"costoTotal" bson:"CostoTotal"`
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

func GetServices() ([]Servicio, error) {
	var servicios []Servicio

	//obtener los servicios de la base de datos
	client, _ := db.ConectarMongoDB()
	serviciosCollection := serviciosCollection(client)

	filter := bson.M{}
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
