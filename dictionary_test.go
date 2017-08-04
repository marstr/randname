package randname

import (
	"fmt"
	"strings"
	"testing"

	"github.com/marstr/collection"
)

func ExampleDictionary_Add() {
	subject := &Dictionary{}

	const example = "hello"
	fmt.Println(subject.Contains(example))
	fmt.Println(subject.Size())
	subject.Add(example)
	fmt.Println(subject.Contains(example))
	fmt.Println(subject.Size())
	// Output:
	// false
	// 0
	// true
	// 1
}

func ExampleDictionary_Enumerate() {
	subject := Dictionary{}
	subject.Add("hello")
	subject.Add("world")

	upperCase := collection.Select(subject, func(x interface{}) interface{} {
		return strings.ToUpper(x.(string))
	})

	for word := range subject.Enumerate(nil) {
		fmt.Println(word)
	}

	for word := range upperCase.Enumerate(nil) {
		fmt.Println(word)
	}
	// Output:
	// hello
	// world
	// HELLO
	// WORLD
}

func TestDictionary_Add(t *testing.T) {
	subject := Dictionary{}

	subject.Add("word")

	if rootChildrenCount := len(subject.root.Children); rootChildrenCount != 1 {
		t.Logf("The root should only have one child, got %d instead.", rootChildrenCount)
		t.Fail()
	}

	if retreived, ok := subject.root.Children['w']; ok {
		leaf := retreived.Navigate("ord")
		if leaf == nil {
			t.Log("Unable to navigate from `w`")
			t.Fail()
		} else if !leaf.IsWord {
			t.Log("leaf shoud have been a word")
			t.Fail()
		}
	} else {
		t.Log("Root doesn't have child for `w`")
		t.Fail()
	}
}

func TestTrieNode_Navigate(t *testing.T) {
	leaf := trieNode{
		IsWord: true,
	}
	subject := trieNode{
		Children: map[rune]*trieNode{
			'a': &trieNode{
				Children: map[rune]*trieNode{
					'b': &trieNode{
						Children: map[rune]*trieNode{
							'c': &leaf,
						},
					},
				},
			},
		},
	}

	testCases := []struct {
		address  string
		expected *trieNode
	}{
		{"abc", &leaf},
		{"abd", nil},
		{"", &subject},
		{"a", subject.Children['a']},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			if result := subject.Navigate(tc.address); result != tc.expected {
				t.Logf("got: %v want: %v", result, tc.expected)
				t.Fail()
			}
		})
	}
}
