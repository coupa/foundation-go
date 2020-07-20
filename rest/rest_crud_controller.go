package rest

import (
	"github.com/coupa/foundation-go/persistence"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type CrudController struct {
	persistenceService persistence.PersistenceService
	httpQueryParser    HttpQueryParser
}

func NewCrudController(persistenceManager persistence.PersistenceService, httpQueryParser HttpQueryParser) *CrudController {
	if httpQueryParser == nil {
		httpQueryParser = &HttpQueryParserRailsActiveAdmin{}
	}
	return &CrudController{
		persistenceService: persistenceManager,
		httpQueryParser:    httpQueryParser,
	}
}

func (self *CrudController) Show(c *gin.Context) {
	id := c.Param("id")
	obj, foundIt, err := self.persistenceService.FindOne(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !foundIt {
		c.JSON(http.StatusNotFound, "")
		return
	}

	c.JSON(http.StatusOK, obj)
}

func (self *CrudController) Index(c *gin.Context) {
	persistenceQuery, err := self.httpQueryParser.Parse(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	objs, err := self.persistenceService.FindMany(persistenceQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if reflect.ValueOf(objs).IsNil() {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	c.JSON(http.StatusOK, objs)
}

func (self *CrudController) Update(c *gin.Context) {
	id := c.Param("id")
	obj := self.persistenceService.NewModelObjPtr()
	foundIt, err := self.persistenceService.FindOneLoader(id, obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if !foundIt {
		c.JSON(http.StatusNotFound, "")
		return
	}

	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	validationErrors, _ := self.persistenceService.Validate(obj)
	if validationErrors.HasErrors() {
		c.JSON(http.StatusUnprocessableEntity, validationErrors)
		return
	}

	_, err = self.persistenceService.UpdateOne(id, obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, obj)
}

func (self *CrudController) Create(c *gin.Context) {
	obj := self.persistenceService.NewModelObjPtr()

	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	validationErrors, _ := self.persistenceService.Validate(obj)
	if validationErrors.HasErrors() {
		c.JSON(http.StatusUnprocessableEntity, validationErrors)
		return
	}

	err := self.persistenceService.CreateOne(obj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, obj)
}

func (self *CrudController) Destroy(c *gin.Context) {
	id := c.Param("id")
	foundIt, err := self.persistenceService.DeleteOne(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if !foundIt {
		c.Status(http.StatusNotFound)
		return
	}
	c.Status(http.StatusNoContent)
}
