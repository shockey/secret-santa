package config

import (
	"errors"
	"fmt"
	"os"
)

// mustValidateConfigDocument validates an already-parsed input document with custom validator logic.
func mustValidateConfigDocument(doc *Document) {
	var errs []error
	if res := validatePersonNameUniqueness(doc); len(*res) > 0 {
		errs = append(errs, *res...)
	}

	if res := validateGroupUniqueness(doc); len(*res) > 0 {
		errs = append(errs, *res...)
	}

	if len(errs) > 0 {
		fmt.Fprint(os.Stderr, "Your input configuration is invalid :(\nHere's what the validation engine said:\n")
		for i, ve := range errs {
			fmt.Fprintf(os.Stderr, "%v: %v\n", i, ve.Error())
		}
		os.Exit(1)
	}
}

func validatePersonNameUniqueness(doc *Document) *[]error {
	errs := []error{}
	seen := make(map[string][]string)

	for _, group := range doc.Groups {
		for groupName, groupContent := range group {
			for _, personName := range groupContent.Members {
				seen[personName] = append(seen[personName], groupName)
			}
		}
	}

	for personName, seenIn := range seen {
		if len(seenIn) > 1 {
			errs = append(errs, errors.New(fmt.Sprintf(`the name "%v" appears more than once (seen %v times, in groups %v)`, personName, len(seenIn), seenIn)))
		}
	}

	return &errs
}

func validateGroupUniqueness(doc *Document) *[]error {
	errs := []error{}
	seen := make(map[string]int)

	for _, group := range doc.Groups {
		for groupName, _ := range group {
			seen[groupName]++
		}
	}

	for groupName, count := range seen {
		if count > 1 {
			errs = append(errs, errors.New(fmt.Sprintf(`the group "%v" appears more than once (seen %v times)`, groupName, count)))
		}
	}

	return &errs
}
