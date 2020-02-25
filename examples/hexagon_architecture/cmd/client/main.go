package main

import (
	"fmt"
	"github.com/coupa/foundation-go/examples/hexagon_architecture/models"
	"github.com/coupa/foundation-go/examples/hexagon_architecture/pkg/services"
	"github.com/coupa/foundation-go/persistence"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

func main() {
	persistenceService, err := persistence.NewPersistenceServiceMySql("root:@/hex_demo", "cars", reflect.TypeOf(models.Car{}))

	apa, _ := persistenceService.FindMany(persistence.QueryParams{})
	fmt.Println(apa)
	if true {
		os.Exit(0)
	}
	//persistenceService, err := rest.NewPersistenceServiceRest("http://localhost:8080/cars", reflect.TypeOf(models.Car{}))
	if err != nil {
		log.Fatal(err)
	}
	dmvService := services.DmvService{PersistenceService: persistenceService}

	ticker3s := time.NewTicker(3 * time.Second)
	ticker5s := time.NewTicker(5 * time.Second)
	ticker10s := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(quit)
	}()

	for {
		select {
		case <-ticker3s.C:
			dmvService.RegisterNewCar()
		case <-ticker5s.C:
			dmvService.CrashRandomCar()
		case <-ticker10s.C:
			dmvService.DeleteCrashedCars()
		case <-quit:
			ticker3s.Stop()
			ticker5s.Stop()
			ticker10s.Stop()
			return
		}
	}
}
