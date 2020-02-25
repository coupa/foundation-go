package main

import (
	"github.com/coupa/foundation-go/examples/hexagon_architecture/models"
	"github.com/coupa/foundation-go/persistence"
	"github.com/coupa/foundation-go/rest"
	"github.com/coupa/foundation-go/server"
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
)

func main() {
	svr := server.Server{Engine: gin.New()}

	carPersistenceService, err := persistence.NewPersistenceServiceMySql("root:@/hex_demo", "cars", reflect.TypeOf(models.Car{}))
	if err != nil {
		log.Fatal(err)
	}

	//Register routes
	v1Cars := svr.Engine.Group("/v1/cars")
	carController := rest.NewCrudController(carPersistenceService, nil)

	v1Cars.GET("/", carController.Index)
	v1Cars.GET("/:id", carController.Show)
	v1Cars.POST("", carController.Create)
	v1Cars.PUT("/:id", carController.Update)
	v1Cars.DELETE("/:id", carController.Destroy)

	svr.Engine.Run(":8080") //svr.Engine.Run() without address parameter will run on ":8080"
}
