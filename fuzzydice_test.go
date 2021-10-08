package fuzzydice_test

import (
	"testing"

	"github.com/elliottcarlson/fuzzydice"
)

type ExampleStruct struct {
	ID   int
	Name string
}

func TestFuzzyDiceStruct(t *testing.T) {
	f := fuzzydice.FuzzyDice{}
	f.Load(ExampleStruct{
		Name: "Jack Jonas",
	}, "Name")
	f.Load(ExampleStruct{
		Name: "Marcus Katerin Thompson",
	}, "Name")
	f.Load(ExampleStruct{
		Name: "Catherine Marcus",
	}, "Name")
	f.Load(ExampleStruct{
		Name: "Katheryne Markus",
	}, "Name")

	results := f.Matches("Catherine Marcus")

	doTest := func(index int, expected string) {
		if results[index].(ExampleStruct).Name != expected {
			t.Errorf("Unexpected result at %d: Expected %s, got %s", index, expected, results[0].(ExampleStruct).Name)
		}
	}

	doTest(0, "Catherine Marcus")
	doTest(1, "Katheryne Markus")
	doTest(2, "Marcus Katerin Thompson")
}

func TestFuzzyDiceStructSlice(t *testing.T) {
	countries := []ExampleStruct{
		{
			ID:   0,
			Name: "American Samoa",
		},
		{
			ID:   1,
			Name: "Antigua & Barbuda",
		},
		{
			ID:   2,
			Name: "Bahamas, The",
		},
		{
			ID:   3,
			Name: "Bosnia & Herzegovina",
		},
		{
			ID:   4,
			Name: "Burma",
		},
		{
			ID:   5,
			Name: "Cambodia",
		},
		{
			ID:   6,
			Name: "Cape Verde",
		},
		{
			ID:   7,
			Name: "Central African Rep.",
		},
		{
			ID:   8,
			Name: "China",
		},
		{
			ID:   9,
			Name: "Congo, Repub. of the",
		},
		{
			ID:   10,
			Name: "Cook Islands",
		},
		{
			ID:   11,
			Name: "Cote d'Ivoire",
		},
		{
			ID:   12,
			Name: "East Timor",
		},
		{
			ID:   13,
			Name: "Gambia, The",
		},
		{
			ID:   14,
			Name: "Hong Kong",
		},
		{
			ID:   15,
			Name: "Macau",
		},
		{
			ID:   16,
			Name: "Macedonia",
		},
		{
			ID:   17,
			Name: "Micronesia, Fed. St.",
		},
		{
			ID:   18,
			Name: "New Caledonia",
		},
		{
			ID:   19,
			Name: "Panama",
		},
		{
			ID:   20,
			Name: "Peru",
		},
		{
			ID:   21,
			Name: "Saint Kitts & Nevis",
		},
		{
			ID:   22,
			Name: "Sao Tome & Principe",
		},
		{
			ID:   23,
			Name: "Sierra Leone",
		},
		{
			ID:   24,
			Name: "South Korea",
		},
		{
			ID:   25,
			Name: "Sudan",
		},
		{
			ID:   26,
			Name: "Tonga",
		},
		{
			ID:   27,
			Name: "Trinidad & Tobago",
		},
		{
			ID:   28,
			Name: "Turkey",
		},
		{
			ID:   29,
			Name: "United States",
		},
		{
			ID:   30,
			Name: "Vatican",
		},
		{
			ID:   31,
			Name: "Zambia",
		},
		{
			ID:   32,
			Name: "Zimbabwe",
		},
	}

	f := fuzzydice.FuzzyDice{}
	f.Load(countries, "Name")

	doTest := func(query string, expected interface{}) {
		result, _ := f.BestMatch(query)

		if result == nil {
			if expected != nil {
				t.Errorf("No matches found, when matches expected. Expected %s, got <nil>", expected)
			}
		} else if result.(ExampleStruct).Name != expected {
			t.Errorf("Unexpected match for '%s'. Expected %s, got %s", query, expected, result.(ExampleStruct).Name)
		}
	}

	doTest("Antigua and Barbuda", "Antigua & Barbuda")
	doTest("Bahamas", "Bahamas, The")
	doTest("Bosnia and Herzegovina", "Bosnia & Herzegovina")
	doTest("Cabo Verde", "Cape Verde")
	doTest("Central African Republic", "Central African Rep.")
	doTest("Cook Islands (NZ)", "Cook Islands")
	doTest("Côte d'Ivoire", "Cote d'Ivoire")
	doTest("Democratic Republic of the Congo", "Congo, Repub. of the")
	doTest("Federated States of Micronesia", "Micronesia, Fed. St.")
	doTest("Myanmar/Burma", "Burma")
	doTest("North Macedonia", "Macedonia")
	doTest("Saint Kitts and Nevis", "Saint Kitts & Nevis")
	doTest("South Sudan", "Sudan")
	doTest("São Tomé and Príncipe", "Sao Tome & Principe")
	doTest("Timor-Leste", "East Timor")
	doTest("Trinidad and Tobago", "Trinidad & Tobago")
	doTest("Vatican City State", "Vatican")
	doTest("zzzzzz", nil)
}
