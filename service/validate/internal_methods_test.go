package validate

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/common/descriptor"
)

func TestValidator(t *testing.T) {

	dm, err := descriptor.GetDescriptorManager()
	arithAddDescriptor, err := dm.GetFunctionDescriptor("arithmatic", "Add")
	test.Assert(t, err == nil)

	// test whether could validate simple struct like such:

	/*struct Request {
	  1: i64 firstArguement (api.query = 'firstArguement')
	  2: i64 secondArguement
	  3: optional string message
	  4: optional map<string, string> extra
	  }

	  Response Add(1: Request request ) ( api.get = "/arith/add" )
	*/
	// all number decoded from json are treated as float64
	var arg1 float64 = 1
	var arg2 float64 = 2

	goodJsonBody := map[string]interface{}{"firstArguement": arg1, "secondArguement": arg2}

	err = validateBody(arithAddDescriptor, goodJsonBody)
	test.Assert(t, err == nil)
	//arguement of wrong types
	badJsonBody := map[string]interface{}{"firstArguement": "1", "secondArguement": "two"}
	err = validateBody(arithAddDescriptor, badJsonBody)
	test.Assert(t, err != nil)

	//missing arguements
	badJsonBody2 := map[string]interface{}{"firstArguement": arg1}
	err = validateBody(arithAddDescriptor, badJsonBody2)
	test.Assert(t, err != nil)

	//test whether could validate nested struct like such:

	/*
						struct testValidator {
					    1: Request recur
					    2: map<string, string> extra}

				Response TestValidator(1: testValidator request)
			}

		where type "Request" is the same as the one above
	*/

	finalTestDescriptor, err := dm.GetFunctionDescriptor("arithmatic", "TestValidator")
	test.Assert(t, err == nil)
	finalTest := map[string]interface{}{"recur": goodJsonBody, "extra": map[string]string{"key": "value"}}
	err = validateBody(finalTestDescriptor, finalTest)
	test.Assert(t, err == nil)

}
