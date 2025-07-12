package evaluator_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/evaluator"
	"github.com/dev-kas/virtlang-go/v4/parser"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
)

func TestNumbers(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "1",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(1),
			},
		},
		{
			input: "123",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(123),
			},
		},
		{
			input: "0",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(0),
			},
		},
		{
			input: "0.1",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(0.1),
			},
		},
		{
			input: "3.14159",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3.14159),
			},
		},
	}

	for _, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("expected no error, got %v", synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("expected no error, got %v", runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("expected %v, got %v", test.output.Type, evaluated.Type)
		}
		if evaluated.Value != test.output.Value {
			t.Errorf("expected %v, got %v", test.output.Value, evaluated.Value)
		}
	}
}

func TestStrings(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: `"hello"`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "hello",
			},
		},
		{
			input: `"hello world"`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "hello world",
			},
		},
		{
			input: `""`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "",
			},
		},
		{
			input: `"123"`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "123",
			},
		},
		{
			input: `'hello world'`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "hello world",
			},
		},
		{
			input: `''`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "",
			},
		},
		{
			input: `'123'`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "123",
			},
		},
	}

	for _, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("expected no error, got %v", synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("expected no error, got %v", runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("expected %v, got %v", test.output.Type, evaluated.Type)
		}
		if evaluated.Value != test.output.Value {
			t.Errorf("expected %v, got %v", test.output.Value.(string), evaluated.Value.(string))
		}
	}
}

