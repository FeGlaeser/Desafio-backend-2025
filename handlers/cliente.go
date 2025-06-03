package handlers

import (
	"api-rest-vendas/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Criar Cliente
func CreateCliente(w http.ResponseWriter, r *http.Request) {
	var cliente models.Cliente
	json.NewDecoder(r.Body).Decode(&cliente)

	err := DB.QueryRow("INSERT INTO clientes (nome, telefone) VALUES ($1, $2) RETURNING id",
		cliente.Nome, cliente.Telefone).Scan(&cliente.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cliente)
}

// Listar Clientes
func ListClientes(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, nome, telefone FROM clientes")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clientes []models.Cliente
	for rows.Next() {
		var c models.Cliente
		rows.Scan(&c.ID, &c.Nome, &c.Telefone)
		clientes = append(clientes, c)
	}

	json.NewEncoder(w).Encode(clientes)
}

// Atualizar Cliente
func UpdateCliente(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var cliente models.Cliente
	json.NewDecoder(r.Body).Decode(&cliente)

	_, err := DB.Exec("UPDATE clientes SET nome=$1, telefone=$2 WHERE id=$3",
		cliente.Nome, cliente.Telefone, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cliente.ID, _ = strconv.Atoi(id)
	json.NewEncoder(w).Encode(cliente)
}

// Deletar Cliente
func DeleteCliente(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := DB.Exec("DELETE FROM clientes WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
