package main //  xsd

//all				规定子元素能够以任意顺序出现，每个子元素可出现零次或一次。
//annotation		annotation 元素是一个顶层元素，规定 schema 的注释。
//any				使创作者可以通过未被 schema 规定的元素来扩展 XML 文档。
//anyAttribute		使创作者可以通过未被 schema 规定的属性来扩展 XML 文档。
//appInfo			规定 annotation 元素中应用程序要使用的信息。
//attribute			定义一个属性。
//attributeGroup	定义在复杂类型定义中使用的属性组。
//choice			仅允许在 <choice> 声明中包含一个元素出现在包含元素中。
//complexContent	定义对复杂类型（包含混合内容或仅包含元素）的扩展或限制。
//complexType		定义复杂类型。
//documentation		定义 schema 中的文本注释。
//element			定义元素。
//extension			扩展已有的 simpleType 或 complexType 元素。
//field				规定 XPath 表达式，该表达式规定用于定义标识约束的值。
//group				定义在复杂类型定义中使用的元素组。
//import			向一个文档添加带有不同目标命名空间的多个 schema。
//include			向一个文档添加带有相同目标命名空间的多个 schema。
//key				指定属性或元素值（或一组值）必须是指定范围内的键。
//keyref			规定属性或元素值（或一组值）对应指定的 key 或 unique 元素的值。
//list				把简单类型定义为指定数据类型的值的一个列表。
//notation			描述 XML 文档中非 XML 数据的格式。
//redefine			重新定义从外部架构文件中获取的简单和复杂类型、组和属性组。
//restriction		定义对 simpleType、simpleContent 或 complexContent 的约束。
//schema			定义 schema 的根元素。
//selector			指定 XPath 表达式，该表达式为标识约束选择一组元素。
//sequence			要求子元素必须按顺序出现。每个子元素可出现 0 到任意次数。
//simpleContent		包含对 complexType 元素的扩展或限制且不包含任何元素。
//simpleType		定义一个简单类型，规定约束以及关于属性或仅含文本的元素的值的信息。
//union				定义多个 simpleType 定义的集合。
//unique			指定属性或元素值（或者属性或元素值的组合）在指定范围内必须是唯一的。

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type out struct {
	dic      map[string]int // 记录elements中的位置
	elements []*xsdElement  // 记录每一种元素
	s        *stack         // 栈，记录当前处理的所有元素
}

type xsdElement struct {
	name          string
	childElements map[string]int    // 子节点
	attr          map[string]string // key记录属性名，value记录属性的数据类型
}

func main() {
	var ot *out = new(out)
	ot.s = NewStack(20)
	ot.dic = make(map[string]int)

	f, _ := os.Open("xsd.xsd")
	defer f.Close()

	p := xml.NewDecoder(f)

	var t xml.Token
	var err error
	for t, err = p.Token(); err == nil; t, err = p.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			name := token.Name.Local
			var index int = 0
			var ok bool

			// s栈中的父节点记录子节点
			if i_top, err1 := ot.s.Top(); err1 == nil {
				ot.elements[i_top].childElements[name] += 1
			}
			// 节点所在elements中的index
			if index, ok = ot.dic[name]; !ok {
				index = len(ot.elements)
				ot.dic[name] = index
				xsdE := new(xsdElement)
				xsdE.name = name
				xsdE.childElements = make(map[string]int)
				xsdE.attr = make(map[string]string)
				ot.elements = append(ot.elements, xsdE)
			}
			// 栈中记录元素在ot.elements中的位置
			ot.s.Push(index)

			for _, attr := range token.Attr {
				attrName := attr.Name.Local

				if _, ok = ot.elements[index].attr[attrName]; !ok {
					attrValue := attr.Value

					if attrValue == "true" || attrValue == "false" {
						ot.elements[index].attr[attrName] = "bool"
					} else if isNumber(attrValue) {
						ot.elements[index].attr[attrName] = "int"
					} else {
						ot.elements[index].attr[attrName] = "string"
					}
				}

			}

		case xml.EndElement:
			ot.s.Pop()

		case xml.CharData:

		default:

		}
	}

	ot.printOut()
	fmt.Println("ok")
}

//
func (me *out) printOut() {
	var str []string = make([]string, 0)
	var a2A byte = 'a' - 'A'

	for i := range me.elements {
		str = append(str, "// "+strconv.Itoa(i)+"\r\ntype "+me.elements[i].name+" struct{\r\n")
		// 输出属性
		for k, v := range me.elements[i].attr {
			b := []byte(k)
			b[0] = b[0] - a2A
			str = append(str, "\t"+string(b)+" "+v+"\t`xml:\""+k+",attr\"`\r\n")
		}
		str = append(str, "\r\n")
		// 输出子节点
		for k, _ := range me.elements[i].childElements {
			b := []byte(k)
			b[0] = b[0] - a2A
			str = append(str, "\t"+string(b)+" []*"+k+"\t`xml:\""+k+"\"`\r\n")
		}
		str = append(str, "}\r\n\r\n")
	}

	ioutil.WriteFile("out.txt", []byte(strings.Join(str, "")), 0666)
}

