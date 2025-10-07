package datasus

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
)

type Info struct {
	UF  string `json:"uf"`
	Ano string `json:"ano"`
	Mes string `json:"mes"`
}

func DownloadDBC(infos []Info) error {
	path := "ftp.datasus.gov.br:21"
	login, password := "anonymous", "anonymous"
	resp, err := ConnectFtp(path, login, password)
	if err != nil {
		panic(err)
	}

	st_eq := []string{"ST", "EQ"}
	for i := range st_eq {
		cadastro := st_eq[i]
		download := "data/dbc/" + cadastro + "/"
		initial := "/dissemin/publicos/CNES/200508_/Dados/" + cadastro + "/"
		for j := range infos {
			archive_name := cadastro + infos[j].UF + infos[j].Ano + infos[j].Mes + ".dbc"
			ftp_path := initial + archive_name
			err = ReadAndDownload(resp, download, archive_name, ftp_path)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

func ConnectFtp(path string, login string, password string) (response *ftp.ServerConn, err error) {
	resp, err := ftp.Dial(path, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = resp.Login(login, password)
	if err != nil {
		log.Fatal(err)
	}

	return resp, nil
}

func ReadAndDownload(resp *ftp.ServerConn, download_path string, name string, ftp_path string) error {
	read, err := resp.Retr(ftp_path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(download_path)
	if err != nil {
		err = os.MkdirAll(download_path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	archive, err := os.Create(download_path + name)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(archive, read)

	defer archive.Close()
	defer read.Close()

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
