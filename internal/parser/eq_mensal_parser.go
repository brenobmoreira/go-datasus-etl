package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/brenobmoreira/go-datasus-etl/internal/entities"
	"github.com/valentin-kaiser/go-dbase/dbase"
)

func EquipamentoParser(archive_name string, blast string, dir string) {
	dbf_path := dir + "/data/dbf/" + archive_name + ".dbf"
	dbc_path := dir + "/data/dbc/" + archive_name + ".dbc"
	outputDir := filepath.Dir(dbf_path)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		panic(err)
	}

	err := DBCtoDBF(dbc_path, dbf_path, blast, dir)
	if err != nil {
		panic(err)
	}

	table, err := ReadDbf(dbf_path)
	if err != nil {
		panic(err)
	}
	defer table.Close()

	path := dir + "/data/csv/" + archive_name + ".csv"
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0755); err != nil {
		panic(err)
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var line uint32

	equipChan := make(chan entities.Equipamentos)
	go WriteEquipamento(file, equipChan)

	for !table.EOF() {
		line++
		row, err := table.Next()
		if err != nil {
			panic(err)
		}

		if row.Deleted {
			fmt.Printf("Deleted row at position: %v \n", row.Position)
			continue
		}

		p, err := RowToEquipamento(row)
		if err != nil {
			fmt.Printf("Error in row: %d, %v", line, err)
			continue
		}

		equipChan <- p
	}
}

func WriteEquipamento(file *os.File, equipChan chan entities.Equipamentos) {
	w := csv.NewWriter(file)
	for r := range equipChan {
		qtdeExistenteStr := strconv.FormatInt(r.QuantidadeExistente, 10)
		qtdeEmUsoStr := strconv.FormatInt(r.QuantidadeUso, 10)
		record := []string{
			r.ID,
			r.CodigoEquipamento,
			qtdeExistenteStr,
			qtdeEmUsoStr,
		}
		if err := w.Write(record); err != nil {
			panic(err)
		}
		fmt.Println(record)
	}
	defer w.Flush()
}

func RowToEquipamento(row *dbase.Row) (entities.Equipamentos, error) {
	p := &entities.Equipamentos{}
	err := row.ToStruct(p)
	if err != nil {
		return entities.Equipamentos{}, err
	}

	return *p, nil
}
