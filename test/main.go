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

/*
	package main

	import (
	"errors"
	"fmt"
	"time"
	)

	type Vehicle struct {
	ID        string
	Type      string // "Motorcycle", "Car", "Truck"
	EntryTime time.Time
	}

	type ParkingSpot struct {
	ID       string
	Type     string // "Motorcycle", "Car", "Truck"
	IsFree   bool
	Assigned *Vehicle
	}

	type ParkingFloor struct {
	Spots map[string]*ParkingSpot
	}

	type ParkingLot struct {
	Floors  map[string]*ParkingFloor
	Revenue float32
	}

	type Payment struct{}

	type SpotManager interface {
	AssignSpot(vehicle Vehicle) (*ParkingSpot, error)
	FreeSpot(spotID string) error
	}

	type VehicleManager interface {
	EnterVehicle(vehicle Vehicle) error
	ExitVehicle(vehicleID string) (float32, error)
	}

	type PaymentProcessor interface {
	CalculateCharge(vehicle Vehicle, hoursParked float32) float32
	ProcessPayment(vehicleID string, amount float32) error
	}

	type ReportGenerator interface {
	GenerateCurrentParkingReport() string
	GenerateDailyRevenueReport() float32
	}

	func NewParkingLot() *ParkingLot {
	return &ParkingLot{
		Floors:  make(map[string]*ParkingFloor),
		Revenue: 0,
	}
	}

	// SpotManager methods
	func (pl *ParkingLot) AssignSpot(vehicle Vehicle) (*ParkingSpot, error) {
	for _, floor := range pl.Floors {
		for _, spot := range floor.Spots {
			if spot.IsFree && spot.Type == vehicle.Type {
				spot.IsFree = false
				spot.Assigned = &vehicle
				return spot, nil
			}
		}
	}
	return nil, errors.New("no available spot for vehicle type")
	}

	func (pl *ParkingLot) FreeSpot(spotID string) error {
	for _, floor := range pl.Floors {
		if spot, exists := floor.Spots[spotID]; exists {
			spot.IsFree = true
			spot.Assigned = nil
			return nil
		}
	}
	return errors.New("spot not found")
	}

	// PaymentProcessor methods
	func (p *Payment) CalculateCharge(vehicle Vehicle, hoursParked float32) float32 {
	rates := map[string]float32{
		"Motorcycle": 500,
		"Car":        1000,
		"Truck":      2000,
	}
	return rates[vehicle.Type] * hoursParked
	}

	func (p *Payment) ProcessPayment(vehicleID string, amount float32) error {
	fmt.Printf("Payment of %.2f processed for vehicle %s\n", amount, vehicleID)
	return nil
	}

	// ReportGenerator methods
	func (pl *ParkingLot) GenerateCurrentParkingReport() string {
	report := "Current Parking Status:\n"
	for floorID, floor := range pl.Floors {
		report += fmt.Sprintf("Floor %s:\n", floorID)
		for spotID, spot := range floor.Spots {
			status := "Free"
			if !spot.IsFree {
				status = fmt.Sprintf("Occupied by %s", spot.Assigned.ID)
			}
			report += fmt.Sprintf("  Spot %s: %s\n", spotID, status)
		}
	}
	return report
	}

	func (pl *ParkingLot) GenerateDailyRevenueReport() float32 {
	return pl.Revenue
 	}

 	func main() {

	parkingLot := NewParkingLot()
	parkingLot.Floors["Floor1"] = &ParkingFloor{
		Spots: map[string]*ParkingSpot{
			"Spot1": {ID: "Spot1", Type: "Car", IsFree: true},
			"Spot2": {ID: "Spot2", Type: "Motorcycle", IsFree: true},
			"Spot3": {ID: "Spot3", Type: "Truck", IsFree: true},
		},
	}
	paymentProcessor := &Payment{}

	car := Vehicle{ID: "CAR123", Type: "Car", EntryTime: time.Now()}
	spot, err := parkingLot.AssignSpot(car)
	if err != nil {
		fmt.Println("Error assigning spot:", err)
		return
	}
	fmt.Printf("Vehicle %s assigned to Spot %s\n", car.ID, spot.ID)

	time.Sleep(2 * time.Second)

	hoursParked := float32(2)
	charge := paymentProcessor.CalculateCharge(car, hoursParked)
	err = paymentProcessor.ProcessPayment(car.ID, charge)
	if err != nil {
		fmt.Println("Error processing payment:", err)
		return
	}
	fmt.Printf("Vehicle %s exited. Charge: %.2f\n", car.ID, charge)

	err = parkingLot.FreeSpot(spot.ID)
	if err != nil {
		fmt.Println("Error freeing spot:", err)
		return
	}
	fmt.Printf("Spot %s is now free\n", spot.ID)

	fmt.Println(parkingLot.GenerateCurrentParkingReport())
	parkingLot.Revenue += charge
	fmt.Printf("Daily Revenue: %.2f\n", parkingLot.GenerateDailyRevenueReport())
	}
*/

package main

func diagonalDifference(arr [][]int32) int32 {
	var sum1 int32
	var sum2 int32

	if len(arr) >= 1 && len(arr) <= 100 {
		for i := 0; i < len(arr); i++ {
			sum1 += arr[i][i]

			sum2 += arr[i][len(arr[i])-1-i]
		}
	}

	n := sum1 - sum2
	if n < 0 {
		n = -n
	}

	return n
}

func main() {
	arr := [][]int32{
		{11, 2, 4},
		{4, 5, 6},
		{10, 8, -12},
	}
	m := diagonalDifference(arr)
	_ = m
}
