// 将某个文件夹下的文件打包成1个文件
// 遍历文件夹，把读取的每个文件写入某个文件，同时记录这个文件的名称、大小等
// 读取完成后，再把文件的名称，大小信息存入到某个文件的最后面

// 读取文件先从最后读取目录大小，再读取目录，根据目录读取文件数据
package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const Path_Separator string = string(os.PathSeparator)

// 目录里存储的信息
type dirInfo struct {
	FileName string // 文件名称
	Start    int64  // 文件开始的位置
	Size     int    // 文件大小byte
	Star     int    // 标星，初始为1
}

// 读取打包的文件
type unPackFile struct {
	f             *os.File
	files         []*dirInfo
	unPacked      map[int]bool // 记录已经读取了的文件的下标
	unPackedFiles []string     // 记录已经读取了的文件

	dirIndex int64 // 目录开始的地方
	bSave    bool  //退出时候，是否需要重新记录目录，重新标记star的时候需要
}

// 打包文件
type packFile struct {
	strDir  string // 打包文件夹的路径
	strSave string // 保存的打包文件名
	f       *os.File
	files   []*dirInfo
	next    int64 // 下一个文件的位置
}

// 打包文件
func PackFile(strPath string, saveName string) {
	if !strings.HasSuffix(strPath, Path_Separator) {
		strPath = strPath + Path_Separator
	}
	pack_File.strDir = strPath
	pack_File.strSave = saveName
	pack_File.files = make([]*dirInfo, 0)
	pack_File.next = 0

	if err := pack_File.packFiles(); err != nil {
		fmt.Println(err)
	}
}

// 读取打包了的文件
func (me *unPackFile) unPackInit(packfile string) (err error) {
	if me.f, err = os.OpenFile(packfile, os.O_RDWR, 0666); err != nil {
		return
	}
	// 读取文件最后几个字节，目录的大小
	finfo, _ := me.f.Stat()
	filesize := finfo.Size()
	fmt.Println("filesize", filesize)

	var dirSize int64
	me.f.Seek(filesize-int64(binary.Size(dirSize)), 0)
	b := make([]byte, binary.Size(dirSize))
	if _, err = me.f.Read(b); err != nil {
		return
	}
	dirSize, _ = binary.Varint(b)
	fmt.Println("dirszie:", dirSize)
	//读取目录，dirSize是包含了目录和表示目标大小的int64
	fmt.Println("读取的目录开始位置：", filesize-dirSize)

	me.dirIndex = filesize - dirSize
	me.f.Seek(me.dirIndex, 0)
	bDir := make([]byte, dirSize-int64(binary.Size(dirSize)))
	if _, err = me.f.Read(bDir); err != nil {
		return
	}
	// 目录是用json存储的
	if err = json.Unmarshal(bDir, &me.files); err != nil {
		return
	}
	return nil
}

func (me *unPackFile) unPackFile(index int) (saveName string, err error) {
	// 已经释放了的就不需要再一次释放
	if _, ok := me.unPacked[index]; ok {
		return
	}
	// 设置读取的位置
	if _, err = me.f.Seek(me.files[index].Start, 0); err != nil {
		return
	}
	// 读取文件的数据
	b := make([]byte, me.files[index].Size)
	if _, err = me.f.Read(b); err != nil {
		return
	}
	// 保存文件	程序路径+index+文件后缀
	saveName = currentPath + strconv.Itoa(index) + filepath.Ext(me.files[index].FileName)
	if err = ioutil.WriteFile(saveName, b, 0666); err != nil {
		return
	}
	// 记录文件名称，最后删除掉
	me.unPackedFiles = append(me.unPackedFiles, saveName)
	return saveName, nil
}

// 程序退出时候，删除已经释放的文件
func (me *unPackFile) deleteUnPackedFIle() {
	for i := range me.unPackedFiles {
		os.Remove(me.unPackedFiles[i])
	}
}

// 排序，按照后缀
func (me *unPackFile) Less(i, j int) bool {
	return filepath.Ext(me.files[i].FileName) < filepath.Ext(me.files[j].FileName)
}
func (me *unPackFile) Swap(i, j int) {
	me.files[i], me.files[j] = me.files[j], me.files[i]
}
func (me *unPackFile) Len() int {
	return len(me.files)
}
func (me *unPackFile) saveDir() {
	//写入目录
	if me.bSave {
		if err := writeDir(me.files, me.f, me.dirIndex); err != nil {
			fmt.Println(err)
		}
	}
	me.f.Close()
}

// 打包文件
func (me *packFile) packFiles() (err error) {
	if me.f, err = os.Create(me.strSave); err != nil {
		return
	}
	defer me.f.Close()
	// 遍历文件打包
	if err = me.scanDir(me.strDir); err != nil {
		return
	}

	//写入目录
	return writeDir(me.files, me.f, me.next)
}

func writeDir(files []*dirInfo, f *os.File, seek int64) (err error) {
	if _, err = f.Seek(seek, 0); err != nil {
		return
	}
	// 目录信息转换为json
	var b []byte
	if b, err = json.Marshal(files); err != nil {
		return
	} else {
		fmt.Println("文件个数：", len(files))

		if _, err = f.Write(b); err != nil {
			return
		}
		seek += int64(len(b))
		// 记录文件目录的大小
		var bDirSize int64
		bDirSize = int64(len(b) + binary.Size(bDirSize))

		var bb = make([]byte, binary.Size(bDirSize))
		binary.PutVarint(bb, bDirSize)
		if _, err = f.Write(bb); err != nil {
			return
		}
	}
	return nil
}

// 遍历文件进行打包
func (me *packFile) scanDir(strDir string) (err error) {
	return filepath.Walk(strDir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		if f.Name() == ".DS_Store" || f.Name() == "Thumbs.db" {
			return nil
		}

		if b, err := ioutil.ReadFile(path); err != nil {
			return errors.New("ReadFile出错。\n" + err.Error())
		} else {
			// 记录文件信息
			dir_info := &dirInfo{FileName: path, Start: me.next, Size: len(b), Star: 1}
			// 写入文件
			if n, err := me.f.Write(b); err != nil || n != dir_info.Size {
				return errors.New("写入文件出错。")
			}
			// 下一个文件的位置
			me.next += int64(dir_info.Size)
			me.files = append(me.files, dir_info)
		}
		return nil
	})
}
