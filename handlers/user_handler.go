package handlers

import (
	"cs-go/db"
	"cs-go/models"
	"encoding/json"
	"log"
	"time"

	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetUsers(rw http.ResponseWriter, r *http.Request) {

	users := models.Users{}
	db.Database.Find(&users)
	sendData(rw, users, http.StatusOK)

}

func getUserById(r *http.Request) (models.User, *gorm.DB) {

	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["id"])
	user := models.User{}

	if err := db.Database.First(&user, userId); err.Error != nil {
		return user, err
	} else {
		return user, nil
	}
}

func GetUser(rw http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(userContextKey).(uint)
	if !ok {
		sendError(rw, http.StatusUnauthorized, "Invalid token")
		return
	}

	if user, err := getUserById(r); err != nil {
		sendError(rw, http.StatusNotFound, "The user was not found")
	} else {
		if userID != user.ID {
			sendError(rw, http.StatusForbidden, "You are not authorized to get this user")
			return
		}
		sendData(rw, user, http.StatusOK)
	}

}

func SignUp(rw http.ResponseWriter, r *http.Request) {
	var validate = validator.New()

	var registrationData models.UserRegistration
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&registrationData); err != nil {
		sendError(rw, http.StatusUnprocessableEntity, "Decoding error")
		return
	}

	if err := validate.Struct(registrationData); err != nil {
		sendValidationError(rw, err)
		return
	}

	existingUser := models.User{}

	err := db.Database.Where("email = ?", registrationData.Email).First(&existingUser).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		sendError(rw, http.StatusInternalServerError, "Database error")
		return
	}

	if err == nil {
		sendError(rw, http.StatusConflict, "User already registered")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationData.PlainPassword), bcrypt.DefaultCost)
	if err != nil {
		sendError(rw, http.StatusInternalServerError, "Error hashing password")
		return
	}

	user := models.User{
		Name:     registrationData.Name,
		Lastname: registrationData.Lastname,
		Password: string(hashedPassword),
		Email:    registrationData.Email,
		Country:  registrationData.Country,
		Role:     registrationData.Role,
	}

	if err := db.Database.Create(&user).Error; err != nil {
		sendError(rw, http.StatusInternalServerError, "Error creating user")
		return
	}

	sendData(rw, user, http.StatusCreated)
}

func LogIn(rw http.ResponseWriter, r *http.Request) {
	var loginData models.UserLogin
	var validate = validator.New()

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		sendError(rw, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validar que los campos Email y PlainPassword no estén vacíos
	if loginData.Email == "" || loginData.PlainPassword == "" {
		sendError(rw, http.StatusBadRequest, "Email and password are required")
		return
	}

	if err := validate.Struct(loginData); err != nil {
		sendValidationError(rw, err)
		return
	}

	var storedUser models.User
	if err := db.Database.Where("email = ?", loginData.Email).First(&storedUser).Error; err != nil {
		sendError(rw, http.StatusUnauthorized, "Invalid email or passworddd")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(loginData.PlainPassword))
	if err != nil {
		sendError(rw, http.StatusUnauthorized, "Invalid email or passwordxd")
		return
	}

	storedUser.Logged = true
	if err := db.Database.Save(&storedUser).Error; err != nil {
		sendError(rw, http.StatusInternalServerError, "Error updating user status")
		return
	}

	userData := map[string]interface{}{
		"ID":       storedUser.ID,
		"name":     storedUser.Name,
		"lastname": storedUser.Lastname,
		"email":    storedUser.Email,
		"country":  storedUser.Country,
		"role":     storedUser.Role,
	}

	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   storedUser.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		sendError(rw, http.StatusInternalServerError, "Error creating token")
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	responseData := map[string]interface{}{
		"token": tokenString,
		"user":  userData,
	}

	sendData(rw, responseData, http.StatusOK)
}

func VerifyToken(rw http.ResponseWriter, r *http.Request) {
	// Obtener el ID de usuario desde el contexto
	userID, ok := r.Context().Value(userContextKey).(uint)
	if !ok {
		log.Println(userID)
		sendError(rw, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Obtener el usuario desde la base de datos
	var user models.User
	if err := db.Database.First(&user, userID).Error; err != nil {
		sendError(rw, http.StatusInternalServerError, "Error retrieving user")
		return
	}

	// Generar un nuevo token
	tokenExpirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		ExpiresAt: tokenExpirationTime.Unix(),
		Subject:   user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		sendError(rw, http.StatusInternalServerError, "Error creating token")
		return
	}

	// Responder con el usuario y el nuevo token
	responseData := map[string]interface{}{
		"success": true,
		"response": map[string]interface{}{
			"user":  user,
			"token": tokenString,
		},
		"message": "Welcome " + user.Name,
	}

	sendData(rw, responseData, http.StatusOK)
}

func LogOut(rw http.ResponseWriter, r *http.Request) {

	user, err := getUserById(r)
	if err != nil {
		sendError(rw, http.StatusNotFound, "The user was not found")
		return
	}
	user.Logged = false
	if err := db.Database.Save(&user).Error; err != nil {
		sendError(rw, http.StatusInternalServerError, "Error updating user status")
		return
	}

	sendData(rw, user, http.StatusOK)

}

func UpdateUser(rw http.ResponseWriter, r *http.Request) {

	user, err := getUserById(r)
	if err != nil {
		sendError(rw, http.StatusNotFound, "The user was not found")
		return
	}

	var updatedUser models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedUser); err != nil {
		sendError(rw, http.StatusUnprocessableEntity, "Decoding error")
		return
	}

	userID, ok := r.Context().Value(userContextKey).(uint)
	if !ok {
		sendError(rw, http.StatusUnauthorized, "Invalid token")
		return
	}

	if userID != user.ID {
		sendError(rw, http.StatusForbidden, "You are not authorized to update this user")
		return
	}

	db.Database.Model(&user).Omit("created_at").Updates(updatedUser)

	sendData(rw, user, http.StatusOK)

}
