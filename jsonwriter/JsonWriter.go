package jsonwriter

import (
	"fmt"
	"os"
)

type JsonWriter struct {
	buffer       []byte
	filepath     string
	file_desc    *os.File
	buffer_index int
	buffer_size  int
	lvl          int
	lvl_created  []bool
}

// Constructor
func NewJsonWriter(FilePath string) *JsonWriter {
	file_desc, err := os.Create(FilePath)
	if err != nil {
		fmt.Println("Some problem occurred while creating the file:", FilePath)
	}

	writer := JsonWriter{
		buffer:       make([]byte, 256),
		filepath:     FilePath,
		file_desc:    file_desc,
		buffer_index: 0,
		buffer_size:  256,
		lvl:          1,
		lvl_created:  make([]bool, 1),
	}

	return &writer
}

// Core I/O Methods
func (writer *JsonWriter) WriteToFile() {
	writer.file_desc.Write(writer.buffer[:writer.buffer_index])
	writer.buffer_index = 0
}

func (writer *JsonWriter) CloseFile() {
	writer.WriteToFile() // flush before closing
	writer.file_desc.Close()
}

// Basic Printing Helpers
func (writer *JsonWriter) Print(pStr string) {
	for i := 0; i < len(pStr); i++ {
		if writer.buffer_index == writer.buffer_size {
			writer.WriteToFile()
		}
		writer.buffer[writer.buffer_index] = pStr[i]
		writer.buffer_index++
	}
}

func (writer *JsonWriter) PrintString(pStr string) {
	writer.Print("\"")
	writer.Print(pStr)
	writer.Print("\"")
}

func (writer *JsonWriter) PrintNewLine() {
	writer.Print("\n")
}

func (writer *JsonWriter) PrintIndent() {
	for i := 0; i <= writer.lvl; i++ {
		writer.Print("\t")
	}
}

// Internal Helper
func (writer *JsonWriter) ResizeLvlArray() {
	if writer.lvl >= len(writer.lvl_created) {
		writer.lvl_created = append(writer.lvl_created, false)
	}
}

// JSON Writing API
func (writer *JsonWriter) StartObject() {
	if writer.lvl_created[writer.lvl-1] {
		writer.Print(",")
	}

	writer.PrintNewLine()
	writer.PrintIndent()
	writer.Print("{\n")

	writer.ResizeLvlArray()
	writer.lvl_created[writer.lvl-1] = true
}

func (writer *JsonWriter) EndObject() {
	writer.PrintNewLine()
	writer.PrintIndent()
	writer.Print("}")
}

func (writer *JsonWriter) AddValue(pKey string, pVal string) {
	writer.PrintIndent()
	writer.PrintString(pKey)
	writer.Print(": ")
	writer.PrintString(pVal)
	writer.Print(",")
	writer.PrintNewLine()
}

func (writer *JsonWriter) StartArray(pKey string) {
	writer.PrintIndent()
	writer.PrintString(pKey)
	writer.Print(": ")
	writer.Print("[")

	writer.ResizeLvlArray()
	writer.lvl++
}

func (writer *JsonWriter) EndArray() {
	writer.Print("]")

	writer.lvl--
	writer.lvl_created[writer.lvl] = false
}
