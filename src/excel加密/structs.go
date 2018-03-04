package main // officeEncryp

// https://msdn.microsoft.com/en-us/library/dd943485(v=office.12).aspx

// 0x06DataSpaces存储：包含所有必需信息的存储器，用于了解应用于给定OLE复合文件中原始文档内容的转换。
// 0x06DataSpaces \ Version流：包含DataSpaceVersionInfo结构的流，如2.1.5节所述。该流指定了文件中使用的数据空间结构的版本。
// 0x06DataSpaces \ DataSpaceMap流：包含2.1.6节中指定的DataSpaceMap结构的流。该流将受保护的内容与用于转换它的数据空间定义相关联。
// 0x06DataSpaces \ DataSpaceInfo storage：包含文件中使用的数据空间定义的存储。这个存储必须包含一个或多个流，每个流包含一个2.1.7 节中规定的DataSpaceDefinition结构。存储必须在\ 0x06DataSpaces \ DataSpaceMap 流（2.2.1节）中为每个DataSpaceMapEntry结构（2.1.6.1节 ）准确包含一个流。每个流的名称必须与包含在\ 0x06DataSpaces \ DataSpaceMap 流中的一个DataSpaceMapEntry结构的DataSpaceName字段相同。
// 转换后的内容流和存储：包含受保护内容的一个或多个存储或流。已转换的内容通过\ 0x06DataSpaces \ DataSpaceMap流中的条目与数据空间定义相关联。
//  0x06DataSpaces \ TransformInfo storage：存储，其中包含 在2.2.2节中指定的存储在\ 0x06DataSpaces \ DataSpaceInfo存储中的数据空间定义中使用的转换的定义。该流包含零个或多个可应用于内容流中的数据的可能转换的定义。

// 2.1.1 File
type File struct {
	DataSpaces0x06
}

// storage
type DataSpaces0x06 struct {
	DataSpaceVersionInfo
	DataSpaceMap
	DataSpaceInfo
	TransformInfo
}
type DataSpaceVersionInfo struct {
	FeatureIdentifier UNICODE_LP_P4 // It MUST be "Microsoft.Container.DataSpaces"
	ReaderVersion     Version       // ReaderVersion.vMajor MUST be 1. ReaderVersion.vMinor MUST be 0.
	UpdaterVersion    Version       // 同上
	WriterVersion     Version       // 同上
}
type Version struct {
	vMajor uint16
	vMinor uint16
}

// 2.1.6
type DataSpaceMap struct {
	HeaderLength uint32 // 一个无符号整数，指定MapEntries数组中第一个条目之前的DataSpaceMap结构中的字节数。它必须等于0x00000008。
	EntryCount   uint32 // MapEntries的数量
	Map_Entries  []MapEntries
}

// 2.1.6.1
type MapEntry struct {
	Length                  uint32 // 指定DataSpaceMapEntry结构的大小（以字节为单位）。
	ReferenceComponentCount uint32 // 指定ReferenceComponents数组中DataSpaceReferenceComponent项的
	ReferenceComponents     []DataSpaceReferenceComponent
	DataSpaceName           UNICODE_LP_P4 // It MUST be equal to the name of a stream in the \0x06DataSpaces\DataSpaceInfo storage as specified in section 2.2.2
}
type DataSpaceReferenceComponent struct {
	ReferenceComponentType uint32 // 指定引用的组件是流还是存储。0x00000000stream 0x00000001storage
	ReferenceComponent     UNICODE_LP_P4
}

// storage
type DataSpaceInfo struct {
	DataSpace_Definition []DataSpaceDefinition
	DataSpaceMapEntry
}

// 2.1.7
type DataSpaceDefinition struct {
	HeaderLength            uint32 // 指定TransformReferences字段之前的DataSpaceDefinition结构中的字节数。它务必是0x00000008。
	TransformReferenceCount uint32 // 指定TransformReferences 数组中的项目数
	TransformReferences     []UNICODE_LP_P4
}

type TransformInfo struct {
	trans_form []transform
}
type transform struct {
	Primary0x06
}

// stream
type Primary0x06 struct {
	IRMDSTransformInfo
}

// 2.1.8
type TransformInfoHeader struct {
	TransformLength uint32 // 指定TransformName 字段之前的此结构中的字节数。
	TransformType   uint32 // 指定要应用的变换的类型。它务必是0x00000001。
	TransformID     UNICODE_LP_P4
	TransformName   UNICODE_LP_P4
	ReaderVersion   Version
	UpdaterVersion  Version
	WriterVersion   Version
}

