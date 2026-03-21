## Deploy da API

### 1) O que a Vercel suporta (importante)
- A Vercel **não executa containers Docker em produção**; o runtime deles é serverless/edge. citeturn0search0
- Para usar Vercel com esta API em Go (que abre WebSocket e depende de Postgres), as opções são:
  1. **Hospedar o contêiner em outra plataforma** (Render, Fly.io, Railway, Cloud Run, etc.) e, se precisar, expor um domínio Vercel apenas para o front‑end/proxy.
  2. Refatorar a API para o modelo serverless da Vercel (Go handlers em `api/`), o que exigiria reescrever o bootstrap do Gin e remover dependências de conexão longa.

O Dockerfile adicionado serve para rodar localmente e para provedores que aceitam containers.

### 2) Usando o Dockerfile em um host de containers
```bash
# Build
docker build -t mychat-api .

# Rodar (apontando para o Postgres existente)
docker run -p 8000:8000 \
  -e API_PORT=8000 \
  -e DATABASE_URL=postgres://user:pass@host:5432/dbname \
  -e GMAIL_EMAIL=... \
  -e GMAIL_APP_PASSWORD=... \
  -e SUPABASE_URL=... \
  -e SUPABASE_KEY=... \
  mychat-api
```

### 3) Se quiser insistir em Vercel
- Coloque o front-end na Vercel e use `rewrites` para o domínio público da API hospedada em outro provedor de containers.
- Caso queira migrar para serverless na Vercel, será preciso criar handlers em `api/*.go` usando o runtime `go` deles, eliminar WebSocket stateful e mover a conexão com Postgres para cada invocação (cold start).
