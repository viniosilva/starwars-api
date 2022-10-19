# Star Wars API

## Requisitos

- [go](https://tip.golang.org/doc/go1.19)
- [sqlboiler](https://github.com/volatiletech/sqlboiler)
- [mockgen](https://github.com/golang/mock)
- [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [swaggo](https://github.com/swaggo/swag)

## Instalação

```bash
$ make install
```

## Configuração

### Criando o arquivo .env

Para o ambiente local apensar crie a `.env` como no exemplo abaixo:

```txt
GIN_MODE=debug
MIGRATE_URL=mysql://luke:xQlpKD95kp20Wa1JAX6O@\(127.0.0.1:3306\)/starwars
MYSQL_PASSWORD=xQlpKD95kp20Wa1JAX6O
```

### Migrate

Após executar o comando `docker-compose`, pode ser necessário esperar alguns segundos até o banco estar apto a receber comandos:

```bash
$ docker-compose up -d
$ make migrate
```

## Executando

Antes de executar a API Rest, é recomendado executar o `feed database` para buscar os dados da API [SWAPI](https://swapi.dev/):

```bash
# Feed database
$ make run/feed-database

# API Rest
$ make run
```

Para visualizar a documentação das rotas localmente, após a API estiver em execução, basta acessar o [swagger](http:localhost:8080/api/swagger/index.html)

## Testes

```bash
# Testes unitários
$ make test

# Cobertura dos testes unitários
$ make test/cov
```

## Arquitetura do projeto

O código está organizado dentro da pasta `internal`, conforme recomendações do [golang-standards](https://github.com/golang-standards/project-layout/blob/master/README_ptBR.md#internal)

- **db**: códigos referentes a banco de dados
    - **migrations**: SQLs para as `migrations`
- **docs**: arquivos swagger
- **log**: arquivos de logs
- **internal**: [golang-standards](https://github.com/golang-standards/project-layout/blob/master/README_ptBR.md#internal)
    - **config**: configurações globais do projeto
    - **controller**: configurações das rotas
    - **dto**: objetos de transferência de dados entre as camadas
    - **exception**: exceções tratadas
    - **model**: representações dos modelos e arquivos gerados pelo `sqlboiler`
    - **request**: abstrações de comunicações com serviços externos
    - **script**: rotinas auxiliares
    - **service**: regras de negócio
- **mock**: arquivos `mock` para dar suporte aos testes unitários

## Atualizando as models com SQLBoiler

Ao executar uma nova migração na base de dados:

- Adicione a senha `xQlpKD95kp20Wa1JAX6O` no mysql password do arquivo `sqlboiler.toml`
- Execute o comando: `make models`
- Troque todos os packages para `model` nos arquivos `internal/models/*`
- Mova todos os arquivos gerados para `internal/model` e delete a pasta `internal/models`
