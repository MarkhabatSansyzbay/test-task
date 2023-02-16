package web

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"
	"regexp"

	"service2/models"

	"github.com/go-chi/chi"
)

const saltURL = "http://localhost:8082/generate-salt"

var errInvalidEmail = errors.New("invalid email")

type Handler struct {
	RpcClient *rpc.Client
}

func (h *Handler) InitRoutes(mux *chi.Mux) {
	mux.Post("/create-user", h.createUser)
	mux.Get("/get-user/{email}", h.getUser)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", saltURL, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	type response struct {
		Salt string `json:"salt"`
	}

	var generatedSalt *response
	if err := json.NewDecoder(res.Body).Decode(&generatedSalt); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	user.Salt = generatedSalt.Salt
	user.Password = hashPassword(user.Password, generatedSalt.Salt)

	if err := h.checkUser(user.Email); err != nil {
		if errors.Is(err, errInvalidEmail) {
			http.Error(w, errInvalidEmail.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var reply string
	if err = h.RpcClient.Call("App.CreateUser", user, &reply); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	var user *models.User

	if err := h.RpcClient.Call("App.GetUser", email, &user); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if user.Email == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Handler) checkUser(email string) error {
	var user *models.User
	if err := h.RpcClient.Call("App.GetUser", email, &user); err != nil {
		return err
	}

	if user.Email != "" ||
		!regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(email) {
		return errInvalidEmail
	}

	return nil
}

func hashPassword(password, salt string) string {
	hasher := md5.New()
	hasher.Write([]byte(salt + password))
	return hex.EncodeToString(hasher.Sum(nil))
}
