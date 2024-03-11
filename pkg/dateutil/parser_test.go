package dateutil

import (
	"testing"
	"time"

	"github.com/emersion/go-vcard"
	"github.com/stretchr/testify/assert"
)

func TestParseDate(t *testing.T) {
	for _, tc := range []struct {
		ExpectedTime time.Time
		Field        *vcard.Field
	}{
		{
			ExpectedTime: time.Date(2000, 2, 5, 0, 0, 0, 0, time.Local),
			Field: &vcard.Field{
				Value: "20000205",
				Params: vcard.Params{
					"VALUE": []string{"DATE"},
				},
			},
		},
		{
			ExpectedTime: time.Date(1, 2, 15, 0, 0, 0, 0, time.Local),
			Field: &vcard.Field{
				Value: "1604-02-15",
				Params: vcard.Params{
					"X-APPLE-OMIT-YEAR": []string{"1604"},
				},
			},
		},
		{
			ExpectedTime: time.Date(1, 11, 2, 0, 0, 0, 0, time.Local),
			Field: &vcard.Field{
				Value: "20221102",
				Params: vcard.Params{
					"X-APPLE-OMIT-YEAR": []string{"2022"},
					"VALUE":             []string{"DATE"},
				},
			},
		},
	} {
		d, err := Parse(tc.Field)
		assert.NoError(t, err, tc.Field)
		assert.Equal(t, tc.ExpectedTime, d)
	}
}
