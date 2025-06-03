package main

import (
	"api-rest-vendas/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Conecta com o banco de dados
	ConnectDB()
	handlers.SetDB(DB)

	// Criar o roteador
	r := mux.NewRouter().StrictSlash(true)

	// ROTAS PRODUTOS
	r.HandleFunc("/produtos", handlers.CreateProduto).Methods("POST")
	r.HandleFunc("/produtos", handlers.ListProdutos).Methods("GET")
	r.HandleFunc("/produtos/{id}", handlers.UpdateProduto).Methods("PUT")
	r.HandleFunc("/produtos/{id}", handlers.DeleteProduto).Methods("DELETE")

	// ROTAS CLIENTES
	r.HandleFunc("/clientes", handlers.CreateCliente).Methods("POST")
	r.HandleFunc("/clientes", handlers.ListClientes).Methods("GET")
	r.HandleFunc("/clientes/{id}", handlers.UpdateCliente).Methods("PUT")
	r.HandleFunc("/clientes/{id}", handlers.DeleteCliente).Methods("DELETE")

	// ROTAS PEDIDOS
	r.HandleFunc("/pedidos", handlers.CreatePedido).Methods("POST")
	r.HandleFunc("/pedidos", handlers.ListPedidos).Methods("GET")
	r.HandleFunc("/pedidos/{id}", handlers.UpdatePedido).Methods("PUT")
	r.HandleFunc("/pedidos/{id}", handlers.DeletePedido).Methods("DELETE")


	// Iniciar servidor
	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
