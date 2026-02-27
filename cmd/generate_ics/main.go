package main

import (
	"fmt"
	"log"
	"os"
	"time"

	ics "github.com/arran4/golang-ical"

	fridagar "github.com/Andrewpk/fridagar-go"
)

func main() {
	now := time.Now()
	currentYear := now.Year()
	years := []int{currentYear - 1, currentYear, currentYear + 1}

	cal := ics.NewCalendarFor("fridagar-go")
	cal.SetName("Icelandic Holidays")
	cal.SetDescription("Icelandic public holidays and special days")
	cal.SetTimezoneId("UTC")
	cal.SetRefreshInterval("P1W")
	cal.SetMethod(ics.MethodPublish)

	for _, year := range years {
		days := fridagar.GetAllDays(year)
		for _, day := range days {
			event := cal.AddEvent(fmt.Sprintf("%s-%d@fridagar-go", day.Key, year))
			event.SetSummary(day.Description)
			event.SetAllDayStartAt(day.Date)
			// All-day events end the next day in ICS (exclusive end date)
			event.SetAllDayEndAt(day.Date.AddDate(0, 0, 1))
			event.SetDtStampTime(now)
			event.SetCreatedTime(now)

			var status string
			if day.Holiday {
				if day.HalfDay {
					status = "Opinber hátíðardagur (hálfdagur)"
				} else {
					status = "Opinber hátíðardagur"
				}
			} else {
				status = "Sérstakur dagur"
			}
			event.SetDescription(status)
			event.SetClass(ics.ClassificationPublic)
		}
	}

	outputFile := "pages/static/icelandic_holidays.ics"
	if len(os.Args) > 1 {
		outputFile = os.Args[1]
	}
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Printf("warning: failed to close output file: %v", cerr)
		}
	}()

	if err := cal.SerializeTo(f); err != nil {
		log.Fatalf("failed to serialize calendar: %v", err)
	}

	fmt.Printf("Calendar written to %s\n", outputFile)
	fmt.Printf("Years included: %d, %d, %d\n", years[0], years[1], years[2])
}
