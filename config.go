package main

var Groups = []*Group{
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

var Exceptions = []*Exception{
	// Deakins kids and their spouses cannot match with Deakins
	{"Robbie Shockey", CannotMatchWith, "Turney Deakins"},
	{"Robbie Shockey", CannotMatchWith, "Shirley Deakins"},
	{"David Shockey", CannotMatchWith, "Turney Deakins"},
	{"David Shockey", CannotMatchWith, "Shirley Deakins"},
	{"Lea Anne Barnett", CannotMatchWith, "Turney Deakins"},
	{"Lea Anne Barnett", CannotMatchWith, "Shirley Deakins"},
	{"Frank Barnett", CannotMatchWith, "Turney Deakins"},
	{"Frank Barnett", CannotMatchWith, "Shirley Deakins"},

	// Barnett parents can't give to 3 youngest Maddox kids
	{"Lea Anne Barnett", CannotGiveTo, "MK Maddox"},
	{"Lea Anne Barnett", CannotGiveTo, "Tripp Maddox"},
	{"Lea Anne Barnett", CannotGiveTo, "Hadley Maddox"},
	{"Frank Barnett", CannotGiveTo, "MK Maddox"},
	{"Frank Barnett", CannotGiveTo, "Tripp Maddox"},
	{"Frank Barnett", CannotGiveTo, "Hadley Maddox"},
}
