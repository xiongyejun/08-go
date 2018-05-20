package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"pkgMyPkg/compoundFile"
)

type MyData struct {
	cf *compoundFile.CompoundFile

	dataSpaceMap              *DataSpaceMap
	strongEncryptionDataSpace *DataSpaceDefinition
	primary                   *IRMDSTransformInfo
}

func (me *MyData) Parse(fileName string) (err error) {
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

	// 读取DataSpaceMap
	if err = me.getDataSpaceMap(); err != nil {
		return
	}
	return nil
}

func (me *MyData) getDataSpaceMap() (err error) {
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
	if string(me.dataSpaceMap.Map_Entries[0].DataSpaceName.Data) == string([]byte{0x53, 0x0, 0x74, 0x0, 0x72, 0x0, 0x6f, 0x0, 0x6e, 0x0, 0x67, 0x0, 0x45, 0x0, 0x6e, 0x0, 0x63, 0x0, 0x72, 0x0, 0x79, 0x0, 0x70, 0x0, 0x74, 0x0, 0x69, 0x0, 0x6f, 0x0, 0x6e, 0x0, 0x44, 0x0, 0x61, 0x0, 0x74, 0x0, 0x61, 0x0, 0x53, 0x0, 0x70, 0x0, 0x61, 0x0, 0x63, 0x0, 0x65, 0x0}) {
		if err = me.getStrongEncryptionDataSpace(); err != nil {
			return
		}
		fmt.Printf("%#v\r\n", me.strongEncryptionDataSpace.TransformReferences[0].Data)
		// The DataSpaceDefinition structure MUST have exactly one TransformReferences entry, which MUST be "StrongEncryptionTransform"
		if string(me.strongEncryptionDataSpace.TransformReferences[0].Data) == string([]byte{0x53, 0x0, 0x74, 0x0, 0x72, 0x0, 0x6f, 0x0, 0x6e, 0x0, 0x67, 0x0, 0x45, 0x0, 0x6e, 0x0, 0x63, 0x0, 0x72, 0x0, 0x79, 0x0, 0x70, 0x0, 0x74, 0x0, 0x69, 0x0, 0x6f, 0x0, 0x6e, 0x0, 0x54, 0x0, 0x72, 0x0, 0x61, 0x0, 0x6e, 0x0, 0x73, 0x0, 0x66, 0x0, 0x6f, 0x0, 0x72, 0x0, 0x6d, 0x0}) {

			fmt.Printf("%#v\r\n", me.dataSpaceMap)
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
func (me *MyData) get06Primary() (err error) {
	var b []byte
	if b, err = me.cf.GetStream(string([]byte{6}) + `DataSpaces\TransformInfo\StrongEncryptionTransform\` + string([]byte{6}) + `Primary`); err != nil {
		return
	}

	me.primary = new(IRMDSTransformInfo)
}

func (me *MyData) getStrongEncryptionDataSpace() (err error) {
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

func byteToUint32(src []byte) (x uint32, err error) {
	if len(src) != 4 {
		return 0, errors.New("转uint32必须是4个字节。")
	}
	return uint32(src[0]) + uint32(src[1]<<8) + uint32(src[2]<<16) + uint32(src[3]<<24), nil
}