func TestObjects(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: `{foo: "bar"}`,
			output: shared.RuntimeValue{
				Type: shared.Object,
				Value: map[string]*shared.RuntimeValue{
					"foo": {
						Type:  shared.String,
						Value: "bar",
					},
				},
			},
		},
		{
			input: `{foo: "bar", bar: "foo"}`,
			output: shared.RuntimeValue{
				Type: shared.Object,
				Value: map[string]*shared.RuntimeValue{
					"foo": {
						Type:  shared.String,
						Value: "bar",
					},
					"bar": {
						Type:  shared.String,
						Value: "foo",
					},
				},
			},
		},
		{
			input: `{}`,
			output: shared.RuntimeValue{
				Type:  shared.Object,
				Value: map[string]*shared.RuntimeValue{},
			},
		},
		{
			input: `{foo: 123}`,
			output: shared.RuntimeValue{
				Type: shared.Object,
				Value: map[string]*shared.RuntimeValue{
					"foo": {
						Type:  shared.Number,
						Value: float64(123),
					},
				},
			},
		},
		{
			input: `{foo: {bar: {bazz: 123}}}`,
			output: shared.RuntimeValue{
				Type: shared.Object,
				Value: map[string]*shared.RuntimeValue{
					"foo": {
						Type: shared.Object,
						Value: map[string]*shared.RuntimeValue{
							"bar": {
								Type: shared.Object,
								Value: map[string]*shared.RuntimeValue{
									"bazz": {
										Type:  shared.Number,
										Value: float64(123),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "let error = 'error not triggered'\ntry {undefinedVariable} catch e {error = e}\nerror.message",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "Cannot resolve variable `undefinedVariable`",
			},
		},
		{
			input: "let obj = { class: 'legendary' }\nobj.fn = 'Function'\nobj.if = 'Branch'\nobj",
			output: shared.RuntimeValue{
				Type: shared.Object,
				Value: map[string]*shared.RuntimeValue{
					"class": {
						Type:  shared.String,
						Value: "legendary",
					},
					"fn": {
						Type:  shared.String,
						Value: "Function",
					},
					"if": {
						Type:  shared.String,
						Value: "Branch",
					},
				},
			},
		},
	}

	for _, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("expected no error, got %v", synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("expected no error, got %v", runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("expected %v, got %v", test.output.Type, evaluated.Type)
		}
		if !reflect.DeepEqual(evaluated.Value, test.output.Value) {
			t.Errorf("value mismatch. expected %#v, got %#v", test.output.Value, evaluated.Value)
		}
	}
}

func TestBinaryExpression(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "1 + 2",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "1 - 2",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(-1),
			},
		},
		{
			input: "1 * 2",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(2),
			},
		},
		{
			input: "10 / 2",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(5),
			},
		},
		{
			input: "10 % 3",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(1),
			},
		},
		{
			input: "3 * (7 + 1)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(24),
			},
		},
		{
			input: "3 * (7 + 1) + 2",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(26),
			},
		},
		{
			input: "1.5 + 2.5",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(4.0),
			},
		},
		{
			input: "5.5 - 2.2",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3.3),
			},
		},
		{
			input: "3.3 * 2.0",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(6.6),
			},
		},
		{
			input: "7.5 / 2.5",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3.0),
			},
		},
		{
			input: "let x = 1.1 let y = 2.2 x + y",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3.3),
			},
		},
		{
			input: "2.5 * (1.5 + 2.5)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10.0),
			},
		},
		{
			input: "10.0 - 3.5",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(6.5),
			},
		},
		{
			input: "4.4 / 2.2",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(2.0),
			},
		},
	}

	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("[%d] expected no error, got %v", i+1, synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("[%d] expected no error, got %v", i+1, runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("[%d] expected %v, got %v", i+1, test.output.Type, evaluated.Type)
		}
		if evaluated.Value != test.output.Value {
			t.Errorf("[%d] expected %v, got %v", i+1, test.output.Value, evaluated.Value)
		}
	}
}
func TestComparisonOperators(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "1 == 1",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "1 != 2",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "1 < 2",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "2 > 1",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "2 <= 2",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "3 >= 2",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "1 == 2",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "1 > 2",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "2 < 1",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "2 >= 3",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "2 <= 1",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "(1==1) == (1==1)",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "(1==1) == (1!=1)",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "{} == []",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "(1!=1) != (1!=1)",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "(1!=1) != (1==1)",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "{} != []",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "{} == {}",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "let a = []\na == a",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let a = {}\n a == a",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let a = []\na != a",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "fn a(){}\nlet b=a\na == b",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let a = fn(){}\na == a",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "(fn(){}) == (fn(){})",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "'string' == 'another'",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "'string' == 'string'",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
	}

	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("[%d] expected no error, got %v, input: %v", i+1, synErr, test.input)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("[%d] expected no error, got %v, input: %v", i+1, runErr, test.input)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("[%d] expected %v, got %v, input: %v", i+1, test.output.Type, evaluated.Type, test.input)
		}
		if evaluated.Value != test.output.Value {
			t.Errorf("[%d] expected %v, got %v, input: %v", i+1, test.output.Value, evaluated.Value, test.input)
		}
	}
}
func TestVariableDeclarationAndAssignment(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "let x = 10\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "let x = 'hello'\nx",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "hello",
			},
		},
		{
			input: "let x = 10\nx = 20\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(20),
			},
		},
		{
			input: "let x = {foo: 'bar'}\nx.foo",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "bar",
			},
		},
		{
			input: "let x = 10\nlet y = x + 5\ny",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(15),
			},
		},
		{
			input: "let x = 10\nx = x + 5\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(15),
			},
		},
		{
			input: "let x = 10\nlet y = x\nx = 20\ny",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "let x = 10\nlet y = x\nx = 20\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(20),
			},
		},
		{
			input: "let x = 10\nlet y = x * 2\nlet z = y + x\nz",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(30),
			},
		},
		{
			input: "let x = {foo: 10}\nx.foo = 20\nx.foo",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(20),
			},
		},
		{
			input: "let x = {foo: {bar: 10}}\nx.foo.bar = 20\nx.foo.bar",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(20),
			},
		},
		{
			input: "let x = 10\nlet y = 20\nlet z = x + y\nz",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(30),
			},
		},
		{
			input: "let x = 10\nlet y = x\nlet z = y\nz",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "let x = 10\nlet y = x\nx = 20\ny",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "let x = 10\nlet y = x\nx = 20\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(20),
			},
		},
		{
			input: "let x = 10\nlet y = x + 5\nlet z = y * 2\nz",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(30),
			},
		},
		{
			input: "let x = 10\nlet y = x + 5\nlet z = y * 2\nx + y + z",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(55),
			},
		},
	}
	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("test %d failed: input=%q, expected type %v, got %v", i, test.input, test.output.Type, evaluated.Type)
		}
		evaluatedStr := fmt.Sprintf("%v", evaluated.Value)
		if evaluated.Value != test.output.Value {
			t.Errorf("test %d failed: input=%q, value mismatch. expected %s, got %s", i, test.input, fmt.Sprintf("%v", test.output.Value), evaluatedStr)
		}
	}
}
func TestFunctions(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "fn myFunc(arg1, arg2) { arg1 }\nmyFunc(10, 20)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "fn myFunc(arg1, arg2) { arg2 }\nmyFunc(10, 20)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(20),
			},
		},
		{
			input: "fn myFunc(arg1, arg2) { arg1 + arg2 }\nmyFunc(10, 20)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(30),
			},
		},
		{
			input: "fn myFunc(arg1, arg2) { arg1 * arg2 }\nmyFunc(10, 20)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(200),
			},
		},
		{
			input: "fn myFunc(arg1) { arg1 * 2 }\nmyFunc(15)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(30),
			},
		},
		{
			input: "fn myFunc() { 42 }\nmyFunc()",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(42),
			},
		},
		{
			input: "fn myFunc() { 'hello' }\nmyFunc()",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "hello",
			},
		},
		{
			input: "fn myFunc(arg1) { let x = arg1 * 2\nx + 10 }\nmyFunc(5)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(20),
			},
		},
		{
			input: "let myFunc = fn(arg1, arg2) { arg1 + arg2 }\nmyFunc(5, 10)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(15),
			},
		},
		{
			input: "let myFunc = fn() { 'anonymous' }\nmyFunc()",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "anonymous",
			},
		},
		{
			input: "fn outerFunc(arg1) { fn innerFunc(arg2) { arg1 + arg2 } }\nlet inner = outerFunc(10)\ninner(5)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(15),
			},
		},
		{
			input: "fn myFunc(arg1) { let x = arg1 * 2\nx }\nmyFunc(7)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(14),
			},
		},
		{
			input: "fn myFunc(arg1) { let x = {foo: arg1}\nx.foo }\nmyFunc(42)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(42),
			},
		},
		{
			input: "fn myFunc() { fn() { 'nested' } }\nlet nestedFunc = myFunc()\nnestedFunc()",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "nested",
			},
		},
		{
			input: "fn myFunc(arg1) { let x = fn() { arg1 * 2 }\nx() }\nmyFunc(8)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(16),
			},
		},
		{
			input: "fn myFunc(arg1, arg2) { let x = arg1 + arg2\nfn() { x * 2 } }\nlet nestedFunc = myFunc(3, 4)\nnestedFunc()",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(14),
			},
		},
		// Function calls with insufficient arguments
		{
			input: "fn myFunc(arg) { return arg }\nmyFunc()",
			output: shared.RuntimeValue{
				Type:  shared.Nil,
				Value: nil,
			},
		},
		// Function calls with too many arguments
		{
			input: "fn myFunc(arg) { return arg }\nmyFunc(1, 2)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(1),
			},
		},
	}

	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("test %d failed: input=%q, expected type %v, got %v", i, test.input, test.output.Type, evaluated.Type)
		}
		if evaluated.Value != test.output.Value {
			t.Errorf("test %d failed: input=%q, value mismatch. expected %v, got %v", i, test.input, test.output.Value, evaluated.Value)
		}
	}
}

