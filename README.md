# ApiMyChat

API de chat em Go (Gin) com WebSocket, PostgreSQL e integrações externas (Supabase Auth, Gmail SMTP para OTP e Firebase Cloud Messaging).

## Principais recursos
- Cadastro/login via Supabase e verificação de e-mail com OTP enviado por Gmail.
- Salas de conversa com membros, criação/listagem e mensagens persistidas em PostgreSQL.
- WebSocket para entrega em tempo real e hub de salas; fallback via endpoint HTTP para envio.
- Notificações push para participantes offline usando FCM; armazenamento de tokens por dispositivo.
- Endpoints REST para usuários, mídias anexadas às mensagens e tokens FCM.

## Requisitos
- Go 1.25+
- Docker + Docker Compose (opcional, recomendado)
- PostgreSQL 16
- Credenciais externas: Gmail app password, Supabase service role key, projeto Firebase (HTTP v1)

## Configuração rápida (Docker Compose)
1. Copie o exemplo: `cp .env.example .env` e preencha valores.
2. Suba tudo: `docker-compose up --build`.
3. API em `http://localhost:8000` (PostgreSQL exposto em 5432). O compose injeta `DB_HOST=postgres` para o container.

## Rodando localmente sem Docker
- Mantenha o PostgreSQL acessível (local ou remoto).
- Instale dependências: `go mod download`.
- Rode: `go run ./cmd/api` (usa `API_PORT` ou 8000 por padrão).

## Variáveis de ambiente
- `API_PORT` – porta HTTP (padrão 8000).
- Conexão com DB: `DB_URL` **ou** `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`.
- OTP por e-mail: `GMAIL_EMAIL`, `GMAIL_APP_PASSWORD`.
- Supabase: `SUPABASE_URL`, `SUPABASE_KEY` (service role).
- FCM: `FCM_PROJECT_ID`, `FCM_CLIENT_EMAIL`, `FCM_PRIVATE_KEY` (ou `FCM_SERVICE_ACCOUNT_JSON`), opcional `FCM_TOKEN_URI`.

## Endpoints principais
- `POST /send-code` – envia OTP por e-mail.
- `POST /verify-email` – confirma OTP.
- `POST /CreateUser` – cria usuário (Supabase + tabela `users`).
- `POST /login` – retorna token do Supabase.
- `GET /GetUserByID/:id`, `GET /GetAll/:id`, `PUT /UpdateUser`.
- `POST /CreateRoom` – cria sala com usuários.
- `GET /GetRoomsByUid/:uid` – salas do usuário.
- `GET /GetRoomByUid/:id` – sala + membros.
- `GET /messages?room=ID&limit=50` – mensagens paginadas.
- `GET /messages/last?room=ID` – última mensagem.
- `POST /CreateMedia` – associa URL de mídia a uma mensagem.
- `GET /GetMediaByMessageId/:messageId` – URLs de mídias.
- `POST /fcm/SaveToken` – salva token de dispositivo.
- `DELETE /fcm/token?token=...` – remove token.
- WebSocket: `GET /ws/connect?id=<jwt|uid>` (usar `dev=1` ou header `X-Dev-Bypass:1` para pular validação), `GET /ws/connected-users` lista online, `POST /ws/message` envia por HTTP; mensagens WS seguem JSON `{ "type": "message", "room": "room-id", "content": "..." }` (também `type: join|leave`).

## Esquema e automação do banco
- DDL em `sql/init.sql` (conteúdo completo abaixo).
- O `docker-compose.yml` monta `./sql` em `/docker-entrypoint-initdb.d`; na primeira subida o Postgres cria o schema automaticamente.
- Se já houver volume `pgdata`, aplique manualmente: `psql "$DB_URL" -f sql/init.sql` ou `psql -h localhost -U $DB_USER -d $DB_NAME -f sql/init.sql`.

### DDL completo
```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    uid UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rooms (
    name TEXT NOT NULL,
    id TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS medias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL,
    message_id UUID NOT NULL,
    uid UUID NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS room_users (
    room_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    joined_at TIMESTAMP DEFAULT NOW(),
    left_at TIMESTAMP NULL,
    PRIMARY KEY (room_id, user_id),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);

CREATE INDEX IF NOT EXISTS idx_room_users_user ON room_users(user_id);
CREATE INDEX IF NOT EXISTS idx_room_users_room ON room_users(room_id);

CREATE TABLE IF NOT EXISTS messages (
   id TEXT PRIMARY KEY,
   sender_id TEXT NOT NULL,
   room_id TEXT NOT NULL,
   content TEXT NOT NULL,
   status VARCHAR(20) NOT NULL DEFAULT 'sent',
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   CONSTRAINT fk_room
       FOREIGN KEY(room_id)
       REFERENCES rooms(id)
       ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_devices (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
 user_id TEXT NOT NULL,
 fcm_token TEXT NOT NULL,
 created_at TIMESTAMP DEFAULT NOW()
);
```

## Testes e debug
- `go test ./...` (configure o banco, se necessário).
- `ws_test.html` oferece um cliente WebSocket simples para inspeção local.

## Deploy
- Dockerfile multi-stage pronto para provedores que aceitam containers. Consulte `DEPLOY.md` para opções como Render, Fly.io, Railway ou Cloud Run.
