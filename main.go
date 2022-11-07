package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
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

type ExceptionType int

const (
	CannotMatchWith = iota // bidirectional
	CannotReceiveFrom
	CannotGiveTo
)

type Exception struct {
	subjectName string
	ruleType    ExceptionType
	targetName  string
}

// TODO: sanity validate to/from names in exceptions against people — protect against misspelling

// CLI flags
var isRealModeFlag = flag.Bool("real", false, "indicates whether filenames will be written as TEST or REAL")

func main() {
	flag.Parse()

	var allGroupedPeople []*GroupedPerson

	for _, group := range Groups {
		groupName := group.name

		for _, person := range group.members {
			allGroupedPeople = append(allGroupedPeople, &GroupedPerson{
				name:      person.name,
				groupName: groupName,
			})
		}
	}

	res := MatchPersons(allGroupedPeople, Exceptions)

	output := "Sender,Recipient\n"

	entries := []string{}

	for _, v := range res {
		entry := fmt.Sprintf("%v,%v", v.sender.name, v.recipient.name)
		entries = append(entries, entry)
	}

	sort.Strings(entries)

	for _, entry := range entries {
		output += entry + "\n"
	}

	dt := time.Now().UTC()
	ts := strings.ReplaceAll(dt.Format(time.RFC3339), ":", "")
	var modestring = "TEST"
	if *isRealModeFlag {
		modestring = "REAL"
	}
	fileName := fmt.Sprintf("Christmas-List-%v-%v.csv", ts, modestring)

	if err := os.WriteFile(fileName, []byte(output), 0666); err != nil {
		log.Fatal(err)
	}
}

func MatchPersons(people []*GroupedPerson, exceptions []*Exception) []*Match {
	senders := people
	shuffleGroupedPersonSlice(&senders)

	recipients := append([]*GroupedPerson{}, people...)
	shuffleGroupedPersonSlice(&recipients)

	matches := []*Match{}

	res, ok := findMatches(recipients, senders, matches, exceptions)

	if !ok {
		panic("Unable to find any resolvable subtree, something is very wrong!")
	}

	return res
}

// findMatches implements a recursive depth-first search strategy to find a solution that results in valid
// matches for all persons
func findMatches(allPeople []*GroupedPerson, pendingSenders []*GroupedPerson, matches []*Match, exceptions []*Exception) ([]*Match, bool) {
	if len(pendingSenders) == 0 {
		// No senders left, we found a good solution set
		return matches, true
	}

	// Always use the next person as the sender
	sender := pendingSenders[0]

	// Find a suitable, unused recipient for the given sender
	for _, person := range allPeople {
		// Check for a match against the static rules
		isMatch := checkForMatch(sender, person, exceptions)

		if !isMatch {
			continue
		}

		// Consider prior search state to disqualify already-tagged recipients
		isAlreadyMatched := false
		for _, match := range matches {
			if match.recipient == person {
				isAlreadyMatched = true
			}
		}

		if isAlreadyMatched {
			continue
		}

		// This person is OK! Continue exploring this subtree...
		newMatches := matches
		newMatches = append(newMatches, &Match{sender, person})
		res, ok := findMatches(allPeople, pendingSenders[1:], newMatches, exceptions)

		if res != nil && ok == true {
			// Found a path that resolves!
			return res, true
		}
	}

	// If we land here, we aren't on a viable path in this subtree
	return matches, false
}

func checkForMatch(sender *GroupedPerson, recipient *GroupedPerson, exceptions []*Exception) bool {
	// Can't match yourself
	if sender == recipient {
		return false
	}

	// Can't match someone in your group
	if sender.groupName == recipient.groupName {
		return false
	}

	// Check exceptions
	for _, exception := range exceptions {
		if exception.ruleType == CannotMatchWith {
			if exception.subjectName == sender.name && exception.targetName == recipient.name {
				return false
			}

			if exception.subjectName == recipient.name && exception.targetName == sender.name {
				return false
			}
		}

		if exception.ruleType == CannotGiveTo {
			if exception.subjectName == sender.name && exception.targetName == recipient.name {
				return false
			}
		}

		if exception.ruleType == CannotReceiveFrom {
			if exception.subjectName == recipient.name && exception.targetName == sender.name {
				return false
			}
		}
	}

	return true
}

func shuffleGroupedPersonSlice(ptr *[]*GroupedPerson) {
	slice := *ptr
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) { slice[i], slice[j] = slice[j], slice[i] })
}
