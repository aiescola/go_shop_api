# Go Shop Backend

Sample backend project made in Go and used in the android Project: [![Android Market App Clean Architecture](https://github.com/aiescola/Android-Market-App-CleanArchitecture)](https://github.com/aiescola/Android-Market-App-CleanArchitecture)

It's deployed in [![Heroku](https://go-shopify-api.herokuapp.com/)](https://go-shopify-api.herokuapp.com/)

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