package handlers

import (
	"database/sql"
	"api-rest-vendas/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreatePedido(w http.ResponseWriter, r *http.Request) {
	var pedido models.Pedido
	json.NewDecoder(r.Body).Decode(&pedido)

	err := DB.QueryRow("INSERT INTO pedidos (cliente_id) VALUES ($1) RETURNING id, data_criacao, valor_total",
		pedido.ClienteID).Scan(&pedido.ID, &pedido.DataCriacao, &pedido.ValorTotal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, item := range pedido.Produtos {
		_, err := DB.Exec("INSERT INTO pedido_produtos (pedido_id, produto_id, quantidade) VALUES ($1, $2, $3)",
			pedido.ID, item.ProdutoID, item.Quantidade)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = DB.QueryRow("SELECT valor_total FROM pedidos WHERE id = $1", pedido.ID).Scan(&pedido.ValorTotal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(pedido)
}

func ListPedidos(w http.ResponseWriter, r *http.Request) {
	clienteIDStr := r.URL.Query().Get("cliente_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	var rows *sql.Rows
	var err error

	if clienteIDStr != "" {
		clienteID, err := strconv.Atoi(clienteIDStr)
		if err != nil {
			http.Error(w, "cliente_id inválido", http.StatusBadRequest)
			return
		}
		rows, err = DB.Query(`
			SELECT id, cliente_id, valor_total, data_criacao 
			FROM pedidos 
			WHERE cliente_id = $1
			ORDER BY id
			LIMIT $2 OFFSET $3
		`, clienteID, limit, offset)
	} else {
		rows, err = DB.Query(`
			SELECT id, cliente_id, valor_total, data_criacao 
			FROM pedidos 
			ORDER BY id
			LIMIT $1 OFFSET $2
		`, limit, offset)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pedidos []models.Pedido
	for rows.Next() {
		var p models.Pedido
		if err := rows.Scan(&p.ID, &p.ClienteID, &p.ValorTotal, &p.DataCriacao); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		prodRows, err := DB.Query("SELECT produto_id, quantidade FROM pedido_produtos WHERE pedido_id = $1", p.ID)
		if err == nil {
			for prodRows.Next() {
				var pp models.PedidoProduto
				if err := prodRows.Scan(&pp.ProdutoID, &pp.Quantidade); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				p.Produtos = append(p.Produtos, pp)
			}
			prodRows.Close()
		}

		pedidos = append(pedidos, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidos)
}

func UpdatePedido(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var pedido models.Pedido
	json.NewDecoder(r.Body).Decode(&pedido)

	for _, item := range pedido.Produtos {
		_, err := DB.Exec(`
			INSERT INTO pedido_produtos (pedido_id, produto_id, quantidade)
			VALUES ($1, $2, $3)
			ON CONFLICT (pedido_id, produto_id) DO UPDATE SET quantidade = EXCLUDED.quantidade
		`, id, item.ProdutoID, item.Quantidade)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err := DB.QueryRow("SELECT valor_total FROM pedidos WHERE id = $1", id).Scan(&pedido.ValorTotal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pedido.ID, _ = strconv.Atoi(id)
	json.NewEncoder(w).Encode(pedido)
}

func DeletePedido(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := DB.Exec("DELETE FROM pedidos WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) 
}
