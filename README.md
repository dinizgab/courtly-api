# Courtly API

API em Go para gerenciar reservas de quadras esportivas, pagamentos via Pix e envio de notificações.

## Requisitos

- Go 1.24+
- PostgreSQL
- [Goose](https://github.com/pressly/goose) para migrações
- Docker (opcional para ambiente local)

## Variáveis de ambiente

| Nome | Descrição |
| --- | --- |
| `API_PORT` | Porta em que a API será exposta |
| `JWT_SECRET` | Chave usada para assinar tokens JWT |
| `DATABASE_URL` | URL de conexão com o banco PostgreSQL |
| `SMTP_EMAIL` | Remetente utilizado para envio de e-mails |
| `SMTP_HOST` | Host do servidor SMTP |
| `SMTP_PORT` | Porta do servidor SMTP |
| `SMTP_USER` | Usuário do servidor SMTP |
| `SMTP_PASS` | Senha do servidor SMTP |
| `OPENPIX_BASE_URL` | URL base da API OpenPix |
| `OPENPIX_APP_ID` | ID de aplicação do OpenPix |
| `STORAGE_PROJECT_URL` | URL do projeto de armazenamento Supabase |
| `STORAGE_API_KEY` | Chave de API para o armazenamento |

## Executando localmente

1. **Inicie os serviços de apoio (opcional):**

   ```bash
   make compose-up
   ```

2. **Execute as migrações do banco:**

   ```bash
   make migration-up
   ```

3. **Inicie a aplicação:**

   ```bash
   make run
   ```

A API estará acessível em `http://localhost:$API_PORT`.

## Estrutura

- `cmd/main.go` – ponto de entrada da aplicação.
- `internal/` – implementação dos módulos de domínio, repositórios, casos de uso e handlers.
- `migrations/` – scripts SQL para criação e alteração do banco de dados.

