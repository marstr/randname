package randname_test

import (
	"fmt"

	"github.com/marstr/randname"
)

func ExampleAdjNoun_Generate() {
	var generator randname.AdjNoun
	fmt.Println(generator.Generate())
}

func ExampleGenerate() {
	fmt.Println(randname.Generate())
}

func ExampleGeneratePascalCase() {
	fmt.Println(randname.GeneratePascalCase())
}

func ExampleGenerateCamelCase() {
	fmt.Println(randname.GenerateCamelCase())
}

func ExampleGenerateHyphened() {
	fmt.Println(randname.GenerateHyphened())
}

func ExampleAdjNoun_GeneratePascalCase() {
	var generator randname.AdjNoun
	fmt.Println(generator.GeneratePascalCase())
}

func ExampleAdjNoun_GenerateCamelCase() {
	var generator randname.AdjNoun
	fmt.Println(generator.GenerateCamelCase())
}

func ExampleAdjNoun_GenerateHyphened() {
	var generator randname.AdjNoun
	fmt.Println(generator.GenerateHyphened())
}
