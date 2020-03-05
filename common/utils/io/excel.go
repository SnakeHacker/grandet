package io

import (
	"fmt"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/tealeg/xlsx"
)

// SaveSimpleXlsx ...
func SaveSimpleXlsx(filename string, dir string, sheetName string, headers []string, data [][]string) (err error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		glog.Error(err)
		return
	}

	// insert header
	r := sheet.AddRow()
	for _, header := range headers {
		cell := r.AddCell()
		cell.Value = header
	}

	// insert data
	for _, row := range data {
		r := sheet.AddRow()
		for _, col := range row {
			cell := r.AddCell()
			cell.Value = col
		}
	}

	err = file.Save(filepath.Join(dir, fmt.Sprintf("%s.xlsx", filename)))
	if err != nil {
		glog.Error(err)
		return
	}

	return
}
