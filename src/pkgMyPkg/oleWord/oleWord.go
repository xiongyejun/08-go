package oleWord

import (
	"errors"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type word struct {
	app *ole.IDispatch
}

var unknown *ole.IUnknown

func init() {
	ole.CoInitialize(0)
	var err error
	unknown, err = oleutil.GetActiveObject("Word.Application")
	if err != nil {
		unknown, _ = oleutil.CreateObject("Word.Application")
	}
}
func UnInit() {

	ole.CoUninitialize()
}

func (me *word) wordInit() {
	me.app = unknown.MustQueryInterface(ole.IID_IDispatch)
	oleutil.PutProperty(me.app, "Visible", true)

}

func (me *word) Msg() string {
	return oleutil.MustGetProperty(me.app, "Name").ToString()
}

func (me *word) Add(saveName string) (err error) {
	var v *ole.VARIANT
	var documents *ole.IDispatch

	if documents, err = me.Documents(""); err != nil {
		return
	}

	if v, err = oleutil.CallMethod(documents, "Add"); err != nil {
		return errors.New("add 添加出错。")
	}
	if saveName != "" {
		if v, err = oleutil.CallMethod(v.ToIDispatch(), "SaveAs", saveName); err != nil {
			return errors.New("SaveAs 保存出错。" + err.Error())
		}
	}
	return
}

func (me *word) Close(index interface{}, saveChanges bool) (err error) {
	var document *ole.IDispatch

	document, err = me.Documents(index)
	if err != nil {
		return
	}
	_, err = oleutil.CallMethod(document, "Close", saveChanges)
	if err != nil {
		return
	}

	return
}

func (me *word) Documents(index interface{}) (dc *ole.IDispatch, err error) {
	var v *ole.VARIANT
	var documents *ole.IDispatch

	switch index.(type) {
	case int:
		if documents, err = me.Documents(""); err != nil {
			return
		}

		if v, err = oleutil.CallMethod(documents, "Item", index.(int)); err != nil {
			err = errors.New("获取Documents Item出错。" + err.Error())
			return
		}

		dc = v.ToIDispatch()
		return
	case string:
		if index == "" {
			v, err = oleutil.GetProperty(me.app, "Documents")
			dc = v.ToIDispatch()
			return
		}
		if documents, err = me.Documents(""); err != nil {
			return
		}

		if v, err = oleutil.CallMethod(documents, "Item", index.(string)); err != nil {
			err = errors.New("获取Documents Item出错。" + err.Error())
			return
		}

		dc = v.ToIDispatch()
		return
	default:
		err = errors.New("index 未知参数类型。")
		return
	}

}

func (me *word) Exists(fileName string) bool {
	_, err := me.Documents(fileName)
	if err != nil {
		return false
	} else {
		return true
	}
}

func (me *word) Fullname(index interface{}) (fullName string, err error) {
	var document *ole.IDispatch
	var v *ole.VARIANT

	document, err = me.Documents(index)
	if err != nil {
		return
	}
	v, err = oleutil.GetProperty(document, "FullName")
	if err != nil {
		err = errors.New("获取属性Fullname出错。")
		return
	}

	fullName = v.ToString()
	return
}

func (me *word) CopyContent(desIndex interface{}, srcIndex interface{}) (err error) {
	var documentSrc *ole.IDispatch
	var documentDes *ole.IDispatch

	var v *ole.VARIANT

	documentSrc, err = me.Documents(srcIndex)
	if err != nil {
		return
	}

	v, err = oleutil.GetProperty(documentSrc, "Content")
	if err != nil {
		return
	}
	content := v.ToIDispatch()
	if _, err = oleutil.CallMethod(content, "Copy"); err != nil {
		return
	}

	documentDes, err = me.Documents(desIndex)
	if err != nil {
		return
	}

	v, err = oleutil.CallMethod(documentDes, "Range", 0, 0)
	if err != nil {
		return
	}
	if _, err = oleutil.CallMethod(v.ToIDispatch(), "Paste"); err != nil {
		return
	}

	return
}

func (me *word) GetContent(index interface{}) (str string, err error) {
	var document *ole.IDispatch
	var v *ole.VARIANT

	document, err = me.Documents(index)
	if err != nil {
		return
	}

	v, err = oleutil.GetProperty(document, "Content")
	if err != nil {
		return
	}
	content := v.ToIDispatch()
	v, err = oleutil.GetProperty(content, "Text")

	str = v.ToString()
	return
}

func (me *word) Quit() {
	oleutil.CallMethod(me.app, "Quit")
}

func (me *word) SetContent(index interface{}, str string) (err error) {
	var document *ole.IDispatch

	document, err = me.Documents(index)
	if err != nil {
		return
	}

	_, err = oleutil.PutProperty(document, "Content", str)
	return
}

func (me *word) SaveAs(index interface{}, saveName string) (err error) {
	var document *ole.IDispatch

	document, err = me.Documents(index)
	if err != nil {
		return
	}

	if saveName != "" {
		if _, err = oleutil.CallMethod(document, "SaveAs", saveName); err != nil {
			return errors.New("SaveAs 保存出错。" + err.Error())
		}
	} else {
		err = errors.New("saveName 不能为空。")
		return
	}
	return
}

func New() *word {
	wd := new(word)
	wd.wordInit()

	return wd
}
