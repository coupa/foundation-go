package persistence

import "reflect"

type PersistenceService interface {
	//Returns nil if record doesnt exist, otherwise returns a struct of same type as returned by GetModelType()
	FindOne(id string) (obj interface{}, foundIt bool, err error)

	//values must be a slice of same type as returned by GetModelType()
	FindMany(params QueryParams) (values interface{}, err error)

	//obj is a pointer to a struct of same type as returned by GetModelType()
	CreateOne(obj interface{}) error

	DeleteOne(id string) (bool, error)

	//obj is a pointer to a struct of same type as returned by GetModelType()
	UpdateOne(id string, obj interface{}) (bool, error)

	Validate(obj interface{}) (ValidationErrors, error)
	GetModelType() reflect.Type

	//Do not use this method directly. It is used as a helper by implementing structs only
	//Returns a struct of same type as returned by GetModelType()
	NewModelObj() interface{}
	//Do not use this method directly. It is used as a helper by implementing structs only
	//Returns a pointer to a struct of same type as returned by GetModelType()
	NewModelObjPtr() interface{}
	//Do not use this method directly, use FindOne() instead. It is used as a helper by implementing structs only
	//Returns a slice of same type as returned by GetModelType(). If no entries loaded, slice is empty
	FindOneLoader(id string, value interface{}) (foundIt bool, err error)
	//Do not use this method directly, use FindMany() instead. It is used as a helper by implementing structs only
	//values must be a pointer to a slice of same type as returned by GetModelType().
	FindManyLoader(params QueryParams, values interface{}) error
}

type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
}

func (self *ValidationErrors) HasErrors() bool {
	return self.Errors != nil && len(self.Errors) > 0
}
