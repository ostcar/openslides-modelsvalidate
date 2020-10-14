package check

import (
	"fmt"
	"strings"

	models "github.com/OpenSlides/openslides-models-to-go"
)

// Check runs some checks on the given models.
func Check(data map[string]models.Model) error {
	validators := []func(map[string]models.Model) error{
		validateTypes,
		validateRelations,
	}

	errors := new(ErrorList)
	for _, v := range validators {
		if err := v(data); err != nil {
			errors.append(err)
		}
	}

	if !errors.empty() {
		return errors
	}
	return nil
}

func validateTypes(models map[string]models.Model) error {
	scalar := scalarTypes()
	special := specialTypes()
	errs := &ErrorList{
		Name:   "type validator",
		intent: 1,
	}
	for modelName, model := range models {
		for attrName, attr := range model.Attributes {
			if scalar[strings.TrimSuffix(attr.Type, "[]")] {
				continue
			}

			if special[attr.Type] {
				continue
			}

			errs.append(fmt.Errorf("Unknown type `%s` in %s/%s", attr.Type, modelName, attrName))
		}
	}
	if errs.empty() {
		return nil
	}
	return errs
}

func validateRelations(models map[string]models.Model) error {
	errs := &ErrorList{
		Name:   "relation validator",
		intent: 1,
	}
	for modelName, model := range models {
	Next:
		for attrName, attr := range model.Attributes {
			r := attr.Relation()
			if r == nil {
				continue
			}

			for _, c := range r.ToCollection() {
				toModel, ok := models[c]
				if !ok {
					errs.append(fmt.Errorf("%s/%s directs to nonexisting model `%s`", modelName, attrName, c))
					continue Next
				}
				_, ok = toModel.Attributes[r.ToField().Name]
				if !ok {
					errs.append(fmt.Errorf("%s/%s directs to nonexisting modelfield `%s/%s`", modelName, attrName, c, r.ToField().Name))
					continue Next
				}

				// TODO: check field type
			}
		}
	}
	if errs.empty() {
		return nil
	}
	return errs
}

// scalarTypes are the main types. All scalarTypes can be used as a list.
// JSON[], timestamp[] etc.
func scalarTypes() map[string]bool {
	s := []string{
		"string",
		"number",
		"boolean",
		"JSON",
		"HTMLPermissive",
		"HTMLStrict",
		"float",
		"decimal(6)",
		"timestamp",
	}
	out := make(map[string]bool)
	for _, t := range s {
		out[t] = true
	}
	return out
}

// specialTypes are realtion types in realtion to other fields.
func specialTypes() map[string]bool {
	s := []string{
		"relation",
		"relation-list",
		"generic-relation",
		"generic-relation-list",
		"template",
	}
	out := make(map[string]bool)
	for _, t := range s {
		out[t] = true
	}
	return out
}
