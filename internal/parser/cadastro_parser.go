package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brenobmoreira/go-datasus-etl/internal/entities"
	"github.com/valentin-kaiser/go-dbase/dbase"
)

func CadastroParser(archive_name string, blast string, dir string, cadastroChan chan entities.EstabelecimentoCadastro) {
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

		p, err := RowToCadastro(row)
		if err != nil {
			fmt.Printf("Error in row: %d, %v", line, err)
			continue
		}

		cadastroChan <- p
	}
}

func WriteCadastro(file *os.File, cadastroChan chan entities.EstabelecimentoCadastro) {
	w := csv.NewWriter(file)
	for r := range cadastroChan {
		if r.Excluido == "0" {
			record := []string{
				r.ID,
				r.Nome,
			}
			if err := w.Write(record); err != nil {
				panic(err)
			}
			fmt.Println(record)
		}
	}
	defer w.Flush()
}

func RowToCadastro(row *dbase.Row) (entities.EstabelecimentoCadastro, error) {
	p := &entities.EstabelecimentoCadastro{}
	err := row.ToStruct(p)
	if err != nil {
		return entities.EstabelecimentoCadastro{}, err
	}

	return *p, nil
}
