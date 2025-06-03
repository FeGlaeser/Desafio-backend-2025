DROP DATABASE IF EXISTS sistema_vendas;

CREATE DATABASE sistema_vendas;

\c sistema_vendas;

CREATE TABLE produtos (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    preco NUMERIC(10, 2) NOT NULL CHECK (preco >= 0),
    estoque INT NOT NULL CHECK (estoque >= 0)
);

CREATE TABLE clientes (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    telefone VARCHAR(20)
);

CREATE TABLE pedidos (
    id SERIAL PRIMARY KEY,
    cliente_id INT NOT NULL REFERENCES clientes(id) ON DELETE CASCADE,
    valor_total NUMERIC(10, 2) DEFAULT 0 NOT NULL,
    data_criacao TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pedido_produtos (
    pedido_id INT NOT NULL REFERENCES pedidos(id) ON DELETE CASCADE,
    produto_id INT NOT NULL REFERENCES produtos(id) ON DELETE CASCADE,
    quantidade INT NOT NULL CHECK (quantidade > 0),
    PRIMARY KEY (pedido_id, produto_id)
);

CREATE OR REPLACE FUNCTION atualizar_valor_total()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE pedidos
    SET valor_total = (
        SELECT COALESCE(SUM(pp.quantidade * p.preco), 0)
        FROM pedido_produtos pp
        JOIN produtos p ON pp.produto_id = p.id
        WHERE pp.pedido_id = NEW.pedido_id
    )
    WHERE id = NEW.pedido_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER atualiza_valor_total_trigger
AFTER INSERT OR UPDATE OR DELETE ON pedido_produtos
FOR EACH ROW
EXECUTE FUNCTION atualizar_valor_total();

CREATE OR REPLACE FUNCTION atualizar_estoque()
RETURNS TRIGGER AS $$
DECLARE
    estoque_atual INT;
BEGIN
    IF TG_OP = 'INSERT' THEN
        SELECT estoque INTO estoque_atual FROM produtos WHERE id = NEW.produto_id;
        IF estoque_atual < NEW.quantidade THEN
            RAISE EXCEPTION 'Estoque insuficiente para o produto %', NEW.produto_id;
        END IF;
        UPDATE produtos SET estoque = estoque - NEW.quantidade WHERE id = NEW.produto_id;

    ELSIF TG_OP = 'DELETE' THEN
        UPDATE produtos SET estoque = estoque + OLD.quantidade WHERE id = OLD.produto_id;

    ELSIF TG_OP = 'UPDATE' THEN
        SELECT estoque INTO estoque_atual FROM produtos WHERE id = NEW.produto_id;
        IF estoque_atual + OLD.quantidade < NEW.quantidade THEN
            RAISE EXCEPTION 'Estoque insuficiente para atualizar o produto %', NEW.produto_id;
        END IF;
        UPDATE produtos
        SET estoque = estoque + OLD.quantidade - NEW.quantidade
        WHERE id = NEW.produto_id;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER atualiza_estoque_trigger
AFTER INSERT OR UPDATE OR DELETE ON pedido_produtos
FOR EACH ROW
EXECUTE FUNCTION atualizar_estoque();

