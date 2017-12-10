package main

import (
	"fmt"
)

type ribbonType struct {
	Type   string
	TypeID int
}

var ribbonTypes []*ribbonType = []*ribbonType{
	{"tab", 1},
	{"group", 2},
	{"button", 3},
}

type dic struct {
	ComplexType    map[string]int
	group          map[string]int
	simpleType     map[string]int
	attributeGroup map[string]int
	control        map[string]int // 记录control是否已经在controls里了，已经存在就不需要继续添加了
}

// 记录控件的属性
type attrType struct {
	name string
	rest []*restriction
}

// 记录每种控件，并且每种控件允许的子控件、属性
type control struct {
	name        string
	children    []string
	dicChildren map[string]int

	attr    []*attrType
	dicAttr map[string]*attrType // 因为有restriction，有时需要减少
}

var d *dic
var controls []*control

func init() {
	d = new(dic)
	d.ComplexType = make(map[string]int)
	d.group = make(map[string]int)
	d.simpleType = make(map[string]int)
	d.attributeGroup = make(map[string]int)
	d.control = make(map[string]int)

	controls = make([]*control, 0)
}

var g_schema = new(schema)

func main() {
	fmt.Println(g_schema.startGetControls())
	//	outStruct()

	//	return

	uiInit()
	//	saveNodeToXml()

	fmt.Println("ok")
}

//func saveNodeToXml() {
//	if b, err := xml.Marshal(ct.treeModle.roots[0]); err != nil {
//		fmt.Println(err)
//	} else {
//		if err := ioutil.WriteFile("xml.txt", b, 0666); err != nil {
//			fmt.Println(err)
//		}
//	}
//}

//func dg(value reflect.Value, pre string) {
//	for i := 0; i < value.NumField(); i++ {
//		fmt.Printf("%s%d Type: %s, Name: %s, Kind: %s, Value: %s\r\n", pre, i, value.Type().String(), value.Field(i).Type().String(), value.Field(i).Kind().String(), value.Field(i))

//		switch value.Field(i).Kind() {
//		case reflect.Ptr:
//			v := value.Field(i).Elem()
//			if v.Kind() == reflect.Struct {
//				dg(v, pre+"  ")
//				fmt.Println()
//			}

//		case reflect.Struct:
//			v := value.Field(i)
//			dg(v, pre+"  ")
//			fmt.Println()
//		default:

//		}
//	}
//}