func TestIfStatements(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "let didExecute = 1!=1\nif (1==1) {didExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (1!=1) {didExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (1<2) {didExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (1>2) {didExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (1==1) {didExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (1!=1) {didExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (2>=2) {didExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (2<=1) {didExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (1==1) {let x = 10\ndidExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (1!=1) {let x = 10\ndidExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (3*2==6) {didExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (3*2!=6) {didExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (5>3) {let y = 20\ndidExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (5<3) {let y = 20\ndidExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (1==1) {let obj = {foo: 'bar'}\ndidExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (1!=1) {let obj = {foo: 'bar'}\ndidExecute=1==1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		// Tests for if+else block
		{
			input: "let didExecute = 1!=1\nif (1==1) {didExecute=1==1} else {didExecute=1!=1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let didExecute = 1!=1\nif (1!=1) {didExecute=1==1} else {didExecute=1!=1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		// Tests for if+elseif block
		{
			input: "let num = 43\nlet didExecute = 1!=1\nif (num > 50) {didExecute=1==1} else if (num < 50) {didExecute=1!=1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "let num = 79\nlet didExecute = 1!=1\nif (num > 50) {didExecute=1==1} else if (num < 50) {didExecute=1!=1}\ndidExecute",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "let num = 50\nlet output = 0\nif (num < 30) {output = 1}\nelse if (num < 60) {output = 2}\nelse if (num < 90) {output = 3}\noutput",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(2),
			},
		},
		// Tests for if+elseif+else block
		{
			input: "let num = 120\nlet output = 0\nif (num < 30) {output = 1}\nelse if (num < 60) {output = 2}\nelse if (num < 90) {output = 3}\nelse {output = 4}\noutput",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(4),
			},
		},
		{
			input: "let x = 2\nif (x == 1) {x = 10}\nelse if (x == 2) {x = 20}\nelse if (x == 3) {x = 30}\nelse {x = 40}\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(20),
			},
		},
		{
			input: "let x = 3\nif (x == 1) {x = 10}\nelse if (x == 2) {x = 20}\nelse if (x == 3) {x = 30}\nelse {x = 40}\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(30),
			},
		},
		{
			input: "let x = 4\nif (x == 1) {x = 10}\nelse if (x == 2) {x = 20}\nelse if (x == 3) {x = 30}\nelse {x = 40}\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(40),
			},
		},
		{
			input: "let x = 5\nif (x == 1) {x = 10}\nelse if (x == 2) {x = 20}\nelse if (x == 3) {x = 30}\nelse {x = 40}\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(40),
			},
		},
	}

	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("test %d failed: input=%q, expected type %v, got %v", i, test.input, test.output.Type, evaluated.Type)
		}
		if evaluated.Value != test.output.Value {
			t.Errorf("test %d failed: input=%q, value mismatch. expected %v, got %v", i, test.input, test.output.Value, evaluated.Value)
		}
	}
}
func TestWhileLoops(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "let x = 0\nwhile (x < 5) { x = x + 1 }\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(5),
			},
		},
		{
			input: "let x = 10\nwhile (x > 0) { x = x - 2 }\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(0),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x <= 5) { sum = sum + x\nx = x + 1 }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(15),
			},
		},
		{
			input: "let x = 1\nlet product = 1\nwhile (x <= 4) { product = product * x\nx = x + 1 }\nproduct",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(24),
			},
		},
		{
			input: "let x = 0\nwhile (x < 3) { x = x + 1 }\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "let x = 0\nlet y = 0\nwhile (x < 3) { y = y + x\nx = x + 1 }\ny",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "let x = 0\nwhile (x < 0) { x = x + 1 }\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(0),
			},
		},
		{
			input: "let x = 5\nwhile (x > 0) { x = x - 1 }\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(0),
			},
		},
		{
			input: "let x = 0\nlet result = 1\nwhile (x < 4) { result = result * 2\nx = x + 1 }\nresult",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(16),
			},
		},
		{
			input: "let x = 0\nlet y = 0\nwhile (x < 5) { if (x % 2 == 0) { y = y + x }\nx = x + 1 }\ny",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(6),
			},
		},
		{
			input: "let x = 0\nlet y = 1\nwhile (x < 3) { y = y * 2\nx = x + 1 }\ny",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(8),
			},
		},
		{
			input: "let x = 10\nwhile (x > 5) { x = x - 1 }\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(5),
			},
		},
		{
			input: "let x = 0\nlet count = 0\nwhile (x < 10) { if (x % 2 == 0) { count = count + 1 }\nx = x + 1 }\ncount",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(5),
			},
		},
		{
			input: "let x = 0\nlet result = 0\nwhile (x < 5) { result = result + x\nx = x + 1 }\nresult",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "let x = 0\nwhile (x < 3) { let y = x * 2\nx = x + 1 }\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
	}

	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("test %d failed: input=%q, expected type %v, got %v", i, test.input, test.output.Type, evaluated.Type)
		}
		if evaluated.Value != test.output.Value {
			t.Errorf("test %d failed: input=%q, value mismatch. expected %v, got %v", i, test.input, test.output.Value, evaluated.Value)
		}
	}
}

