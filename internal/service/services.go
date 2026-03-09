package service

type Services struct {
	Auth     *AuthService
	Product  *ProductService
	Cart     *CartService
	Order    *OrderService
	Category *CategoryService
}

func NewServices(
	auth *AuthService,
	product *ProductService,
	cart *CartService,
	order *OrderService,
	category *CategoryService,
) *Services {
	return &Services{
		Auth:     auth,
		Product:  product,
		Cart:     cart,
		Order:    order,
		Category: category,
	}
}
