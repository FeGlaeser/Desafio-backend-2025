package handlers

import (
	"api-rest-vendas/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var DB *sql.DB

func SetDB(db *sql.DB) {
	DB = db
}

func CreateProduto(w http.ResponseWriter, r *http.Request) {
	var produto models.Produto
	json.NewDecoder(r.Body).Decode(&produto)

	err := DB.QueryRow("INSERT INTO produtos (nome, preco, estoque) VALUES ($1, $2, $3) RETURNING id",
		produto.Nome, produto.Preco, produto.Estoque).Scan(&produto.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(produto)
}

func ListProdutos(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, nome, preco, estoque FROM produtos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var produtos []models.Produto
	for rows.Next() {
		var p models.Produto
		rows.Scan(&p.ID, &p.Nome, &p.Preco, &p.Estoque)
		produtos = append(produtos, p)
	}

	json.NewEncoder(w).Encode(produtos)
}

func UpdateProduto(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var produto models.Produto
	json.NewDecoder(r.Body).Decode(&produto)

	_, err := DB.Exec("UPDATE produtos SET nome=$1, preco=$2, estoque=$3 WHERE id=$4",
		produto.Nome, produto.Preco, produto.Estoque, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	produto.ID, _ = strconv.Atoi(id)
	json.NewEncoder(w).Encode(produto)
}

func DeleteProduto(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := DB.Exec("DELETE FROM produtos WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
