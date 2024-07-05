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

type record struct {
	Key         string
	Name        string
	Description string
	Logo        string
	Levels      []byte
	Tags        pq.StringArray
}

func toSkill(r record, lvl []Level) Skill {
	return Skill{
		Key:         r.Key,
		Name:        r.Name,
		Description: r.Description,
		Logo:        r.Logo,
		Tags:        r.Tags,
		Levels:      lvl,
	}
}

func unmarshalLevels() {

}

func findSkillByKey(db *sql.DB, key string) (Skill, error) {
	// query the database for the skill with the given key
	row := db.QueryRow("SELECT key, name, description, logo, levels, tags FROM skill WHERE key = $1", key)

	// scan data from row into a Skill struct
	r := record{}
	if err := row.Scan(&r.Key, &r.Name, &r.Description, &r.Logo, &r.Levels, &r.Tags); err != nil {
		return Skill{}, err
	}
	lvl := []Level{}
	if err := json.Unmarshal(r.Levels, &lvl); err != nil {
		return Skill{}, err
	}

	s := toSkill(r, lvl)

	return s, nil
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
