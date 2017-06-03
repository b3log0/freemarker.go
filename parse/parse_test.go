// freemarker.go - FreeMarker template engine in golang.
// Copyright (C) 2017, b3log.org & hacpai.com
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package parse

import (
	"flag"
	"fmt"
	"testing"
)

func mkItem(typ itemType, text string) item {
	return item{
		typ: typ,
		val: text,
	}
}

func TestStack(t *testing.T) {
	e1 := mkItem(itemNumber, "1")
	e2 := mkItem(itemAdd, "+")
	e3 := mkItem(itemNumber, "2")

	s := &stack{}
	s.push(&e1)
	s.push(&e2)
	s.push(&e3)

	if "2" != s.pop().(*item).val {
		t.Log("unexpected stack item")
	}

	if "+" != s.pop().(*item).val {
		t.Log("unexpected stack item")
	}

	if "1" != s.peek().(*item).val {
		t.Log("unexpected stack item")
	}

	if "1" != s.pop().(*item).val {
		t.Log("unexpected stack item")
	}
}

var debug = flag.Bool("debug", true, "show the errors produced by the main tests")

type parseTest struct {
	name   string
	input  string
	ok     bool
	result string // what the user would see in an error message.
}

const (
	noError  = true
	hasError = false
)

var parseTests = []parseTest{
	//	{"empty", "", noError, ``},
	//	{"comment", "hello-<#--\n\n\n-->-world", noError, `"hello-""-world"`},
	//	{"spaces", " \t\n", noError, `" \t\n"`},
	//	{"text", "some text", noError, `"some text"`},
	{"interpolation iden", "hello${abc}world", noError, `"hello"${abc}"world"`},
	{"interpolation iden.", "hello${a.b}world", noError, `"hello"${a.b}"world"`},
	//	{"true if", "<#if true></#if>", noError, `<#if true></#if>`},
	//	{"simple if", "<#if a == b>true content</#if>following content", noError,
	//		`<#if a==b>"true content"</#if>"following content"`},
}

func testParse(doCopy bool, t *testing.T) {
	textFormat = "%q"
	defer func() { textFormat = "%s" }()
	for _, test := range parseTests {
		tmpl, err := New(test.name).Parse(test.input, make(map[string]*Tree))
		switch {
		case err == nil && !test.ok:
			t.Errorf("%q: expected error; got none", test.name)
			continue
		case err != nil && test.ok:
			t.Errorf("%q: unexpected error: %v", test.name, err)
			continue
		case err != nil && !test.ok:
			// expected error, got one
			if *debug {
				fmt.Printf("%s: %s\n\t%s\n", test.name, test.input, err)
			}
			continue
		}
		var result string
		if doCopy {
			result = tmpl.Root.Copy().String()
		} else {
			result = tmpl.Root.String()
		}
		if result != test.result {
			t.Errorf("%s=(%q): got\n\t%v\nexpected\n\t%v", test.name, test.input, result, test.result)
		}
	}
}

func TestParse(t *testing.T) {
	testParse(false, t)
}

type isEmptyTest struct {
	name  string
	input string
	empty bool
}

var isEmptyTests = []isEmptyTest{
	{"empty", "", true},
	{"nonempty", "hello", false},
	{"spaces only", " \t\n \t\n", true},
}

func TestIsEmpty(t *testing.T) {
	if !IsEmptyTree(nil) {
		t.Errorf("nil tree is not empty")
	}
	for _, test := range isEmptyTests {
		tree, err := New("root").Parse(test.input, make(map[string]*Tree))
		if err != nil {
			t.Errorf("%q: unexpected error: %v", test.name, err)
			continue
		}
		if empty := IsEmptyTree(tree.Root); empty != test.empty {
			t.Errorf("%q: expected %t got %t", test.name, test.empty, empty)
		}
	}
}

func TestErrorContextWithTreeCopy(t *testing.T) {
	tree, err := New("root").Parse("{{if true}}{{end}}", make(map[string]*Tree))
	if err != nil {
		t.Fatalf("unexpected tree parse failure: %v", err)
	}
	treeCopy := tree.Copy()
	wantLocation, wantContext := tree.ErrorContext(tree.Root.Nodes[0])
	gotLocation, gotContext := treeCopy.ErrorContext(treeCopy.Root.Nodes[0])
	if wantLocation != gotLocation {
		t.Errorf("wrong error location want %q got %q", wantLocation, gotLocation)
	}
	if wantContext != gotContext {
		t.Errorf("wrong error location want %q got %q", wantContext, gotContext)
	}
}
