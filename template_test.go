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

package template_test

import (
	"log"
	"os"

	"github.com/b3log/freemarker.go"
)

func ExampleTemplate() {
	const ftl = "It's a simple ${thing}."

	t, err := template.New("ftl").Parse(ftl)
	if nil != err {
		panic(err)
	}

	dataModel := map[string]interface{}{}
	dataModel["thing"] = "template engine"
	err = t.Execute(os.Stdout, dataModel)
	if err != nil {
		log.Println("executing template:", err)
	}

	// Output:
	// It's a simple template engine.
}
