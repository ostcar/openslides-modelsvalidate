package check_test

import (
	"errors"
	"strings"
	"testing"

	models "github.com/OpenSlides/openslides-models-to-go"
	"github.com/OpenSlides/openslides-modelsvalidate/internal/check"
)

func TestCheck(t *testing.T) {
	for _, tt := range []struct {
		name string
		yaml string
		err  string
	}{
		{
			"unknown type",
			yamlUnknownAttrType,
			"Unknown type `unknown` in some_model/field",
		},
		{
			"invalid relation",
			yamlInvalidRelation,
			"some_model/no_other_model directs to nonexisting model `not_existing`",
		},
		{
			"invalid relation",
			yamlInvalidRelation,
			"some_model/no_other_field directs to nonexisting modelfield `other_model/bar`",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			data, err := models.Unmarshal(strings.NewReader(tt.yaml))
			if err != nil {
				t.Fatalf("Can not unmarshal yaml: %v", err)
			}
			gotErr := check.Check(data)
			if tt.err == "" {
				if gotErr != nil {
					t.Errorf("Models.Check() returned an unexepcted error: %v", err)
				}
				return
			}

			if gotErr == nil {
				t.Fatalf("Models.Check() did not return an error, expected: %v", tt.err)
			}

			var errList *check.ErrorList
			if !errors.As(gotErr, &errList) {
				t.Fatalf("Models.Check() did not return a ListError, got: %v", gotErr)
			}

			var found bool
			for _, err := range errList.Errs {
				var errList *check.ErrorList
				if !errors.As(err, &errList) {
					continue
				}

				for _, err := range errList.Errs {
					if err.Error() == tt.err {
						found = true
					}
				}
			}

			if !found {
				t.Errorf("Models.Check() returned %v, expected %v", gotErr, tt.err)
			}
		})
	}
}

const yamlUnknownAttrType = `---
some_model:
  field: unknown
`

const yamlInvalidRelation = `---
some_model:
  no_other_model:
    type: relation
    to: not_existing/field
  no_other_field:
    type: relation
    to: other_model/bar
other_model:
  foo: string
`
