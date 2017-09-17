// 解压dir流，获取信息
package compdocFile

import (
	"encoding/binary"
	"pkgMySelf/ucs2T0utf8"
)

type dirInfo struct {
	name       string
	textOffset int32 // stream中开始的位置
	moduleType int32
}

type projectModules struct {
	Id             int16 // 必须是0x000f
	Size           int32 // 必须是 0x00000002
	Count          int16
	Project_Cookie projectCookie // 8 bytes
	//Modules
}

type projectCookie struct {
	Id     int16 // 必须是0x0013
	Size   int32 // 必须是 0x00000002
	Cookie int16 // MUST be ignored on read. MUST be 0xFFFF on write
}
type moduleName struct {
	Id               int16 // 必须是0x0019
	SizeOfModuleName int32
	// Dim ModuleName() As Byte
}
type moduleNameUnicode struct {
	Id                      int16 // 必须是0x0047
	SizeOfModuleNameUnicode int32
	// Dim ModuleNameUnicode() As Byte
}
type moduleStreamName struct {
	Id               int16 // 必须是0x001A
	SizeOfStreamName int32
	// Dim StreamName() As Byte
}
type moduleStreamNameUnicode struct {
	Reserved                int16
	SizeOfStreamNameUnicode int32
	// Dim StreamNameUnicode() As Byte
}
type moduleDocString struct {
	Id              int16 // 必须是0x001C
	SizeOfDocString int32
	// DocString() As Byte
}

type moduleDocStringUnicode struct {
	Reserved               int16
	SizeOfDocStringUnicode int32
	// Dim DocStringUnicode() As Byte
}
type moduleOffset struct {
	Id         int16 // 必须是0x0031
	Size       int32
	TextOffset int32
}
type moduleHelpContext struct {
	Id          int16 // 必须是0x001E
	Size        int32
	HelpContext int32
}
type moduleCookie struct {
	Id     int16 // 必须是0x002C
	Size   int32 // 必须是 0x00000002
	Cookie int16 // MUST be 0xFFFF on write
}

func getModuleInfo(dirStream []byte) (arrDirInfo []dirInfo) {
	project_Modules := projectModules{}

	var pDirStream int32 = 0
	// 找到Project_Modules开始的地址
	for project_Modules.Id != 0xf ||
		project_Modules.Size != 2 ||
		project_Modules.Project_Cookie.Id != 0x13 ||
		project_Modules.Project_Cookie.Size != 2 {

		pDirStream++
		byte2struct(dirStream[pDirStream:], &project_Modules)
	}

	pDirStream += int32(binary.Size(project_Modules))
	// 读取模块个数
	arrDirInfo = make([]dirInfo, project_Modules.Count)

	var i int16 = 0
	var module_Name moduleName
	for ; i < project_Modules.Count; i++ {
		// 读取模块名称
		module_Name = moduleName{}

		// 找到Module_Name ID =0x0019的地方，
		//		buf := bytes.NewBuffer(dirStream[pDirStream:])
		//		binary.Read(buf, binary.LittleEndian, &module_Name)

		byte2struct(dirStream[pDirStream:], &module_Name)
		//因为有2个结构不一定有，所有不能保证第2后的moduleName位置
		for module_Name.Id != 0x0019 {
			pDirStream++
			byte2struct(dirStream[pDirStream:], &module_Name)
		}
		pDirStream += int32(binary.Size(module_Name))
		pDirStream += module_Name.SizeOfModuleName

		// 读取模块Unicode名称
		module_NameUnicode := moduleNameUnicode{}
		byte2struct(dirStream[pDirStream:], &module_NameUnicode)
		pDirStream += int32(binary.Size(module_NameUnicode))
		bName := dirStream[pDirStream : pDirStream+module_NameUnicode.SizeOfModuleNameUnicode]
		bName, _ = ucs2T0utf8.UCS2toUTF8(bName)
		pDirStream += module_NameUnicode.SizeOfModuleNameUnicode
		arrDirInfo[i].name = string(bName)

		// 读取模块stream名称
		module_StreamName := moduleStreamName{}
		byte2struct(dirStream[pDirStream:], &module_StreamName)
		pDirStream += int32(binary.Size(module_StreamName))
		pDirStream += module_StreamName.SizeOfStreamName

		// 读取模块streamUnicode名称
		module_StreamNameUnicode := moduleStreamNameUnicode{}
		byte2struct(dirStream[pDirStream:], &module_StreamNameUnicode)
		pDirStream += int32(binary.Size(module_StreamNameUnicode))
		pDirStream += module_StreamNameUnicode.SizeOfStreamNameUnicode

		// 读取模块DocString
		module_DocString := moduleDocString{}
		byte2struct(dirStream[pDirStream:], &module_DocString)
		pDirStream += int32(binary.Size(module_DocString))
		pDirStream += module_DocString.SizeOfDocString

		// 读取模块ModuleDocStringUnicode
		module_DocStringUnicode := moduleDocStringUnicode{}
		byte2struct(dirStream[pDirStream:], &module_DocStringUnicode)
		pDirStream += int32(binary.Size(module_DocStringUnicode))
		pDirStream += module_DocStringUnicode.SizeOfDocStringUnicode

		// 读取ModuleOffset
		module_Offset := moduleOffset{}
		byte2struct(dirStream[pDirStream:], &module_Offset)
		pDirStream += int32(binary.Size(module_Offset))
		arrDirInfo[i].textOffset = module_Offset.TextOffset

		// 跳过ModuleHelpContext
		module_HelpContext := moduleHelpContext{}
		pDirStream += int32(binary.Size(module_HelpContext))

		// 跳过ModuleCookie
		module_Cookie := moduleCookie{}
		pDirStream += int32(binary.Size(module_Cookie))

		//            '这2个不一定有！
		//            ''跳过ModuleReadonly
		//            'Dim Module_Readonly As ModuleReadonly = Nothing
		//            'i_start += Marshal.SizeOf(Module_Readonly)
		//            ''跳过ModulePrivate
		//            'Dim Module_Private As ModulePrivate = Nothing
		//            'i_start += Marshal.SizeOf(Module_Private)

	}

	return
}
