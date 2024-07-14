// forum/controllers/user_controller.go

package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum-backend/models"
	"forum-backend/services"
	"forum-backend/utils"

	"github.com/gorilla/mux"
)

func validateEmailHandler(email string) bool {
	if email == "" {
		return false
	}

	if strings.Contains(email, "@") {
		fmt.Println("Geçerli bir e-mail adresi.")
		return true
	} else {
		fmt.Println("Geçersiz bir e-mail adresi.")
		return false
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if user.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}
	if !validateEmailHandler(user.Email) {
		fmt.Println("tt")
		http.Error(w, "Email Valid değil", http.StatusInternalServerError)
		return
	}

	err := services.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := services.LoginUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := services.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract the id from the URL path
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	fmt.Println("id : ", id)

	user, err := services.GetUser(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func GetUserFromToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	userId, _ := utils.GetUserId(token)

	user, err := services.GetUser(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := services.UpdateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get("id"))

	err := services.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetUserLikes(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get("id"))

	likes, err := services.GetUserLikes(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(likes)
}

func GetUserDislikes(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get("id"))

	likes, err := services.GetUserDislikes(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(likes)
}

func GetUserPosts(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get("id"))

	posts, err := services.GetUserPosts(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func GetUserComments(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get("id"))

	comments, err := services.GetUserComments(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

func MostPostedUser(w http.ResponseWriter, r *http.Request) {
	user, err := services.MostPostedUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func ProfileData(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, _ := strconv.Atoi(params.Get("id"))

	user, err := services.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	posts, err := services.GetUserPosts(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	comments, err := services.GetUserComments(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	likes, err := services.GetUserLikes(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	dislikes, err := services.GetUserDislikes(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"user": user, "posts": posts, "comments": comments, "likes": likes, "dislikes": dislikes})
}

func ThirdPartyLoginAndRegister(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var token string

	db := utils.GetDB()

	var userId int

	err := db.QueryRow("SELECT id FROM USERS WHERE email = ?", user.Email).Scan(&userId)
	if err == sql.ErrNoRows {
		// kullanıcı yok
		if user.Username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}
		if user.Email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		err := services.CreateThirdPartyUser(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err = services.GetThirdPartyUserToken(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	} else {
		// kullanıcı var
		var err error
		token, err = services.GetThirdPartyUserToken(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
