package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/shockey/secret-santa/configreader"
	"github.com/shockey/secret-santa/rules"
)

// TODO: tests :) ideally integration tests that fuzz MatchPersons

type Person struct {
	name string
}

type Group struct {
	name    string
	members []*Person
}

type GroupedPerson struct {
	name      string
	groupName string
}

type Match struct {
	sender    *GroupedPerson
	recipient *GroupedPerson
}

const (
	CannotMatchWith = iota // bidirectional
	CannotReceiveFrom
	CannotGiveTo
)

// CLI flags
var isRealModeFlag = flag.Bool("real", false, "indicates whether filenames will be written as TEST or REAL")
var inputNameFlag = flag.String("input", "", "the name of the input config profile to use")

func main() {
	flag.Parse()

	if *inputNameFlag == "" {
		fmt.Fprintf(os.Stderr, "missing required argument flag `input`")
		os.Exit(2) // the same exit code flag.Parse uses
	}

	inputDocument := configreader.MustLoadConfigDocument(*inputNameFlag)

	var allGroupedPeople []*GroupedPerson

	for _, groupRecord := range inputDocument.Groups {
		for groupName, group := range groupRecord {
			for _, personName := range group.Members {
				if _, ok := inputDocument.InactiveGroupMembersMap[[2]string{groupName, personName}]; ok {
					// Skip inactive members
					continue
				}
				allGroupedPeople = append(allGroupedPeople, &GroupedPerson{
					name:      personName,
					groupName: groupName,
				})
			}
		}
	}

	res := MatchPersons(allGroupedPeople, inputDocument.Rules)

	entries := []string{}

	for _, v := range res {
		entry := fmt.Sprintf("%v,%v", v.sender.name, v.recipient.name)
		entries = append(entries, entry)
	}

	slices.Sort(entries)

	var builder strings.Builder
	builder.WriteString("Sender,Recipient\n")
	for _, entry := range entries {
		builder.WriteString(entry)
		builder.WriteString("\n")
	}
	output := builder.String()

	dt := time.Now().UTC()
	ts := strings.ReplaceAll(dt.Format(time.RFC3339), ":", "")
	var modestring = "TEST"
	if *isRealModeFlag {
		modestring = "REAL"
	}
	fileName := fmt.Sprintf("output/List-%v-%v.csv", ts, modestring)

	if err := os.WriteFile(fileName, []byte(output), 0666); err != nil {
		log.Fatal(err)
	}
}

func MatchPersons(people []*GroupedPerson, rules []*rules.Rule) []*Match {
	senders := people
	shuffleGroupedPersonSlice(&senders)

	recipients := append([]*GroupedPerson{}, people...)
	shuffleGroupedPersonSlice(&recipients)

	matches := []*Match{}

	res, ok := findMatches(recipients, senders, matches, rules)

	if !ok {
		panic("Unable to find any resolvable subtree, something is very wrong!")
	}

	return res
}

// findMatches implements a recursive depth-first search strategy to find a solution that results in valid
// matches for all persons
func findMatches(allPeople []*GroupedPerson, pendingSenders []*GroupedPerson, matches []*Match, rules []*rules.Rule) ([]*Match, bool) {
	if len(pendingSenders) == 0 {
		// No senders left, we found a good solution set
		return matches, true
	}

	// Always use the next person as the sender
	sender := pendingSenders[0]

	// Find a suitable, unused recipient for the given sender
	for _, person := range allPeople {
		// Check for a match against the static rules
		var isMatchForAllRules bool = true

		for _, rule := range rules {
			if isMatch := rule.IsPairMatchable(sender.name, sender.groupName, person.name, person.groupName); !isMatch {
				isMatchForAllRules = false
				break
			}
		}

		if !isMatchForAllRules {
			continue
		}

		// Consider prior search state to disqualify already-tagged recipients
		if slices.ContainsFunc(matches, func(m *Match) bool {
			return m.recipient == person
		}) {
			continue
		}

		// This person is OK! Continue exploring this subtree...
		newMatches := matches
		newMatches = append(newMatches, &Match{sender, person})
		res, ok := findMatches(allPeople, pendingSenders[1:], newMatches, rules)

		if res != nil && ok == true {
			// Found a path that resolves!
			return res, true
		}
	}

	// If we land here, we aren't on a viable path in this subtree
	return matches, false
}

func shuffleGroupedPersonSlice(ptr *[]*GroupedPerson) {
	slice := *ptr
	rand.Shuffle(len(slice), func(i, j int) { slice[i], slice[j] = slice[j], slice[i] })
}
