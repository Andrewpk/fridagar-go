# Frídagar

A Go library for looking up Icelandic public holidays, commonly celebrated
"special" days, and resolving business days before/after a given date.

This is a Go port of [fridagar-node](https://github.com/gaui/fridagar-node).

All returned dates are in the UTC timezone and set to midnight.

## Installation

```bash
go get fridagar
```

## Usage

```go
package main

import (
	"fmt"
	"time"

	"fridagar"
)

func main() {
	// Get all holidays and special days for 2025
	days := fridagar.GetAllDays(2025)
	for _, d := range days {
		fmt.Printf("%s  %s\n", d.Date.Format("2006-01-02"), d.Description)
	}

	// Get only official public holidays
	holidays := fridagar.GetHolidays(2025)

	// Get only unofficial "special" days
	other := fridagar.GetOtherDays(2025)

	// Get all days as a keyed map
	keyed := fridagar.GetAllDaysKeyed(2025)
	fmt.Println(keyed[fridagar.KeyJola].Description) // "Jóladagur"

	// Check if a specific date is a holiday
	day, ok := fridagar.IsHoliday(time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC))
	if ok {
		fmt.Println(day.Description) // "Jóladagur"
	}

	// Check if a date is any special day (holiday or not)
	day, ok = fridagar.IsSpecialDay(time.Date(2025, 3, 3, 0, 0, 0, 0, time.UTC))
	if ok {
		fmt.Println(day.Description) // "Bolludagur"
	}

	// Find the next workday (skipping weekends and holidays)
	workday := fridagar.WorkdaysFromDate(1, time.Date(2025, 12, 24, 0, 0, 0, 0, time.UTC), false)
	fmt.Println(workday.Format("2006-01-02")) // "2025-12-29"

	// Find workday with half-days treated as workdays
	workday = fridagar.WorkdaysFromDate(1, time.Date(2025, 12, 23, 0, 0, 0, 0, time.UTC), true)
	fmt.Println(workday.Format("2006-01-02")) // "2025-12-24" (Aðfangadagur is a half-day)

	_ = holidays
	_ = other
}
```

## API

### `GetAllDays(year int) []Day`

Returns all Icelandic public holidays and commonly celebrated "special" days for a given year, sorted by date.

### `GetAllDaysForMonth(year int, month int) []Day`

Returns all days for a given year and 1-based month (January = 1, December = 12).

### `GetAllDaysKeyed(year int) map[DayKey]Day`

Returns a map of all days for a given year, keyed by their `DayKey`.

### `GetHolidays(year int) []Day`

Returns only official public holidays for a given year.

### `GetHolidaysForMonth(year int, month int) []Day`

Returns official public holidays for a given year and 1-based month.

### `GetOtherDays(year int) []Day`

Returns only unofficial "special" days for a given year.

### `GetOtherDaysForMonth(year int, month int) []Day`

Returns unofficial "special" days for a given year and 1-based month.

### `IsSpecialDay(date time.Time) (Day, bool)`

Checks if a given date is either a holiday or a special day.

### `IsHoliday(date time.Time) (Day, bool)`

Checks if a given date is an official public holiday.

### `WorkdaysFromDate(days int, refDate time.Time, includeHalfDays bool) time.Time`

Returns the date that is the given number of business days before (negative) or after (positive) the reference date. Weekends and official holidays are skipped. If `includeHalfDays` is true, half-day holidays (Aðfangadagur, Gamlársdagur) are treated as workdays.

## Types

### `Day`

```go
type Day struct {
    Date        time.Time // UTC midnight date
    Description string    // Icelandic name
    Key         DayKey    // Stable identifier for translations
    Holiday     bool      // true = official public holiday
    HalfDay     bool      // true = half-day holiday
}
```

### Day Keys

**Holiday keys:** `nyars`, `skir`, `foslangi`, `paska`, `paska2`, `sumar1`, `uppst`, `mai1`, `hvitas`, `hvitas2`, `jun17`, `verslm`, `adfanga`, `jola`, `jola2`, `gamlars`

**Special day keys:** `þrettand`, `bonda`, `bollu`, `sprengi`, `osku`, `valent`, `konu`, `sjomanna`, `sumsolst`, `jonsm`, `vetur1`, `hrekkja`, `isltungu`, `fullv`, `vetsolst`, `thorl`

## Supported Days

| Key | Icelandic Name | Holiday |
|---|---|---|
| `nyars` | Nýársdagur | ✅ |
| `þrettand` | Þrettándinn | |
| `bonda` | Bóndadagur | |
| `bollu` | Bolludagur | |
| `sprengi` | Sprengidagur | |
| `osku` | Öskudagur | |
| `valent` | Valentínusardagur | |
| `konu` | Konudagur | |
| `skir` | Skírdagur | ✅ |
| `foslangi` | Föstudagurinn langi | ✅ |
| `paska` | Páskadagur | ✅ |
| `paska2` | Annar í páskum | ✅ |
| `sumar1` | Sumardagurinn fyrsti | ✅ |
| `mai1` | Verkalýðsdagurinn | ✅ |
| `uppst` | Uppstigningardagur | ✅ |
| `hvitas` | Hvítasunnudagur | ✅ |
| `hvitas2` | Annar í Hvítasunnu | ✅ |
| `sjomanna` | Sjómannadagurinn | |
| `jun17` | Þjóðhátíðardagurinn | ✅ |
| `sumsolst` | Sumarsólstöður | |
| `jonsm` | Jónsmessa | |
| `verslm` | Frídagur verslunarmanna | ✅ |
| `vetur1` | Fyrsti vetrardagur | |
| `hrekkja` | Hrekkjavaka | |
| `isltungu` | Dagur íslenskrar tungu | |
| `fullv` | Fullveldisdagurinn | |
| `vetsolst` | Vetrarsólstöður | |
| `thorl` | Þorláksmessa | |
| `adfanga` | Aðfangadagur | ½ |
| `jola` | Jóladagur | ✅ |
| `jola2` | Annar í Jólum | ✅ |
| `gamlars` | Gamlársdagur | ½ |

## CLI

A simple command-line tool is included:

```bash
go run ./cmd/fridagar [year]
```

## License

ISC

