package compiler

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"path/filepath"
	"strings"
)

func loadDataTable(file util.TableFile, fileName string, types *model.TypeManager) (ret []*model.DataTable) {
	for _, sheet := range file.Sheets() {

		tab := model.NewDataTable()

		// 单元测试或者某些特殊情况, Sheet名为空, 默认用文件名填充
		if sheet.Name() == "" {
			tab.HeaderType = strings.TrimSuffix(fileName, filepath.Ext(fileName))
		} else {
			tab.HeaderType = sheet.Name()
		}

		if types.ObjectExists(tab.HeaderType) {
			util.ReportError("DuplicateHeaderType", tab.HeaderType)
		}

		tab.FileName = fileName

		ret = append(ret, tab)

		maxCol := parseHeader(sheet, tab, types)
		checkHeaderTypes(tab, types)

		// 遍历所有数据行
		for row := maxHeaderRow; ; row++ {

			if sheet.IsRowEmpty(row, maxCol+1) {
				break
			}

			firstCol := sheet.GetValue(row, 0, nil)
			if strings.HasPrefix(firstCol, "#") {
				continue
			}

			tgtRow := tab.AddRow()

			// 读取每一行
			readOneRow(sheet, tab, row, tgtRow)
		}
	}

	for _, tab := range ret {
		checkRepeat(tab)
	}

	return
}

func checkRepeat(tab *model.DataTable) {
	for _, header := range tab.Headers {
		if header.TypeInfo.MakeIndex {
			checker := map[string]*model.Cell{}
			for row := 0; row < len(tab.Rows); row++ {
				cell := tab.GetCell(row, header.Col)
				if cell == nil {
					continue
				}

				if cell.Value == "" {
					continue
				}

				if _, ok := checker[cell.Value]; ok {
					util.ReportError("DuplicateValueInMakingIndex", cell.String())
				} else {
					checker[cell.Value] = cell
				}
			}
		}
	}
}

func readOneRow(sheet util.TableSheet, tab *model.DataTable, srcRow int, tgtRow *model.DataRow) bool {

	arrayCellByName := map[string]*model.Cell{}

	for _, header := range tab.Headers {

		if header.TypeInfo == nil {
			continue
		}

		// 浮点数用库取时，需要特殊处理
		isFloat := util.LanguagePrimitive(header.TypeInfo.FieldType, "go") == "float32"

		value := sheet.GetValue(srcRow, header.Col, &util.ValueOption{ValueAsFloat: isFloat})
		cell := tgtRow.AddCell()
		cell.Value = value

		if header.TypeInfo.IsArray() {
			preCell := arrayCellByName[header.TypeInfo.FieldName]
			if preCell == nil {
				arrayCellByName[header.TypeInfo.FieldName] = cell
				preCell = cell
			}

			for _, element := range strings.Split(value, header.TypeInfo.ArraySplitter) {
				preCell.ValueList = append(preCell.ValueList, element)
			}
		}

	}

	return true
}
