package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Level struct {
	Key          string   `json:"key"`
	Name         string   `json:"name"`
	Brief        string   `json:"brief"`
	Descriptions []string `json:"descriptions"`
	Level        int      `json:"level"`
}

type Skill struct {
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Logo        string   `json:"logo"`
	Levels      []Level  `json:"levels"`
	Tags        []string `json:"tags"`
}

func main() {
	r := gin.Default()
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r.GET("/skills/:key", func(c *gin.Context) {
		key := c.Param("key")
		row := db.QueryRow("SELECT key, name, description, logo, levels, tags FROM skills WHERE key = $1", key)

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
	})

	r.GET("/skills", func(c *gin.Context) {
		rows, err := db.Query("SELECT key, name, description, logo, levels, tags FROM skills")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		skills := []Skill{}
		for rows.Next() {
			var skill Skill
			var levels []byte
			var tags pq.StringArray
			if err := rows.Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, &levels, &tags); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if err := json.Unmarshal(levels, &skill.Levels); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			skill.Tags = tags
			skills = append(skills, skill)
		}

		c.JSON(http.StatusOK, gin.H{"data": skills})
	})

	r.POST("/skills", func(c *gin.Context) {
		var skill Skill
		if err := c.ShouldBindJSON(&skill); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stmt, err := db.Prepare("INSERT INTO skills (key, name, description, logo, levels, tags) VALUES ($1, $2, $3, $4, $5, $6)")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		levels, err := json.Marshal(skill.Levels)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tags := pq.StringArray(skill.Tags)
		_, err = stmt.Exec(skill.Key, skill.Name, skill.Description, skill.Logo, levels, tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": skill})
	})
	r.Run(":8080")
}