// 判断字符串是否都是数字
func isNumber(str string) bool {
	for i := range str {
		if str[i] < '0' || str[i] > '9' {
			return false
		}
	}
	return true
}

//type xsdStruct struct {
//	XsdSchema            xml.Name `xml:"schema"`
//	TargetNamespace      string   `xml:"targetNamespace,attr"`
//	Xmlns                string   `xml:"xmlns,attr"`
//	ElementFormDefault   string   `xml:"elementFormDefault,attr"`
//	attributeFormDefault string   `xml:"attributeFormDefault,attr"`

//	AttributeGroup []*attributeGroup `xml:"attributeGroup"`
//	SimpleType     []*simpleType     `xml:"simpleType"`
//	Annotation     []*annotation     `xml:"annotation"`
//	ComplexType    []*complexType    `xml:"complexType"`
//}

//// attributeGroup	定义在复杂类型定义中使用的属性组。
//type attributeGroup struct {
//	Name string `xml:"name,attr"`
//	Ref  string `xml:"ref,attr"`

//	Annotation     *annotation       `xml:"annotation"`
//	Attribute      []*attribute      `xml:"attribute"`
//	attributeGroup []*attributeGroup `xml:"attributeGroup"`
//}

//// attribute	定义一个属性。
//type attribute struct {
//	Name string `xml:"name,attr"`
//	Type string `xml:"type,attr"`
//	Use  string `xml:"use,attr"`

//	Annotation     []*annotation     `xml:"annotation"`
//	attributeGroup []*attributeGroup `xml:"attributeGroup"`
//}

//// simpleType	定义一个简单类型，规定约束以及关于属性或仅含文本的元素的值的信息。
//type simpleType struct {
//	Name string `xml:"name,attr"`

//	Annotation  []*annotation  `xml:"annotation"`
//	Restriction []*restriction `xml:"restriction"`
//}

//// restriction	定义对 simpleType、simpleContent 或 complexContent 的约束。
//type restriction struct {
//	Base       string `xml:"base,attr"`
//	MinLength  int    `xml:"minLength,attr"`
//	MaxLength  int    `xml:"maxLength,attr"`
//	WhiteSpace string `xml:"whiteSpace,attr"`

//	Enumeration []*enumeration `xml:"enumeration,attr"`
//	Attribute   []*attribute   `xml:"attribute"`
//}

//type enumeration struct {
//	Value string `xml:"value,attr"`
//}

//// annotation 元素是一个顶层元素，规定 schema 的注释。
//type annotation struct {
//	Documentation string `xml:"documentation"`
//}

//// complexType	定义复杂类型。
//type complexType struct {
//	Name  string `xml:"name,attr"`
//	Mixed bool   `xml:"mixed,attr"`

//	Annotation     *annotation       `xml:"annotation"`
//	attributeGroup []*attributeGroup `xml:"attributeGroup"`
//}

//// complexContent	定义对复杂类型（包含混合内容或仅包含元素）的扩展或限制。
//type complexContent struct {
//	Extension *extension `xml:"extension"`

//	Restriction *restriction `xml:"restriction"`
//}

//// extension	扩展已有的 simpleType 或 complexType 元素。
//type extension struct {
//	Base string `xml:"base,attr"`

//	AttributeGroup []*attributeGroup `xml:"attributeGroup"`
//	Restriction    *restriction      `xml:"restriction"`
//	Attribute      []*attribute      `xml:"attribute"`
//}

//// A push-type button
//func main() {
//	b, _ := ioutil.ReadFile("xsd.xsd")

//	xsd := new(xsdStruct)
//	xml.Unmarshal(b, xsd)

//	//	fmt.Printf("%#v\r\n", xsd)
//	fmt.Printf("%#v\r\n AttributeGroup==%d\r\n", xsd.AttributeGroup[0], len(xsd.AttributeGroup))
//	fmt.Printf("%#v\r\n SimpleType==%d\r\n", xsd.SimpleType[0], len(xsd.SimpleType))
//	fmt.Printf("%#v\r\n Annotation==%d\r\n", xsd.Annotation[0], len(xsd.Annotation))
//	fmt.Printf("%#v\r\n ComplexType==%d\r\n", xsd.ComplexType[0], len(xsd.ComplexType))
//	//	fmt.Printf("%#v\r\n", xsd.Annotation.Documentation)
//}
