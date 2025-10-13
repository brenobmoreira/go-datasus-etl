package entities

type Estabelecimento struct {
	ID              string `dbase:"CNES"`
	CodigoMunicipio string `dbase:"CODUFMUN"`
}

type Equipamentos struct {
	ID                  string `dbase:"CNES"`
	CodigoEquipamento   string `dbase:"TIPEQUIP"`
	QuantidadeExistente int64  `dbase:"QT_EXIST"`
	QuantidadeUso       int64  `dbase:"QT_USO"`
}

type EstabelecimentoCadastro struct {
	ID         string `dbase:"CNES"`
	Nome       string `dbase:"FANTASIA"`
	Logradouro string `dbase:"NO_LOGRAD"`
	Numero     string `dbase:"NUMERO_END"`
	Bairro     string `dbase:"NO_BAIRRO"`
	CEP        string `dbase:"CO_CEP"`
}

type EquipamentoDescricao struct {
	CodigoEquipamento string `dbase:"CHAVE"`
	Descricao         string `dbase:"DS_TPEQUIP"`
}
