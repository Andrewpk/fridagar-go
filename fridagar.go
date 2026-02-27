// Package fridagar provides functions for looking up Icelandic public holidays
// and commonly celebrated "special" days, as well as resolving business days
// before/after a given date.
//
// It is a Go port of https://github.com/gaui/fridagar-node
//
// All returned dates are in UTC and set to midnight.
package fridagar

import (
	"sort"
	"time"
)

// HolidayKey is a stable identifier for an Icelandic public holiday.
type HolidayKey = string

// SpecialDayKey is a stable identifier for an Icelandic "special" day.
type SpecialDayKey = string

// DayKey is a stable identifier for any Icelandic holiday or special day.
type DayKey = string

// Predefined HolidayKey values.
const (
	KeyNyars    HolidayKey = "nyars"
	KeySkir     HolidayKey = "skir"
	KeyFosLangi HolidayKey = "foslangi"
	KeyPaska    HolidayKey = "paska"
	KeyPaska2   HolidayKey = "paska2"
	KeySumar1   HolidayKey = "sumar1"
	KeyUppst    HolidayKey = "uppst"
	KeyMai1     HolidayKey = "mai1"
	KeyHvitas   HolidayKey = "hvitas"
	KeyHvitas2  HolidayKey = "hvitas2"
	KeyJun17    HolidayKey = "jun17"
	KeyVerslm   HolidayKey = "verslm"
	KeyAdfanga  HolidayKey = "adfanga"
	KeyJola     HolidayKey = "jola"
	KeyJola2    HolidayKey = "jola2"
	KeyGamlars  HolidayKey = "gamlars"
)

// Predefined SpecialDayKey values.
const (
	KeyThrettand SpecialDayKey = "þrettand"
	KeyBonda     SpecialDayKey = "bonda"
	KeyBollu     SpecialDayKey = "bollu"
	KeySprengi   SpecialDayKey = "sprengi"
	KeyOsku      SpecialDayKey = "osku"
	KeyValent    SpecialDayKey = "valent"
	KeyKonu      SpecialDayKey = "konu"
	KeySjomanna  SpecialDayKey = "sjomanna"
	KeySumsolst  SpecialDayKey = "sumsolst"
	KeyJonsm     SpecialDayKey = "jonsm"
	KeyVetur1    SpecialDayKey = "vetur1"
	KeyHrekkja   SpecialDayKey = "hrekkja"
	KeyIslTungu  SpecialDayKey = "isltungu"
	KeyFullv     SpecialDayKey = "fullv"
	KeyVetsolst  SpecialDayKey = "vetsolst"
	KeyThorl     SpecialDayKey = "thorl"
)

// Day represents an Icelandic public holiday or commonly celebrated "special" day.
type Day struct {
	// Date is the UTC midnight date of this day.
	Date time.Time
	// Description is the Icelandic name of the day.
	Description string
	// Key is a stable identifier useful for translations.
	Key DayKey
	// Holiday indicates whether this is an official public holiday (non-working day).
	Holiday bool
	// HalfDay indicates if this is only a half-day holiday (e.g. Aðfangadagur, Gamlársdagur).
	HalfDay bool
}

// GetAllDays returns all Icelandic public holidays and commonly celebrated
// "special" days for a given year, sorted by date.
func GetAllDays(year int) []Day {
	return calcSpecialDays(year)
}

// GetAllDaysForMonth returns all days for a given year and 1-based month.
func GetAllDaysForMonth(year int, month int) []Day {
	if month < 1 || month > 12 {
		return nil
	}
	days := calcSpecialDays(year)
	var result []Day
	for _, d := range days {
		if int(d.Date.Month()) == month {
			result = append(result, d)
		}
	}
	return result
}

// GetAllDaysKeyed returns a map of all days for a given year, keyed by their DayKey.
func GetAllDaysKeyed(year int) map[DayKey]Day {
	days := GetAllDays(year)
	m := make(map[DayKey]Day, len(days))
	for _, d := range days {
		m[d.Key] = d
	}
	return m
}

