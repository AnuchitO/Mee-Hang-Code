package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) handler {
	return handler{db: db}
}

func (h handler) GetSkillByKey(c *gin.Context) {
	key := c.Param("key")
	row := h.db.QueryRow("SELECT key, name, description, logo, levels, tags FROM skill WHERE key = $1", key)

	var skill Skill
	var levels []byte
	var tags pq.StringArray
	if err := row.Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, &levels, &tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := json.Unmarshal(levels, &skill.Levels); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	skill.Tags = tags

	c.JSON(http.StatusOK, gin.H{"data": skill})
}