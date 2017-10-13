package main

import (
	"fmt"

	"github.com/xuri/excelize"
)

func main() {
	xlsx, err := excelize.OpenFile("./2.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	cell := xlsx.GetCellValue("月报表", "A1")
	fmt.Println(cell, "111")

}
