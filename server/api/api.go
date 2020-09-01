package api

import (
	"text/template"

	"github.com/aiescola/go_shop_api/api/discounts"
	"github.com/aiescola/go_shop_api/api/login"
	"github.com/aiescola/go_shop_api/api/products"
	"github.com/aiescola/go_shop_api/middleware"
	"github.com/go-redis/redis/v8"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Api struct {
	authMiddleWare     middleware.Middleware
	productController  products.ProductController
	loginController    login.LoginController
	discountController discounts.DiscountController
}

// Initializes the api
func New(database *mongo.Database, redis *redis.Client, cookieStore *sessions.CookieStore, logger *log.Logger) Api {
	templates := template.Must(template.ParseGlob("templates/*.html"))
	productDataSource := products.NewMongoDataSource(database)
	loginDataSource := login.NewMongoDataSource(database)
	discountDataSource := discounts.NewMongoDataSource(database)

	productsCache := products.NewProductsRedisCache(redis)
	discountsCache := discounts.NewDiscountsRedisCache(redis)

	productController := products.NewController(productDataSource, productsCache, logger)
	loginController := login.NewController(loginDataSource, cookieStore, templates, logger)
	discountController := discounts.NewController(discountDataSource, discountsCache, logger)

	return Api{
		middleware.NewAuthMiddleware(cookieStore, logger),
		productController,
		loginController,
		discountController,
	}
}

func (api Api) CreateRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/register", api.loginController.Register).Methods("POST")
	router.HandleFunc("/login", api.loginController.Login).Methods("POST")
	router.HandleFunc("/login", api.loginController.LoginForm).Methods("GET")

	router.HandleFunc("/api/products", api.productController.GetProducts).Methods("GET")
	router.HandleFunc("/api/products/{code}", api.productController.GetProduct).Methods("GET")
	router.HandleFunc("/api/products", api.authMiddleWare.Intercept(api.productController.AddProduct)).Methods("POST")

	router.HandleFunc("/api/discounts", api.discountController.GetDiscounts).Methods("GET")
	router.HandleFunc("/api/discounts/{code}", api.discountController.GetDiscount).Methods("GET")
	router.HandleFunc("/api/discounts", api.authMiddleWare.Intercept(api.discountController.AddDiscount)).Methods("POST")

	return router
}
