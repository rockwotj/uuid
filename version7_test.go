package uuid

import (
	"testing"
	"time"
)

func TestNewV7WithTime(t *testing.T) {
	testCases := map[string]string{
		"test with current date":                      time.Now().Format(time.RFC3339),                                // now
		"test with past date":                         time.Now().Add(-1 * time.Hour * 24 * 365).Format(time.RFC3339), // 1 year ago
		"test with future date":                       time.Now().Add(time.Hour * 24 * 365).Format(time.RFC3339),      // 1 year from now
		"test with different timezone":                "2021-09-01T12:00:00+04:00",
		"test with negative timezone":                 "2021-09-01T12:00:00-12:00",
		"test with future date in different timezone": "2124-09-23T12:43:30+09:00",
	}

	for testName, inputTime := range testCases {
		t.Run(testName, func(t *testing.T) {
			customTime, err := time.Parse(time.RFC3339, inputTime)
			if err != nil {
				t.Errorf("time.Parse returned unexpected error %v", err)
			}
			id, err := NewV7WithTime(&customTime)
			if err != nil {
				t.Errorf("NewV7WithTime returned unexpected error %v", err)
			}
			if id.Version() != 7 {
				t.Errorf("got %d, want version 7", id.Version())
			}
			unixTime := time.Unix(id.Time().UnixTime())
			// Compare the times in UTC format, since the input time might have different timezone,
			// and the result is always in system timezone
			if customTime.UTC().Format(time.RFC3339) != unixTime.UTC().Format(time.RFC3339) {
				t.Errorf("got %s, want %s", unixTime.Format(time.RFC3339), customTime.Format(time.RFC3339))
			}
		})
	}
}

func TestNewV7FromTimeGeneratesUniqueUUIDs(t *testing.T) {
	now := time.Now()
	ids := make([]string, 0)
	runs := 26000

	for i := 0; i < runs; i++ {
		id, err := NewV7WithTime(&now)
		if err != nil {
			t.Errorf("NewV7WithTime returned unexpected error %v", err)
		}
		if id.Version() != 7 {
			t.Errorf("got %d, want version 7", id.Version())
		}

		// Make sure we add only unique values
		if !contains(t, ids, id.String()) {
			ids = append(ids, id.String())
		}
	}

	// Check we added all the UIDs
	if len(ids) != runs {
		t.Errorf("got %d UUIDs, want %d", len(ids), runs)
	}
}
