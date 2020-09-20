package models

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

// Unmarshal parses the content of models.yml to a datastruct.q
func Unmarshal(r io.Reader) (map[string]Model, error) {
	var m map[string]Model
	if err := yaml.NewDecoder(r).Decode(&m); err != nil {
		return nil, fmt.Errorf("decoding models: %w", err)
	}
	return m, nil
}

// Model replresents one model from models.yml.
type Model struct {
	Attributes map[string]Attribute
}

// UnmarshalYAML decodes a yaml model to models.Model.
func (m *Model) UnmarshalYAML(node *yaml.Node) error {
	return node.Decode(&m.Attributes)
}

// Relation represents some kind of relation between fields.
type Relation interface {
	toCollection() []string
	toField() ToField
}

// Attribute is a field of a model.
type Attribute struct {
	Type     string
	relation Relation
	template *AttributeTemplate
}

// Relation returns the relation object if the Attribute is a relation or a
// template with a relation. In other case, it returns nil.
func (a *Attribute) Relation() Relation {
	if a.relation != nil {
		return a.relation
	}

	if a.template != nil && a.template.Fields.relation != nil {
		return a.template.Fields.relation
	}
	return nil
}

// AttributeRelation is a relation or relation-list field.
type AttributeRelation struct {
	To To `yaml:"to"`
}

func (r AttributeRelation) toCollection() []string {
	return []string{r.To.Collection}
}

func (r AttributeRelation) toField() ToField {
	return r.To.Field
}

// AttributeGenericRelation is a generic-relation or generic-relation-list fiedl.
type AttributeGenericRelation struct {
	To ToGeneric `yaml:"to"`
}

func (r AttributeGenericRelation) toCollection() []string {
	return r.To.Collection
}

func (r AttributeGenericRelation) toField() ToField {
	return r.To.Field
}

// AttributeTemplate represents a template field.
type AttributeTemplate struct {
	Replacement string    `yaml:"replacement"`
	Fields      Attribute `yaml:"fields"`
}

// UnmarshalYAML decodes a model attribute from yaml.
func (a *Attribute) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err == nil {
		a.Type = s
		return nil
	}

	var typer struct {
		Type string `yaml:"type"`
	}
	if err := value.Decode(&typer); err != nil {
		return fmt.Errorf("field object without type: %w", err)
	}

	a.Type = typer.Type
	switch typer.Type {
	case "relation":
		fallthrough
	case "relation-list":
		var relation AttributeRelation
		if err := value.Decode(&relation); err != nil {
			return fmt.Errorf("invalid object of type %s at line %d object: %w", typer.Type, value.Line, err)
		}
		a.relation = &relation
	case "generic-relation":
		fallthrough
	case "generic-relation-list":
		var relation AttributeGenericRelation
		if err := value.Decode(&relation); err != nil {
			return fmt.Errorf("invalid object of type %s at line %d object: %w", typer.Type, value.Line, err)
		}
		a.relation = &relation
	case "template":
		var template AttributeTemplate
		if err := value.Decode(&template); err != nil {
			return fmt.Errorf("invalid object of type template object in line %d: %w", value.Line, err)
		}
		a.template = &template
	}
	return nil
}

// To is shows a Relation where to point to.
type To struct {
	Collection string
	Field      ToField
}

// UnmarshalYAML decodes the models.yml to a To object.
func (t *To) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err == nil {
		cf := strings.Split(s, "/")
		if len(cf) != 2 {
			return fmt.Errorf("invalid value of `to` in line %d, expected one `/`: %s", value.Line, s)
		}
		t.Collection = cf[0]
		t.Field.Name = cf[1]
		return nil
	}

	var d struct {
		Collection string  `yaml:"collection"`
		Field      ToField `yaml:"field"`
	}
	if err := value.Decode(&d); err != nil {
		return fmt.Errorf("decoding to field at line %d: %w", value.Line, err)
	}

	t.Collection = d.Collection
	t.Field = d.Field
	return nil
}

// ToGeneric is like a To object, but for generic relations.
type ToGeneric struct {
	Collection []string
	Field      ToField
}

// ToField is the field part of a To object.
type ToField struct {
	Name string
	Type string
}

// UnmarshalYAML decodes the models.yml to a ToField object.
func (t *ToField) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err == nil {
		t.Name = s
		t.Type = "normal"
		return nil
	}

	var d struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
	}
	if err := value.Decode(&d); err != nil {
		return fmt.Errorf("decoding to field at line %d: %w", value.Line, err)
	}

	t.Name = d.Name
	t.Type = d.Type
	return nil
}
