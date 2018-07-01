package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"pkgMyPkg/compoundFile"
)

type myData struct {
	cf             *compoundFile.CompoundFile
	encryptionInfo *EncryptionInfo

	dataSpaceMap              *DataSpaceMap
	strongEncryptionDataSpace *DataSpaceDefinition
	primary                   *IRMDSTransformInfo
}

func Parse(fileName string) (iEncryptedType IEncryptedType, err error) {
	var me *myData = &myData{}

	var b []byte
	if b, err = ioutil.ReadFile(fileName); err != nil {
		return
	}

	if me.cf, err = compoundFile.NewCompoundFile(b); err != nil {
		return
	}
	// 解析复合文档
	if err = me.cf.Parse(); err != nil {
		return
	}

	// 读取EncryptionInfo，判断使用的是什么加密方式
	if iEncryptedType, err = me.getEncryptionInfo(); err != nil {
		return
	}
	if err = iEncryptedType.initData(); err != nil {
		return
	}

	//	// 读取DataSpaceMap
	//	if err = me.getDataSpaceMap(); err != nil {
	//		return
	//	}

	return
}

// 读取EncryptionInfo -- 用来判断使用的是哪种加密方式
// 这个是ECMA-376（也就是zip文档加密后形成的复合文档，里面有EncryptionInfo流）的加密
func (me *myData) getEncryptionInfo() (iEncryptedType IEncryptedType, err error) {
	var b []byte
	if b, err = me.cf.GetStream(`EncryptionInfo`); err != nil {
		return
	}
	me.cf = nil

	me.encryptionInfo = new(EncryptionInfo)

	var startIndex int = 0
	if startIndex, err = readVersion(&me.encryptionInfo.EncryptionVersionInfo, b, startIndex); err != nil {
		return
	}
	if me.encryptionInfo.EncryptionVersionInfo.vMajor == 0x0004 &&
		me.encryptionInfo.EncryptionVersionInfo.vMinor == 0x0004 {
		// Agile敏捷 Encryption
		fmt.Println("ECMA-376 Agile Encryption")
		agl := &agile{}
		agl.b = b
		return agl, nil

	} else if (me.encryptionInfo.EncryptionVersionInfo.vMajor == 0x0002 ||
		me.encryptionInfo.EncryptionVersionInfo.vMajor == 0x0003 ||
		me.encryptionInfo.EncryptionVersionInfo.vMajor == 0x0004) &&
		me.encryptionInfo.EncryptionVersionInfo.vMinor == 0x0002 {
		// Standard Encryption
		fmt.Println("ECMA-376 Encryption")
		r := &rc4{}
		r.b = b
		return r, nil

	} else if (me.encryptionInfo.EncryptionVersionInfo.vMajor == 0x0003 ||
		me.encryptionInfo.EncryptionVersionInfo.vMajor == 0x0004) &&
		me.encryptionInfo.EncryptionVersionInfo.vMinor == 0x0003 {
		// Extensible Encryption
		fmt.Println("Extensible Encryption")
		return nil, errors.New("未实现的加密类型。")
	} else {
		return nil, errors.New("未知加密类型。")
	}

	return nil, nil
}

