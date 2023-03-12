package configreader

import (
	"errors"
	"fmt"
	"os"

	"github.com/shockey/secret-santa/rules"
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

	if res := validateNameReferentialIntegrity(doc); len(*res) > 0 {
		errs = append(errs, *res...)
	}

	if len(errs) > 0 {
		fmt.Fprint(os.Stderr, "Your input configuration is invalid :(\nHere's what the validation engine said:\n")
		for _, ve := range errs {
			fmt.Fprintf(os.Stderr, "- %v\n", ve.Error())
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
		for groupName := range group {
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

func validateNameReferentialIntegrity(doc *Document) *[]error {
	errs := []error{}
	knownGroups := map[string]bool{}
	knownPeople := map[string]bool{}

	for _, group := range doc.Groups {
		for groupName, groupContent := range group {
			knownGroups[groupName] = true
			for _, personName := range groupContent.Members {
				knownPeople[personName] = true
			}
		}
	}

	for i, rule := range doc.Rules {
		if rule.NoMatchTo != nil && rule.NoMatchTo.From != nil {
			res := subvalidateEntityMatcherIntegrity(&knownGroups, &knownPeople, rule.NoMatchTo.From, fmt.Sprintf(`rules.%v.from`, i))
			errs = append(errs, res...)
		}

		if rule.NoMatchTo != nil && rule.NoMatchTo.To != nil {
			res := subvalidateEntityMatcherIntegrity(&knownGroups, &knownPeople, rule.NoMatchTo.To, fmt.Sprintf(`rules.%v.to`, i))
			errs = append(errs, res...)
		}

		if rule.NoMatchBetween != nil && rule.NoMatchBetween[0] != nil {
			res := subvalidateEntityMatcherIntegrity(&knownGroups, &knownPeople, rule.NoMatchBetween[0], fmt.Sprintf(`rules.%v.0`, i))
			errs = append(errs, res...)
		}

		if rule.NoMatchBetween != nil && rule.NoMatchBetween[1] != nil {
			res := subvalidateEntityMatcherIntegrity(&knownGroups, &knownPeople, rule.NoMatchBetween[1], fmt.Sprintf(`rules.%v.1`, i))
			errs = append(errs, res...)
		}
	}

	return &errs
}

func subvalidateEntityMatcherIntegrity(knownGroups *map[string]bool, knownPeople *map[string]bool, em *rules.EntityMatcher, locationPrefix string) []error {
	errs := []error{}
	if em.Groups != nil {
		for j, groupName := range *em.Groups {
			if _, exists := (*knownGroups)[groupName]; !exists {
				errs = append(errs, errors.New(fmt.Sprintf(`%v.groups.%v references nonexistent group name "%v"`, locationPrefix, j, groupName)))
			}
		}
	}

	if em.People != nil {
		for j, personName := range *em.People {
			if _, exists := (*knownPeople)[personName]; !exists {
				errs = append(errs, errors.New(fmt.Sprintf(`%v.people.%v references nonexistent person name "%v"`, locationPrefix, j, personName)))
			}
		}
	}

	return errs
}
