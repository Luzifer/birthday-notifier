package dateutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const timeDay = 24 * time.Hour

func TestProjectToNextBirthday(t *testing.T) {
	// Now should stay in the year
	assert.Equal(
		t,
		time.Now().Year(),
		ProjectToNextBirthday(time.Now()).Year(),
	)

	// Start-of-day should stay in the year
	assert.Equal(
		t,
		time.Now().Year(),
		ProjectToNextBirthday(TodayStartOfDay()).Year(),
	)

	// Tomorrow should stay in the year
	assert.Equal(
		t,
		time.Now().Year(),
		ProjectToNextBirthday(time.Now().Add(timeDay)).Year(),
	)

	// Yesterday should go to next year
	assert.Equal(
		t,
		time.Now().Year()+1,
		ProjectToNextBirthday(time.Now().Add(-timeDay)).Year(),
	)

	// Yesterday, thirty years ago should go to next year
	assert.Equal(
		t,
		time.Now().Year()+1,
		ProjectToNextBirthday(time.Date(time.Now().Year()-30, time.Now().Month(), time.Now().Day()-1, 0, 0, 0, 0, time.Local)).Year(),
	)

	// Tomorrow, thirty years ago should go to this year
	assert.Equal(
		t,
		time.Now().Year(),
		ProjectToNextBirthday(time.Date(time.Now().Year()-30, time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.Local)).Year(),
	)
}