// GetHolidays returns only the official public holidays (non-working days) for a given year.
func GetHolidays(year int) []Day {
	days := calcSpecialDays(year)
	var result []Day
	for _, d := range days {
		if d.Holiday {
			result = append(result, d)
		}
	}
	return result
}

// GetHolidaysForMonth returns official public holidays for a given year and 1-based month.
func GetHolidaysForMonth(year int, month int) []Day {
	if month < 1 || month > 12 {
		return nil
	}
	days := calcSpecialDays(year)
	var result []Day
	for _, d := range days {
		if d.Holiday && int(d.Date.Month()) == month {
			result = append(result, d)
		}
	}
	return result
}

// GetOtherDays returns only the unofficial "special" days for a given year.
func GetOtherDays(year int) []Day {
	days := calcSpecialDays(year)
	var result []Day
	for _, d := range days {
		if !d.Holiday {
			result = append(result, d)
		}
	}
	return result
}

// GetOtherDaysForMonth returns unofficial "special" days for a given year and 1-based month.
func GetOtherDaysForMonth(year int, month int) []Day {
	if month < 1 || month > 12 {
		return nil
	}
	days := calcSpecialDays(year)
	var result []Day
	for _, d := range days {
		if !d.Holiday && int(d.Date.Month()) == month {
			result = append(result, d)
		}
	}
	return result
}

// IsSpecialDay checks if the given date is a holiday or a special day, returning
// the Day info and true if found, or an empty Day and false otherwise.
func IsSpecialDay(date time.Time) (Day, bool) {
	date = truncateToUTCDay(date)
	days := calcSpecialDays(date.Year())
	for _, d := range days {
		if d.Date.Equal(date) {
			return d, true
		}
	}
	return Day{}, false
}

// IsHoliday checks if the given date is an official public holiday, returning
// the Day info and true if found, or an empty Day and false otherwise.
func IsHoliday(date time.Time) (Day, bool) {
	day, ok := IsSpecialDay(date)
	if ok && day.Holiday {
		return day, true
	}
	return Day{}, false
}

// WorkdaysFromDate returns the date that is the given number of business days
// before (negative) or after (positive) the reference date.
// If includeHalfDays is true, half-day holidays are treated as workdays.
// A zero offset returns the reference date itself.
func WorkdaysFromDate(days int, refDate time.Time, includeHalfDays bool) time.Time {
	date := truncateToUTCDay(refDate)
	if days == 0 {
		return date
	}

	delta := 1
	if days < 0 {
		delta = -1
	}
	count := days
	if count < 0 {
		count = -count
	}

	var holidays []Day
	holidayYear := -1

	for count > 0 {
		date = date.AddDate(0, 0, delta)
		wDay := date.Weekday()
		dateYear := date.Year()

		if dateYear != holidayYear {
			holidayYear = dateYear
			all := calcSpecialDays(dateYear)
			holidays = holidays[:0]
			for _, d := range all {
				if d.Holiday {
					holidays = append(holidays, d)
				}
			}
		}

		notWorkDay := wDay == time.Sunday || wDay == time.Saturday
		if !notWorkDay {
			for _, h := range holidays {
				if h.Date.Equal(date) && !(includeHalfDays && h.HalfDay) {
					notWorkDay = true
					break
				}
			}
		}
		if notWorkDay {
			continue
		}
		count--
	}

	return date
}

// truncateToUTCDay truncates a time to midnight UTC.
func truncateToUTCDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// sortDays sorts days by date, with holidays sorted before special days on the same date.
func sortDays(days []Day) {
	sort.SliceStable(days, func(i, j int) bool {
		if !days[i].Date.Equal(days[j].Date) {
			return days[i].Date.Before(days[j].Date)
		}
		// Sort holidays ahead of special days on the same date
		if days[i].Holiday != days[j].Holiday {
			return days[i].Holiday
		}
		return false
	})
}
