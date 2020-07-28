package login

import (
	"encoding/json"
	"net/http"
	"shopify/util"
	"text/template"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"

	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	dataSource  LoginDataSource
	cookieStore *sessions.CookieStore
	templates   *template.Template
	logger      *log.Entry
}

func NewController(registerDataSource LoginDataSource,
	cookieStore *sessions.CookieStore,
	templates *template.Template,
	logger *log.Logger) LoginController {
	return LoginController{
		registerDataSource,
		cookieStore,
		templates,
		logger.WithFields(log.Fields{
			"file": "LoginController",
		}),
	}
}

func (r LoginController) Register(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	user := request.URL.Query()["user"][0]
	password := request.URL.Query()["password"][0]

	encodedPassword, _ := bcrypt.GenerateFromPassword([]byte(string(password)), bcrypt.DefaultCost)

	credentials := Credentials{string(user), string(encodedPassword)}

	if err := r.dataSource.createUser(credentials); err != nil {
		r.logger.Error(err.Error())

		util.EncodeError(response, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(response).Encode("User created")
}

func (lc LoginController) LoginForm(w http.ResponseWriter, r *http.Request) {
	lc.templates.ExecuteTemplate(w, "login.html", nil)
}

func (r LoginController) Login(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	request.ParseForm()
	username := request.PostForm.Get("username")
	password := request.PostForm.Get("password")

	storedApiKey, err := r.dataSource.getPassword(username)

	if err != nil {
		r.logger.Error(err.Error())
		util.EncodeError(response, http.StatusNotFound, err.Error())
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*storedApiKey), []byte(password)); err != nil {
		r.logger.Error("Invalid password")
		util.EncodeError(response, http.StatusForbidden, "Invalid password")
		return
	}

	session, err := r.cookieStore.Get(request, "session")
	if err != nil {
		r.logger.Error("Error getting session ", err)
		util.EncodeError(response, http.StatusInternalServerError, "Session couldn't be decoded")
		return
	}

	session.Values["username"] = username
	if err = session.Save(request, response); err != nil {
		r.logger.Error("Error saving session: ", err)
	}

	r.logger.Info("User ", username, "logged in")
	json.NewEncoder(response).Encode("Login successful")
}
