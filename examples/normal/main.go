package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/dzonerzy/quickjs-go"
)

func check(err error) {
	if err != nil {
		var evalErr *quickjs.Error
		if errors.As(err, &evalErr) {
			fmt.Println(evalErr.Cause)
			fmt.Println(evalErr.Stack)
		}
		panic(err)
	}
}

func main() {
	runtime := quickjs.NewRuntime()
	defer runtime.Free()

	context := runtime.NewContext()
	defer context.Free()

	globals := context.Globals()

	// Test evaluating template strings.

	result, err := context.Eval("`Hello world! 2 ** 8 = ${2 ** 8}.`", quickjs.EVAL_GLOBAL)
	check(err)
	defer result.Free()

	fmt.Println(result.String())
	fmt.Println()

	// Test evaluating numeric expressions.

	/* test bytecode */
	/*ret, err := context.Binary(`
	(function(global) {
		global.Deaf = function() {
			return 1;
		};
	})(this)
	`, "<code>")
	check(err)
	fmt.Println(ret)*/
	ret := []byte{2, 2, 12, 60, 99, 111, 100, 101, 62, 8, 68, 101, 97, 102, 14, 0, 6, 0, 160, 1, 0, 2, 0, 2, 0, 1, 8, 2, 162, 1, 0, 0, 0, 16, 0, 1, 0, 8, 202, 192, 0, 198, 239, 205, 40, 194, 3, 1, 3, 0, 2, 10, 14, 67, 6, 0, 0, 1, 0, 1, 2, 0, 1, 9, 1, 220, 1, 0, 1, 0, 209, 192, 0, 67, 226, 0, 0, 0, 41, 194, 3, 2, 3, 3, 9, 38, 14, 67, 6, 0, 0, 0, 0, 0, 1, 0, 0, 2, 0, 182, 40, 194, 3, 3, 1, 3}
	result, err = context.EvalBinary(ret)
	check(err)
	fmt.Println(result)

	/* end */

	result, err = context.Eval(`1 + 2 * 100 - 3 + Math.sin(10)`, quickjs.EVAL_GLOBAL)
	check(err)
	defer result.Free()

	fmt.Println(result.Int64())
	fmt.Println()

	// Test evaluating big integer expressions.

	result, err = context.Eval(`128n ** 16n`, quickjs.EVAL_GLOBAL)
	check(err)
	defer result.Free()

	fmt.Println(result.BigInt())
	fmt.Println()

	// Test evaluating big decimal expressions.

	result, err = context.Eval(`128l ** 12l`, quickjs.EVAL_GLOBAL)
	check(err)
	defer result.Free()

	fmt.Println(result.BigFloat())
	fmt.Println()

	// Test evaluating boolean expressions.

	result, err = context.Eval(`false && true`, quickjs.EVAL_GLOBAL)
	check(err)
	defer result.Free()

	fmt.Println(result.Bool())
	fmt.Println()

	// Test setting and calling functions.

	A := func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		fmt.Println("A got called!")
		return ctx.Null()
	}

	B := func(ctx *quickjs.Context, this quickjs.Value, args []quickjs.Value) quickjs.Value {
		fmt.Println("B got called!")
		return ctx.Null()
	}

	globals.Set("A", context.Function(A))
	globals.Set("B", context.Function(B))

	_, err = context.Eval(`for (let i = 0; i < 10; i++) { if (i % 2 === 0) A(); else B(); }`, quickjs.EVAL_GLOBAL)
	check(err)

	fmt.Println()

	// Test setting global variables.

	_, err = context.Eval(`HELLO = "world"; TEST = false;`, quickjs.EVAL_GLOBAL)
	check(err)

	names, err := globals.PropertyNames()
	check(err)

	fmt.Println("Globals:")
	for _, name := range names {
		val := globals.GetByAtom(name.Atom)
		defer val.Free()

		fmt.Printf("'%s': %s\n", name, val)
	}
	fmt.Println()

	// Test evaluating arbitrary expressions from flag arguments.

	flag.Parse()
	if flag.NArg() == 0 {
		return
	}

	result, err = context.Eval(strings.Join(flag.Args(), " "), quickjs.EVAL_GLOBAL)
	check(err)
	defer result.Free()

	if result.IsObject() {
		names, err := result.PropertyNames()
		check(err)

		fmt.Println("Object:")
		for _, name := range names {
			val := result.GetByAtom(name.Atom)
			defer val.Free()

			fmt.Printf("'%s': %s\n", name, val)
		}
	} else {
		fmt.Println(result.String())
	}
}
