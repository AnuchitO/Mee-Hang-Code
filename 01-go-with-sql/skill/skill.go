package skill

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

type handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) handler {
	return handler{db: db}
}

func findSkillByKey(db *sql.DB, key string) (Skill, error) {
	// query the database for the skill with the given key
	row := db.QueryRow("SELECT key, name, description, logo, levels, tags FROM skill WHERE key = $1", key)

	// scan data from row into a Skill struct
	var Key string
	var Name string
	var Description string
	var Logo string
	var Levels []byte
	var Tags pq.StringArray

	var skill Skill
	if err := row.Scan(&Key, &Name, &Description, &Logo, &Levels, &Tags); err != nil {
		return Skill{}, err
	}
	lvl := []Level{}
	if err := json.Unmarshal(Levels, &lvl); err != nil {
		return Skill{}, err
	}

	skill.Key = Key
	skill.Name = Name
	skill.Description = Description
	skill.Logo = Logo
	skill.Tags = Tags
	skill.Levels = lvl

	return skill, nil
}

func (h handler) GetSkillByKey(c *gin.Context) {
	// get the key from the URL path param
	key := c.Param("key")

	skill, err := findSkillByKey(h.db, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": skill})
}
