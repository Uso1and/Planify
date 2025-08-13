package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"planify/internal/domain/models"
	"planify/internal/domain/repo"
)

type NoteHandler struct {
	noteRepo repo.NoteRepoInterface
}

func NewNoteHandler(noteHandler repo.NoteRepoInterface) *NoteHandler {
	return &NoteHandler{noteRepo: noteHandler}
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	userID := c.GetInt("userID")

	var note models.Note

	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internal request note"})
		return
	}

	note.UserID = userID

	if err := h.noteRepo.CreateNewNote(c.Request.Context(), &note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}
	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) GetNotes(c *gin.Context) {
	userID := c.GetInt("userID")

	notes, err := h.noteRepo.GetNotesByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to featch notes"})
		return
	}

	c.JSON(http.StatusOK, notes)
}
