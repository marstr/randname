package randname

import "fmt"

func ExampleAdjNoun_Generate() {
	generator := AdjNoun{}

	fmt.Println(generator.Generate())
	// Output: Foo
}
