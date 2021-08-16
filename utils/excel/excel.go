package excel

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

func GetStyle(bold bool) (centerStyle *xlsx.Style, borderStyle *xlsx.Style, boldStyle *xlsx.Style, boderCenterStyle *xlsx.Style) {
	centerStyle = xlsx.NewStyle()
	centerStyle.Alignment = xlsx.Alignment{
		Horizontal:  "center",
		Vertical:    "center",
		WrapText:    true,
		ShrinkToFit: true,
	}

	borderStyle = xlsx.NewStyle()
	borderStyle.Border = xlsx.Border{
		Bottom:      "thin",
		BottomColor: "black",
		Top:         "thin",
		TopColor:    "black",
		Left:        "thin",
		LeftColor:   "black",
		Right:       "thin",
		RightColor:  "black",
	}
	borderStyle.ApplyBorder = true

	boldStyle = xlsx.NewStyle()
	boldStyle.Alignment = centerStyle.Alignment
	borderStyle.ApplyAlignment = true
	boldStyle.Font.Bold = bold
	boldStyle.ApplyFont = true
	boldStyle.Border = borderStyle.Border
	boldStyle.ApplyBorder = true

	boderCenterStyle = borderStyle
	boderCenterStyle.Alignment = centerStyle.Alignment
	return
}

func CreateSheetSample(xcl *xlsx.File, nameSheet string) (*xlsx.Sheet, error) {
	return xcl.AddSheet(nameSheet)
}

func CreateSheet(xcl *xlsx.File, styleHeader *xlsx.Style, title, sub, nameSheet string, cellCount int) (*xlsx.Sheet, error) {
	sheet, err := xcl.AddSheet(nameSheet)
	if err != nil {
		return nil, err
	}
	if title != "" {
		{
			sheet.Row(0).SetHeight(20)
			titleCell := sheet.Row(0).AddCell()
			titleCell.SetString(title)
			titleCell.HMerge = cellCount + 5
			style := xlsx.NewStyle()
			style.Alignment = styleHeader.Alignment
			style.ApplyAlignment = true
			style.Font.Bold = true
			style.ApplyFont = true
			if styleHeader != nil {
				titleCell.SetStyle(styleHeader)
			}
		}

		{
			sheet.Row(1).SetHeight(20)
			cell := sheet.Row(1).AddCell()
			cell.SetString(sub)
			cell.HMerge = cellCount + 5
			if styleHeader != nil {
				cell.SetStyle(styleHeader)
			}
		}
	}

	return sheet, nil
}

//Vertical and horizontal

func CreateHeader(objStruct interface{}, sheet *xlsx.Sheet, style *xlsx.Style, row, left int, isHorizontal, isAddStt bool) {
	var t = reflect.TypeOf(objStruct)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	sheet.Row(row) //.SetHeight(15)
	{
		lenField := t.NumField()
		for i := 0; i <= lenField; i++ {
			var name = "STT"
			if i > 0 {
				name = t.Field(i - 1).Name
			} else if !isAddStt {
				continue
			}
			cell := sheet.Cell(row, left)
			name = strings.ToLower(name)
			cell.SetString(name)
			if style != nil {
				cell.SetStyle(style)
			}
			if isHorizontal {
				left++
			} else {
				row++
			}
		}
	}
}

func CreateHeaderArr(headers []string, sheet *xlsx.Sheet, style *xlsx.Style, row, left int, isHorizontal, isAddStt bool) {
	sheet.Row(row) //.SetHeight(15)
	{
		lenField := len(headers)
		for i := 0; i <= lenField; i++ {
			var name = "STT"
			if i > 0 {
				name = headers[i]
			} else if !isAddStt {
				continue
			}
			cell := sheet.Cell(row, left)
			cell.SetString(name)
			if style != nil {
				cell.SetStyle(style)
			}
			if isHorizontal {
				left++
			} else {
				row++
			}
		}
	}
}

func CreateBody(sheet *xlsx.Sheet, style *xlsx.Style, objs []interface{}, row, left int, isHorizontal, isAddStt bool) {
	sheet.Row(row).SetHeight(15)
	{
		for m, val := range objs {
			var t = reflect.TypeOf(val)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			v := reflect.ValueOf(val)
			var lenField = t.NumField()
			for i := 0; i <= lenField; i++ {
				cellNum := left
				rowCell := row
				var value = ""
				if isHorizontal {

					cellNum = left + i
					if !isAddStt {
						cellNum = cellNum - 1
					}
					value = strconv.Itoa(m + 1)
				} else if isAddStt {
					rowCell = row + i
					value = strconv.Itoa(m + 1)
				}
				if i > 0 {
					var k = reflect.Indirect(v)
					value = k.FieldByName(t.Field(i - 1).Name).String()
				} else if !isAddStt {
					continue
				}

				cell := sheet.Cell(rowCell, cellNum)
				cell.SetString(value)
				if style != nil {
					cell.SetStyle(style)
				}
			}
			if isHorizontal {
				row++
			} else {
				left++
			}
		}
	}
}

func CreateHeaderCSV(objStruct interface{}, sheet *xlsx.Sheet, row, left int) {
	var t = reflect.TypeOf(objStruct)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	sheet.Row(row).SetHeight(60)
	{
		var name = "STT"
		lenField := t.NumField()
		for i := 0; i < lenField; i++ {
			name += "," + t.Field(i).Name
		}
	}
}

func CreateHeaderArrCSV(vals []string, sheet *xlsx.Sheet, row, left int) {
	sheet.Row(row) //.SetHeight(60)
	{
		var name = "STT"
		for _, val := range vals {
			name += "," + val
		}
		cell := sheet.Cell(row, left)
		cell.SetString(name)
	}
}

func CreateBodyCSV(sheet *xlsx.Sheet, objs []interface{}, row, left int) {
	// sheet.Row(row).SetHeight(20)
	// {
	for m, val := range objs {
		var t = reflect.TypeOf(val)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		v := reflect.ValueOf(val)
		var lenField = t.NumField()
		var value = strconv.Itoa(m + 1)
		for i := 0; i < lenField; i++ {
			if i > 0 {
				var k = reflect.Indirect(v)
				value = k.FieldByName(t.Field(i).Name).String()
			}

		}
		cell := sheet.Cell(row, left)
		cell.SetString(value)
		row++
	}
	//}
}
