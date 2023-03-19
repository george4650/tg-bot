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

var Cart []Model.Product //Корзина 

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
	for _, product := range Cart {
		if product.Product_id == p.Product_id {
			return fmt.Errorf("Товар - %s уже в корзине!", product.Product_name)
		}
	}
	Cart = append(Cart, p)
	return nil
}

//Вывести информацию о корзине пользователя
func ReturnCart() []Model.Product {
	return Cart
}


//Получить данные о продукте
func ReadOne(id string) (*Model.Product, error) {

	sess := Connection.NewSession(nil)
	product_id, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("неверно введён параметр id: %v", err)
	}
	rows, err := sess.Select("*").From("product").Where("product_id = ?", product_id).Rows()
	if err != nil {
		return nil, err
	}

	var p Model.Product
	for rows.Next() {
		err = rows.Scan(&p.Product_id, &p.Product_image, &p.Product_name, &p.Product_price)
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
		err := rows.Scan(&p.Product_id, &p.Product_image, &p.Product_name, &p.Product_price)
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
	rows, err := sess.Select("order_id, product_name, product_koll, product_price,order_time,order_status, customer_name, customer_address,customer_phone,customer_email").From("orders").Where("user_name = ?", user_id).OrderAsc("order_id").Rows()
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
