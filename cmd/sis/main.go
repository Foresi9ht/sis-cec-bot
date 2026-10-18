// Command sis is a small demo of the complaint-processing core: it creates
// several complaints, processes them and prints the results.
package main

import (
	"errors"
	"fmt"

	"github.com/Foresi9ht/sis-cec-bot/internal/complaint"
)

func main() {
	svc := complaint.NewService(complaint.NewKeywordProcessor())

	inputs := []complaint.Input{
		{
			Category:     complaint.CategoryOrg,
			Organization: "Student Union",
			Message:      "Руководство организации игнорирует студентов, есть жалоба на несправедливое решение",
		},
		{
			Category:    complaint.CategoryStaff,
			School:      "ШИТиИ",
			TeacherName: "Иванов Иван Иванович",
			Subject:     "Алгоритмы",
			Message:     "Преподаватель вымогает взятку за пересдачу экзамена",
		},
		{
			Category: complaint.CategoryCafeteria,
			Message:  "Суп в столовой холодный, порции маленькие, меню однообразное",
		},
		{
			Category: complaint.CategoryOther,
			Message:  "Хочу рассказать про ситуацию на кафедре, там всё очень странно",
		},
		{
			Category: complaint.CategoryOther,
			Message:  "коротко", // too short: rejected by validation
		},
	}

	fmt.Println("== Creating and processing complaints ==")
	for _, in := range inputs {
		c, err := svc.Create(in)
		if err != nil {
			fmt.Printf("rejected input: %v\n", err)
			continue
		}
		fmt.Printf("created   #%d [%s] status=%s\n", c.ID, c.Input.Category, c.Status)

		c, err = svc.Process(c.ID)
		if err != nil {
			fmt.Printf("process error: %v\n", err)
			continue
		}
		fmt.Printf("processed #%d status=%s validity=%s priority=%s (valid hits %d, invalid hits %d)\n",
			c.ID, c.Status, c.Result.Validity, c.Result.Priority, c.Result.ValidHits, c.Result.InvalidHits)
	}

	fmt.Println("\n== Error handling ==")
	if _, err := svc.Get(999); errors.Is(err, complaint.ErrNotFound) {
		fmt.Printf("get unknown id: %v\n", err)
	}
	if _, err := svc.Process(1); errors.Is(err, complaint.ErrInvalidTransition) {
		fmt.Printf("process twice:  %v\n", err)
	}

	fmt.Println("\n== All complaints ==")
	for _, c := range svc.List() {
		fmt.Printf("#%d %-9s %-10s %s/%s\n", c.ID, c.Input.Category, c.Status, c.Result.Validity, c.Result.Priority)
	}
}
