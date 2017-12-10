package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// 读取CustomUI.xsd，能达到的目的：
//	1 读取到各种控件，如button等
//		CT_Tab	CT_Group
//  2 控件需要满足的attribute
//  3 控件所能包含的子控件，如果Box 包含EG_Controls
//  4 属性的各种限制

// simpleType 也就是各种数据类型
// complexType 各种控件，如button等

// attribute 中的属性use
// prohibited 禁止
// optional
// required

// 读取控件、控件允许的子件、控件的属性
func (me *schema) startGetControls() (err error) {
	b, _ := ioutil.ReadFile("CustomUI.xsd")
	xml.Unmarshal(b, me)
	fmt.Println("startGetControls")

	// 记录ComplexType在schema中的下标
	for i := range me.ComplexType {
		d.ComplexType[me.ComplexType[i].Name] = i
	}
	// 记录group在schema中的下标
	for i := range me.Group {
		d.group[me.Group[i].Name] = i
	}
	// 记录simpleType在schema中的下标
	for i := range me.SimpleType {
		d.simpleType[me.SimpleType[i].Name] = i
	}
	// 记录attributeGroup在schema中的下标
	for i := range me.AttributeGroup {
		d.attributeGroup[me.AttributeGroup[i].Name] = i
	}

	// complexType中，先找到CT_CustomUI，再一路找下去
	if err = getControls("CT_CustomUI"); err != nil {
		return
	}
	// 清除Regular后缀的
	var p int = 0
	for i := range controls {
		if strings.HasSuffix(controls[i].name, "Regular") {

		} else {
			controls[p] = controls[i]
			p++
		}
	}
	controls = controls[:p]
	// 将记录在字典里的属性放到切片中
	for i := range controls {
		for _, v := range controls[i].dicAttr {
			controls[i].attr = append(controls[i].attr, v)
		}

		for k, _ := range controls[i].dicChildren {
			controls[i].children = append(controls[i].children, k)
		}
	}

	return
}

// 根据控件的名称来查找控件的children
// ctlName	控件的名称
func getControls(ctlName string) (err error) {
	if index, ok := d.ComplexType[ctlName]; ok {
		// 读取controls
		if _, ok2 := d.control[ctlName]; !ok2 {
			// 如果controls已经有了这个控件，就不用再处理了
			ctl := newControl(ctlName)
			controls = append(controls, ctl)
			d.control[ctlName] = len(controls) - 1
			// 读取children
			return g_schema.ComplexType[index].getControls(ctl)
		}
	} else {
		return errors.New("xsd 文件中未找到 " + ctlName)
	}
	return nil
}

func (me *complexType) getControls(ctl *control) (err error) {
	for i := range me.Sequence {
		if err = me.Sequence[i].getControls(ctl); err != nil {
			return
		}
	}

	for i := range me.ComplexContent {
		if err = me.ComplexContent[i].getControls(ctl); err != nil {
			return
		}
	}

	for i := range me.Group {
		if err = me.Group[i].getControls(ctl); err != nil {
			return
		}
	}

	for i := range me.All {
		if err = me.All[i].getControls(ctl); err != nil {
			return
		}
	}

	// 获取属性
	for i := range me.Attribute {
		if err = me.Attribute[i].getAttr(ctl, true); err != nil {
			return
		}
	}

	for i := range me.AttributeGroup {
		if err = me.AttributeGroup[i].getAttr(ctl); err != nil {
			return
		}
	}

	return nil

}
func (me *all) getControls(ctl *control) (err error) {
	for i := range me.Element {
		ctlChildren := newControl(me.Element[i].Name)
		ctl.appendChildren(ctlChildren.name)
		// 子控件也要继续去获取子控件
		if err = getControls(me.Element[i].Type); err != nil {
			return
		}
	}

	return nil
}

func (me *attributeGroup) getAttr(ctl *control) (err error) {
	// 找到指向的那个AttributeGroup
	if index, ok := d.attributeGroup[me.Ref]; ok {
		if err = g_schema.AttributeGroup[index].getAttr(ctl); err != nil {
			return
		}
	} else {
		err = errors.New("xsd 文件中未找到 " + me.Name)
	}
	// 包含的Attribute
	for i := range me.Attribute {
		if err = me.Attribute[i].getAttr(ctl, true); err != nil {
			return
		}
	}

	for i := range me.AttributeGroup {
		if err = me.AttributeGroup[i].getAttr(ctl); err != nil {
			return
		}
	}

	return nil
}

// bAdd	是增加属性还是减少属性
func (me *attribute) getAttr(ctl *control, bAdd bool) (err error) {

	if bAdd {
		at := newAttr(me.Name)
		if index, ok := d.simpleType[me.Type]; ok {
			at.rest = g_schema.SimpleType[index].Restriction
		}
		ctl.dicAttr[me.Name] = at
	} else {
		delete(ctl.dicAttr, me.Name)
	}

	return nil
}

