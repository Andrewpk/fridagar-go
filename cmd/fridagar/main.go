// Command fridagar prints Icelandic holidays for a given year.
//
// Usage:
//
//	fridagar [year]
//
// If no year is given, the current year is used.
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Andrewpk/fridagar-go"
)

func main() {
	year := time.Now().Year()

	if len(os.Args) > 1 {
		y, err := strconv.Atoi(os.Args[1])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "invalid year: %s\n", os.Args[1])
			os.Exit(1)
		}
		year = y
	}

	days := fridagar.GetAllDays(year)
	fmt.Printf("Icelandic holidays and special days for %d:\n\n", year)

	for _, d := range days {
		holiday := " "
		if d.Holiday {
			holiday = "*"
			if d.HalfDay {
				holiday = "½"
			}
		}
		fmt.Printf("  %s %s  %-30s [%s]\n",
			holiday,
			d.Date.Format("2006-01-02"),
			d.Description,
			d.Key)
	}

	fmt.Printf("\n  * = official public holiday, ½ = half-day holiday\n")
}
