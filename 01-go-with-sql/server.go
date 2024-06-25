package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

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

type SkillStore struct {
	Key         string
	Name        string
	Description string
	Logo        string
	Levels      []byte
	Tags        pq.StringArray
}

func findSkillByKey(db *sql.DB, key string) (Skill, error) {
	row := db.QueryRow("SELECT key, name, description, logo, levels, tags FROM skills WHERE key = $1", key)

	s := SkillStore{}
	if err := row.Scan(&s.Key, &s.Name, &s.Description, &s.Logo, &s.Levels, &s.Tags); err != nil {
		return Skill{}, err
	}
	var levels []Level
	if err := json.Unmarshal(s.Levels, &levels); err != nil {
		return Skill{}, err
	}

	skill := Skill{
		Key:         s.Key,
		Name:        s.Name,
		Description: s.Description,
		Logo:        s.Logo,
		Levels:      levels,
		Tags:        s.Tags,
	}

	return skill, nil
}

func GetSkillByKey(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")

		skill, err := findSkillByKey(db, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": skill})
	}
}

func main() {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	uri := os.Getenv("POSTGRES_URI")
	db, err := sql.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r.GET("/skills/:key", GetSkillByKey(db))

	r.GET("/skills", func(c *gin.Context) {
		rows, err := db.Query("SELECT key, name, description, logo, levels, tags FROM skill")
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
	r.Run("127.0.0.1:8080")
}
