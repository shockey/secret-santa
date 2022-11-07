package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"
)

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

var Data = []*Group{
	{"Shockey", []*Person{
		{"Kyle Shockey"},
		{"Erin Shockey"},
		{"Rachel Lawrimore"},
		{"Blaze Lawrimore"},
		{"Robbie Shockey"},
		{"David Shockey"},
		// {"Chase Foster"},
	}},
	{"Barnett", []*Person{
		{"Austin Barnett"},
		{"Devan Barnett"},
		{"Alyssa Barnett"},
		{"Lea Anne Barnett"},
		{"Frank Barnett"},
	}},
	{"Maddox", []*Person{
		{"Steve Maddox"},
		{"Candice Maddox"},
		{"Tripp Maddox"},
		{"MK Maddox"},
		{"Hadley Maddox"},
		// {"Anne Claire Gray"},
		// {"David Gray"},
	}},
	{"Deakins", []*Person{
		{"Turney Deakins"},
		{"Shirley Deakins"},
	}},
}

var cannotMatchDeakins = []string{
	"Robbie Shockey",
	"David Shockey",
	"Lea Anne Barnett",
	"Frank Barnett",
}

func main() {
	var allGroupedPeople []*GroupedPerson

	for _, group := range Data {
		groupName := group.name

		for _, person := range group.members {
			allGroupedPeople = append(allGroupedPeople, &GroupedPerson{
				name:      person.name,
				groupName: groupName,
			})
		}
	}

	res := MatchPersons(allGroupedPeople)

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

	dt := time.Now()
	fileName := fmt.Sprintf("Christmas-List-%v.csv", dt.Format(time.RFC3339))

	if err := os.WriteFile(fileName, []byte(output), 0666); err != nil {
		log.Fatal(err)
	}
}

func MatchPersons(people []*GroupedPerson) []*Match {
	senders := people
	shuffleGroupedPersonSlice(&senders)

	recipients := append([]*GroupedPerson{}, people...)
	shuffleGroupedPersonSlice(&recipients)

	matches := []*Match{}

	res, ok := findMatches(recipients, senders, matches)

	if !ok {
		panic("Unable to find any resolvable subtree, something is very wrong!")
	}

	return res
}

// findMatches implements a recursive depth-first search strategy to find a solution that results in valid
// matches for all persons
func findMatches(allPeople []*GroupedPerson, pendingSenders []*GroupedPerson, matches []*Match) ([]*Match, bool) {
	if len(pendingSenders) == 0 {
		// No senders left, we found a good solution set
		return matches, true
	}

	// Always use the next person as the sender
	sender := pendingSenders[0]

	// Find a suitable, unused recipient for the given sender
	for _, person := range allPeople {
		// Check for a match against the static rules
		isMatch := checkForMatch(sender, person)

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
		res, ok := findMatches(allPeople, pendingSenders[1:], newMatches)

		if res != nil && ok == true {
			// Found a path that resolves!
			return res, true
		}
	}

	// If we land here, we aren't on a viable path in this subtree
	return matches, false
}

func checkForMatch(sender *GroupedPerson, recipient *GroupedPerson) bool {
	// Can't match yourself
	if sender == recipient {
		return false
	}

	// Can't match someone in your group
	if sender.groupName == recipient.groupName {
		return false
	}

	// Special case: Frank, Lea, Robbie, David can't match a Deakins in any direction
	if sender.groupName == "Deakins" || recipient.groupName == "Deakins" {
		for _, name := range cannotMatchDeakins {
			if sender.name == name || recipient.name == name {
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