// 2.1.9
type EncryptionTransformInfo struct {
	EncryptionName      UTF_8_LP_P4 // 甲UTF-8-LP-P4 结构（部分2.1.3），指定的加密算法的名称。名称必须是加密算法的名称，例如“AES 128”，“AES 192”或“AES 256”。与可扩展加密一起使用时，此值由可扩展加密模块指定。
	EncryptionBlockSize uint32      // 指定由EncryptionName指定的加密算法的块大小 。它必须是由高级加密标准（AES）指定的0x00000010 。与可扩展加密一起使用时，此值由可扩展加密模块指定。
	CipherMode          uint32      // 必须为0x00000000的值，除非与可扩展加密一起使用。与可扩展加密一起使用时，此值由可扩展加密模块指定。
	Reserved            uint32      // 必须是0x00000004的值。
}

// 2.2.5
type ExtensibilityHeader struct {
	Length uint32 // It MUST be 0x00000004.
}

// 2.2.6
type IRMDSTransformInfo struct {
	TransformInfoHeader
	ExtensibilityHeader
	XrMLLicense UTF_8_LP_P4
}

// 2.2.9
type EndUserLicenseHeader struct {
	Length    uint32
	ID_String UTF_8_LP_P4
}

// 2.2.10
// ECMA-376	MUST be named "EncryptedPackage"
// Other	MUST be named "\0x09DRMContent"
type ProtectedContentStream struct {
	Length   uint64
	Contents []byte
}

// 2.2.11
type ViewerContentStream struct {
	Length   uint64
	Contents []byte
}

// 2.3.2
type EncryptionHeader struct {
	Flags     int32 // EncryptionHeaderFlags
	SizeExtra int32 // MUST be 0x00000000
	AlgID     int32
	//	0x00000000	Determined by Flags
	//	0x00006801	RC4
	//	0x0000660E	128-bit AES
	//	0x0000660F	192-bit AES
	//	0x00006610	256-bit AES

	//	Flags.fCryptoAPI	Flags.fAES	Flags.fExternal	AlgID		Algorithm
	//	0					0			1				0x00000000	Determined by the application
	//	1					0			0				0x00000000	RC4
	//	1					0			0				0x00006801	RC4
	//	1					1			0				0x00000000	128-bit AES
	//	1					1			0				0x0000660E	128-bit AES
	//	1					1			0				0x0000660F	192-bit AES
	//	1					1			0				0x00006610	256-bit AES
	AlgIDHash int32
	//	AlgIDHash	Flags.fExternal	Algorithm
	//	0x00000000				1	Determined by the application
	//	0x00000000				0	SHA-1
	//	0x00008004				0	SHA-1
	KeySize uint32
	//	Algorithm	Value									Comment
	//	Any			0x00000000								Determined by Flags
	//	RC4			0x00000028 – 0x00000080 (inclusive)		8-bit increments
	//	AES			0x00000080, 0x000000C0, 0x00000100		128-bit, 192-bit, or 256-bit
	ProviderType int32

	Reserved1 int32
	Reserved2 int32
	CSPName   string // (variable)
}

// 2.3.3
type EncryptionVerifier struct {
	SaltSize              uint32 // It MUST be 0x00000010
	Salt                  [16]byte
	EncryptedVerifier     [16]byte
	VerifierHashSize      uint32
	EncryptedVerifierHash []byte //  (variable) RC4-20个字节。AES-32个字节
}

// 2.3.4.4
type EncryptedPackage struct {
	StreamSize    uint8
	EncryptedData []byte // (variable)
}

// 2.3.4.5
type EncryptionInfo struct {
	EncryptionVersionInfo Version
	EncryptionHeaderFlags int32
	EncryptionHeaderSize  uint32
	EncryptionHeader
	//Field			Value
	//Flags			The fCryptoAPI and fAES bits MUST be set. The fDocProps bit MUST be 0.
	//SizeExtra		This value MUST be 0x00000000.
	//AlgID			This value MUST be 0x0000660E (AES-128), 0x0000660F (AES-192), or 0x00006610 (AES-256).
	//AlgIDHash		This value MUST be 0x00008004 (SHA-1).
	//KeySize		This value MUST be 0x00000080 (AES-128), 0x000000C0 (AES-192), or 0x00000100 (AES-256).
	//ProviderType	This value SHOULD<10> be 0x00000018 (AES).
	//Reserved1		This value is undefined and MUST be ignored.
	//Reserved2		This value MUST be 0x00000000 and MUST be ignored.
	//CSPName		This value SHOULD<11> be set to either "Microsoft Enhanced RSA and AES Cryptographic Provider" or "Microsoft Enhanced RSA and AES Cryptographic Provider (Prototype)" as a null-terminated Unicode string.
	EncryptionVerifier
}

type UNICODE_LP_P4 struct {
	Length  uint32 // 它必须是2字节的倍数。
	Data    []byte
	Padding []byte // 结构体必须是4bytes的倍数，这个是填充
}
type UTF_8_LP_P4 struct {
	Length  uint32 // 它必须是4字节的倍数。
	Data    []byte
	Padding []byte // 结构体必须是4bytes的倍数，这个是填充
}