func TestTryCatch(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
		check  func(t *testing.T, result shared.RuntimeValue) bool
	}{
		{
			input: "let error = 'error not triggered'\ntry {undefinedVariable} catch e {error = e}\nerror.message",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "Cannot resolve variable `undefinedVariable`",
			},
			check: func(t *testing.T, result shared.RuntimeValue) bool {
				if result.Type != shared.String {
					t.Errorf("Expected error to be a string, got %v", result.Type)
					return false
				}

				if !strings.Contains(result.Value.(string), "Cannot resolve variable") {
					t.Errorf("Unexpected error message: %v", result.Value)
					return false
				}

				return true
			},
		},
		{
			input: "let error = 'error not triggered'\ntry {1 + 1} catch e {error = e}\nerror",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "error not triggered",
			},
		},
		{
			input: "let error = 'error not triggered'\ntry {let x = 10\nx} catch e {error = e}\nerror",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "error not triggered",
			},
		},
		{
			input: "let error = 'error not triggered'\ntry {let x = y} catch e {error = e}\nerror.message",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "Cannot resolve variable `y`",
			},
		},
		{
			input: "let error = 'error not triggered'\ntry {let x = 10\nx = x + y} catch e {error = e}\nerror.message",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "Cannot resolve variable `y`",
			},
		},
		{
			input: "let error = 'error not triggered'\ntry {let obj = {foo: 'bar'}\nobj.bar} catch e {error = e}\nerror",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "error not triggered",
			},
		},
		{
			input: "let error = 'error not triggered'\ntry {let obj = {foo: 'bar'}\nobj.foo()} catch e {error = e}\nerror.message",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "Cannot invoke a non-function (attempted to call a string).",
			},
		},
		{
			input: "let error = 'error not triggered'\ntry {let arr = [1, 2, 3]\narr[5]} catch e {error = e}\nerror",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "error not triggered",
			},
		},
		{
			input: "let error = 'error not triggered'\ntry {let arr = [1, 2, 3]\narr.foo()} catch e {error = e}\nerror.message",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "Cannot access property of array by non-number (attempting to access properties by Identifier).",
			},
		},
		{
			input: "let error = 'error not triggered'\ntry {let x = 10\nx = x + 5} catch e {error = e}\nerror",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "error not triggered",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			p := parser.New("test")
			env := environment.NewEnvironment(nil)
			program, synErr := p.ProduceAST(test.input)
			if synErr != nil {
				t.Fatalf("parser error: %v", synErr)
			}

			evaluated, runErr := evaluator.Evaluate(program, env, nil)
			if runErr != nil {
				t.Fatalf("evaluation error: %v", runErr)
			}

			if test.check != nil {
				if !test.check(t, *evaluated) {
					t.Errorf("test %d failed custom check for input: %q", i, test.input)
				}
			} else {
				if evaluated.Type != test.output.Type {
					t.Errorf("test %d failed: input=%q, expected type %v, got %v",
						i, test.input, test.output.Type, evaluated.Type)
				}
				if !reflect.DeepEqual(evaluated.Value, test.output.Value) {
					t.Errorf("test %d failed: input=%q, value mismatch. expected %+v, got %+v",
						i, test.input, test.output.Value, evaluated.Value)
				}
			}
		})
	}
}
func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "fn myFunc() { return 42 }\nmyFunc()",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(42),
			},
		},
		{
			input: "fn myFunc() { return 'hello' }\nmyFunc()",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "hello",
			},
		},
		{
			input: "fn myFunc() { return {foo: 'bar'} }\nmyFunc()",
			output: shared.RuntimeValue{
				Type: shared.Object,
				Value: map[string]*shared.RuntimeValue{
					"foo": {
						Type:  shared.String,
						Value: "bar",
					},
				},
			},
		},
		{
			input: "fn myFunc() { if (1 == 1) { return 10 } return 20 }\nmyFunc()",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "fn myFunc() { let x = 5\nreturn x * 2 }\nmyFunc()",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "fn myFunc() { return fn() { return 'nested' } }\nlet nestedFunc = myFunc()\nnestedFunc()",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "nested",
			},
		},
		{
			input: "fn myFunc() { let x = 10\nif (x > 5) { return x + 5 } return x - 5 }\nmyFunc()",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(15),
			},
		},
		{
			input: "fn myFunc() { let x = 10\nwhile (x > 0) { return x } }\nmyFunc()",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "fn myFunc() { return fn(arg) { return arg * 2 } }\nlet double = myFunc()\ndouble(5)",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "fn myFunc() { let x = {foo: 'bar'}\nreturn x.foo }\nmyFunc()",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "bar",
			},
		},
	}

	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("test %d failed: input=%q, expected type %v, got %v", i, test.input, test.output.Type, evaluated.Type)
		}
		if !reflect.DeepEqual(evaluated.Value, test.output.Value) {
			t.Errorf("test %d failed: input=%q, value mismatch. expected %v, got %v", i, test.input, test.output.Value, evaluated.Value)
		}
	}
}
func TestContinueKeyword(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 5) { x = x + 1\nif (x == 3) { continue }\nsum = sum + x }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(12),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 5) { x = x + 1\nif (x % 2 == 0) { continue }\nsum = sum + x }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(9),
			},
		},
		{
			input: "let x = 0\nlet count = 0\nwhile (x < 10) { x = x + 1\nif (x < 5) { continue }\ncount = count + 1 }\ncount",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(6),
			},
		},
		{
			input: "let x = 0\nlet product = 1\nwhile (x < 5) { x = x + 1\nif (x % 2 == 0) { continue }\nproduct = product * x }\nproduct",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(15),
			},
		},
		{
			input: "let x = 0\nlet result = 0\nwhile (x < 5) { x = x + 1\nif (x == 3) { continue }\nresult = result + x }\nresult",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(12),
			},
		},
		{
			input: "let x = 0\nlet evenSum = 0\nwhile (x < 10) { x = x + 1\nif (x % 2 != 0) { continue }\nevenSum = evenSum + x }\nevenSum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(30),
			},
		},
		{
			input: "let x = 0\nlet y = 0\nwhile (x < 5) { x = x + 1\nif (x == 2) { continue }\ny = y + x }\ny",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(13),
			},
		},
		{
			input: "let x = 0\nlet count = 0\nwhile (x < 7) { x = x + 1\nif (x % 3 == 0) { continue }\ncount = count + 1 }\ncount",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(5),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 6) { x = x + 1\nif (x == 4) { continue }\nsum = sum + x }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(17),
			},
		},
		{
			input: "let x = 0\nlet result = 1\nwhile (x < 5) { x = x + 1\nif (x%2 == 1) { continue }\nresult = result * x }\nresult",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(8),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 8) { x = x + 1\nif (x % 2 == 1) { continue }\nsum = sum + x }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(20),
			},
		},
		{
			input: "let x = 0\nlet count = 0\nwhile (x < 10) { x = x + 1\nif (x > 7) { continue }\ncount = count + 1 }\ncount",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(7),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 5) { x = x + 1\nif (x == 3) { continue }\nsum = sum + x }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(12),
			},
		},
		{
			input: "let x = 0\nlet product = 1\nwhile (x < 6) { x = x + 1\nif (x == 2) { continue }\nif (x == 5) { continue }\nproduct = product * x }\nproduct",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(72),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 4) { x = x + 1\nif (x == 2) { continue }\nsum = sum + x }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(8),
			},
		},
		{
			input: "let x = 0\nlet count = 0\nwhile (x < 10) { x = x + 1\nif (x % 2 == 0) { continue }\ncount = count + 1 }\ncount",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(5),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 7) { x = x + 1\nif (x == 6) { continue }\nsum = sum + x }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(22),
			},
		},
		{
			input: "let x = 0\nlet result = 1\nwhile (x < 4) { x = x + 1\nif (x == 2) { continue }\nresult = result * x }\nresult",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(12),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 5) { x = x + 1\nif (x == 3) { continue }\nsum = sum + x }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(12),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 5) { x = x + 1\nif (x == 3) { continue }\nsum = sum + x }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(12),
			},
		},
		{
			input: "let err = ''\ntry {continue} catch e {err = e}\nerr.message",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "`continue` statement used outside of a loop context.",
			},
		},
	}

	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("test %d failed: input=%q, expected type %v, got %v", i, test.input, test.output.Type, evaluated.Type)
		}
		if !reflect.DeepEqual(evaluated.Value, test.output.Value) {
			t.Errorf("test %d failed: input=%q, value mismatch. expected %v, got %v", i, test.input, test.output.Value, evaluated.Value)
		}
	}
}

