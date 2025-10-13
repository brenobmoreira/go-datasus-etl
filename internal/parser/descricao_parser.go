package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brenobmoreira/go-datasus-etl/internal/entities"
	"github.com/valentin-kaiser/go-dbase/dbase"
)

func DescricaoParser(archive_name string, blast string, dir string) {
	dbf_path := dir + "/assets/" + archive_name + ".dbf"
	outputDir := filepath.Dir(dbf_path)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
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

	descricaoChan := make(chan entities.EquipamentoDescricao)
	go WriteDescricao(file, descricaoChan)

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

		p, err := RowToDescricao(row)
		if err != nil {
			fmt.Printf("Error in row: %d, %v", line, err)
			continue
		}

		descricaoChan <- p
	}
}

func WriteDescricao(file *os.File, descricaoChan chan entities.EquipamentoDescricao) {
	w := csv.NewWriter(file)
	for r := range descricaoChan {
		record := []string{
			r.CodigoEquipamento,
			r.Descricao,
		}
		if err := w.Write(record); err != nil {
			panic(err)
		}
		fmt.Println(record)
	}
	defer w.Flush()
}

func RowToDescricao(row *dbase.Row) (entities.EquipamentoDescricao, error) {
	p := &entities.EquipamentoDescricao{}
	err := row.ToStruct(p)
	if err != nil {
		return entities.EquipamentoDescricao{}, err
	}

	return *p, nil
}
