// Package dateutil contains a helper to parse vcard dates
package dateutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-vcard"
)

// Parse parses a vcard.Field into a time.Time
func Parse(field *vcard.Field) (d time.Time, err error) {
	if field == nil {
		return d, fmt.Errorf("nil-field given")
	}

	rawDate := field.Value

	if field.Params.Get("X-APPLE-OMIT-YEAR") != "" {
		// Yay, Apple bullshit. They don't use the proper way defined in
		// the RFC to omit the year but specify an invalid format with a
		// replace. As the year 0001 is the "zero-time" in Go we replace
		// the defined year (likely 1604) and move it to the year 1.
		//
		// Field should be something like this:
		// &{1604-09-13 map[X-APPLE-OMIT-YEAR:[1604]] }

		rawDate = strings.Replace(rawDate, field.Params.Get("X-APPLE-OMIT-YEAR"), "0001", 1)
	}

	// And now as we can't rely on `VALUE=DATE` being set (thanks Sabre)
	// we're trying to walk possible formats until we found a matching
	// one…

	for _, fmtCandidate := range []string{
		"20060102",        // Most likely, test first
		"2006-01-02",      // Invalid as of RFC, used by Apple, see above
		"060102",          // RFC compliant with 2-digit year
		"--0102",          // RFC-compliant omit-year
		"2006-01",         // Shouldn't exist as no day present but is valid
		"2006",            // Somewhere in the year
		"20060102T150405", // Full DATE-TIME in RFC
		"20060102T1504",   // DATE-TIME in RFC without seconds
		"20060102T15",     // DATE-TIME in RFC without minutes & seconds
	} {
		d, err = time.ParseInLocation(fmtCandidate, rawDate, time.Local)
		if err == nil {
			// Yay! It matched. Or at least it had the right length and numbers
			// at the right places so it SHOULD have matched.
			return d, nil
		}
	}

	// Well. We found no matching format.
	return d, fmt.Errorf("no format defined for %q", rawDate)
}
