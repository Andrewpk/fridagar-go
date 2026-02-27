package fridagar

import (
	"math"
	"time"
)

// rimspillir returns 1 if the given year is a "Rímspilliár" year, 0 otherwise.
// This value is used to shift the base/reference date for certain special days.
func rimspillir(year int) int {
	// Precalculated values for 1900-2100
	precalc := map[int]int{
		1911: 1, 1939: 1, 1967: 1, 1995: 1,
		2023: 1, 2051: 1, 2079: 1,
	}
	if v, ok := precalc[year]; ok {
		return v
	}
	nextYear := year + 1
	nextIsLeapYear := nextYear%4 == 0 && (nextYear%100 != 0 || nextYear%400 == 0)
	// Dec 31 of (year-1)
	dec31 := time.Date(year-1, 12, 31, 0, 0, 0, 0, time.UTC)
	isRimspilliar := nextIsLeapYear && dec31.Weekday() == time.Saturday
	if isRimspilliar {
		return 1
	}
	return 0
}

// findNextWeekDay finds the next occurrence of targetWDay on or after the given date.
// month is 1-based (time.Month).
func findNextWeekDay(year int, month time.Month, day int, targetWDay time.Weekday) time.Time {
	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	offset := (int(targetWDay) - int(date.Weekday()) + 7) % 7
	return date.AddDate(0, 0, offset)
}

// solstice calculates the approximate date of a solstice for a given year.
// season should be "summer" or "winter".
func solstice(year int, season string) time.Time {
	// Interval: 365 days + 5 hours + 47 minutes + 56.5 seconds in milliseconds
	solsticeIntervalMs := (56.5 + 47*60 + 5*3600 + 365*86400) * 1000

	// Base dates (in milliseconds since epoch)
	// Summer solstice base: 2016-06-20T22:34:00Z
	baseSummer := time.Date(2016, 6, 20, 22, 34, 0, 0, time.UTC)
	// Winter solstice base: 2016-12-21T10:44:00Z
	baseWinter := time.Date(2016, 12, 21, 10, 44, 0, 0, time.UTC)

	var baseMs float64
	if season == "winter" {
		baseMs = float64(baseWinter.UnixMilli())
	} else {
		baseMs = float64(baseSummer.UnixMilli())
	}

	timeMs := baseMs + solsticeIntervalMs*float64(year-2016)
	// Truncate to day
	dayMs := float64(24 * 3600 * 1000)
	timeMs = timeMs - math.Mod(timeMs, dayMs)

	return time.UnixMilli(int64(timeMs)).UTC()
}