func TestBreakKeyword(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "let x = 0\nwhile (x < 5) { if (x == 3) { break }\nx = x + 1 }\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 5) { if (x == 3) { break }\nsum = sum + x\nx = x + 1 }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "let x = 0\nwhile (x < 5) { break\nx = x + 1 }\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(0),
			},
		},
		{
			input: "let x = 0\nlet y = 0\nwhile (x < 5) { if (x == 2) { break }\ny = y + x\nx = x + 1 }\ny",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(1),
			},
		},
		{
			input: "let x = 0\nwhile (x < 5) { if (x == 0) { break }\nx = x + 1 }\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(0),
			},
		},
		{
			input: "let x = 0\nlet product = 1\nwhile (x < 5) { if (x == 3) { break }\nproduct = product * (x + 1)\nx = x + 1 }\nproduct",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(6),
			},
		},
		{
			input: "let x = 0\nwhile (x < 5) { if (x == 4) { break }\nx = x + 1 }\nx",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(4),
			},
		},
		{
			input: "let x = 0\nlet count = 0\nwhile (x < 10) { if (x == 5) { break }\ncount = count + 1\nx = x + 1 }\ncount",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(5),
			},
		},
		{
			input: "let x = 0\nlet result = 0\nwhile (x < 5) { if (x == 2) { break }\nresult = result + x\nx = x + 1 }\nresult",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(1),
			},
		},
		{
			input: "let x = 0\nlet evenSum = 0\nwhile (x < 10) { if (x == 6) { break }\nif (x % 2 == 0) { evenSum = evenSum + x }\nx = x + 1 }\nevenSum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(6),
			},
		},
		{
			input: "let x = 0\nlet y = 0\nwhile (x < 5) { if (x == 3) { break }\ny = y + x\nx = x + 1 }\ny",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "let x = 0\nlet count = 0\nwhile (x < 7) { if (x == 4) { break }\ncount = count + 1\nx = x + 1 }\ncount",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(4),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 6) { if (x == 3) { break }\nsum = sum + x\nx = x + 1 }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "let x = 0\nlet result = 1\nwhile (x < 5) { if (x == 2) { break }\nresult = result * (x + 1)\nx = x + 1 }\nresult",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(2),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 5) { if (x == 4) { break }\nsum = sum + x\nx = x + 1 }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(6),
			},
		},
		{
			input: "let x = 0\nlet count = 0\nwhile (x < 10) { if (x == 7) { break }\ncount = count + 1\nx = x + 1 }\ncount",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(7),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 5) { if (x == 3) { break }\nsum = sum + x\nx = x + 1 }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "let x = 0\nlet product = 1\nwhile (x < 6) { if (x == 2) { break }\nproduct = product * (x + 1)\nx = x + 1 }\nproduct",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(2),
			},
		},
		{
			input: "let x = 0\nlet sum = 0\nwhile (x < 4) { if (x == 2) { break }\nsum = sum + x\nx = x + 1 }\nsum",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(1),
			},
		},
		{
			input: "let x = 0\nlet count = 0\nwhile (x < 10) { if (x == 5) { break }\ncount = count + 1\nx = x + 1 }\ncount",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(5),
			},
		},
		{
			input: "let err = ''\ntry {break} catch e {err = e}\nerr.message",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "`break` statement used outside of a loop context.",
			},
		},
	}

	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("test %d failed: input=%q, expected type %v, got %v", i, test.input, test.output.Type, evaluated.Type)
		}
		if !reflect.DeepEqual(evaluated.Value, test.output.Value) {
			t.Errorf("test %d failed: input=%q, value mismatch. expected %v, got %v", i, test.input, test.output.Value, evaluated.Value)
		}
	}
}
func TestArrays(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "let arr = [1, 2, 3]\narr[0]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(1),
			},
		},
		{
			input: "let arr = [1, 2, 3]\narr[1]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(2),
			},
		},
		{
			input: "let arr = [1, 2, 3]\narr[2]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "let arr = [1, 2, 3]\narr[3]",
			output: shared.RuntimeValue{
				Type:  shared.Nil,
				Value: nil,
			},
		},
		{
			input: "let arr = [1, 2, 3]\narr[0] = 10\narr[0]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "let arr = [1, 2, 3]\narr[1] = 20\narr[1]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(20),
			},
		},
		{
			input: "let arr = [1, 2, 3]\narr[2] = 30\narr[2]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(30),
			},
		},
		{
			input: "let arr = [1, 2, 3]\narr[3] = 40\narr[3]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(40),
			},
		},
		{
			input: "let arr = [1, 2, 3]\narr[4] = 50\narr[4]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(50),
			},
		},
		{
			input: "let arr = [1, 2, 3]\narr[5] = 60\narr[5]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(60),
			},
		},
		{
			input: "let arr = [1, 2, 3]\nlet i = 1\narr[i]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(2),
			},
		},
		{
			input: "let arr = [1, 2, 3]\nlet i = 2\narr[i]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "[1, 2, 3][2]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "let obj = {arr: [1, 2, 3]}\nobj.arr[2]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		/*
			Expressions like f()[2] and (f())[2] are evaluated differently due to explicit grouping rules.
			Without parentheses, f()[2] is parsed as two separate expressions: a function call f() followed by a standalone array literal [2].
			Since they are not syntactically connected, the array does not operate on the function result.
			To index into the result of a function, parentheses must be used to group the call expression: (f())[2].
			This ensures the function is evaluated first, and the result is then indexed as expected.
		*/
		{
			input: "fn foo() { return [1, 2, 3] }\n(foo())[2]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		{
			input: "let obj = {arr: [1, 2, 3]}\nlet foo = fn() { return obj.arr }\n(foo())[2]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(3),
			},
		},
		// ISSUE#5: Array Assignment in Object
		{
			input: "let obj = {arr: [1, 2, 3]}\nobj.arr[0] = 10\nobj.arr[0]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "let obj = {arr: [1, \"hello\"]}\nobj.arr[1] = 10\nobj.arr[1]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "let obj = {prop:{arr: [1, 2, 3]}}\nobj.prop.arr[0] = 10\nobj.prop.arr[0]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "let obj = {prop:{arr: [1, \"hello\"]}}\nobj.prop.arr[1] = 10\nobj.prop.arr[1]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "let obj = {prop:{arr: [1, {prop2:{arr2: [1, 2, 3]}}]}}\nobj.prop.arr[1].prop2.arr2[0] = 10\nobj.prop.arr[1].prop2.arr2[0]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
		{
			input: "let obj = {prop:{arr: [1, {prop2:{arr2: [1, {prop3:{arr3: [1, 2, 3]}}]}}]}}\nobj.prop.arr[1].prop2.arr2[1].prop3.arr3[0] = 10\nobj.prop.arr[1].prop2.arr2[1].prop3.arr3[0]",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(10),
			},
		},
	}

	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Errorf("test %d failed: input=%q, expected no error, got %v", i, test.input, runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("test %d failed: input=%q, expected type %v, got %v", i, test.input, test.output.Type, evaluated.Type)
		}
		if !reflect.DeepEqual(evaluated.Value, test.output.Value) {
			t.Errorf("test %d failed: input=%q, value mismatch. expected %v, got %v", i, test.input, test.output.Value, evaluated.Value)
		}
	}
}

