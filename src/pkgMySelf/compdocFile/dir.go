// 解压dir流，获取信息
package compdocFile

import (
	"fmt"
	"pkgAPI/kernel32"
	"unsafe"
)

type dirInfo struct {
	name       string
	textOffset int32 // stream中开始的位置
	moduleType int32
}

type projectModules struct {
	id    int16 // 必须是0x000F
	size  int32 // 必须是 0x00000002
	count int16
	//	projectCookieRecord projectCookie // 8 bytes
	projectCookieID     int16
	projectCookieSize   int32
	projectCookieCookie int16

	//Modules
}

type projectCookie struct {
	id     int16 // 必须是0x0013
	size   int32 // 必须是 0x00000002
	cookie int16 // MUST be ignored on read. MUST be 0xFFFF on write
}
type moduleName struct {
	id               int16 // 必须是0x0019
	sizeOfModuleName int32
	// Dim ModuleName() As Byte
}
type moduleNameUnicode struct {
	id                      int16 // 必须是0x0047
	sizeOfModuleNameUnicode int32
	// Dim ModuleNameUnicode() As Byte
}
type moduleStreamName struct {
	id               int16 // 必须是0x001A
	sizeOfStreamName int32
	// Dim StreamName() As Byte
}
type moduleStreamNameUnicode struct {
	reserved                int16
	sizeOfStreamNameUnicode int32
	// Dim StreamNameUnicode() As Byte
}
type moduleDocString struct {
	id              int16 // 必须是0x001C
	sizeOfDocString int32
	// DocString() As Byte
}

type moduleDocStringUnicode struct {
	reserved               int16
	sizeOfDocStringUnicode int32
	// Dim DocStringUnicode() As Byte
}
type moduleOffset struct {
	id         int16 // 必须是0x0031
	size       int32
	textOffset int32
}
type moduleHelpContext struct {
	id          int16 // 必须是0x001E
	size        int32
	helpContext int32
}
type moduleCookie struct {
	id     int16 // 必须是0x002C
	size   int32 // 必须是 0x00000002
	cookie int16 // MUST be 0xFFFF on write
}

