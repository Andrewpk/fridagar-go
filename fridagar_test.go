package fridagar

import (
	"fmt"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

func expectDate(t *testing.T, label string, got time.Time, year, month, day int) {
	t.Helper()
	if got.Year() != year || int(got.Month()) != month || got.Day() != day {
		t.Errorf("%s: expected %d-%02d-%02d, got %s",
			label, year, month, day, got.Format("2006-01-02"))
	}
}

func expectDateStr(t *testing.T, label string, got time.Time, expected string) {
	t.Helper()
	gotStr := got.Format("2006-01-02")
	if gotStr != expected {
		t.Errorf("%s: expected %s, got %s", label, expected, gotStr)
	}
}

func d(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(fmt.Sprintf("invalid date %q: %v", s, err))
	}
	return t
}

// ---------------------------------------------------------------------------
// getAllDays / getHolidays / getOtherDays shared behavior
// ---------------------------------------------------------------------------

func TestGetAllDays_SharedBehavior(t *testing.T) {
	t.Run("all dates are at midnight UTC", func(t *testing.T) {
		for _, day := range GetAllDays(2024) {
			if day.Date.Hour() != 0 || day.Date.Minute() != 0 || day.Date.Second() != 0 || day.Date.Nanosecond() != 0 {
				t.Errorf("expected midnight UTC for %s, got %s", day.Key, day.Date.Format(time.RFC3339Nano))
			}
			if day.Date.Location() != time.UTC {
				t.Errorf("expected UTC location for %s", day.Key)
			}
		}
	})

	t.Run("all dates are in the requested year", func(t *testing.T) {
		days := GetAllDays(2021)
		if len(days) == 0 {
			t.Fatal("expected non-empty result")
		}
		for _, day := range days {
			if day.Date.Year() != 2021 {
				t.Errorf("expected year 2021 for %s, got %d", day.Key, day.Date.Year())
			}
		}
		// First day should be in January, last in December
		if days[0].Date.Month() != time.January {
			t.Errorf("expected first day in January, got %s", days[0].Date.Month())
		}
		if days[len(days)-1].Date.Month() != time.December {
			t.Errorf("expected last day in December, got %s", days[len(days)-1].Date.Month())
		}
	})
}

func TestGetHolidays_SharedBehavior(t *testing.T) {
	t.Run("all dates are at midnight UTC", func(t *testing.T) {
		for _, day := range GetHolidays(2024) {
			if day.Date.Hour() != 0 || day.Date.Minute() != 0 || day.Date.Second() != 0 {
				t.Errorf("expected midnight UTC for %s, got %s", day.Key, day.Date.Format(time.RFC3339))
			}
		}
	})

	t.Run("all dates are in the requested year", func(t *testing.T) {
		days := GetHolidays(2021)
		if len(days) == 0 {
			t.Fatal("expected non-empty result")
		}
		for _, day := range days {
			if day.Date.Year() != 2021 {
				t.Errorf("expected year 2021 for %s, got %d", day.Key, day.Date.Year())
			}
		}
		if days[0].Date.Month() != time.January {
			t.Errorf("expected first holiday in January, got %s", days[0].Date.Month())
		}
		if days[len(days)-1].Date.Month() != time.December {
			t.Errorf("expected last holiday in December, got %s", days[len(days)-1].Date.Month())
		}
	})

	t.Run("only returns holidays", func(t *testing.T) {
		for _, day := range GetHolidays(2021) {
			if !day.Holiday {
				t.Errorf("expected holiday=true for %s", day.Key)
			}
		}
		// Also check with month filtering
		for _, day := range GetHolidaysForMonth(2021, 12) {
			if !day.Holiday {
				t.Errorf("expected holiday=true for %s (month filter)", day.Key)
			}
		}
	})
}

