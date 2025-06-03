# 🛒 API RESTful - Sistema de Vendas

Esta API RESTful foi desenvolvida em Go para gerenciar produtos, clientes e pedidos de uma loja fictícia. Utiliza PostgreSQL como banco de dados e segue boas práticas de organização de código.

## 🚀 Funcionalidades

### 📦 Produtos
- Criar, listar, atualizar e remover produtos
- Cada produto tem: `id`, `nome`, `preço`, `estoque`

### 👤 Clientes
- Criar, listar, atualizar e remover clientes
- Cada cliente tem: `id`, `nome`, `telefone`

### 🧾 Pedidos
- Criar, listar, atualizar e remover pedidos
- Cada pedido contém:
  - `id`
  - `cliente_id`
  - `lista de produtos` (com `produto_id` e `quantidade`)
  - `valor_total` (calculado automaticamente)
  - `data_criacao` (gerado automaticamente)
- Atualização automática de estoque
- Recalcula o valor total ao atualizar produtos no pedido

---

## 🛠️ Tecnologias e Ferramentas

- Go
- PostgreSQL
- `gorilla/mux` para rotas
- SQL puro (sem ORM)
- Funções e triggers no banco para lógica de negócio

---

## 📁 Estrutura do Projeto

```
api-rest-vendas/
├── main.go
├── db.go
├── handlers/
│   ├── produto.go
│   ├── cliente.go
│   └── pedido.go
├── models/
│   ├── produto.go
│   ├── cliente.go
│   └── pedido.go
├── utils/
│   └── json.go
├── go.mod
├── go.sum
├── schema.sql
└── .env (opcional)
```

---

## ⚙️ Como Executar o Projeto

### ✅ Pré-requisitos

- Go 1.18 ou superior
- PostgreSQL instalado e em execução

### 📦 1. Clone o projeto

```bash
git clone https://github.com/FeGlaeser/api-rest-vendas.git
cd api-rest-vendas
```

### 🛠️ 2. Configure o Banco de Dados

Execute o script SQL para criar as tabelas e triggers:

```bash
psql -U seu_usuario -f schema.sql
```

> Altere o `seu_usuario` para o usuário do seu PostgreSQL.

### 🔐 3. Configure as variáveis de ambiente (opcional)

Crie um arquivo `.env`:

```env
DB_USER=seu_usuario
DB_PASSWORD=sua_senha
DB_NAME=sistema_vendas
DB_HOST=localhost
DB_PORT=5432
```

Ou configure diretamente no `db.go`.

### 🧰 4. Instale as dependências

```bash
go mod tidy
```

### ▶️ 5. Execute a aplicação

```bash
go run main.go db.go
```

Servidor disponível em: [http://localhost:8080](http://localhost:8080)

---

## 🔄 Exemplos de Requisições

### Criar Produto

```http
POST /produtos
Content-Type: application/json

{
  "nome": "Teclado Mecânico",
  "preco": 250.00,
  "estoque": 10
}
```

### Criar Pedido

```http
POST /pedidos
Content-Type: application/json

{
  "cliente_id": 1,
  "produtos": [
    { "produto_id": 2, "quantidade": 1 },
    { "produto_id": 3, "quantidade": 2 }
  ]
}
```

---

## ✅ Extras Implementados

- [x] Atualização de valor total do pedido com trigger
- [x] Controle automático de estoque
- [x] Relacionamentos entre tabelas
- [x] Filtro por cliente nos pedidos
- [x] Paginação
- [x] Documentação Swagger (opcional)

---

## 📄 Paginação

A listagem de pedidos suporta paginação através dos seguintes parâmetros de query:

### 🔹 Parâmetros

| Parâmetro   | Tipo   | Descrição                               |
|-------------|--------|-----------------------------------------|
| `limit`     | int    | Quantidade máxima de itens por página.  |
| `offset`    | int    | Quantidade de itens a serem ignorados (deslocamento). |
| `cliente_id`| int    | (Opcional) Filtra pedidos de um cliente específico. |

### 🔹 Exemplo de requisição
GET /pedidos?limit=3&offset=6
### 🔹 Regras

- `limit` padrão: `10` (se não informado).
- `limit` máximo recomendado: `100`.
- `offset` mínimo: `0`.
- Se `limit` ou `offset` forem inválidos (ex.: negativos), a API retornará um erro `400 Bad Request`.
- Quando não houver resultados, será retornado: `[]`.

### 🔹 Exemplo de resposta

```json
[
  {
    "id": 7,
    "cliente_id": 2,
    "produtos": [{"produto_id": 1, "quantidade": 1}],
    "valor_total": 10.5,
    "data_criacao": "2025-06-02T21:47:30.742979Z"
  },
  {
    "id": 8,
    "cliente_id": 1,
    "produtos": [{"produto_id": 1, "quantidade": 1}],
    "valor_total": 10.5,
    "data_criacao": "2025-06-02T21:47:30.769867Z"
  }
] 
```

## 📅 Prazo

Entrega até **segunda-feira, 2 de junho, às 23:59**.

---

## 👨‍💻 Autor

Fernando Glaeser da Silva  
https://github.com/FeGlaeser

---
