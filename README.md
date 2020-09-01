# Go Shop Backend

Sample backend project made in Go and used in the android Project: [![Android Market App Clean Architecture](https://github.com/aiescola/Android-Market-App-CleanArchitecture)](https://github.com/aiescola/Android-Market-App-CleanArchitecture)

It's deployed in [![Heroku](https://go-shopify-api.herokuapp.com/)](https://go-shopify-api.herokuapp.com/), although it uses a free dyno, so the first time it will delay a bit or even timeout.


## Register/Login

POST methods require authentication, for it, you must first `/register` using the appropiate `username` and `password` in the body (urlencoded)

As for `/login` you can either use a POST with the encoded credentials or a GET to fill the credentials in an html form.

## Products

### GET Products
```sh
/api/products
```
```json
{
    "products": [
        {
            "code": "MUG",
            "name": "Mug",
            "price": 5.4
        },
        {
            "code": "TSHIRT",
            "name": "T-Shirt",
            "price": 12.5
        },
        {
            "code": "PEN",
            "name": "Pen",
            "price": 3.2
        }
    ]
}
```

### GET Product
```sh
/api/products/{Code}
```
```json
{
    "code": "MUG",
    "name": "Mug",
    "price": 5.4
}
```

### POST Product
```sh
/api/products/
```
Raw body:
```json
{ "code": "CODE", "name": "NAME", "price": 99.99 }

```

## Discounts

### GET Discounts
```sh
/api/discounts/
```
```json
{
    "discounts": [
        {
            "code": "PRD02",
            "type": "product",
            "name": "One Product discount (voucher)",
            "description": "Mugs at 9.95 for a limited time",
            "productCodes": [
                "MUG"
            ],
            "price": 9.95
        }
    ]
}

```

### GET Discounts
```sh
/api/discounts/{code}
```
```json

{
    "code": "PRD02",
    "type": "product",
    "name": "One Product discount (voucher)",
    "description": "Mugs at 9.95 for a limited time",
    "productCodes": [
        "MUG"
    ],
    "price": 9.95
}

```

### POST Discount

```sh
/api/discounts/
```

Raw body:
```json
{ "code": "TSHIRT", "name": "T-Shirt", "price": 15.90 }

```
