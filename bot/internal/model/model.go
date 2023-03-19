package Model

type Product struct {
	Product_id    int
	Product_image string
	Product_name  string
	Product_price int
}

type Order struct {
	Order_Id       int
	User           string
	Product_Id     int
	Product_Image  string
	Product_Name   string
	Product_Price  string
	Product_Koll   int
	Product_amount int
	Order_time     string
	Order_status   string
	Total_price    int
	//Данные покупателя
	Customer_Name    string
	Customer_Address string
	Customer_Email   string
	Customer_Phone   string
}

type UserCart struct {
	Product_Id           int
	Product_Image        string
	Product_Name         string
	Product_Manufacturer string
	Product_Category     string
	Product_Description  string
	Product_Price        string
	Product_Koll         int
	Product_amount       int
}
