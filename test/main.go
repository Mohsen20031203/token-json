/*
	package main

	import (
	"fmt"
	)

	type Product struct {
	Name     string
	Price    float32
	Number   uint
	Category string
	}

	type Customer struct {
	Name     string
	Email    string
	ListShop []Product
	}

	type Shop struct {
	Users    map[string]*Customer
	Products map[string]*Product
	}

	type Cliant interface {
	AddProduct(name, category string, price float32, number uint) error
	ListPrintProduct()
	Shopping(nameUser, nameProduct string, number uint) error
	ListShopUser(nameUser string) error
	AddUser(name, email string) error
	}

	func newShopRoshd() Cliant {
	return &Shop{}
	}
	func newShopKetab() Cliant {
	return &Shop{}
	}

	func (s *Shop) AddProduct(name, category string, price float32, number uint) error {
	if _, exists := s.Products[name]; exists {
		return fmt.Errorf("product %s already exists", name)
	}
	p := &Product{Name: name, Price: price, Number: number, Category: category}
	s.Products[name] = p
	return nil
	}

	func (s *Shop) ListPrintProduct() {
	for _, p := range s.Products {
		if p.Number > 0 {
			fmt.Println(*p)
		}
	}
	}

	func (s *Shop) Shopping(nameUser, nameProduct string, number uint) error {
	customer, exists := s.Users[nameUser]
	if !exists {
		return fmt.Errorf("user %s not found", nameUser)
	}

	product, exists := s.Products[nameProduct]
	if !exists {
		return fmt.Errorf("product %s not found", nameProduct)
	}

	if product.Number < number {
		return fmt.Errorf("not enough stock for product %s", nameProduct)
	}

	// Add product to customer's shopping list
	customer.ListShop = append(customer.ListShop, *product)

	// Decrease the product's number in stock
	product.Number -= number
	return nil
	}

	func (s *Shop) ListShopUser(nameUser string) error {
	customer, exists := s.Users[nameUser]
	if !exists {
		return fmt.Errorf("user %s not found", nameUser)
	}
	fmt.Println("Shopping list for", nameUser)
	for _, product := range customer.ListShop {
		fmt.Println(product)
	}
	return nil
	}

	func (s *Shop) AddUser(name, email string) error {
	if _, exists := s.Users[name]; exists {
		return fmt.Errorf("user %s already exists", name)
	}
	c := &Customer{Name: name, Email: email}
	s.Users[name] = c
	return nil
	}

	func main() {
	var varShop Cliant

	varShop = newShopKetab()

	varShop.ListPrintProduct()

	}
*/

package main
