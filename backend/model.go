package main

type Pedidos struct {
	Id_Pedido    int     `json:"id_pedido"`
	Nombre       string  `json:"nombre"`
	Precio       float64 `json:"precio"`
	Imagen       string  `json:"imagen"`
	Variante     string  `json:"variante"`
	Cantidad     int     `json:"cantidad"`
	Id_Comprador int     `json:"id_comprador"`
}

type PedidosDTO struct {
	Id_Pedido    int     `json:"id_pedido"`
	Nombre       string  `json:"nombre"`
	Precio       float64 `json:"precio"`
	Imagen       string  `json:"imagen"`
	Variante     string  `json:"variante"`
	Cantidad     int     `json:"cantidad"`
	Comprador    string  `json:"comprador"`
	Id_Comprador int     `json:"id_comprador"`
}

type Comprador struct {
	Id_Comprador int    `json:"id_comprador"`
	Nombre       string `json:"nombre"`
}

type AsignarPedidoRequest struct {
	Id_Pedido    int `json:"id_pedido"`
	Id_Comprador int `json:"id_comprador"`
}