func (me *myData) getDataSpaceMap() (err error) {
	var b []byte
	if b, err = me.cf.GetStream(string([]byte{6}) + `DataSpaces\DataSpaceMap`); err != nil {
		return
	}

	me.dataSpaceMap = new(DataSpaceMap)
	var startIndex int = 0
	// 读取HeaderLength
	if me.dataSpaceMap.HeaderLength, err = byteToUint32(b[startIndex:binary.Size(me.dataSpaceMap.HeaderLength)]); err != nil {
		return
	}
	startIndex += binary.Size(me.dataSpaceMap.HeaderLength)
	// 读取EntryCount
	if me.dataSpaceMap.EntryCount, err = byteToUint32(b[startIndex : startIndex+binary.Size(me.dataSpaceMap.EntryCount)]); err != nil {
		return
	}
	startIndex += binary.Size(me.dataSpaceMap.EntryCount)
	// 读取Map_Entries
	me.dataSpaceMap.Map_Entries = make([]MapEntry, me.dataSpaceMap.EntryCount)

	var i uint32 = 0
	for ; i < me.dataSpaceMap.EntryCount; i++ {
		// 读取MapEntries[i]的长度Length
		if me.dataSpaceMap.Map_Entries[i].Length, err = byteToUint32(b[startIndex : startIndex+binary.Size(me.dataSpaceMap.Map_Entries[i].Length)]); err != nil {
			return
		}
		startIndex += binary.Size(me.dataSpaceMap.Map_Entries[i].Length)
		// 读取MapEntries[i]的ReferenceComponentCount个数
		if me.dataSpaceMap.Map_Entries[i].ReferenceComponentCount, err = byteToUint32(b[startIndex : startIndex+binary.Size(me.dataSpaceMap.Map_Entries[i].ReferenceComponentCount)]); err != nil {
			return
		}
		startIndex += binary.Size(me.dataSpaceMap.Map_Entries[i].ReferenceComponentCount)

		me.dataSpaceMap.Map_Entries[i].ReferenceComponents = make([]DataSpaceReferenceComponent, me.dataSpaceMap.Map_Entries[i].ReferenceComponentCount)

		var j uint32 = 0
		for ; j < me.dataSpaceMap.Map_Entries[i].ReferenceComponentCount; j++ {
			// 读取MapEntries[i]的ReferenceComponentType, 0表示是个stream,1表示storage
			if me.dataSpaceMap.Map_Entries[i].ReferenceComponents[j].ReferenceComponentType, err = byteToUint32(b[startIndex : startIndex+binary.Size(me.dataSpaceMap.Map_Entries[i].ReferenceComponents[j].ReferenceComponentType)]); err != nil {
				return
			}
			startIndex += binary.Size(me.dataSpaceMap.Map_Entries[i].ReferenceComponents[j].ReferenceComponentType)
			if me.dataSpaceMap.Map_Entries[i].ReferenceComponents[j].ReferenceComponentType == 0 {
				if startIndex, err = readUNICODE_LP_P4(&me.dataSpaceMap.Map_Entries[i].ReferenceComponents[j].ReferenceComponent, b, startIndex); err != nil {
					return
				}
			} else {
				// storage怎么处理？
			}

		}

		// 读取me.dataSpaceMap.Map_Entries[i].DataSpaceName
		if startIndex, err = readUNICODE_LP_P4(&me.dataSpaceMap.Map_Entries[i].DataSpaceName, b, startIndex); err != nil {
			return
		}

	}
	// The \0x06DataSpaces\DataSpaceInfo storage MUST contain a stream named "StrongEncryptionDataSpace"
	if bytes.Compare(me.dataSpaceMap.Map_Entries[0].DataSpaceName.Data, []byte{0x53, 0x0, 0x74, 0x0, 0x72, 0x0, 0x6f, 0x0, 0x6e, 0x0, 0x67, 0x0, 0x45, 0x0, 0x6e, 0x0, 0x63, 0x0, 0x72, 0x0, 0x79, 0x0, 0x70, 0x0, 0x74, 0x0, 0x69, 0x0, 0x6f, 0x0, 0x6e, 0x0, 0x44, 0x0, 0x61, 0x0, 0x74, 0x0, 0x61, 0x0, 0x53, 0x0, 0x70, 0x0, 0x61, 0x0, 0x63, 0x0, 0x65, 0x0}) == 0 {
		if err = me.getStrongEncryptionDataSpace(); err != nil {
			return
		}

		// The DataSpaceDefinition structure MUST have exactly one TransformReferences entry, which MUST be "StrongEncryptionTransform"
		if bytes.Compare(me.strongEncryptionDataSpace.TransformReferences[0].Data, []byte{0x53, 0x0, 0x74, 0x0, 0x72, 0x0, 0x6f, 0x0, 0x6e, 0x0, 0x67, 0x0, 0x45, 0x0, 0x6e, 0x0, 0x63, 0x0, 0x72, 0x0, 0x79, 0x0, 0x70, 0x0, 0x74, 0x0, 0x69, 0x0, 0x6f, 0x0, 0x6e, 0x0, 0x54, 0x0, 0x72, 0x0, 0x61, 0x0, 0x6e, 0x0, 0x73, 0x0, 0x66, 0x0, 0x6f, 0x0, 0x72, 0x0, 0x6d, 0x0}) == 0 {

			// The "StrongEncryptionTransform" storage MUST contain a stream named "0x06Primary"
			if err = me.get06Primary(); err != nil {
				return
			}

			fmt.Printf("%#v\r\n", me.primary)
			fmt.Println(string(me.dataSpaceMap.Map_Entries[0].DataSpaceName.Data))
			fmt.Println(string(me.dataSpaceMap.Map_Entries[0].ReferenceComponents[0].ReferenceComponent.Data))

			fmt.Println(string(me.strongEncryptionDataSpace.TransformReferences[0].Data))
		} else {
			return errors.New("没有找到StrongEncryptionTransform")
		}
	} else {
		return errors.New("没有找到StrongEncryptionDataSpace")
	}

	return nil
}
func (me *myData) get06Primary() (err error) {
	var b []byte
	if b, err = me.cf.GetStream(string([]byte{6}) + `DataSpaces\TransformInfo\StrongEncryptionTransform\` + string([]byte{6}) + `Primary`); err != nil {
		return
	}

	me.primary = new(IRMDSTransformInfo)
	var startIndex int = 0
	// 第1部分 TransformInfoHeader
	if startIndex, err = me.primary.getTransformInfoHeader(b, startIndex); err != nil {
		return
	}
	// 第2部分 TransformInfoHeader
	if startIndex, err = me.primary.getExtensibilityHeader(b, startIndex); err != nil {
		return
	}
	// 第3部分
	if startIndex, err = readUTF_8_LP_P4(&me.primary.XrMLLicense, b, startIndex); err != nil {
		return
	}

	return
}

// 第1部分 TransformInfoHeader
func (me *IRMDSTransformInfo) getTransformInfoHeader(b []byte, startIndex int) (endIndex int, err error) {
	me.TransformInfoHeader = TransformInfoHeader{}
	// 读取Length
	if me.TransformInfoHeader.TransformLength, err = byteToUint32(b[startIndex : startIndex+binary.Size(me.TransformInfoHeader.TransformLength)]); err != nil {
		return
	}
	startIndex += binary.Size(me.TransformInfoHeader.TransformLength)
	// 读取TransformType
	if me.TransformInfoHeader.TransformType, err = byteToUint32(b[startIndex : startIndex+binary.Size(me.TransformInfoHeader.TransformType)]); err != nil {
		return
	}
	startIndex += binary.Size(me.TransformInfoHeader.TransformType)
	// 读取TransformID
	if startIndex, err = readUNICODE_LP_P4(&me.TransformInfoHeader.TransformID, b, startIndex); err != nil {
		return
	}
	// 读取TransformName
	if startIndex, err = readUNICODE_LP_P4(&me.TransformInfoHeader.TransformName, b, startIndex); err != nil {
		return
	}
	// 读取ReaderVersion
	if startIndex, err = readVersion(&me.TransformInfoHeader.ReaderVersion, b, startIndex); err != nil {
		return
	}
	// 读取UpdaterVersion
	if startIndex, err = readVersion(&me.TransformInfoHeader.UpdaterVersion, b, startIndex); err != nil {
		return
	}
	// 读取WriterVersion
	if startIndex, err = readVersion(&me.TransformInfoHeader.WriterVersion, b, startIndex); err != nil {
		return
	}

	return startIndex, nil
}

// 第2部分 ExtensibilityHeader
func (me *IRMDSTransformInfo) getExtensibilityHeader(b []byte, startIndex int) (endIndex int, err error) {
	me.ExtensibilityHeader = ExtensibilityHeader{}
	// 读取Length
	if me.ExtensibilityHeader.Length, err = byteToUint32(b[startIndex : startIndex+binary.Size(me.ExtensibilityHeader.Length)]); err != nil {
		return
	}
	startIndex += binary.Size(me.ExtensibilityHeader.Length)
	// It MUST be 0x00000004
	if me.ExtensibilityHeader.Length != 4 {
		return 0, errors.New("ExtensibilityHeader.Length MUST be 0x00000004")
	}
	return startIndex, nil
}

func (me *myData) getStrongEncryptionDataSpace() (err error) {
	var b []byte
	if b, err = me.cf.GetStream(string([]byte{6}) + `DataSpaces\DataSpaceInfo\StrongEncryptionDataSpace`); err != nil {
		return
	}

	me.strongEncryptionDataSpace = new(DataSpaceDefinition)
	var startIndex int = 0
	// 读取Length
	if me.strongEncryptionDataSpace.HeaderLength, err = byteToUint32(b[startIndex:binary.Size(me.strongEncryptionDataSpace.HeaderLength)]); err != nil {
		return
	}
	startIndex += binary.Size(me.strongEncryptionDataSpace.HeaderLength)
	// 读取TransformReferenceCount
	if me.strongEncryptionDataSpace.TransformReferenceCount, err = byteToUint32(b[startIndex : startIndex+binary.Size(me.strongEncryptionDataSpace.TransformReferenceCount)]); err != nil {
		return
	}
	startIndex += binary.Size(me.strongEncryptionDataSpace.TransformReferenceCount)

	me.strongEncryptionDataSpace.TransformReferences = make([]UNICODE_LP_P4, me.strongEncryptionDataSpace.TransformReferenceCount)

	var i uint32 = 0
	for ; i < me.strongEncryptionDataSpace.TransformReferenceCount; i++ {
		if startIndex, err = readUNICODE_LP_P4(&me.strongEncryptionDataSpace.TransformReferences[i], b, startIndex); err != nil {
			return
		}
	}
	return
}

// 读取UNICODE_LP_P4结构
func readUNICODE_LP_P4(p *UNICODE_LP_P4, b []byte, startIndex int) (endIndex int, err error) {
	// 读取长度
	if p.Length, err = byteToUint32(b[startIndex : startIndex+binary.Size(p.Length)]); err != nil {
		return
	}
	startIndex += binary.Size(p.Length)
	// 读取data
	p.Data = b[startIndex : startIndex+int(p.Length)]
	startIndex += int(p.Length)
	// Length必须是4的倍数，不够的补足
	if p.Length%4 != 0 {
		startIndex += int(4 - p.Length%4)
	}
	return startIndex, nil
}

// 读取UTF_8_LP_P4结构
func readUTF_8_LP_P4(p *UTF_8_LP_P4, b []byte, startIndex int) (endIndex int, err error) {
	return readUNICODE_LP_P4(&p.UNICODE_LP_P4, b, startIndex)
}

// 读取Version结构
func readVersion(p *Version, b []byte, startIndex int) (endIndex int, err error) {
	if p.vMajor, err = byteToUint16(b[startIndex : startIndex+binary.Size(p.vMajor)]); err != nil {
		return
	}
	startIndex += binary.Size(p.vMajor)

	if p.vMinor, err = byteToUint16(b[startIndex : startIndex+binary.Size(p.vMinor)]); err != nil {
		return
	}
	startIndex += binary.Size(p.vMinor)

	return startIndex, nil
}
