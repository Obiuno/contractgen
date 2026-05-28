package normalise

import (
	"contractgen/internal/parser"
	"strings"
)

func normaliseName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.Join(strings.Fields(s), "_")
}

func Normalise(c *parser.Contract) {
	for i := range c.Tables {
		c.Tables[i].Name = normaliseName(c.Tables[i].Name)
		for j := range c.Tables[i].Columns {
			c.Tables[i].Columns[j].Name = normaliseName(c.Tables[i].Columns[j].Name)
			if ref := c.Tables[i].Columns[j].References; ref != nil {
				ref.Table = normaliseName(ref.Table)
				ref.Column = normaliseName(ref.Column)
			}
		}
		for j, pk := range c.Tables[i].PrimaryKey {
			c.Tables[i].PrimaryKey[j] = normaliseName(pk)
		}
		for j, group := range c.Tables[i].Unique {
			for k, col := range group {
				c.Tables[i].Unique[j][k] = normaliseName(col)
			}
		}
	}
}