func TestGetOtherDays_SharedBehavior(t *testing.T) {
	t.Run("all dates are at midnight UTC", func(t *testing.T) {
		for _, day := range GetOtherDays(2024) {
			if day.Date.Hour() != 0 || day.Date.Minute() != 0 || day.Date.Second() != 0 {
				t.Errorf("expected midnight UTC for %s, got %s", day.Key, day.Date.Format(time.RFC3339))
			}
		}
	})

	t.Run("all dates are in the requested year", func(t *testing.T) {
		days := GetOtherDays(2021)
		if len(days) == 0 {
			t.Fatal("expected non-empty result")
		}
		for _, day := range days {
			if day.Date.Year() != 2021 {
				t.Errorf("expected year 2021 for %s, got %d", day.Key, day.Date.Year())
			}
		}
	})

	t.Run("only returns non-holidays", func(t *testing.T) {
		for _, day := range GetOtherDays(2021) {
			if day.Holiday {
				t.Errorf("expected holiday=false for %s", day.Key)
			}
		}
		// Also check with month filtering
		for _, day := range GetOtherDaysForMonth(2021, 12) {
			if day.Holiday {
				t.Errorf("expected holiday=false for %s (month filter)", day.Key)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Month filtering
// ---------------------------------------------------------------------------

func TestGetAllDaysForMonth(t *testing.T) {
	t.Run("returns days in the requested month", func(t *testing.T) {
		days := GetAllDaysForMonth(2021, 12)
		if len(days) == 0 {
			t.Fatal("expected non-empty days for December 2021")
		}
		for _, day := range days {
			if day.Date.Month() != time.December {
				t.Errorf("expected December, got %s for %s", day.Date.Month(), day.Key)
			}
			if day.Date.Year() != 2021 {
				t.Errorf("expected year 2021, got %d for %s", day.Date.Year(), day.Key)
			}
		}
	})

	t.Run("out of bound months return nil", func(t *testing.T) {
		if days := GetAllDaysForMonth(2012, -1); days != nil {
			t.Errorf("expected nil for month -1, got %d days", len(days))
		}
		if days := GetAllDaysForMonth(2012, 0); days != nil {
			t.Errorf("expected nil for month 0, got %d days", len(days))
		}
		if days := GetAllDaysForMonth(2012, 13); days != nil {
			t.Errorf("expected nil for month 13, got %d days", len(days))
		}
	})
}

func TestGetHolidaysForMonth(t *testing.T) {
	t.Run("returns holidays in the requested month", func(t *testing.T) {
		days := GetHolidaysForMonth(2021, 12)
		if len(days) == 0 {
			t.Fatal("expected non-empty holidays for December 2021")
		}
		for _, day := range days {
			if day.Date.Month() != time.December {
				t.Errorf("expected December, got %s for %s", day.Date.Month(), day.Key)
			}
			if !day.Holiday {
				t.Errorf("expected holiday=true for %s", day.Key)
			}
		}
	})

	t.Run("out of bound months return nil", func(t *testing.T) {
		if days := GetHolidaysForMonth(2012, -1); days != nil {
			t.Errorf("expected nil for month -1, got %d days", len(days))
		}
		if days := GetHolidaysForMonth(2012, 0); days != nil {
			t.Errorf("expected nil for month 0, got %d days", len(days))
		}
		if days := GetHolidaysForMonth(2012, 13); days != nil {
			t.Errorf("expected nil for month 13, got %d days", len(days))
		}
	})
}

func TestGetOtherDaysForMonth(t *testing.T) {
	t.Run("returns other days in the requested month", func(t *testing.T) {
		days := GetOtherDaysForMonth(2021, 12)
		if len(days) == 0 {
			t.Fatal("expected non-empty other days for December 2021")
		}
		for _, day := range days {
			if day.Date.Month() != time.December {
				t.Errorf("expected December, got %s for %s", day.Date.Month(), day.Key)
			}
			if day.Holiday {
				t.Errorf("expected holiday=false for %s", day.Key)
			}
		}
	})

	t.Run("out of bound months return nil", func(t *testing.T) {
		if days := GetOtherDaysForMonth(2012, -1); days != nil {
			t.Errorf("expected nil for month -1, got %d days", len(days))
		}
		if days := GetOtherDaysForMonth(2012, 0); days != nil {
			t.Errorf("expected nil for month 0, got %d days", len(days))
		}
		if days := GetOtherDaysForMonth(2012, 13); days != nil {
			t.Errorf("expected nil for month 13, got %d days", len(days))
		}
	})
}

// ---------------------------------------------------------------------------
// getAllDays — sorting
// ---------------------------------------------------------------------------

func TestGetAllDays_Sorted(t *testing.T) {
	for _, year := range []int{2011, 2023, 2024} {
		t.Run(fmt.Sprintf("sorted_%d", year), func(t *testing.T) {
			days := GetAllDays(year)
			for i := 1; i < len(days); i++ {
				if days[i].Date.Before(days[i-1].Date) {
					t.Errorf("year %d: not sorted — %s (%s) before %s (%s)",
						year,
						days[i-1].Key, days[i-1].Date.Format("2006-01-02"),
						days[i].Key, days[i].Date.Format("2006-01-02"))
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// getAllDays — day count and key completeness
// ---------------------------------------------------------------------------

func TestGetAllDays_CompleteDayCount(t *testing.T) {
	days := GetAllDays(2024)

	// The Node.js source defines exactly 32 days
	if len(days) != 32 {
		t.Errorf("expected 32 days total, got %d", len(days))
	}

	holidays := GetHolidays(2024)
	if len(holidays) != 16 {
		t.Errorf("expected 16 holidays, got %d", len(holidays))
	}

	other := GetOtherDays(2024)
	if len(other) != 16 {
		t.Errorf("expected 16 other/special days, got %d", len(other))
	}
}

func TestGetAllDays_AllKeysPresent(t *testing.T) {
	expectedHolidayKeys := []string{
		"nyars", "skir", "foslangi", "paska", "paska2", "sumar1",
		"uppst", "mai1", "hvitas", "hvitas2", "jun17", "verslm",
		"adfanga", "jola", "jola2", "gamlars",
	}
	expectedSpecialKeys := []string{
		"þrettand", "bonda", "bollu", "sprengi", "osku", "valent",
		"konu", "sjomanna", "sumsolst", "jonsm", "vetur1", "hrekkja",
		"isltungu", "fullv", "vetsolst", "thorl",
	}

	keyed := GetAllDaysKeyed(2024)

	for _, key := range expectedHolidayKeys {
		day, ok := keyed[key]
		if !ok {
			t.Errorf("missing holiday key: %s", key)
			continue
		}
		if !day.Holiday {
			t.Errorf("key %s should be a holiday", key)
		}
	}
	for _, key := range expectedSpecialKeys {
		day, ok := keyed[key]
		if !ok {
			t.Errorf("missing special day key: %s", key)
			continue
		}
		if day.Holiday {
			t.Errorf("key %s should NOT be a holiday", key)
		}
	}
}

// ---------------------------------------------------------------------------
// getAllDays — rímspilliár handling
// ---------------------------------------------------------------------------

func TestGetAllDays_Rimspillir(t *testing.T) {
	// Bóndadagur dates
	bonda2023 := GetAllDaysKeyed(2023)[KeyBonda]
	expectDate(t, "bonda 2023", bonda2023.Date, 2023, 1, 20)

	bonda2024 := GetAllDaysKeyed(2024)[KeyBonda]
	expectDate(t, "bonda 2024", bonda2024.Date, 2024, 1, 26)

	bonda2025 := GetAllDaysKeyed(2025)[KeyBonda]
	expectDate(t, "bonda 2025", bonda2025.Date, 2025, 1, 24)

	// Fyrsti vetrardagur dates
	vetur2023 := GetAllDaysKeyed(2023)[KeyVetur1]
	expectDate(t, "vetur1 2023", vetur2023.Date, 2023, 10, 28)

	vetur2024 := GetAllDaysKeyed(2024)[KeyVetur1]
	expectDate(t, "vetur1 2024", vetur2024.Date, 2024, 10, 26)

	vetur2025 := GetAllDaysKeyed(2025)[KeyVetur1]
	expectDate(t, "vetur1 2025", vetur2025.Date, 2025, 10, 25)
}

// ---------------------------------------------------------------------------
// getAllDays — sjómanna shift when hvítasunnu is first Sunday in June
// ---------------------------------------------------------------------------

func TestGetAllDays_SjomannaShift(t *testing.T) {
	// Hvítasunnudagur is first Sunday of June in 2022 → sjómanna shifts to 2nd Sunday
	sjomanna22 := GetAllDaysKeyed(2022)[KeySjomanna]
	expectDate(t, "sjomanna 2022", sjomanna22.Date, 2022, 6, 12)

	// Hvítasunnudagur is in May in 2023 → sjómanna is first Sunday of June
	sjomanna23 := GetAllDaysKeyed(2023)[KeySjomanna]
	expectDate(t, "sjomanna 2023", sjomanna23.Date, 2023, 6, 4)

	// Hvítasunnudagur is 2nd Sunday of June in 2025 → sjómanna is first Sunday of June
	sjomanna25 := GetAllDaysKeyed(2025)[KeySjomanna]
	expectDate(t, "sjomanna 2025", sjomanna25.Date, 2025, 6, 1)
}

// ---------------------------------------------------------------------------
// getAllDays — solstice days
// ---------------------------------------------------------------------------

func TestGetAllDays_Solstice(t *testing.T) {
	days := GetAllDaysKeyed(2023)

	sumsolst := days[KeySumsolst]
	expectDate(t, "sumsolst 2023", sumsolst.Date, 2023, 6, 21)

	vetsolst := days[KeyVetsolst]
	expectDate(t, "vetsolst 2023", vetsolst.Date, 2023, 12, 22)
}

// ---------------------------------------------------------------------------
// getAllDaysKeyed
// ---------------------------------------------------------------------------

func TestGetAllDaysKeyed(t *testing.T) {
	t.Run("returns keyed map with correct descriptions", func(t *testing.T) {
		days := GetAllDaysKeyed(2024)

		bonda, ok := days[KeyBonda]
		if !ok {
			t.Fatal("expected bonda key in keyed result")
		}
		if bonda.Description != "Bóndadagur" {
			t.Errorf("expected description 'Bóndadagur', got %q", bonda.Description)
		}

		adfanga, ok := days[KeyAdfanga]
		if !ok {
			t.Fatal("expected adfanga key in keyed result")
		}
		if adfanga.Description != "Aðfangadagur" {
			t.Errorf("expected description 'Aðfangadagur', got %q", adfanga.Description)
		}
	})
}

// ---------------------------------------------------------------------------
// isSpecialDay
// ---------------------------------------------------------------------------

func TestIsSpecialDay(t *testing.T) {
	t.Run("special days in 2023-12", func(t *testing.T) {
		// Christmas
		day, ok := IsSpecialDay(time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC))
		if !ok {
			t.Fatal("expected 2023-12-25 to be a special day")
		}
		if day.Key != KeyJola {
			t.Errorf("expected key %s, got %s", KeyJola, day.Key)
		}

		// Winter solstice
		_, ok = IsSpecialDay(time.Date(2023, 12, 22, 0, 0, 0, 0, time.UTC))
		if !ok {
			t.Fatal("expected 2023-12-22 to be a special day (Vetrarsólstöður)")
		}

		// Non-midnight should still work
		_, ok = IsSpecialDay(time.Date(2023, 12, 22, 13, 0, 0, 0, time.UTC))
		if !ok {
			t.Fatal("expected 2023-12-22T13:00 to be a special day")
		}

		// New Year's Eve
		_, ok = IsSpecialDay(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC))
		if !ok {
			t.Fatal("expected 2023-12-31 to be a special day")
		}
	})

	t.Run("returns false for a non-special day", func(t *testing.T) {
		_, ok := IsSpecialDay(time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC))
		if ok {
			t.Error("expected 2024-03-15 to NOT be a special day")
		}
	})
}

// ---------------------------------------------------------------------------
// isHoliday
// ---------------------------------------------------------------------------

func TestIsHoliday(t *testing.T) {
	t.Run("holidays in 2023-12", func(t *testing.T) {
		// Vetrarsólstöður (Dec 22) is NOT an official holiday
		_, ok := IsHoliday(time.Date(2023, 12, 22, 0, 0, 0, 0, time.UTC))
		if ok {
			t.Error("expected 2023-12-22 NOT to be a holiday")
		}

		// Christmas IS a holiday
		day, ok := IsHoliday(time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC))
		if !ok {
			t.Fatal("expected 2023-12-25 to be a holiday")
		}
		if day.Key != KeyJola {
			t.Errorf("expected key %s, got %s", KeyJola, day.Key)
		}

		// New Year's Eve IS a holiday
		_, ok = IsHoliday(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC))
		if !ok {
			t.Fatal("expected 2023-12-31 to be a holiday")
		}
	})

	t.Run("returns false for a non-holiday date", func(t *testing.T) {
		_, ok := IsHoliday(time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC))
		if ok {
			t.Error("expected 2024-03-15 to NOT be a holiday")
		}
	})
}

// ---------------------------------------------------------------------------
// workdaysFromDate
// ---------------------------------------------------------------------------

func TestWorkdaysFromDate(t *testing.T) {
	t.Run("finds future days", func(t *testing.T) {
		expectDateStr(t, "workday +2 from Wed",
			WorkdaysFromDate(2, d("2023-12-06"), false), "2023-12-08")
		expectDateStr(t, "workday +2 from Fri",
			WorkdaysFromDate(2, d("2023-12-08"), false), "2023-12-12")
		expectDateStr(t, "workday +2 from Mon (Christmas)",
			WorkdaysFromDate(2, d("2023-12-25"), false), "2023-12-28")
		expectDateStr(t, "workday +60 from Wed",
			WorkdaysFromDate(60, d("2023-06-01"), false), "2023-08-25")
	})

	t.Run("treats halfDays as non-workdays by default", func(t *testing.T) {
		expectDateStr(t, "halfday default +1 from Wed",
			WorkdaysFromDate(1, d("2021-12-22"), false), "2021-12-23")
		expectDateStr(t, "halfday default +1 from Thu (Þorláksmessa)",
			WorkdaysFromDate(1, d("2021-12-23"), false), "2021-12-27")
		expectDateStr(t, "halfday included +1 from Thu (Þorláksmessa)",
			WorkdaysFromDate(1, d("2021-12-23"), true), "2021-12-24")
	})

	t.Run("returns the refDate if offset is 0", func(t *testing.T) {
		expectDateStr(t, "zero offset Wed",
			WorkdaysFromDate(0, d("2023-12-06"), false), "2023-12-06")
		expectDateStr(t, "zero offset Fri",
			WorkdaysFromDate(0, d("2023-12-08"), false), "2023-12-08")
	})

	t.Run("finds previous workdays", func(t *testing.T) {
		expectDateStr(t, "workday -1 from Mon",
			WorkdaysFromDate(-1, d("2020-12-21"), false), "2020-12-18")
		expectDateStr(t, "workday -1 from Fri (Christmas)",
			WorkdaysFromDate(-1, d("2020-12-25"), false), "2020-12-23")
		expectDateStr(t, "workday -1 from Fri (Christmas) halfday",
			WorkdaysFromDate(-1, d("2020-12-25"), true), "2020-12-24")
	})

	t.Run("searches into adjacent year", func(t *testing.T) {
		expectDateStr(t, "workday +1 from Fri Dec 29",
			WorkdaysFromDate(1, d("2023-12-29"), false), "2024-01-02")
		expectDateStr(t, "workday -1 from Tue Jan 2",
			WorkdaysFromDate(-1, d("2025-01-02"), false), "2024-12-30")
	})
}

// ---------------------------------------------------------------------------
// Easter calculation
// ---------------------------------------------------------------------------

func TestEaster(t *testing.T) {
	tests := []struct {
		year  int
		month time.Month
		day   int
	}{
		{2020, time.April, 12},
		{2021, time.April, 4},
		{2022, time.April, 17},
		{2023, time.April, 9},
		{2024, time.March, 31},
		{2025, time.April, 20},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("easter_%d", tt.year), func(t *testing.T) {
			got := easter(tt.year)
			if got.Year() != tt.year || got.Month() != tt.month || got.Day() != tt.day {
				t.Errorf("easter(%d) = %s, want %d-%02d-%02d",
					tt.year, got.Format("2006-01-02"), tt.year, tt.month, tt.day)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Half-day flag
// ---------------------------------------------------------------------------

func TestHalfDayFlags(t *testing.T) {
	keyed := GetAllDaysKeyed(2024)

	// Aðfangadagur (Christmas Eve) should be a half-day holiday
	adfanga := keyed[KeyAdfanga]
	if !adfanga.Holiday {
		t.Error("adfanga should be a holiday")
	}
	if !adfanga.HalfDay {
		t.Error("adfanga should be a half-day")
	}

	// Gamlársdagur (New Year's Eve) should be a half-day holiday
	gamlars := keyed[KeyGamlars]
	if !gamlars.Holiday {
		t.Error("gamlars should be a holiday")
	}
	if !gamlars.HalfDay {
		t.Error("gamlars should be a half-day")
	}

	// Jóladagur should NOT be a half-day
	jola := keyed[KeyJola]
	if !jola.Holiday {
		t.Error("jola should be a holiday")
	}
	if jola.HalfDay {
		t.Error("jola should NOT be a half-day")
	}

	// Special (non-holiday) days should never be half-days
	for _, day := range GetOtherDays(2024) {
		if day.HalfDay {
			t.Errorf("special day %s should not be a half-day", day.Key)
		}
	}
}

// ---------------------------------------------------------------------------
// Verify specific well-known dates for a given year
// ---------------------------------------------------------------------------

func TestKnownDates2024(t *testing.T) {
	keyed := GetAllDaysKeyed(2024)

	expected := map[string]string{
		"nyars":    "2024-01-01",
		"þrettand": "2024-01-06",
		"bonda":    "2024-01-26",
		"bollu":    "2024-02-12",
		"sprengi":  "2024-02-13",
		"osku":     "2024-02-14",
		"valent":   "2024-02-14",
		"konu":     "2024-02-25",
		"skir":     "2024-03-28",
		"foslangi": "2024-03-29",
		"paska":    "2024-03-31",
		"paska2":   "2024-04-01",
		"sumar1":   "2024-04-25",
		"mai1":     "2024-05-01",
		"uppst":    "2024-05-09",
		"hvitas":   "2024-05-19",
		"hvitas2":  "2024-05-20",
		"sjomanna": "2024-06-02",
		"jun17":    "2024-06-17",
		"jonsm":    "2024-06-24",
		"verslm":   "2024-08-05",
		"hrekkja":  "2024-10-31",
		"isltungu": "2024-11-16",
		"fullv":    "2024-12-01",
		"thorl":    "2024-12-23",
		"adfanga":  "2024-12-24",
		"jola":     "2024-12-25",
		"jola2":    "2024-12-26",
		"gamlars":  "2024-12-31",
	}

	for key, wantDate := range expected {
		day, ok := keyed[key]
		if !ok {
			t.Errorf("missing key %s", key)
			continue
		}
		gotDate := day.Date.Format("2006-01-02")
		if gotDate != wantDate {
			t.Errorf("key %s: expected %s, got %s", key, wantDate, gotDate)
		}
	}
}
