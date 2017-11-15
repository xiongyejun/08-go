package oleIE

import (
	"errors"
	"strings"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type ieType struct {
	unknown     *ole.IUnknown
	app         *ole.IDispatch
	ie          *ole.IDispatch // iexplore.exe
	ieDocuments []*ole.IDispatch
}

func (me *ieType) InitIE() (err error) {
	ole.CoInitialize(0)

	me.unknown, err = oleutil.CreateObject("Shell.Application")
	if err != nil {
		return
	}

	if me.app, err = me.unknown.QueryInterface(ole.IID_IDispatch); err != nil {
		return
	}

	var v *ole.VARIANT
	if v, err = oleutil.CallMethod(me.app, "Windows"); err != nil {
		return
	}
	windows := v.ToIDispatch()
	defer windows.Release()

	if v, err = oleutil.GetProperty(windows, "Count"); err != nil {
		return
	}
	count := int(v.Val)

	for i := 0; i < count; i++ {
		if v, err = oleutil.CallMethod(windows, "Item", i); err != nil {
			return
		}

		me.ie = v.ToIDispatch()
		fullName := oleutil.MustGetProperty(me.ie, "FullName")
		if strings.HasSuffix(fullName.ToString(), "iexplore.exe") {
			if v, err = oleutil.GetProperty(me.ie, "Document"); err != nil {
				err = errors.New("getIE Document 出错." + err.Error())
				return
			}
			me.ieDocuments = append(me.ieDocuments, v.ToIDispatch())

			me.enumIE(me.ie)
		}
		me.ie.Release()
	}
	return
}

func (me *ieType) UnInit() {
	for i := range me.ieDocuments {
		me.ieDocuments[i].Release()
	}
	me.app.Release()
	me.unknown.Release()
	ole.CoUninitialize()
}

func (me *ieType) GetIEHtml(urlFind string) (strHtml string, err error) {
	// ie.Document.body.innerHTML
	return me.GetAttr(urlFind, "body.innerHTML")
}

func (me *ieType) GetCookie(urlFind string) (strCookie string, err error) {
	return me.GetAttr(urlFind, "Cookie")
}

// 获取属性
// attr	以.分开的属性，如body.innerHTML
func (me *ieType) GetAttr(urlFind string, attr string) (str string, err error) {
	if len(attr) == 0 {
		err = errors.New("attr 为空.")
		return
	}
	var rtIEDocument *ole.IDispatch

	if rtIEDocument, err = me.getIE(urlFind); err != nil {
		err = errors.New("获取IE出错." + err.Error())
		return
	}

	attrArr := strings.Split(attr, ".")

	var v *ole.VARIANT
	var IDispatchArr []*ole.IDispatch = make([]*ole.IDispatch, len(attrArr))
	IDispatchArr[0] = rtIEDocument
	for i := range attrArr {
		if v, err = oleutil.GetProperty(IDispatchArr[i], attrArr[i]); err != nil {
			err = errors.New(attrArr[i] + " 出错." + err.Error())
			return
		}
		if i == len(attrArr)-1 {
			str = v.ToString()
			return
		}
		IDispatchArr[i+1] = v.ToIDispatch()
	}
	return
}

// 获取me.ieDocuments里url包含urlFind的document
func (me *ieType) getIE(urlFind string) (rtIEDocument *ole.IDispatch, err error) {
	var v *ole.VARIANT

	for i := range me.ieDocuments {
		rtIEDocument = me.ieDocuments[i]

		if v, err = oleutil.GetProperty(rtIEDocument, "url"); err != nil {
			err = errors.New("getIE url 出错." + err.Error())
			return
		}

		if strings.Index(v.ToString(), urlFind) > -1 {
			return
		}
	}
	err = errors.New("getIE 出错.")
	return nil, err
}

func (me *ieType) enumIE(ie *ole.IDispatch) {
	var v *ole.VARIANT
	var err error

	if v, err = oleutil.GetProperty(ie, "Document"); err != nil {
		return
	}
	ieDocument := v.ToIDispatch()

	if v, err = oleutil.GetProperty(ieDocument, "frames"); err != nil {
		return
	}
	obj_frame := v.ToIDispatch()

	if v, err = oleutil.GetProperty(obj_frame, "Length"); err != nil {
		return
	}
	iLen := int(v.Val)

	for i := 0; i < iLen; i++ {
		if v, err = oleutil.CallMethod(obj_frame, "Item", i); err != nil {
			return
		}
		obj_ie := v.ToIDispatch()
		if v, err = oleutil.GetProperty(obj_ie, "Document"); err != nil {
			err = errors.New("getIE Document 出错." + err.Error())
			return
		}
		me.ieDocuments = append(me.ieDocuments, v.ToIDispatch())

		me.enumIE(obj_ie)
	}
	obj_frame.Release()
}

func New() *ieType {
	return new(ieType)
}
