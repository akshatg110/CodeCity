/* Copyright 2017 Google Inc.
 * https://github.com/NeilFraser/CodeCity
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package interpreter implements a JavaScript interpreter.
package interpreter

import (
	"fmt"

	"CodeCity/server/interpreter/ast"
	"CodeCity/server/interpreter/object"
)

// Interpreter implements a JavaScript interpreter.
type Interpreter struct {
	state   state
	value   *cval
	Verbose bool
}

// New takes a JavaScript program, in the form of an JSON-encoded
// ESTree, and creates a new Interpreter that will execute that
// program.
func New(astJSON string) *Interpreter {
	var this = new(Interpreter)

	tree, err := ast.NewFromJSON(astJSON)
	if err != nil {
		panic(err)
	}
	s := newScope(nil, nil)
	initArrayProto(s)
	// FIXME: insert global names into s
	s.populate(tree)
	this.state = newState(nil, s, tree)
	return this
}

// Step performs the next step in the evaluation of program.  Returns
// true if a step was executed; false if the program has terminated.
func (intrp *Interpreter) Step() bool {
	if intrp.state == nil {
		switch intrp.value.typ {
		case BREAK:
			panic(fmt.Errorf("illegal break to %s", intrp.value.targ))
		case CONTINUE:
			panic(fmt.Errorf("illegal continue of %s", intrp.value.targ))
		case RETURN:
			panic(fmt.Errorf("illegal return of %s", intrp.value.val))
		case THROW:
			panic(fmt.Errorf("unhandled exception %s", intrp.value.val))
		}
		return false
	}
	if intrp.Verbose {
		fmt.Printf("Next step is %T.step(%#v)\n", intrp.state, intrp.value)
	}
	intrp.state, intrp.value = intrp.state.step(intrp.value)
	return true
}

// Run runs the program to completion.
func (intrp *Interpreter) Run() {
	for intrp.Step() {
	}
}

// Value returns the final value computed by the last statement
// expression of the program.
func (intrp *Interpreter) Value() object.Value {
	if intrp.value == nil {
		return nil
	}
	return intrp.value.val
}
