DROP TABLE IF EXISTS estabelecimento;
DROP TABLE IF EXISTS equipamento;
DROP TABLE IF EXISTS estabelecimento_cadastro;
DROP TABLE IF EXISTS equipamento_descricao;

CREATE TABLE estabelecimento_cadastro (
    cnes TEXT PRIMARY KEY,
    nome TEXT
);

CREATE TABLE equipamento_descricao (
    codigo TEXT PRIMARY KEY,
    descricao TEXT
);

CREATE TABLE estabelecimento (
    cnes TEXT REFERENCES estabelecimento_cadastro(cnes),
    codigo_municipio TEXT,
    competencia DATE,
    PRIMARY KEY (cnes, competencia)
);

CREATE TABLE equipamento (
    cnes TEXT REFERENCES estabelecimento_cadastro(cnes),
    codigo_equipamento TEXT,
    quantidade_existente INTEGER,
    quantidade_uso INTEGER,
    competencia DATE,
    PRIMARY KEY (cnes, competencia)
);
