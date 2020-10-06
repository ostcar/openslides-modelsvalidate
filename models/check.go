package models

import (
	"fmt"
	"strings"
)

// Check runs some checks on the given models.
func Check(models map[string]Model) error {
	validators := []func(map[string]Model) error{
		validateTypes,
		validateRelations,
	}

	errors := new(ErrorList)
	for _, v := range validators {
		if err := v(models); err != nil {
			errors.append(err)
		}
	}

	if !errors.empty() {
		return errors
	}
	return nil
}

func validateTypes(models map[string]Model) error {
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

func validateRelations(models map[string]Model) error {
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

			for _, c := range r.toCollection() {
				toModel, ok := models[c]
				if !ok {
					errs.append(fmt.Errorf("%s/%s directs to nonexisting model `%s`", modelName, attrName, c))
					continue Next
				}
				_, ok = toModel.Attributes[r.toField().Name]
				if !ok {
					errs.append(fmt.Errorf("%s/%s directs to nonexisting modelfield `%s/%s`", modelName, attrName, c, r.toField().Name))
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

func scalarTypes() map[string]bool {
	s := []string{
		"string",
		"number",
		"boolean",
		"JSON",
		"HtmlPermissive",
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