func getModuleInfo(dirStream []byte) (arrDirInfo []dirInfo) {
	project_Modules := projectModules{}

	var pDirStream int32 = 0
	// 找到Project_Modules开始的地址
	for project_Modules.id != 0xf ||
		project_Modules.size != 2 ||
		project_Modules.projectCookieID != 0x13 ||
		project_Modules.projectCookieSize != 2 {

		pDirStream++
		kernel32.MoveMemory(unsafe.Pointer(&project_Modules.id), unsafe.Pointer(&dirStream[pDirStream]), uintptr(unsafe.Sizeof(project_Modules)))

		fmt.Println(project_Modules)
	}
	pDirStream += int32(unsafe.Sizeof(project_Modules))
	// 读取模块个数
	arrDirInfo = make([]dirInfo, project_Modules.count)

	var i int16 = 0
	var module_Name *moduleName
	for ; i < project_Modules.count; i++ {
		// 读取模块名称
		module_Name = new(moduleName)
		kernel32.MoveMemory(unsafe.Pointer(&module_Name.id), unsafe.Pointer(&dirStream[pDirStream]), uintptr(unsafe.Sizeof(module_Name)))
		pDirStream += int32(unsafe.Sizeof(module_Name))
		strName := string(dirStream[pDirStream : pDirStream+module_Name.sizeOfModuleName])
		pDirStream += module_Name.sizeOfModuleName
		arrDirInfo[i].name = strName

		// 读取模块Unicode名称
		module_NameUnicode := new(moduleNameUnicode)
		kernel32.MoveMemory(unsafe.Pointer(&module_NameUnicode.id), unsafe.Pointer(&dirStream[pDirStream]), uintptr(unsafe.Sizeof(module_NameUnicode)))
		pDirStream += int32(unsafe.Sizeof(module_NameUnicode))
		//		strName = string(dirStream[pDirStream : pDirStream+module_NameUnicode.sizeOfModuleNameUnicode])
		pDirStream += module_NameUnicode.sizeOfModuleNameUnicode

		// 读取模块stream名称
		module_StreamName := new(moduleStreamName)
		kernel32.MoveMemory(unsafe.Pointer(&module_StreamName.id), unsafe.Pointer(&dirStream[pDirStream]), uintptr(unsafe.Sizeof(module_StreamName)))
		pDirStream += int32(unsafe.Sizeof(module_StreamName))
		//		strName = string(dirStream[pDirStream : pDirStream+module_StreamName.sizeOfStreamName])
		pDirStream += module_StreamName.sizeOfStreamName

		// 读取模块streamUnicode名称
		module_StreamNameUnicode := new(moduleStreamNameUnicode)
		kernel32.MoveMemory(unsafe.Pointer(&module_StreamNameUnicode.reserved), unsafe.Pointer(&dirStream[pDirStream]), uintptr(unsafe.Sizeof(module_StreamNameUnicode)))
		pDirStream += int32(unsafe.Sizeof(module_StreamNameUnicode))
		//		strName = string(dirStream[pDirStream : pDirStream+module_StreamNameUnicode.sizeOfStreamNameUnicode])
		pDirStream += module_StreamNameUnicode.sizeOfStreamNameUnicode

		// 读取模块DocString
		module_DocString := new(moduleDocString)
		kernel32.MoveMemory(unsafe.Pointer(&module_DocString.id), unsafe.Pointer(&dirStream[pDirStream]), uintptr(unsafe.Sizeof(module_DocString)))
		pDirStream += int32(unsafe.Sizeof(module_DocString))
		//		strName = string(dirStream[pDirStream : pDirStream+module_DocString.sizeOfDocString])
		pDirStream += module_DocString.sizeOfDocString

		// 读取模块ModuleDocStringUnicode
		module_DocStringUnicode := new(moduleDocStringUnicode)
		kernel32.MoveMemory(unsafe.Pointer(&module_DocStringUnicode.reserved), unsafe.Pointer(&dirStream[pDirStream]), uintptr(unsafe.Sizeof(module_DocStringUnicode)))
		pDirStream += int32(unsafe.Sizeof(module_DocStringUnicode))
		//		strName = string(dirStream[pDirStream : pDirStream+module_DocStringUnicode.sizeOfDocStringUnicode])
		pDirStream += module_DocStringUnicode.sizeOfDocStringUnicode

		// 读取ModuleOffset
		module_Offset := new(moduleOffset)
		kernel32.MoveMemory(unsafe.Pointer(&module_Offset.id), unsafe.Pointer(&dirStream[pDirStream]), uintptr(unsafe.Sizeof(module_Offset)))
		pDirStream += int32(unsafe.Sizeof(module_Offset))
		arrDirInfo[i].textOffset = module_Offset.textOffset

		// 跳过ModuleHelpContext
		module_HelpContext := new(moduleHelpContext)
		pDirStream += int32(unsafe.Sizeof(module_HelpContext))

		// 跳过ModuleCookie
		module_Cookie := new(moduleCookie)
		pDirStream += int32(unsafe.Sizeof(module_Cookie))

		//            '这2个不一定有！
		//            ''跳过ModuleReadonly
		//            'Dim Module_Readonly As ModuleReadonly = Nothing
		//            'i_start += Marshal.SizeOf(Module_Readonly)
		//            ''跳过ModulePrivate
		//            'Dim Module_Private As ModulePrivate = Nothing
		//            'i_start += Marshal.SizeOf(Module_Private)

		// 找到下一个Module_Name ID =0x0019的地方
		kernel32.MoveMemory(unsafe.Pointer(&module_Name.id), unsafe.Pointer(&dirStream[pDirStream]), uintptr(unsafe.Sizeof(module_Name)))
		for module_Name.id != 0x19 {
			pDirStream++
			kernel32.MoveMemory(unsafe.Pointer(&module_Name.id), unsafe.Pointer(&dirStream[pDirStream]), uintptr(unsafe.Sizeof(module_Name)))
		}
	}

	return
}
