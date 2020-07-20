package rest_test

import (
	"fmt"
	"github.com/coupa/foundation-go/persistence"
	"github.com/coupa/foundation-go/rest"
	"github.com/coupa/foundation-go/server"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

var testServer TestServer
var restService persistence.PersistenceService

func TestRest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rest Suite")
}

var _ = BeforeSuite(func() {
	testServer = TestServer{}
	testServer.start()
})

var _ = AfterSuite(func() {
	testServer.stop()
})

type TestModel struct {
	ID      int64  `db:"id" json:"id"`
	StringA string `db:"string_a" json:"string_a"`
	IntA    int    `db:"int_a" json:"int_a"`
}

type TestServer struct {
	baseUrl                  string
	svr                      server.Server
	mockDb                   *PersistenceServiceMock
	crudController           *rest.CrudController
	listener                 net.Listener
	lastestRawQuery          string
	lastestQuery             url.Values
	latestResponseStatusCode int
}

var _ persistence.PersistenceService = (*PersistenceServiceMock)(nil)

func (self *TestServer) start() {
	self.svr = server.Server{Engine: gin.New()}
	var err error
	self.mockDb = &PersistenceServiceMock{}
	self.mockDb.PersistenceServiceCommon = persistence.NewPersistenceServiceCommon(reflect.TypeOf(TestModel{}), self.mockDb.FindOneLoader, self.mockDb.FindManyLoader)
	self.crudController = rest.NewCrudController(self.mockDb, nil)
	self.svr.Engine.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			self.latestResponseStatusCode = -1
			self.lastestRawQuery = c.Request.URL.RawQuery
			self.lastestQuery = c.Request.URL.Query()
			c.Next()
			self.latestResponseStatusCode = c.Writer.Status()
		}
	}())
	self.svr.Engine.GET("/", self.crudController.Index)
	self.svr.Engine.GET("/:id", self.crudController.Show)
	self.svr.Engine.POST("", self.crudController.Create)
	self.svr.Engine.PUT("/:id", self.crudController.Update)
	self.svr.Engine.DELETE("/:id", self.crudController.Destroy)
	self.listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	self.baseUrl = fmt.Sprintf("http://127.0.0.1:%d", self.listener.Addr().(*net.TCPAddr).Port)
	restService, err = rest.NewPersistenceServiceRest(testServer.baseUrl, reflect.TypeOf(TestModel{}))
	if err != nil {
		panic(err)
	}
	go self.svr.Engine.RunListener(self.listener)
}

func (self *TestServer) stop() {
	self.listener.Close()
}

func (self *TestServer) resetMockDb() {
	self.mockDb.db = map[int64]*TestModel{}
}

type PersistenceServiceMock struct {
	persistence.PersistenceServiceCommon
	db        map[int64]*TestModel
	currentId int64
}

func (self *PersistenceServiceMock) FindOneLoader(idStr string, value interface{}) (bool, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return false, err
	}
	obj, foundIt := self.db[id]
	if !foundIt {
		return false, nil
	}
	*value.(*TestModel) = *obj
	return true, nil
}

func (self *PersistenceServiceMock) FindManyLoader(params persistence.QueryParams, values interface{}) error {
	//as of now returns everything without honoring any query params filters
	ret := values.(*[]TestModel)
	for _, v := range self.db {
		*ret = append(*ret, *v)
	}
	return nil
}

func (self *PersistenceServiceMock) CreateOne(obj interface{}) error {
	testModel := obj.(*TestModel)
	testModel.ID = self.nextId()
	self.db[testModel.ID] = testModel
	return nil
}

func (self *PersistenceServiceMock) DeleteOne(idStr string) (bool, error) {
	id, _ := strconv.ParseInt(idStr, 10, 64)
	_, exists := self.db[id]
	delete(self.db, id)
	return exists, nil
}

func (self *PersistenceServiceMock) UpdateOne(idStr string, obj interface{}) (bool, error) {
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if _, exists := self.db[id]; !exists {
		return false, nil
	}
	self.db[id] = obj.(*TestModel)
	return true, nil
}

func (self *PersistenceServiceMock) Validate(obj interface{}) (persistence.ValidationErrors, error) {
	return persistence.Validate(obj)
}

func (self *PersistenceServiceMock) nextId() int64 {
	self.currentId += 1
	return self.currentId
}
