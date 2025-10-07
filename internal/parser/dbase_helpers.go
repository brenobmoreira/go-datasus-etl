package parser

import (
	"fmt"
	"os/exec"

	"github.com/valentin-kaiser/go-dbase/dbase"
)

func ReadDbf(path string) (*dbase.File, error) {
	table, err := dbase.OpenTable(&dbase.Config{
		Filename:   path,
		TrimSpaces: true,
		Untested:   true,
	})
	if err != nil {
		return nil, err
	}

	return table, nil
}

func DBCtoDBF(dbc string, dbf string, blast string, dir string) error {
	cmd := exec.Command(blast, dbc, dbf)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running blast-dbf: %v\nOutput: %s\n", err, string(output))
		return err
	}

	fmt.Println("Successfully converted .dbc to .dbf")
	return err
}
