package service

type Services struct {
	Auth    *AuthService
	Product *ProductService
	Cart    *CartService
	Order   *OrderService
}

func NewServices(
	auth *AuthService,
	product *ProductService,
	cart *CartService,
	order *OrderService,
) *Services {
	return &Services{
		Auth:    auth,
		Product: product,
		Cart:    cart,
		Order:   order,
	}
}
