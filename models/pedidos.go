package models

import "time"

type PedidoProduto struct {
	ProdutoID int `json:"produto_id"`
	Quantidade int `json:"quantidade"`
}

type Pedido struct {
	ID          int             `json:"id"`
	ClienteID   int             `json:"cliente_id"`
	Produtos    []PedidoProduto `json:"produtos"`
	ValorTotal  float64         `json:"valor_total"`
	DataCriacao time.Time       `json:"data_criacao"`
}
