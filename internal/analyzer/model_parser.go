package analyzer

import (
	"fmt"
	"regexp"
	"strings"
)

type DatabaseModel struct {
	TableName  string
	Fields     []ModelField
	PrimaryKey string
}

type ModelField struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

func ParseTableSchema(ddl string) (*DatabaseModel, error) {
	model := &DatabaseModel{}

	tableRegex := regexp.MustCompile(`CREATE\s+TABLE\s+` + "`" + `?(\w+)` + "`" + `?`)
	tableMatch := tableRegex.FindStringSubmatch(ddl)
	if len(tableMatch) > 1 {
		model.TableName = tableMatch[1]
	} else {
		return nil, fmt.Errorf("no table name found in DDL")
	}

	fieldRegex := regexp.MustCompile("`" + `(\w+)` + "`" + `\s+(\w+(?:\(\d+\))?)\s*([^,\n]*)?`)
	fieldMatches := fieldRegex.FindAllStringSubmatch(ddl, -1)

	for _, match := range fieldMatches {
		if len(match) > 2 {
			field := ModelField{
				Name: match[1],
				Type: match[2],
			}
			if len(match) > 3 {
				field.Comment = strings.TrimSpace(match[3])
				if strings.Contains(strings.ToUpper(field.Comment), "PRIMARY KEY") {
					model.PrimaryKey = field.Name
				}
			}
			model.Fields = append(model.Fields, field)
		}
	}

	return model, nil
}

func (m *DatabaseModel) GetFieldNames() []string {
	names := make([]string, len(m.Fields))
	for i, field := range m.Fields {
		names[i] = field.Name
	}
	return names
}
