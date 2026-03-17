CREATE TABLE IF NOT EXISTS pessoas (
    id VARCHAR PRIMARY KEY,
    apelido VARCHAR(32),
    nome VARCHAR(100),
    nascimento VARCHAR(10),
    stack VARCHAR,
    busca_termos VARCHAR
);
