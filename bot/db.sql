CREATE TABLE products (
	product_id serial PRIMARY KEY,
	product_image character varying NOT NULL,
	product_name character varying NOT NULL,
	product_price integer NOT NULL
);

CREATE TABLE order_status (
	order_status_id serial PRIMARY KEY,
	order_status_name character varying NOT NULL UNIQUE
);

CREATE TABLE orders (
	order_id serial PRIMARY KEY,
	user_id character varying NOT NULL,
	product_id integer NOT NULL,
	product_name  character varying NOT NULL,
	product_koll integer NOT NULL,
	product_price integer NOT NULL,
	order_time timestamp NOT NULL,
	order_status character varying NOT NULL,
	customer_name character varying NOT NULL,
	customer_address character varying NOT NULL,
	customer_email character varying NOT NULL,
	customer_phone character varying NOT NULL,
	customer_comment character varying,
	FOREIGN KEY (order_status) REFERENCES order_status (order_status_name),
	FOREIGN KEY (product_id) REFERENCES products (product_id)
);

INSERT INTO order_status (
	order_status_name
)
VALUES (
	'Ожидает подтверждения'
);

INSERT INTO order_status (
	order_status_name
)
VALUES (
	'Принято в работу'
);

INSERT INTO order_status (
	order_status_name
)
VALUES (
	'Заказ доставлен'
);

INSERT INTO order_status (
	order_status_name
)
VALUES (
	'Отказ'
);


INSERT INTO products (
	product_image, product_name, product_price)
VALUES (
	'images/sushi/meksika.jpeg', 'Суши Мексика', '1000'
);

INSERT INTO products (
	product_image, product_name, product_price)
VALUES (
	'images/sushi/midori.jpg', 'Суши Мидори', '950'
);

INSERT INTO products (
	product_image, product_name, product_price)
VALUES (
	'images/sushi/losos.jpg', 'Суши Лосось', '1200'
);