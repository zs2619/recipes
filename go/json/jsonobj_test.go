package recipesgo

import (
	"encoding/json"
	"fmt"
	"recipes-go/test"
	"testing"

	"github.com/stretchr/testify/suite"
)

type JsonObjTestSuite struct {
	test.BaseTestSuite
}

func TestJsonObj(t *testing.T) {
	suite.Run(t, new(JsonObjTestSuite))
}

func (t *JsonObjTestSuite) TestDiffJsonObj() {
	src := `{"z":{"ss":345,"z":true},"a":[2,1]}`
	dest := `{"z":{"ss":123,"z":false},"a":[2,2]}`

	var srcObj map[string]interface{}
	err := json.Unmarshal([]byte(src), &srcObj)
	t.Nil(err)

	var destObj map[string]interface{}
	err = json.Unmarshal([]byte(dest), &destObj)
	t.Nil(err)
	json.Marshal(destObj)

	diff := diffJSObj(&srcObj, &destObj, "$")
	t.NotNil(diff)
	fmt.Printf(" shuai %v", diff)

}

func (t *JsonObjTestSuite) TestDiffJsonObj2() {
	diff := DiffJSON(nil, nil)
	t.Nil(diff)

	js := `{"z":{"ss":345,"z":true},"a":[2,1]}`
	var jsObj map[string]interface{}
	err := json.Unmarshal([]byte(js), &jsObj)
	t.Nil(err)

	diff = DiffJSON(nil, &jsObj)
	t.NotNil(diff)

	diff = DiffJSON(&jsObj, nil)
	t.NotNil(diff)
}

func (t *JsonObjTestSuite) TestJsonPath() {
	src := `{"z":{"ss":345,"z":true},"a":[2,1]}`
	var srcObj map[string]interface{}
	err := json.Unmarshal([]byte(src), &srcObj)
	t.Nil(err)

	v, err := JsonPathGet(&srcObj, "$.a.a.a.a")
	t.NotNil(err)
	v, err = JsonPathGet(&srcObj, "$.a.")
	t.NotNil(err)
	v, err = JsonPathGet(&srcObj, ".a.")
	t.NotNil(err)

	v, err = JsonPathGet(&srcObj, "$.z.ss")
	t.Nil(err)
	t.Equal(v, float64(345))

	v, err = JsonPathGet(&srcObj, "$.z.z")
	t.Nil(err)
	t.Equal(v, true)

	v, err = JsonPathGet(&srcObj, "$.a")
	t.Nil(err)
	t.Equal(2, len(v.([]interface{})))

}

func (t *JsonObjTestSuite) TestJsonPathOP() {
	src := `{"z":{"ss":345,"z":true},"a":[2,1]}`
	var srcObj map[string]interface{}
	err := json.Unmarshal([]byte(src), &srcObj)
	t.Nil(err)

	retObj, err := JsonPathOP(srcObj, "$.name", SetOP, 1)
	t.Nil(err)
	t.Equal((*retObj)["name"], float64(1))

	name, err := JsonPathGet(retObj, "$.name")
	t.Nil(err)
	t.Equal(name, float64(1))

	retObj, err = JsonPathOP(srcObj, "$.name", IncOP, 1)
	t.Nil(err)
	t.Equal((*retObj)["name"], float64(2))

	name, err = JsonPathGet(retObj, "$.name")
	t.Nil(err)
	t.Equal(name, float64(2))

	retObj, err = JsonPathOP(srcObj, "$.name.name", IncOP, 1)
	t.Nil(err)
	t.Equal((*retObj)["name"].(map[string]interface{})["name"], float64(1))

	name, err = JsonPathGet(retObj, "$.name.name")
	t.Nil(err)
	t.Equal(name, float64(1))

	retObj, err = JsonPathOP(srcObj, "$.name", UnSetOP, nil)
	t.Nil(err)
	_, ok := (*retObj)["name"]
	t.Equal(ok, false)

	_, err = JsonPathGet(retObj, "$.name")
	t.NotNil(err)

	retObj, err = JsonPathOP(srcObj, "$.name.name", SetOP, 1)
	t.Nil(err)
	t.Equal((*retObj)["name"].(map[string]interface{})["name"], float64(1))

	name, err = JsonPathGet(retObj, "$.name.name")
	t.Nil(err)
	t.Equal(name, float64(1))

	retObj, err = JsonPathOP(srcObj, "$.name.name", SetOP, 2)
	t.Nil(err)
	t.Equal((*retObj)["name"].(map[string]interface{})["name"], float64(2))

	name, err = JsonPathGet(retObj, "$.name.name")
	t.Nil(err)
	t.Equal(name, float64(2))

	retObj, err = JsonPathOP(srcObj, "$.name", UnSetOP, nil)
	t.Nil(err)

	_, err = JsonPathGet(retObj, "$.name")
	t.NotNil(err)

	retObj, err = JsonPathOP(srcObj, "$.list", PushOP, 10)
	t.Nil(err)

	l, ok := (*retObj)["list"]
	t.Equal(ok, true)
	t.Equal(1, len(l.([]interface{})))
	t.Equal(float64(10), (l.([]interface{}))[0])

	retObj, err = JsonPathOP(*retObj, "$.list", PushOP, 11)
	t.Nil(err)

	l, ok = (*retObj)["list"]
	t.Equal(ok, true)
	t.Equal(2, len(l.([]interface{})))
	t.Equal(float64(11), (l.([]interface{}))[1])

	ll, err := JsonPathGet(retObj, "$.list")
	t.Nil(err)
	t.Equal(2, len(ll.([]interface{})))
	t.Equal(float64(10), (ll.([]interface{}))[0])
	t.Equal(float64(11), (ll.([]interface{}))[1])

	retObj, err = JsonPathOP(*retObj, "$.list", PullOP, 11)
	t.Nil(err)
	l, ok = (*retObj)["list"]
	t.Equal(ok, true)
	t.Equal(1, len(l.([]interface{})))
	t.Equal(float64(10), (l.([]interface{}))[0])

	ll, err = JsonPathGet(retObj, "$.list")
	t.Nil(err)
	t.Equal(1, len(ll.([]interface{})))
	t.Equal(float64(10), (ll.([]interface{}))[0])

	retObj, err = JsonPathOP(*retObj, "$.list", PullOP, 10)
	t.Nil(err)
	l, ok = (*retObj)["list"]
	t.Equal(ok, true)
	t.Equal(0, len(l.([]interface{})))

}

func (t *JsonObjTestSuite) TestCheckJsonOpJs() {
	s := `{"ps":{"items":{"configID":123,"name":"Knife"}},"s":{"name":"Bohu"}}`
	err := CheckJsonOpJs(&s)
	t.Nil(err)

	s = `{"ps":{"name":{"configID":123,"name":"Knife"}},"s":{"name":"Bohu"}}`
	err = CheckJsonOpJs(&s)
	t.NotNil(err)

	s = `{"ps1":{"items":{"configID":123,"name":"Knife"}},"s":{"name":"Bohu"}}`
	err = CheckJsonOpJs(&s)
	t.NotNil(err)
}
