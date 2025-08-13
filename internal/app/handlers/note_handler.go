package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

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

func (h *NoteHandler) GetNoteByID(c *gin.Context) {
	userID := c.GetInt("userID")
	noteID := c.Param("id")

	id, err := strconv.Atoi(noteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note ID"})
		return
	}

	note, err := h.noteRepo.GetNoteByID(c.Request.Context(), id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {

	userID := c.GetInt("userID")
	noteID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}
	var note models.Note

	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	note.ID = noteID
	note.UserID = userID

	if err := h.noteRepo.UpdateNote(c.Request.Context(), &note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error update note"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "note Update"})
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	userID := c.GetInt("userID")
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}
	if err := h.noteRepo.DeleteNote(c.Request.Context(), noteID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delite note"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "note delete"})
}
