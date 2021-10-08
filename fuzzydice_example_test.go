package fuzzydice_test

import (
	"fmt"

	"github.com/elliottcarlson/fuzzydice"
)

func ExampleFuzzyDice_Load() {
	type Fruit struct {
		ID	     int
		Name	     string
		Translations []string
	}

	type Vegetable struct {
		ID	     int
		Name	     string
		Translations []string
	}

	fruits := []Fruit{
		{
			ID:	      1,
			Name:	      "Apple",
			Translations: []string{"Pomme", "Manzana"},
		},
		{
			ID:	      2,
			Name:	      "Banana",
			Translations: []string{"Banane", "Plátano"},
		},
	}

	vegetables := []Vegetable{
		{
			ID:	      1,
			Name:	      "Broccoli",
			Translations: []string{"Brocoli", "Brócoli"},
		},
		{
			ID:	      2,
			Name:	      "Carrot",
			Translations: []string{"Carotte", "Zanahoria"},
		},
	}

	FuzzyDice := fuzzydice.FuzzyDice{}

	// Load a Slice of Structs at once.
	FuzzyDice.Load(fruits, "Name", "Translations")

	// Load in each Struct individually
	for _, vegetable := range vegetables {
		FuzzyDice.Load(vegetable, "Name", "Translations")
	}

	// Load in an additional Struct
	FuzzyDice.Load(Vegetable{
		ID:	     3,
		Name:	     "Zucchini",
		Translations: []string{"Courgette", "Calabacín"},
	})
}

func ExampleFuzzyDice_BestMatch() {
	type Fruit struct {
		ID	     int
		Name	     string
		Translations []string
	}

	type Vegetable struct {
		ID	     int
		Name	     string
		Translations []string
	}

	fruits := []Fruit{
		{
			ID:	      1,
			Name:	      "Apple",
			Translations: []string{"Pomme", "Manzana"},
		},
		{
			ID:	      2,
			Name:	      "Banana",
			Translations: []string{"Banane", "Plátano"},
		},
	}

	vegetables := []Vegetable{
		{
			ID:	      1,
			Name:	      "Broccoli",
			Translations: []string{"Brocoli", "Brócoli"},
		},
		{
			ID:	      2,
			Name:	      "Carrot",
			Translations: []string{"Carotte", "Zanahoria"},
		},
	}

	FuzzyDice := fuzzydice.FuzzyDice{}
	FuzzyDice.Load(fruits, "Name", "Translations")
	FuzzyDice.Load(vegetables, "Name", "Translations")

	result, similarity := FuzzyDice.BestMatch("plantain")

	if result != nil {
		if fruit, isFruit := result.(Fruit); isFruit {
			fmt.Printf("[Fruit] Best Match: id=%d name=%s similarity=%f\n", fruit.ID, fruit.Name, similarity)
		}
		if vegetable, isVegetable := result.(Vegetable); isVegetable {
			fmt.Printf("[Vegetable] Best Match: id=%d name=%s similarity=%f\n", vegetable.ID, vegetable.Name, similarity)
		}
	}

	// output: [Fruit] Best Match: id=2 name=Banana similarity=0.428571
}

func ExampleFuzzyDice_Matches() {
	type Employee struct {
		Name	     string
	}

	employees := []Employee{
		{
			Name: "Jack Jonas",
		},
		{
			Name: "Marcus Katerin Thompson",
		},
		{
			Name: "Catherine Marcus",
		},
		{
			Name: "Katheryne Markus",
		},
	}

	FuzzyDice := fuzzydice.FuzzyDice{}
	FuzzyDice.Load(employees, "Name")

	results := FuzzyDice.Matches("Kathy Marx")

	for _, result := range results {
		if employee, isEmployee := result.(Employee); isEmployee {
			fmt.Printf("Employee: name=%s\n", employee.Name)
		}
	}

	// output: Employee: name=Katheryne Markus
	// Employee: name=Catherine Marcus
	// Employee: name=Marcus Katerin Thompson
}

func ExampleFuzzyDice_Rank() {
	type Server struct {
		IP   string
		Host string
	}

	servers := []Server{
		{
			IP:	      "192.168.1.2",
			Host:	      "prod.us-east01.example.com",
		},
		{
			IP:	      "192.168.2.2",
			Host:	      "prod.us-west01.example.com",
		},
		{
			IP:	      "192.168.1.3",
			Host:	      "prod.us-east02.example.com",
		},
		{
			IP:	      "192.168.2.3",
			Host:	      "prod.us-west02.example.com",
		},
		{
			IP:	      "192.168.1.4",
			Host:	      "prod.us-east03.example.com",
		},
		{
			IP:	      "192.168.2.4",
			Host:	      "prod.us-west03.example.com",
		},
		{
			IP:	      "192.168.1.5",
			Host:	      "prod.us-east04.example.com",
		},
		{
			IP:	      "192.168.2.5",
			Host:	      "prod.us-west04.example.com",
		},
		{
			IP:	      "192.168.1.6",
			Host:	      "prod.us-east05.example.com",
		},
		{
			IP:	      "192.168.2.6",
			Host:	      "prod.us-west05.example.com",
		},
	}

	FuzzyDice := fuzzydice.FuzzyDice{}
	FuzzyDice.Load(servers, "Host", "IP")

	results := FuzzyDice.Rank("prod east 05")

	for _, result := range results {
		if server, isServer := result.Source.(Server); isServer {
			fmt.Printf("Server: IP=%s Host=%s Similarity=%f\n", server.IP, server.Host, result.Rank)
		}
	}

	// output: Server: IP=192.168.1.6 Host=prod.us-east05.example.com Similarity=0.470588
	// Server: IP=192.168.1.2 Host=prod.us-east01.example.com Similarity=0.411765
	// Server: IP=192.168.1.4 Host=prod.us-east03.example.com Similarity=0.411765
	// Server: IP=192.168.1.5 Host=prod.us-east04.example.com Similarity=0.411765
	// Server: IP=192.168.1.3 Host=prod.us-east02.example.com Similarity=0.411765
	// Server: IP=192.168.2.6 Host=prod.us-west05.example.com Similarity=0.352941
	// Server: IP=192.168.2.2 Host=prod.us-west01.example.com Similarity=0.294118
	// Server: IP=192.168.2.4 Host=prod.us-west03.example.com Similarity=0.294118
	// Server: IP=192.168.2.5 Host=prod.us-west04.example.com Similarity=0.294118
	// Server: IP=192.168.2.3 Host=prod.us-west02.example.com Similarity=0.294118
}
