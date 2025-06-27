# Processador de Logs

Este é um projeto de processamento de logs que utiliza Go, MySQL e RabbitMQ.
O objetivo é processar logs recebidos via RabbitMQ e armazená-los em um banco de dados MySQL.

## Sumário

- [Configuração do Ambiente para o Projeto de Logs](#configuração-do-ambiente-para-o-projeto-de-logs)

## Configuração do Ambiente para o Projeto de Logs

### Crie o banco `logs_db` no MySQL

```bash
CREATE DATABASE IF NOT EXISTS logs_db;
```

### Crie a tabela `logs` no banco de dados MySQL

````bash
CREATE TABLE IF NOT EXISTS logs (
id VARCHAR(255) PRIMARY KEY,
message TEXT,
severity VARCHAR(50),
source VARCHAR(255),
timestamp DATETIME,
processed_at DATETIME
);
````

### Crie o usuário `user_go` com senha `user1234.`

- Lembre de atribuir as permissões necessárias para o banco de dados.

### Certifique-se de que o MySQL está disponível em `localhost` ou `127.0.0.1` e rodando na porta `3306`.
- Essas configurações estão padronizadas, mas você pode ajustá-las conforme necessário no arquivo `main.go`