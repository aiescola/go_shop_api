package api

import (
	"shopify/api/discounts"
	"shopify/api/login"
	"shopify/api/products"
	"shopify/middleware"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Api struct {
	authMiddleWare     middleware.Middleware
	ProductController  products.ProductController
	LoginController    login.LoginController
	DiscountController discounts.DiscountController
}

// Initializes the api
func New(database *mongo.Database, cookieStore *sessions.CookieStore, logger *log.Logger) Api {
	templates := template.Must(template.ParseGlob("templates/*.html"))
	productDataSource := products.NewMongoDataSource(database)
	loginDataSource := login.NewMongoDataSource(database)
	discountDataSource := discounts.NewMongoDataSource(database)

	productController := products.NewController(productDataSource, logger)
	loginController := login.NewController(loginDataSource, cookieStore, templates, logger)
	discountController := discounts.NewController(discountDataSource, logger)

	return Api{
		middleware.NewAuthMiddleware(cookieStore, logger),
		productController,
		loginController,
		discountController,
	}
}

func (api Api) CreateRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/products", api.ProductController.GetProducts).Methods("GET")
	router.HandleFunc("/api/products/{code}", api.ProductController.GetProduct).Methods("GET")
	router.HandleFunc("/api/products", api.authMiddleWare.Intercept(api.ProductController.AddProduct)).Methods("POST")

	router.HandleFunc("/register", api.LoginController.Register).Methods("POST")
	router.HandleFunc("/login", api.LoginController.Login).Methods("POST")
	router.HandleFunc("/login", api.LoginController.LoginForm).Methods("GET")

	router.HandleFunc("/discounts", api.DiscountController.GetDiscounts).Methods("GET")
	router.HandleFunc("/discounts/{code}", api.DiscountController.GetDiscount).Methods("GET")
	router.HandleFunc("/discounts", api.authMiddleWare.Intercept(api.DiscountController.AddDiscount)).Methods("POST")

	return router
}
