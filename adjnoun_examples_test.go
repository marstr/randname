package randname_test

import (
	"fmt"

	"github.com/marstr/randname"
)

func ExampleAdjNoun_Generate() {
	generator := randname.NewAdjNoun()
	fmt.Println(generator.Generate())
}

func ExampleGenerate() {
	fmt.Println(randname.Generate())
}
