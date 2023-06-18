package descriptor

import (
	"testing"

	"github.com/cloudwego/thriftgo/pkg/test"
	"github.com/yiwen101/CardWizards/common"
)

func TestDescriptorKeeper(t *testing.T) {
	filename := "http.thrift"
	includeDir := common.RelativePathToIDL
	d, err := buildDescriptorKeeperFromPath(filename, includeDir)
	test.Assert(t, err == nil)

	des := d.get()
	test.Assert(t, des != nil)

	err = d.validateMethodName("BizMethod1")
	test.Assert(t, err == nil)

	err = d.validateMethodName("fake")
	test.Assert(t, err != nil)

	/*
		body := map[string]interface{}{
			"text": "text",
			"some": map[string]interface{}{
				"id":   1,
				"text": "text",
			},
			"req_items_map": map[string]interface{}{
				"1": map[string]interface{}{
					"id":   1,
					"text": "text",
				},
			},
		}
		data, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		url := "http://example.com/1/1?v_int64=1&req_items=item1,item2,itme3&cids=1,2,3&vids=1,2,3"
		req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}
		req.Header.Set("token", "1")
		customReq, err := generic.FromHTTPRequest(req)
		test.Assert(t, err == nil)

		b := d.matchedRouter(customReq)
		test.Assert(t, b == true)
	*/
}
