package repository

import (
	"database/sql"
	"fmt"
	Model "myapp/internal/model"
	"strconv"

	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "4650"
	dbname   = "postgres"
	sslmode  = "disable"
)

var connectionString string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

var Connection *dbr.Connection

var Cart map[int]Model.UserCart //Корзина

func init() {
	Cart = make(map[int]Model.UserCart)
}

func OpenTable() error {
	var err error
	Connection, err = dbr.Open("postgres", connectionString, nil)
	if err != nil {
		return err
	}
	return nil
}

//Добавить товар в корзину
func AddToCart(p Model.Product) error {
	for Product_Id := range Cart {
		if Product_Id == p.Product_Id {
			return fmt.Errorf(`Товар - "%s" уже в корзине`, p.Product_Name)
		}
	}
	userCart := Model.UserCart{
		Product_Id:    p.Product_Id,
		Product_Image: p.Product_Image,
		Product_Name:  p.Product_Name,
		Product_Price: p.Product_Price,
		Product_Koll:  1,
	}
	Cart[p.Product_Id] = userCart
	return nil
}

//Увеличить товар в корзине
func IncrementKoll(id string) error {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("неверно введён параметр id: %v", err)
	}
	if p, ok := Cart[product_id]; ok {
		p.Product_Koll++

		Cart[product_id] = p
	}
	return nil
}

//Уменьшить товар в корзине
func DeincrementKoll(id string) error {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("неверно введён параметр id: %v", err)
	}

	if p, ok := Cart[product_id]; ok {
		if p.Product_Koll == 1 {
			DeleteFromCart(id)
			return nil
		}
		p.Product_Koll--

		Cart[product_id] = p
	}
	return nil
}

//Удалить товар из корзины
func DeleteFromCart(id string) error {
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("неверно введён параметр id: %v", err)
	}
	delete(Cart, product_id)
	return nil
}

//Вывести информацию о корзине пользователя
func ReturnCart() map[int]Model.UserCart {
	return Cart
}

//Получить данные о продукте
func ReadOne(id string) (*Model.Product, error) {

	sess := Connection.NewSession(nil)
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("неверно введён параметр id: %v", err)
	}
	rows, err := sess.Select("*").From("products").Where("product_id = ?", product_id).Rows()
	if err != nil {
		return nil, err
	}

	var p Model.Product
	for rows.Next() {
		err = rows.Scan(&p.Product_Id, &p.Product_Image, &p.Product_Name, &p.Product_Price)
		if err != nil {
			return nil, err
		}
	}
	product := Model.Product{}
	if product == p {
		return nil, sql.ErrNoRows
	}
	return &p, nil
}

//Получить данные о всех продуктах
func GetAllProducts() ([]Model.Product, error) {
	sess := Connection.NewSession(nil)
	personInfo := []Model.Product{}
	rows, err := sess.Select("*").From("products").OrderAsc("product_id").Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p Model.Product
		err := rows.Scan(&p.Product_Id, &p.Product_Image, &p.Product_Name, &p.Product_Price)
		if err != nil {
			return nil, err
		}
		personInfo = append(personInfo, p)
	}
	return personInfo, nil
}

//Получить данные о заказах
func GetOrders(user_id string) ([]Model.Order, error) {
	sess := Connection.NewSession(nil)
	OrderInfo := []Model.Order{}
	rows, err := sess.Select("order_id, product_name, product_koll, product_price,order_time,order_status, user_id, customer_name, customer_address,customer_phone,customer_email").From("orders").Where("user_id = ?", user_id).OrderAsc("order_id").Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var o Model.Order
		err := rows.Scan(&o.Order_Id, &o.Product_Id, &o.Product_Name, &o.Product_Koll, &o.Product_Price, &o.Order_time, &o.Order_status, &o.Customer_Name, &o.Customer_Address, &o.Customer_Phone, &o.Customer_Email)
		if err != nil {
			return nil, err
		}
		OrderInfo = append(OrderInfo, o)
	}
	return OrderInfo, nil
}

/*func CreateOrder(o Model.Order) error {
	sess := Connection.NewSession(nil)

	sqlStatement := `INSERT INTO person ("product_name, product_koll, product_price, order_time, order_status, user_id, customer_name, customer_address, customer_phone, customer_email") VALUES ($1, $2, $3, $4)`

	var id int
	err := sess.QueryRow( sqlStatement, o.Product_Name, o.Product_Koll,o.Product_Price, o.LastName).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
*/