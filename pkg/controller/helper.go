package controller

import (
	"fmt"
	"reflect"

	restful "github.com/emicklei/go-restful/v3"
)

const (
	KeyOpenAPITags = "openapi.tags"
)

var SecurityHeader = restful.HeaderParameter("Authorization", "访问令牌Bearer Token").DataType("string")

func writeErrorAsJSON(response *restful.Response, statusCode int, err error) {
	response.WriteHeaderAndJson(statusCode, err, restful.MIME_JSON)
}

//CopyProperties  copy src struct to dst struct, dst must be struct pointer
func CopyProperties(src, dst interface{}) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%v", p)
		}
	}()

	srcType, srcVal := reflect.TypeOf(src), reflect.ValueOf(src)
	dstType, dstVal := reflect.TypeOf(dst), reflect.ValueOf(dst)

	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dst type must be a struct pointer")
	}
	if srcType.Kind() != reflect.Struct {
		return fmt.Errorf("src type must be a struct or a struct pointer")
	}
	if srcType.Kind() == reflect.Ptr {
		_, srcVal = srcType.Elem(), srcVal.Elem()
	}
	dstType, dstVal = dstType.Elem(), dstVal.Elem()
	for i := 0; i < dstType.NumField(); i++ {
		property := dstType.Field(i)
		dstValue := dstVal.Field(i)
		srcValue := srcVal.FieldByName(property.Name)
		if !srcValue.IsValid() || srcValue.Type() != property.Type {
			continue
		}
		if dstValue.CanSet() {
			dstValue.Set(srcValue)
		}

	}
	return
}

type BoolRes struct {
	Result bool `json:"result" default:"true"`
}
