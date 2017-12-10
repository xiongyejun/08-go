// 读取xsd用的结构，然后可以获取ribbon的结构

package main

import (
	"encoding/xml"
)

// 0
type schema struct {
	XsdSchema            xml.Name `xml:"schema"`
	Xsd                  string   `xml:"xsd,attr"`
	Version              string   `xml:"version,attr"`
	TargetNamespace      string   `xml:"targetNamespace,attr"`
	Xmlns                string   `xml:"xmlns,attr"`
	ElementFormDefault   string   `xml:"elementFormDefault,attr"`
	AttributeFormDefault string   `xml:"attributeFormDefault,attr"`

	Annotation     []*annotation     `xml:"annotation"`
	SimpleType     []*simpleType     `xml:"simpleType"`
	AttributeGroup []*attributeGroup `xml:"attributeGroup"`
	ComplexType    []*complexType    `xml:"complexType"`
	Group          []*group          `xml:"group"`
	Element        []*element        `xml:"element"`
}

// 1
type annotation struct {
	Documentation []*documentation `xml:"documentation"`
}

// 2
type documentation struct {
}

// 3
type simpleType struct {
	Name string `xml:"name,attr"`

	Annotation  []*annotation  `xml:"annotation"`
	Restriction []*restriction `xml:"restriction"`
}

// 4
type restriction struct {
	Base string `xml:"base,attr"`

	MinLength    []*minLength    `xml:"minLength"`
	MaxLength    []*maxLength    `xml:"maxLength"`
	MinInclusive []*minInclusive `xml:"minInclusive"`
	MaxInclusive []*maxInclusive `xml:"maxInclusive"`
	Enumeration  []*enumeration  `xml:"enumeration"`
	WhiteSpace   []*whiteSpace   `xml:"whiteSpace"`
	Attribute    []*attribute    `xml:"attribute"`
}

// 5
type minLength struct {
	Value int `xml:"value,attr"`
}

// 6
type maxLength struct {
	Value int `xml:"value,attr"`
}

// 7
type minInclusive struct {
	Value int `xml:"value,attr"`
}

// 8
type maxInclusive struct {
	Value int `xml:"value,attr"`
}

// 9
type enumeration struct {
	Value string `xml:"value,attr"`
}

// 10
type whiteSpace struct {
	Value string `xml:"value,attr"`
}

// 11
type attributeGroup struct {
	Ref  string `xml:"ref,attr"`
	Name string `xml:"name,attr"`

	Annotation     []*annotation     `xml:"annotation"`
	Attribute      []*attribute      `xml:"attribute"`
	AttributeGroup []*attributeGroup `xml:"attributeGroup"`
}

// 12
type attribute struct {
	Use  string `xml:"use,attr"`
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`

	Annotation []*annotation `xml:"annotation"`
}

// 13
type complexType struct {
	Name  string `xml:"name,attr"`
	Mixed bool   `xml:"mixed,attr"`

	All            []*all            `xml:"all"`
	Annotation     []*annotation     `xml:"annotation"`
	AttributeGroup []*attributeGroup `xml:"attributeGroup"`
	ComplexContent []*complexContent `xml:"complexContent"`
	Attribute      []*attribute      `xml:"attribute"`
	Sequence       []*sequence       `xml:"sequence"`
	Group          []*group          `xml:"group"`
}

// 14
type complexContent struct {
	Extension   []*extension   `xml:"extension"`
	Restriction []*restriction `xml:"restriction"`
}

// 15
type extension struct {
	Base string `xml:"base,attr"`

	AttributeGroup []*attributeGroup `xml:"attributeGroup"`
	Attribute      []*attribute      `xml:"attribute"`
	Sequence       []*sequence       `xml:"sequence"`
}

// 16
type sequence struct {
	MinOccurs int `xml:"minOccurs,attr"`

	Element  []*element  `xml:"element"`
	Choice   []*choice   `xml:"choice"`
	Sequence []*sequence `xml:"sequence"`
}

// 17
type element struct {
	Name      string `xml:"name,attr"`
	Type      string `xml:"type,attr"`
	MinOccurs int    `xml:"minOccurs,attr"`
	MaxOccurs int    `xml:"maxOccurs,attr"`

	Annotation []*annotation `xml:"annotation"`
	Unique     []*unique     `xml:"unique"`
}

// 18
type group struct {
	Name      string `xml:"name,attr"`
	Ref       string `xml:"ref,attr"`
	MinOccurs int    `xml:"minOccurs,attr"`
	MaxOccurs int    `xml:"maxOccurs,attr"`

	Annotation []*annotation `xml:"annotation"`
	Choice     []*choice     `xml:"choice"`
}

// 19
type choice struct {
	MinOccurs int `xml:"minOccurs,attr"`
	MaxOccurs int `xml:"maxOccurs,attr"`

	Element []*element `xml:"element"`
	Group   []*group   `xml:"group"`
}

// 20
type all struct {
	Element []*element `xml:"element"`
}

// 21
type unique struct {
	Name string `xml:"name,attr"`

	Annotation []*annotation `xml:"annotation"`
	Selector   []*selector   `xml:"selector"`
	Field      []*field      `xml:"field"`
}

// 22
type selector struct {
	Xpath string `xml:"xpath,attr"`
}

// 23
type field struct {
	Xpath string `xml:"xpath,attr"`
}
