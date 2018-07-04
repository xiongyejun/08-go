// permutation and combination
package permuCombin

import (
	"errors"
	"strings"
)

// 排列
func Permu(src []byte, selectCount uint, ch chan []byte) (resultCount uint32, err error) {
	if selectCount == 0 {
		return 0, errors.New("选择的个数为0。")
	}

	iCount := uint(len(src))
	if iCount == 0 {
		return 0, errors.New("数据源为空。")
	}

	// src 0-4		selectCount=3			  i
	// 程序首先都选择第0个数据源			000
	// 然后									    001	002	003	004
	// 到了4也就是到达了src的最后，进位		010	011 012 013 014
	//										020 021 022 023 024
	//										030 031 032 033 034
	//										040 041 042 043 044
	//										100 101 102 103 104
	// ……
	//										440 441 442 443 444
	arrPointer := make([]uint, selectCount)
	var j uint = 0

	// 从最后1个位置开始变化
	index := int(selectCount - 1)
	lastIndex := index
	for {
		// 永远是最后1个位置增加，然后进位过去
		for ; arrPointer[lastIndex] < iCount; arrPointer[lastIndex]++ {

			// 生成1个结果
			item := make([]byte, selectCount)
			for j = 0; j < selectCount; j++ {
				item[j] = src[arrPointer[j]]
			}
			ch <- item
		}

		index = lastIndex
		// 进位
		// 找到还没有超过iCount的index
		for arrPointer[index] == iCount {
			arrPointer[index] = 0
			index--
			if index >= 0 {
				arrPointer[index]++
			} else {
				return 1, nil
			}
		}
	}

	return 1, nil
}
func PermuString(src []string, selectCount uint, ch chan []byte) (resultCount uint32, err error) {
	if selectCount == 0 {
		return 0, errors.New("选择的个数为0。")
	}

	iCount := uint(len(src))
	if iCount == 0 {
		return 0, errors.New("数据源为空。")
	}

	arrPointer := make([]uint, selectCount)
	var j uint = 0

	// 从最后1个位置开始变化
	index := int(selectCount - 1)
	lastIndex := index
	for {
		// 永远是最后1个位置增加，然后进位过去
		for ; arrPointer[lastIndex] < iCount; arrPointer[lastIndex]++ {

			// 生成1个结果
			item := make([]string, selectCount)
			for j = 0; j < selectCount; j++ {
				item[j] = src[arrPointer[j]]
			}
			ch <- []byte(strings.Join(item, ""))
		}

		index = lastIndex
		// 进位
		// 找到还没有超过iCount的index
		for arrPointer[index] == iCount {
			arrPointer[index] = 0
			index--
			if index >= 0 {
				arrPointer[index]++
			} else {
				return 1, nil
			}
		}
	}

	return 1, nil
}

// 组合
func Combin(src []byte, selectCount uint, ch chan []byte) (resultCount uint32, err error) {
	if selectCount == 0 {
		return 0, errors.New("选择的个数为0。")
	}

	iCount := uint(len(src))
	if iCount < selectCount {
		return 0, errors.New("数据源的个数少于选择的个数。")
	}

	return 1, nil
}
