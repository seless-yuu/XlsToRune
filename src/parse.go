package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	SRuneType = "RuneType"
	SType     = "type"
	SEnum     = "enum"
	SString   = "string"
	SInt      = "int"
	SInt2     = "int2"
	SInt3     = "int3"
	SFloat    = "float"
	SFloat2   = "float2"
	SFloat3   = "float3"
	SFloat4   = "float4"
	SComment  = "#"
)

const value_separator = ","

type RuneTypeName struct {
	Kind  string
	Value string
}

func (c *RuneTypeName) Print() {
	fmt.Println("--- RuneTypeName ---")
	fmt.Printf("kind  : %s\n", c.Kind)
	fmt.Printf("value : %s\n", c.Value)
}

type RuneTypeValue struct {
	Values []string
}

type RuneTypeType struct {
	TypeName RuneTypeName
	colIndex int
}

func (c *RuneTypeType) Print() {
	fmt.Println("--- RuneTypeValue ---")
	fmt.Printf("col_index   : %d\n", c.colIndex)
	c.TypeName.Print()
}

type RuneTypeTable struct {
	Name        string
	Types       []RuneTypeType
	Values      []RuneTypeValue
	typeIndex   int
	ignoreIndex []int
}

func (c *RuneTypeTable) Print() {
	fmt.Println("--- RuneTypeTable ---")
	fmt.Printf("name        : %s\n", c.Name)
	fmt.Printf("type  count : %d\n", len(c.Types))
	fmt.Printf("value count : %d\n", len(c.Values))
	for _, t := range c.Types {
		t.Print()
	}
}

func (c *RuneTypeTable) IsIgnoreIndex(index int) bool {
	for _, v := range c.ignoreIndex {
		if v == index {
			return true
		}
	}

	return false
}

type RuneTypeSheet struct {
	Name   string
	Tables []RuneTypeTable
}

func (c *RuneTypeSheet) Print() {
	fmt.Println("--- RuneTypeSheet ---")
	fmt.Printf("name : %s\n", c.Name)
	for _, v := range c.Tables {
		v.Print()
	}
}

type RuneTypeBook struct {
	Name   string
	Sheets []RuneTypeSheet
}

func (c *RuneTypeBook) Print() {
	fmt.Println("--- RuneTypeBook ---")
	fmt.Printf("name : %s\n", c.Name)
	for _, v := range c.Sheets {
		v.Print()
	}
}

type contextType int

const (
	ContextNone contextType = iota
	ContextTable
)

type context struct {
	runeBook      RuneTypeBook
	currentType   contextType
	currentSheet  RuneTypeSheet
	pcurrentTable *RuneTypeTable
	rowIndex      int
	colIndex      int
}

var gctx context

func ParseXls(path string) (RuneTypeBook, error) {
	file, err := excelize.OpenFile(path)
	if err != nil {
		return gctx.runeBook, err
	}

	name := strings.Split(filepath.Base(path), ".")[0]
	gctx.runeBook = RuneTypeBook{}
	gctx.runeBook.Name = name

	sheets := file.GetSheetList()
	for i := 0; i < len(sheets); i++ {
		name := sheets[i]

		rows, err := file.Rows(name)
		if err != nil {
			return gctx.runeBook, err
		}

		gctx.rowIndex = 0

		gctx.currentSheet = RuneTypeSheet{}
		gctx.currentSheet.Name = name
		for rows.Next() {
			cols, err := rows.Columns()
			if err != nil {
				return gctx.runeBook, err
			}

			err = parseCols(cols)
			if err != nil {
				return gctx.runeBook, err
			}
			gctx.rowIndex++
		}

		if len(gctx.currentSheet.Tables) > 0 {
			gctx.runeBook.Sheets = append(gctx.runeBook.Sheets, gctx.currentSheet)
		}
	}

	return gctx.runeBook, nil
}

func parseCols(cols []string) error {
	col_len := len(cols)

	switch gctx.currentType {
	case ContextNone:
		return parseColsForNone(cols)
	case ContextTable:
		if col_len > 0 {
			if cols[0] == SRuneType {
				return parseColsForNone(cols)
			} else {
				return parseColsForTable(cols)
			}
		}
	}

	return nil
}

func parseColsForNone(cols []string) error {
	if len(cols) <= 0 {
		return nil
	}

	gctx.colIndex = 0
	col := cols[0]
	if col == SRuneType {
		err := newCurrentTable(cols)
		if err != nil {
			return err
		}
		gctx.currentType = ContextTable
	}

	return nil
}

func parseColsForTable(cols []string) error {
	result := RuneTypeValue{}

	table_len := len(gctx.currentSheet.Tables)
	if table_len <= 0 {
		return nil
	}

	current_table := &gctx.currentSheet.Tables[table_len-1]
	col_len := len(cols)
	type_len := len(current_table.Types)
	if type_len <= 0 || col_len <= 0 {
		return nil
	}

	gctx.colIndex = 0

	begin_index := current_table.typeIndex + 1
	end_index := begin_index + type_len - 1 // -1 はテーブル型名分
	for i := begin_index; i <= end_index; i++ {
		if current_table.IsIgnoreIndex(i) {
			continue
		}

		if i < col_len {
			col_type := current_table.Types[i-2]
			col := cols[i]
			col = strings.TrimSpace(col)

			switch col_type.TypeName.Kind {
			case SInt:
				err := checkInt(col)
				if err != nil {
					return err
				}

			case SInt2:
			case SInt3:
				tokens := strings.Split(col, value_separator)
				for _, token := range tokens {
					err := checkInt(token)
					if err != nil {
						return err
					}
				}

			case SFloat:
				err := checkFloat(col)
				if err != nil {
					return err
				}

			case SFloat2:
			case SFloat3:
			case SFloat4:
				tokens := strings.Split(col, value_separator)
				for _, token := range tokens {
					err := checkFloat(token)
					if err != nil {
						return err
					}
				}
			}

			result.Values = append(result.Values, col)
		} else {
			result.Values = append(result.Values, "")
		}

		gctx.colIndex++
	}
	current_table.Values = append(current_table.Values, result)

	return nil
}