func TestClasses(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: `
				let i = 1
				class Person {
					public name
					public age
					private id
					public constructor(n, a) {
						name = n
						age = a
						id = i
						i = i + 1
					}
				
					private joinStrings(str1, str2) {
						return str1 + ": " + str2
					}

					public speak(str) {
						return joinStrings(name, str)
					}
				}

				let john = Person("John", 30)
				john.speak("Hello World!")
			`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "John: Hello World!",
			},
		},
		{
			input: `
				let i = 1
				class Person {
					public name
					public age
					private id
					public constructor(n, a) {
						name = n
						age = a
						id = i
						i = i + 1
					}
				}
				let john = Person("John", 30)
				john.id
			`,
			output: shared.RuntimeValue{
				Type:  shared.Nil,
				Value: nil,
			},
		},
		{
			input: `
				let i = 1
				class Person {
					public name
					public age
					private id
					public constructor(n, a) {
						name = n
						age = a
						id = i
						i = i + 1
					}
					public test() { name = "Jedi" }
					public getName() { return name }
				}
				let john = Person("John", 30)
				john.getName()
				john.test()
				john.getName()
			`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "Jedi",
			},
		},
		{
			input: `
				let i = 1
				class Person {
					public name
					public age
					private id
					public constructor(n, a) {
						name = n
						age = a
						id = i
						i = i + 1
					}
				}
				let john = Person("John", 30)
				john.name
			`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "John",
			},
		},
		{
			input: `
				class Test {
					private secret() { return "shhh" }
				}
				let t = Test()
				t.secret
			`,
			output: shared.RuntimeValue{
				Type:  shared.Nil,
				Value: nil,
			},
		},
		{
			input: `
				class Data {
					public value
					public constructor(v) { value = v }
				}
				let d = Data(42)
				d["value"]
			`,
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 42.0,
			},
		},
		{
			input: `
				class Data {
					private hidden
					public constructor() { hidden = 999 }
				}
				let d = Data()
				let key = "hidden"
				d[key]
			`,
			output: shared.RuntimeValue{
				Type:  shared.Nil,
				Value: nil,
			},
		},
		{
			input: `
				class Chain {
					public one() { return two() }
					public two() { return "called two" }
				}
				let c = Chain()
				c.one()
			`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "called two",
			},
		},
		{
			input: `
				class Mutate {
					public val
					public constructor() { val = 5 }
					public inc() { val = val + 1 }
					public get() { return val }
				}
				let m = Mutate()
				m.inc()
				m.inc()
				m.get()
			`,
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 7.0,
			},
		},
		{
			input: `
				class Mathy {
					public double
					public constructor(x) {
						double = x * 2
					}
				}
				let m = Mathy(4)
				m.double
			`,
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 8.0,
			},
		},
		{
			input: `
				class Weird {
					public prop = fn (){"This is a function property!"}
					public constructor() {}
				}
				let weirdObj = Weird()
				weirdObj.prop()
			`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "This is a function property!",
			},
		},
		{
			input: `
				class Foo {
					public bar
					public constructor() {
						bar = 42
					}
				}
				class Baz {
					public foo
					public constructor() {
						foo = Foo()
					}
				}
				let baz = Baz()
				baz.foo.bar
			`,
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(42),
			},
		},
		{
			input: `
				class Factorial {
					public compute(n) {
						if (n == 0) {
							return 1
						}
						return n * compute(n - 1)
					}
				}
				let fact = Factorial()
				fact.compute(5)
			`,
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: float64(120),
			},
		},
		{
			input: `
				class Secret {
					private secret() { return "This is private!" }
				}
				let s = Secret()
				s.secret
			`,
			output: shared.RuntimeValue{
				Type:  shared.Nil,
				Value: nil,
			},
		},
		{
			input: `
				class Transformer {
					public value
					public transform = fn () { value = "Transformed" }
					public constructor() {
						value = "Original"
					}
				}
				let obj = Transformer()
				obj.transform()
				obj.value
			`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "Transformed",
			},
		},
		{
			input: `
				class Test {
					public arg = "i am not initialized"
					public constructor(arg_param) {
						arg = arg_param
					}
				}

				let instance = Test()
				instance.arg
			`,
			output: shared.RuntimeValue{
				Type:  shared.Nil,
				Value: nil,
			},
		},
		{
			input: `
				class Test {
					public arg = "i am not initialized"
					public constructor(arg_param) {
						arg = arg_param
					}
				}

				let instance = Test("i am initialized", "i am extra")
				instance.arg
			`,
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "i am initialized",
			},
		},
	}

	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Fatalf("test %d failed: input=%q, expected no error, got %v", i, test.input, synErr)
		}

		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Fatalf("test %d failed: input=%q, expected no error, got %v", i, test.input, runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("test %d failed: input=%q, expected type %v, got %v", i, test.input, test.output.Type, evaluated.Type)
		}
		if !reflect.DeepEqual(evaluated.Value, test.output.Value) {
			t.Errorf("test %d failed: input=%q, value mismatch. expected %v, got %v", i, test.input, test.output.Value, evaluated.Value)
		}
	}
}

