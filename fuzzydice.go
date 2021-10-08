/*
Package fuzzydice provides a mechanism for searching fields within structures using the Sørensen–Dice coefficient.
*/
package fuzzydice

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// FuzzyDice describes a FuzzyDice instance with loaded objects. A single FuzzyDice instance can contain structures of
// varying types, as long as they share common fields for searching.
type FuzzyDice struct {
	objects []object
}

// Load will load searchable data into the instance. Searchable data can either be a singular Struct, or a Slice of
// Structs. Not all objects fed into the search need to be of the same type, however, they must share common field names
// that can be searched on.
func (f *FuzzyDice) Load(source interface{}, fields ...string) error {
	value := reflect.ValueOf(source)

	switch value.Kind() {
	case reflect.Struct:
		for _, field := range fields {
			if !value.FieldByName(field).IsValid() {
				panic("Field doesn't exist " + field)
			}
		}

		f.objects = append(f.objects, object{
			source: source,
			fields: fields,
		})
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			item := value.Index(i)

			if item.Kind() == reflect.Struct {
				f.Load(item.Interface(), fields...)
			}
		}
	default:
		return errors.New("Unsupported type provided.")
	}

	return nil
}

// BestMatch will find the top ranked match for a given query. It will return the matching interface{} and a float32
// similiarity coefficient. If there are no matches available, the match will return nil as the interface{}.
func (f *FuzzyDice) BestMatch(query string) (interface{}, float32) {
	results := f.Rank(query)

	if len(results) > 0 {
		return results[0].source, results[0].rank
	}

	return nil, 0
}

// Matches will return all matches for a given query. The resulting interface{} will contain each of the matching
// objects that were loaded in that were matched.
func (f *FuzzyDice) Matches(query string) []interface{} {
	ranks := f.Rank(query)

	results := make([]interface{}, len(ranks))
	for i, r := range ranks {
		results[i] = r.source
	}

	return results
}

// Rank will perform the search for a given query. It will return a []rankedObject that contains a reference to the
// source object that matched, as well as the similarity coefficient of the match to the query.
func (f *FuzzyDice) Rank(query string) []rankedObject {
	ranks := []rankedObject{}
	for _, o := range f.objects {
		distance := o.rank(normalize(query))
		if distance <= 0 {
			continue
		}

		ranks = append(ranks, rankedObject{
			source: o.source,
			rank:   distance,
		})
	}

	if len(ranks) == 0 {
		return []rankedObject{}
	}

	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i].rank > ranks[j].rank
	})

	return ranks
}

type object struct {
	source interface{}
	fields []string
}

type rankedObject struct {
	source interface{}
	rank   float32
}

func (o object) rank(query string) float32 {
	var highestRank float32
	highestRank = -1
	for _, field := range o.fields {
		for _, value := range valuesForField(o, field) {
			rank := calcCoefficient(query, value)
			if rank > highestRank {
				highestRank = rank
			}
		}
	}
	return highestRank
}

func calcCoefficient(source, target string) float32 {
	if value := returnEarlyIfPossible(source, target); value >= 0 {
		return value
	}

	bigrams := make(map[string]int)
	for i := 0; i < len(source)-1; i++ {
		var count int
		bigram := makeBigram(source, i)
		if value, ok := bigrams[bigram]; ok {
			count = value + 1
		} else {
			count = 1
		}

		bigrams[bigram] = count
	}

	var intersectionSize float32
	intersectionSize = 0

	for i := 0; i < len(target)-1; i++ {
		var count int
		bigram := makeBigram(target, i)
		if value, ok := bigrams[bigram]; ok {
			count = value
		} else {
			count = 0
		}

		if count > 0 {
			bigrams[bigram] = count - 1
			intersectionSize = intersectionSize + 1
		}
	}

	return (2.0 * intersectionSize) / (float32(len(source)) + float32(len(target)) - 2)
}

func makeBigram(source string, index int) string {
	a := fmt.Sprintf("%c", source[index])
	b := fmt.Sprintf("%c", source[index+1])

	return a + b
}

func normalize(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "", -1)
}

func returnEarlyIfPossible(source, target string) float32 {
	if len(source) == 0 && len(target) == 0 {
		return 1
	}

	if len(source) == 0 || len(target) == 0 {
		return 0
	}

	if source == target {
		return 1
	}

	if len(source) == 1 && len(target) == 1 {
		return 0
	}

	if len(source) < 2 || len(target) < 2 {
		return 0
	}

	return -1
}

func valuesForField(o object, fieldName string) []string {
	field := reflect.ValueOf(o.source).FieldByName(fieldName)

	if strSlice, isStrSlice := field.Interface().([]string); isStrSlice {
		lc := make([]string, len(strSlice))
		for i, s := range strSlice {
			lc[i] = normalize(s)
		}

		return lc
	}

	return []string{normalize(field.String())}
}