func checkInt(token string) error {
	token = strings.TrimSpace(token)
	_, err := strconv.Atoi(token)
	if err != nil {
		return fmt.Errorf("値:%s は整数ではありません", token)
	} else {
		return nil
	}
}

func checkFloat(token string) error {
	token = strings.TrimSpace(token)
	_, err := strconv.ParseFloat(token, 32)
	if err != nil {
		return fmt.Errorf("値:%s は浮動小数ではありません", token)
	} else {
		return nil
	}
}

func newCurrentTable(cols []string) error {
	table := RuneTypeTable{}
	gctx.colIndex = 0
	gctx.pcurrentTable = &table

	for i, col := range cols {
		value := strings.Trim(col, " ")

		if isComment(value) {
			table.ignoreIndex = append(table.ignoreIndex, i)
		} else if strings.Contains(value, SEnum) {
			v := parseSEnum()
			table.Types = append(table.Types, v)
		} else {
			err := checkTypeValidity(value)
			if err != nil {
				return err
			}

			if strings.Contains(value, SRuneType) {
				// RuneType宣言なので何もしない
				continue
			}

			type_string := strings.TrimSpace(strings.Split(value, ":")[1])
			if strings.Contains(type_string, SType) {
				err = parseSType(col)
				table.typeIndex = i
			} else if strings.Contains(type_string, SString) {
				v := parseSString(col)
				table.Types = append(table.Types, v)
			} else if strings.Contains(type_string, SInt2) {
				v := parseSInt2(col)
				table.Types = append(table.Types, v)
			} else if strings.Contains(type_string, SInt3) {
				v := parseSInt3(col)
				table.Types = append(table.Types, v)
			} else if strings.Contains(type_string, SInt) {
				v := parseSInt(col)
				table.Types = append(table.Types, v)
			} else if strings.Contains(type_string, SFloat2) {
				v := parseSFloat2(col)
				table.Types = append(table.Types, v)
			} else if strings.Contains(type_string, SFloat3) {
				v := parseSFloat3(col)
				table.Types = append(table.Types, v)
			} else if strings.Contains(type_string, SFloat4) {
				v := parseSFloat4(col)
				table.Types = append(table.Types, v)
			} else if strings.Contains(type_string, SFloat) {
				v := parseSFloat(col)
				table.Types = append(table.Types, v)
			}

			if err != nil {
				return err
			}
		}

		gctx.colIndex++
	}

	if len(table.Types) > 0 {
		gctx.currentSheet.Tables = append(gctx.currentSheet.Tables, table)
	}

	return nil
}

func checkTypeValidity(str string) error {
	strs := strings.Split(str, ":")
	str_len := len(strs)

	validity := true
	if str_len == 1 {
		str := strs[0]
		if strings.Contains(str, SRuneType) {
			validity = true
		} else {
			validity = false
		}
	} else if str_len == 2 {
		validity = true
	} else {
		validity = false
	}

	if validity {
		return nil
	} else {
		return fmt.Errorf("%s : invalid type string : %s", makeRuneErrorPrefix(), str)
	}
}

func makeRuneErrorPrefix() string {
	return fmt.Sprintf("row, %d, col %d", gctx.rowIndex, gctx.colIndex)
}

func parseSType(str string) error {
	str_array := parseTypeString(str)
	if len(str_array) > 1 {
		name := str_array[0]
		if str_array[1] != SType {
			return fmt.Errorf("row %d, col %d, not found RuneType", gctx.rowIndex, gctx.colIndex)
		}
		gctx.pcurrentTable.Name = name
	}

	return nil
}

func parseSEnum() RuneTypeType {
	result := RuneTypeType{}
	result.TypeName.Kind = SEnum
	result.colIndex = gctx.colIndex

	return result
}

func parseSString(str string) RuneTypeType {
	return parseTypeValue(str, SString)
}

func parseSInt(str string) RuneTypeType {
	return parseTypeValue(str, SInt)
}

func parseSInt2(str string) RuneTypeType {
	return parseTypeValue(str, SInt2)
}

func parseSInt3(str string) RuneTypeType {
	return parseTypeValue(str, SInt3)
}

func parseSFloat(str string) RuneTypeType {
	return parseTypeValue(str, SFloat)
}

func parseSFloat2(str string) RuneTypeType {
	return parseTypeValue(str, SFloat2)
}

func parseSFloat3(str string) RuneTypeType {
	return parseTypeValue(str, SFloat3)
}

func parseSFloat4(str string) RuneTypeType {
	return parseTypeValue(str, SFloat4)
}

func isComment(str string) bool {
	return strings.Contains(str, SComment)
}

func parseTypeValue(str string, type_name string) RuneTypeType {
	result := RuneTypeType{}
	result.TypeName.Kind = type_name

	strs := parseTypeString(str)
	result.TypeName.Value = strs[0]
	result.colIndex = gctx.colIndex

	return result
}

func parseTypeString(str string) []string {
	str = strings.Trim(str, " ")
	strs := strings.Split(str, ":")
	for i, s := range strs {
		strs[i] = strings.Trim(s, " ")
	}

	return strs
}
