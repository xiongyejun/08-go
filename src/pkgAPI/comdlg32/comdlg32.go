// GetOpenFileName 选择单个文件
// GetOpenFileNames 选择多个文件

package comdlg32

import (
	"path/filepath"
	"syscall"
	"unsafe"
)

type HWND uintptr
type HINSTANCE uintptr
type LPOFNHOOKPROC uintptr

type fileDialog struct {
	Title          string
	FilePath       string
	FilePaths      []string
	InitialDirPath string
	Filter         string
	FilterIndex    int
}
type oPENFILENAME struct {
	LStructSize       uint32
	HwndOwner         HWND
	HInstance         HINSTANCE
	LpstrFilter       *uint16
	LpstrCustomFilter *uint16
	NMaxCustFilter    uint32
	NFilterIndex      uint32
	LpstrFile         *uint16
	NMaxFile          uint32
	LpstrFileTitle    *uint16
	NMaxFileTitle     uint32
	LpstrInitialDir   *uint16
	LpstrTitle        *uint16
	Flags             uint32
	NFileOffset       uint16
	NFileExtension    uint16
	LpstrDefExt       *uint16
	LCustData         uintptr
	LpfnHook          LPOFNHOOKPROC
	LpTemplateName    *uint16
	PvReserved        unsafe.Pointer
	DwReserved        uint32
	FlagsEx           uint32
}

// GetOpenFileName and GetSaveFileName flags
const (
	OFN_ALLOWMULTISELECT     = 0x00000200
	OFN_CREATEPROMPT         = 0x00002000
	OFN_DONTADDTORECENT      = 0x02000000
	OFN_ENABLEHOOK           = 0x00000020
	OFN_ENABLEINCLUDENOTIFY  = 0x00400000
	OFN_ENABLESIZING         = 0x00800000
	OFN_ENABLETEMPLATE       = 0x00000040
	OFN_ENABLETEMPLATEHANDLE = 0x00000080
	OFN_EXPLORER             = 0x00080000
	OFN_EXTENSIONDIFFERENT   = 0x00000400
	OFN_FILEMUSTEXIST        = 0x00001000
	OFN_FORCESHOWHIDDEN      = 0x10000000
	OFN_HIDEREADONLY         = 0x00000004
	OFN_LONGNAMES            = 0x00200000
	OFN_NOCHANGEDIR          = 0x00000008
	OFN_NODEREFERENCELINKS   = 0x00100000
	OFN_NOLONGNAMES          = 0x00040000
	OFN_NONETWORKBUTTON      = 0x00020000
	OFN_NOREADONLYRETURN     = 0x00008000
	OFN_NOTESTFILECREATE     = 0x00010000
	OFN_NOVALIDATE           = 0x00000100
	OFN_OVERWRITEPROMPT      = 0x00000002
	OFN_PATHMUSTEXIST        = 0x00000800
	OFN_READONLY             = 0x00000001
	OFN_SHAREAWARE           = 0x00004000
	OFN_SHOWHELP             = 0x00000010
)

// GetOpenFileName and GetSaveFileName extended flags
const (
	OFN_EX_NOPLACESBAR = 0x00000001
)

var comdlg32 syscall.Handle

func init() {
	comdlg32, _ = syscall.LoadLibrary("comdlg32.dll")
}

func getOpenFileName(lpofn *oPENFILENAME) bool {
	getOpenFileName, _ := syscall.GetProcAddress(comdlg32, "GetOpenFileNameW")

	ret, _, _ := syscall.Syscall(getOpenFileName,
		1,
		uintptr(unsafe.Pointer(lpofn)),
		0,
		0,
	)
	return ret != 0
}

func commDlgExtendedError() uint32 {
	defer syscall.FreeLibrary(comdlg32)

	cdee, _ := syscall.GetProcAddress(comdlg32, "CommDlgExtendedError")
	ret, _, _ := syscall.Syscall(cdee,
		0,
		0,
		0,
		0)

	return uint32(ret)
}

func getSaveFileName(lpofn *oPENFILENAME) bool {
	gsfn, _ := syscall.GetProcAddress(comdlg32, "GetSaveFileNameW")

	ret, _, _ := syscall.Syscall(gsfn,
		1,
		uintptr(unsafe.Pointer(lpofn)),
		0,
		0)
	return ret != 0
}
func (dlg *fileDialog) show(fun func(ofn *oPENFILENAME) bool, flags int32) (bool, error) {
	ofn := new(oPENFILENAME)

	ofn.LStructSize = uint32(unsafe.Sizeof(*ofn))
	filter := make([]uint16, len(dlg.Filter)+2)
	copy(filter, syscall.StringToUTF16(dlg.Filter))
	// replace '|' with ten expected '\0'
	for i, c := range filter {
		if byte(c) == '|' {
			filter[i] = uint16(0)
		}
	}

	ofn.LpstrFilter = &filter[0]
	ofn.NFilterIndex = uint32(dlg.FilterIndex)

	ofn.LpstrInitialDir = syscall.StringToUTF16Ptr(dlg.InitialDirPath)
	ofn.LpstrTitle = syscall.StringToUTF16Ptr(dlg.Title)
	ofn.Flags = uint32(OFN_FILEMUSTEXIST | flags)

	var fileBuf []uint16
	if flags&OFN_ALLOWMULTISELECT > 0 {
		fileBuf = make([]uint16, 65536)
	} else {
		fileBuf = make([]uint16, 1024)
		copy(fileBuf, syscall.StringToUTF16(dlg.FilePath))
	}

	ofn.LpstrFile = &fileBuf[0]
	ofn.NMaxFile = uint32(len(fileBuf))

	if !fun(ofn) {
		errno := commDlgExtendedError()
		if errno != 0 {
			abort("CommDlgExtendedError", int(errno))
		}
		return false, nil
	}

	if flags&OFN_ALLOWMULTISELECT > 0 {
		split := func() [][]uint16 {
			var parts [][]uint16

			from := 0
			for i, c := range fileBuf {
				if c == 0 {
					if i == from {
						return parts
					}

					parts = append(parts, fileBuf[from:i])
					from = i + 1
				}
			}
			return parts
		}
		parts := split()

		if len(parts) == 1 {
			dlg.FilePaths = []string{syscall.UTF16ToString(parts[0])}
		} else {
			dirPath := syscall.UTF16ToString(parts[0])
			dlg.FilePaths = make([]string, len(parts)-1)

			for i, fp := range parts[1:] {
				dlg.FilePaths[i] = filepath.Join(dirPath, syscall.UTF16ToString(fp))
			}
		}
	} else {
		dlg.FilePath = syscall.UTF16ToString(fileBuf)
	}

	return true, nil
}

func (dlg *fileDialog) GetOpenFileName() (bool, error) {
	return dlg.show(getOpenFileName, 0)
}
func (dlg *fileDialog) GetOpenFileNames() (bool, error) {
	return dlg.show(getOpenFileName, OFN_ALLOWMULTISELECT|OFN_EXPLORER)
}
func (dlg *fileDialog) GetSaveFileName() (bool, error) {
	return dlg.show(getSaveFileName, 0)
}

func NewFileDialog() *fileDialog {
	return new(fileDialog)
}

func abort(funcName string, err int) {
	panic(funcName + " failed:" + syscall.Errno(err).Error())
}