func (me *complexContent) getControls(ctl *control) (err error) {
	for i := range me.Extension {
		if err = me.Extension[i].getControls(ctl); err != nil {
			return
		}
	}

	// restriction，先要获取base的属性，再减去restriction下的attr
	for i := range me.Restriction {
		if err = me.Restriction[i].getAttr(ctl); err != nil {
			return
		}
	}
	return nil
}
func (me *restriction) getAttr(ctl *control) (err error) {
	if index, ok := d.ComplexType[me.Base]; ok {
		g_schema.ComplexType[index].getControls(ctl)
	} else {
		err = errors.New("xsd 文件中未找到 " + me.Base)
		return
	}
	// 减掉属性
	for i := range me.Attribute {
		if err = me.Attribute[i].getAttr(ctl, false); err != nil {
			return
		}
	}
	return nil
}
func (me *extension) getControls(ctl *control) (err error) {
	if err = getControls(me.Base); err != nil {
		return
	} else {
		// extension的子控件，返回给ctl
		index := d.control[me.Base]
		for i := range controls[index].children {
			ctl.appendChildren(controls[index].children[i])
		}
	}

	// base指向的complexType的属性
	if index, ok := d.ComplexType[me.Base]; ok {
		g_schema.ComplexType[index].getControls(ctl)
	} else {
		err = errors.New("xsd 文件中未找到 " + me.Base)
		return
	}

	// Extension包含的attributeGroup
	for i := range me.AttributeGroup {
		fmt.Println(me.AttributeGroup[i].Ref)
		if err = me.AttributeGroup[i].getAttr(ctl); err != nil {
			return
		}
	}

	for i := range me.Sequence {
		if err = me.Sequence[i].getControls(ctl); err != nil {
			return
		}
	}

	return nil
}

func (me *sequence) getControls(ctl *control) (err error) {
	for i := range me.Element {
		ctlChildren := newControl(me.Element[i].Name)
		ctl.appendChildren(ctlChildren.name)

		// 子控件也要继续去获取子控件
		if err = getControls(me.Element[i].Type); err != nil {
			return
		}
	}

	for i := range me.Choice {
		if err = me.Choice[i].getControls(ctl); err != nil {
			return
		}
	}

	for i := range me.Sequence {
		if err = me.Sequence[i].getControls(ctl); err != nil {
			return
		}
	}
	return nil
}

func (me *choice) getControls(ctl *control) (err error) {
	for i := range me.Element {
		ctlChildren := newControl(me.Element[i].Name)
		ctl.appendChildren(ctlChildren.name)
		// 子控件也要继续去获取子控件
		if err = getControls(me.Element[i].Type); err != nil {
			return
		}
	}

	for i := range me.Group {
		if err = me.Group[i].getControls(ctl); err != nil {
			return
		}
	}
	return nil
}

func (me *group) getControls(ctl *control) (err error) {
	if index, ok := d.group[me.Ref]; ok {
		if err = g_schema.Group[index].getControls(ctl); err != nil {
			return
		}
	}

	for i := range me.Choice {
		if err = me.Choice[i].getControls(ctl); err != nil {
			return
		}
	}
	return nil
}

// 添加子控件名字
func (me *control) appendChildren(childrenName string) {
	me.dicChildren[childrenName] = 0
}

func newControl(ctlName string) *control {
	ctl := new(control)
	ctl.name = ctlName
	ctl.dicAttr = make(map[string]*attrType)
	ctl.dicChildren = make(map[string]int)
	return ctl
}

func newAttr(attrName string) *attrType {
	at := new(attrType)
	at.name = attrName

	return at
}

// 读取xsd文件后，输出结构体
func outStruct() {
	var str []string = make([]string, 0)
	var a2A byte = 'a' - 'A'

	for i := range controls {
		controls[i].name = controls[i].name[3:]                                                 // [3:] 要去掉CT_
		controls[i].name = "rb" + strings.ToLower(controls[i].name[0:1]) + controls[i].name[1:] // 加rb(要不然容易重复) 首字母小写

		str = append(str, "// "+strconv.Itoa(i)+"\r\ntype "+controls[i].name+" struct{\r\n")
		// 输出属性
		for j := range controls[i].attr {
			b := []byte(controls[i].attr[j].name)
			b[0] = b[0] - a2A
			str = append(str, "\t"+string(b)+" string\t`xml:\""+controls[i].attr[j].name+",attr\"`\r\n")
		}
		str = append(str, "\r\n")
		// 输出子节点
		for j := range controls[i].children {
			b := []byte(controls[i].children[j])
			b[0] = b[0] - a2A
			str = append(str, "\t"+string(b)+" []*rb"+controls[i].children[j]+"\t`xml:\""+controls[i].children[j]+"\"`\r\n")
		}
		str = append(str, "}\r\n\r\n")
	}

	ioutil.WriteFile("outStruct.txt", []byte(strings.Join(str, "")), 0666)
}
