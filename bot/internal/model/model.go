package Model

import "time"

type Product struct {
	Product_Id    int
	Product_Image string
	Product_Name  string
	Product_Price int
}

type Order struct {
	Order_Id       int
	User_id        int
	Product_Id     int
	Product_Image  string
	Product_Name   string
	Product_Price  int
	Product_Koll   int
	Order_time     time.Time
	Order_status   string
	//Данные покупателя
	Customer_Name    string
	Customer_Address string
	Customer_Email   string
	Customer_Phone   string
}

type UserCart struct {
	Product_Id    int
	Product_Image string
	Product_Name  string
	Product_Price int
	Product_Koll  int
}