// calcSpecialDays computes all holidays and special days for a given year.
func calcSpecialDays(year int) []Day {
	bondadagur := findNextWeekDay(year, time.January, 19+rimspillir(year-1), time.Friday)
	easterSunday := easter(year)

	whitsunday := easterSunday.AddDate(0, 0, 49)
	// Check if Whitsunday falls on the first Sunday of June
	withsun1stJuneSun := whitsunday.Month() == time.June && whitsunday.Day() < 8

	// Determine Sjómannadagurinn start day
	sjomannaStartDay := 1
	if withsun1stJuneSun {
		sjomannaStartDay = 8
	}

	days := []Day{
		{
			Date:        time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC),
			Description: "Nýársdagur",
			Key:         KeyNyars,
			Holiday:     true,
		},
		{
			Date:        time.Date(year, time.January, 6, 0, 0, 0, 0, time.UTC),
			Description: "Þrettándinn",
			Key:         KeyThrettand,
			Holiday:     false,
		},
		{
			Date:        bondadagur,
			Description: "Bóndadagur",
			Key:         KeyBonda,
			Holiday:     false,
		},
		{
			Date:        easterSunday.AddDate(0, 0, -48),
			Description: "Bolludagur",
			Key:         KeyBollu,
			Holiday:     false,
		},
		{
			Date:        easterSunday.AddDate(0, 0, -47),
			Description: "Sprengidagur",
			Key:         KeySprengi,
			Holiday:     false,
		},
		{
			Date:        easterSunday.AddDate(0, 0, -46),
			Description: "Öskudagur",
			Key:         KeyOsku,
			Holiday:     false,
		},
		{
			Date:        time.Date(year, time.February, 14, 0, 0, 0, 0, time.UTC),
			Description: "Valentínusardagur",
			Key:         KeyValent,
			Holiday:     false,
		},
		{
			Date:        bondadagur.AddDate(0, 0, 30),
			Description: "Konudagur",
			Key:         KeyKonu,
			Holiday:     false,
		},
		{
			Date:        easterSunday.AddDate(0, 0, -3),
			Description: "Skírdagur",
			Key:         KeySkir,
			Holiday:     true,
		},
		{
			Date:        easterSunday.AddDate(0, 0, -2),
			Description: "Föstudagurinn langi",
			Key:         KeyFosLangi,
			Holiday:     true,
		},
		{
			Date:        easterSunday,
			Description: "Páskadagur",
			Key:         KeyPaska,
			Holiday:     true,
		},
		{
			Date:        easterSunday.AddDate(0, 0, 1),
			Description: "Annar í páskum",
			Key:         KeyPaska2,
			Holiday:     true,
		},
		{
			Date:        findNextWeekDay(year, time.April, 19, time.Thursday),
			Description: "Sumardagurinn fyrsti",
			Key:         KeySumar1,
			Holiday:     true,
		},
		{
			Date:        time.Date(year, time.May, 1, 0, 0, 0, 0, time.UTC),
			Description: "Verkalýðsdagurinn",
			Key:         KeyMai1,
			Holiday:     true,
		},
		{
			Date:        easterSunday.AddDate(0, 0, 39),
			Description: "Uppstigningardagur",
			Key:         KeyUppst,
			Holiday:     true,
		},
		{
			Date:        whitsunday,
			Description: "Hvítasunnudagur",
			Key:         KeyHvitas,
			Holiday:     true,
		},
		{
			Date:        easterSunday.AddDate(0, 0, 50),
			Description: "Annar í Hvítasunnu",
			Key:         KeyHvitas2,
			Holiday:     true,
		},
		{
			Date:        findNextWeekDay(year, time.June, sjomannaStartDay, time.Sunday),
			Description: "Sjómannadagurinn",
			Key:         KeySjomanna,
			Holiday:     false,
		},
		{
			Date:        time.Date(year, time.June, 17, 0, 0, 0, 0, time.UTC),
			Description: "Þjóðhátíðardagurinn",
			Key:         KeyJun17,
			Holiday:     true,
		},
		{
			Date:        solstice(year, "summer"),
			Description: "Sumarsólstöður",
			Key:         KeySumsolst,
			Holiday:     false,
		},
		{
			Date:        time.Date(year, time.June, 24, 0, 0, 0, 0, time.UTC),
			Description: "Jónsmessa",
			Key:         KeyJonsm,
			Holiday:     false,
		},
		{
			Date:        findNextWeekDay(year, time.August, 1, time.Monday),
			Description: "Frídagur verslunarmanna",
			Key:         KeyVerslm,
			Holiday:     true,
		},
		{
			Date:        findNextWeekDay(year, time.October, 21+rimspillir(year), time.Saturday),
			Description: "Fyrsti vetrardagur",
			Key:         KeyVetur1,
			Holiday:     false,
		},
		{
			Date:        time.Date(year, time.October, 31, 0, 0, 0, 0, time.UTC),
			Description: "Hrekkjavaka",
			Key:         KeyHrekkja,
			Holiday:     false,
		},
		{
			Date:        time.Date(year, time.November, 16, 0, 0, 0, 0, time.UTC),
			Description: "Dagur íslenskrar tungu",
			Key:         KeyIslTungu,
			Holiday:     false,
		},
		{
			Date:        time.Date(year, time.December, 1, 0, 0, 0, 0, time.UTC),
			Description: "Fullveldisdagurinn",
			Key:         KeyFullv,
			Holiday:     false,
		},
		{
			Date:        solstice(year, "winter"),
			Description: "Vetrarsólstöður",
			Key:         KeyVetsolst,
			Holiday:     false,
		},
		{
			Date:        time.Date(year, time.December, 23, 0, 0, 0, 0, time.UTC),
			Description: "Þorláksmessa",
			Key:         KeyThorl,
			Holiday:     false,
		},
		{
			Date:        time.Date(year, time.December, 24, 0, 0, 0, 0, time.UTC),
			Description: "Aðfangadagur",
			Key:         KeyAdfanga,
			Holiday:     true,
			HalfDay:     true,
		},
		{
			Date:        time.Date(year, time.December, 25, 0, 0, 0, 0, time.UTC),
			Description: "Jóladagur",
			Key:         KeyJola,
			Holiday:     true,
		},
		{
			Date:        time.Date(year, time.December, 26, 0, 0, 0, 0, time.UTC),
			Description: "Annar í Jólum",
			Key:         KeyJola2,
			Holiday:     true,
		},
		{
			Date:        time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC),
			Description: "Gamlársdagur",
			Key:         KeyGamlars,
			Holiday:     true,
			HalfDay:     true,
		},
	}

	sortDays(days)

	return days
}
