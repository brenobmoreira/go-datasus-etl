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
	ID       string `dbase:"CNES"`
	Nome     string `dbase:"FANTASIA"`
	Excluido string `dbase:"EXCLUIDO"`
}

type EquipamentoDescricao struct {
	CodigoEquipamento string `dbase:"CHAVE"`
	Descricao         string `dbase:"DS_TPEQUIP"`
}
