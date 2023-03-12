package rules

// TODO: sumtype instead? https://github.com/BurntSushi/go-sumtype
type Rule struct {
	NoMatchBetween *NoMatchNondirectionalCondition
	NoMatchTo      *NoMatchDirectionalCondition
}

func (r *Rule) DoesPairMatch(senderPersonName string, senderGroupName string, recipientPersonName string, recipientGroupName string) bool {
	// Can't match yourself
	if senderPersonName == recipientPersonName && senderGroupName == recipientGroupName {
		return false
	}

	// Can't match someone in your group
	if senderGroupName == recipientGroupName {
		return false
	}

	if r.NoMatchBetween != nil {
		criteria := r.NoMatchBetween

		doesSenderMatch := criteria[0].DoesPersonMatch(senderPersonName, senderGroupName) || criteria[1].DoesPersonMatch(senderPersonName, senderGroupName)
		doesRecipientMatch := criteria[0].DoesPersonMatch(recipientPersonName, recipientGroupName) || criteria[1].DoesPersonMatch(recipientPersonName, recipientGroupName)

		if doesSenderMatch && doesRecipientMatch {
			return false
		}
	}

	if r.NoMatchTo != nil {
		criteria := r.NoMatchTo

		if criteria.From.DoesPersonMatch(senderPersonName, senderGroupName) && criteria.To.DoesPersonMatch(recipientPersonName, recipientGroupName) {
			return false
		}
	}

	return true
}

type NoMatchNondirectionalCondition [2]*EntityMatcher

type NoMatchDirectionalCondition struct {
	From *EntityMatcher
	To   *EntityMatcher
}

type EntityMatcher struct {
	Groups *[]string
	People *[]string
}

func (e *EntityMatcher) DoesPersonMatch(personName string, groupName string) bool {
	if e.People != nil {
		for _, matchablePerson := range *e.People {
			if matchablePerson == personName {
				return true
			}
		}
	}

	if e.Groups != nil {
		for _, matchableGroup := range *e.Groups {
			if matchableGroup == groupName {
				return true
			}
		}
	}

	return false
}