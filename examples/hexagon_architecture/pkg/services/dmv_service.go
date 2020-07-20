package services

import (
	"fmt"
	"github.com/coupa/foundation-go/examples/hexagon_architecture/models"
	"github.com/coupa/foundation-go/persistence"
	"log"
	"math/rand"
	"time"
)

type DmvService struct {
	PersistenceService persistence.PersistenceService
}

func (self *DmvService) RegisterNewCar() {
	car := models.Car{
		LicensePlate: randSeq(6),
		Make:         makes[rand.Intn(len(makes))],
		Year:         1990 + rand.Intn(30),
	}
	err := self.PersistenceService.CreateOne(&car)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("new car registration %+v\n", car)
}

func (self *DmvService) CrashRandomCar() {
	var cars []models.Car
	err := self.PersistenceService.FindManyLoader(persistence.QueryParams{
		Operands: []persistence.QueryExpression{{
			Key:      "crashed",
			Operator: persistence.QUERY_OPERATOR_EQ,
			Value:    "false",
		}},
	}, &cars)
	if err != nil {
		log.Fatal(err)
	}
	if len(cars) > 0 {
		carToCrash := cars[rand.Intn(len(cars))]
		carToCrash.Crashed = true
		fmt.Printf("about to crash car %v\n", carToCrash.ID)
		_, err := self.PersistenceService.UpdateOne(fmt.Sprint(carToCrash.ID), &carToCrash)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (self *DmvService) DeleteCrashedCars() {
	var cars []models.Car
	err := self.PersistenceService.FindManyLoader(persistence.QueryParams{
		Operands: []persistence.QueryExpression{{
			Key:      "crashed",
			Operator: persistence.QUERY_OPERATOR_EQ,
			Value:    "true",
		}},
	}, &cars)
	if err != nil {
		log.Fatal(err)
	}
	for _, crashedCar := range cars {
		fmt.Printf("about to delete car %v\n", crashedCar.ID)
		_, err := self.PersistenceService.DeleteOne(fmt.Sprint(crashedCar.ID))
		if err != nil {
			log.Fatal(err)
		}
	}
}

var makes = []string{"Tesla", "Volvo", "Honda", "Dodge"}
var letters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
