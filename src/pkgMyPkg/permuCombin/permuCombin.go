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

// 排列
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
	index := int(selectCount - 1)
	lastIndex := index
	for {
		for ; arrPointer[lastIndex] < iCount; arrPointer[lastIndex]++ {
			item := make([]string, selectCount)
			for j = 0; j < selectCount; j++ {
				item[j] = src[arrPointer[j]]
			}
			ch <- []byte(strings.Join(item, ""))
		}

		index = lastIndex
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
func Combin(src []byte, selectCount int, ch chan []byte) (resultCount int, err error) {
	if selectCount <= 0 {
		return 0, errors.New("选择的个数必须为正数。")
	}

	iCount := len(src)
	if iCount < selectCount {
		return 0, errors.New("数据源的个数少于选择的个数。")
	}
	// src 0-4		selectCount=3			  i
	// 程序首先都选择						012
	// 然后									    013	014
	// 到了4也就是到达了src的最后，进位		023	024
	//										034
	//										123 124
	//										134
	//										234

	arrPointer := make([]int, selectCount)
	var j int = 0
	// 最后1个往后移动，到达最后，进位
	// 从最后1个位置开始变化
	index := selectCount - 1
	lastIndex := index
	// 初始为012
	for j = 0; j < selectCount; j++ {
		arrPointer[j] = j
	}

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

		index = lastIndex - 1
		// 进位
		// 找到还允许增加的那个位置（最后1个的下标最大是lastIndex，倒数第2个就是lastIndex-1）
		for arrPointer[index] >= (iCount - selectCount + index) {
			//            2        5   -     3      +2=4
			//            1        5         3       1=3
			index--
			if index < 0 {
				return getCombinCount(iCount, selectCount), nil
			}
		}
		// 后面的都比前1个+1
		arrPointer[index]++
		for i := index + 1; i <= lastIndex; i++ {
			arrPointer[i] = arrPointer[i-1] + 1
		}
	}

	return getCombinCount(iCount, selectCount), nil
}

func getPermuCount(m, n int) (icount int) {
	icount = 1
	for i := 0; i < n; i++ {
		icount *= m
	}
	return
}

func fact(n int) (icount int) {
	if n <= 0 {
		return 0
	}
	icount = 1
	for i := n; i > 0; i-- {
		icount *= i
	}
	return
}

func getCombinCount(m, n int) (icount int) {
	if n <= 0 || m <= 0 || m < n {
		return 0
	}

	return fact(m) / fact(n) / fact(m-n)
}
