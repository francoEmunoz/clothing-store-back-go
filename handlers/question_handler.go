package handlers

import (
	"cs-go/db"
	"cs-go/models"
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetQuestions(rw http.ResponseWriter, r *http.Request) {

	var questions []models.Question
	if err := db.Database.Preload("User").Order("id DESC").Find(&questions).Error; err != nil {
		sendError(rw, http.StatusInternalServerError, "Database error")
		return
	}

	sendData(rw, questions, http.StatusOK)

}

func getQuestionsById(r *http.Request) (models.Question, *gorm.DB) {

	vars := mux.Vars(r)
	questionId, _ := strconv.Atoi(vars["id"])
	question := models.Question{}

	if err := db.Database.First(&question, questionId); err.Error != nil {
		return question, err
	} else {
		return question, nil
	}
}

func CreateQuestion(rw http.ResponseWriter, r *http.Request) {
	var validate = validator.New()

	// Decodificar la solicitud en una estructura Question
	var question models.Question
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&question); err != nil {
		sendError(rw, http.StatusUnprocessableEntity, "Decoding error")
		return
	}

	// Validar los datos de la pregunta
	if err := validate.StructPartial(question, "Content", "UserID", "ProductID"); err != nil {
		sendValidationError(rw, err)
		return
	}

	// Obtener el userID del contexto
	userID, ok := r.Context().Value(userContextKey).(uint)
	if !ok {
		sendError(rw, http.StatusUnauthorized, "Invalid token")
		return
	}

	if userID != question.UserID {
		sendError(rw, http.StatusForbidden, "You are not authorized to create this comment")
		return
	}

	// Asignar el userID desde el contexto a la pregunta
	question.UserID = userID

	// Crear la pregunta en la base de datos
	if err := db.Database.Create(&question).Error; err != nil {
		sendError(rw, http.StatusInternalServerError, "Error saving question")
		return
	}

	if err := db.Database.Preload("User").First(&question, question.ID).Error; err != nil {
		sendError(rw, http.StatusInternalServerError, "Error loading user")
		return
	}

	sendData(rw, question, http.StatusCreated)
}

func UpdateComment(rw http.ResponseWriter, r *http.Request) {

	comment, err := getQuestionsById(r)
	if err != nil {
		sendError(rw, http.StatusNotFound, "The comment was not found")
		return
	}

	// Obtener el userID del contexto
	userID, ok := r.Context().Value(userContextKey).(uint)
	if !ok {
		sendError(rw, http.StatusUnauthorized, "Invalid token")
		return
	}

	// Verificar que el userID del token coincide con el userID del comentario
	if userID != comment.UserID {
		sendError(rw, http.StatusForbidden, "You are not authorized to update this comment")
		return
	}

	var updatedComment models.Question
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedComment); err != nil {
		sendError(rw, http.StatusUnprocessableEntity, "Decoding error")
		return
	}

	db.Database.Model(&comment).Omit("created_at").Updates(updatedComment)

	// Preload the User after updating the comment
	if err := db.Database.Preload("User").First(&comment, comment.ID).Error; err != nil {
		sendError(rw, http.StatusInternalServerError, "Error loading user")
		return
	}

	sendData(rw, comment, http.StatusOK)
}

func DeleteQuestion(rw http.ResponseWriter, r *http.Request) {

	if question, err := getQuestionsById(r); err != nil {
		sendError(rw, http.StatusNotFound, "The question was not found")
	} else {
		userID, ok := r.Context().Value(userContextKey).(uint)

		if !ok {
			sendError(rw, http.StatusUnauthorized, "Invalid token")
			return
		}

		if userID != question.UserID {
			sendError(rw, http.StatusForbidden, "You are not authorized to delete this question")
			return
		}

		db.Database.Delete(&question)
		sendData(rw, question, http.StatusOK)
	}

}