func TestLogicalOperators(t *testing.T) {
	tests := []struct {
		input  string
		output shared.RuntimeValue
	}{
		{
			input: "'apple' && 'banana'",
			output: shared.RuntimeValue{
				Type:  shared.String,
				Value: "banana",
			},
		},
		{
			input: "1 && 0",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 0.0,
			},
		},
		{
			input: "0 && 1",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 0.0,
			},
		},
		{
			input: "0 && 0",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 0.0,
			},
		},
		{
			input: "false && true",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "true && false",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "false && false",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "nil && 123",
			output: shared.RuntimeValue{
				Type:  shared.Nil,
				Value: nil,
			},
		},
		{
			input: "1 || 1",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 1.0,
			},
		},
		{
			input: "1 || 0",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 1.0,
			},
		},
		{
			input: "0 || 1",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 1.0,
			},
		},
		{
			input: "0 || 0",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 0.0,
			},
		},
		{
			input: "true || false",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "false || true",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
		{
			input: "false || false",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "nil || 123",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 123.0,
			},
		},
		{
			input: "nil || nil",
			output: shared.RuntimeValue{
				Type:  shared.Nil,
				Value: nil,
			},
		},
		{
			input: "0 ?? 1",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 0.0,
			},
		},
		{
			input: "1 ?? 0",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 1.0,
			},
		},
		{
			input: "nil ?? 1",
			output: shared.RuntimeValue{
				Type:  shared.Number,
				Value: 1.0,
			},
		},
		{
			input: "false ?? false",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "!true",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: false,
			},
		},
		{
			input: "!false",
			output: shared.RuntimeValue{
				Type:  shared.Boolean,
				Value: true,
			},
		},
	}

	for i, test := range tests {
		p := parser.New("test")
		env := environment.NewEnvironment(nil)
		env.DeclareVar("true", values.MK_BOOL(true), false)
		env.DeclareVar("false", values.MK_BOOL(false), false)
		env.DeclareVar("nil", values.MK_NIL(), false)
		program, synErr := p.ProduceAST(test.input)
		if synErr != nil {
			t.Fatalf("test %d failed: input=%q, expected no error, got %v", i, test.input, synErr)
		}
		evaluated, runErr := evaluator.Evaluate(program, env, nil)
		if runErr != nil {
			t.Fatalf("test %d failed: input=%q, expected no error, got %v", i, test.input, runErr)
		}
		if evaluated.Type != test.output.Type {
			t.Errorf("test %d failed: input=%q, expected type %v, got %v", i, test.input, test.output.Type, evaluated.Type)
		}
		if !reflect.DeepEqual(evaluated.Value, test.output.Value) {
			t.Errorf("test %d failed: input=%q, value mismatch. expected %v, got %v", i, test.input, test.output.Value, evaluated.Value)
		}
	}
}
