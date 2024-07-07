package skill

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestRecord(t *testing.T) {
	t.Run("toSkill: should mapping record to Skill model", func(t *testing.T) {
		r := record{
			Key:         "key",
			Name:        "name",
			Description: "description",
			Logo:        "logo",
			Levels:      []byte(`[{"name":"name","description":"description"}]`),
			Tags:        pq.StringArray{"tag1", "tag2"},
		}
		lvl := []Level{{Name: "name", Descriptions: []string{"description"}}}

		got := r.toSkill(lvl)

		want := Skill{
			Key:         "key",
			Name:        "name",
			Description: "description",
			Logo:        "logo",
			Levels:      lvl,
			Tags:        pq.StringArray{"tag1", "tag2"},
		}
		assert.Equal(t, want, got)
	})

	t.Run("unmarshalLevels: should be able to unmarshal levels", func(t *testing.T) {
		r := record{
			Levels: []byte(`[{"name":"name","descriptions":["description"]}]`),
		}

		got, err := r.unmarshalLevels()

		want := []Level{{Name: "name", Descriptions: []string{"description"}}}
		assert.Equal(t, want, got)
		assert.Nil(t, err)
	})

	t.Run("unmarshalLevels: should return error when unmarshal failed", func(t *testing.T) {
		r := record{Levels: []byte(`invalid json`)}

		_, err := r.unmarshalLevels()

		assert.NotNil(t, err)
	})
}