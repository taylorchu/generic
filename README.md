# Generic, a code generation tool to enable generic in go.

`go install github.com/taylorchu/generic/cmd/generic`

This is an experiment to enable generic with code generation in the most elegant way.

## Existing approaches to generic in go

  - code generation
    - output
      - file
      - package (*)
    - rewrite method
      - simple string replacement
      - ast-based replacement (*)
    - type placeholder
      - text/template, i.e. `{{ .Type }}`
      - special types from a package, i.e. `generic.Type`
      - in-package type declaration, i.e. `type A int`
        - above-declaration comment, i.e. `// template type Vector(A, N)`
        - without comment
          - any type name
          - type name with certain pattern (*)
  - language change
  - interface{}
  - reflect
  - marshal everything to bytes/string
  - copy&paste

## What does `generic` do?

![](http://i.imgur.com/X07XInF.png)

`generic` does the followings if you put the following comments in your go code:

```go
// THINK: generate from a generic package and save result as a new package,
// with a list of rewrite rules!

//go:generate generic github.com/go/sort int_sort Type->int
```

and then run `go generate`:

1. Run `go get -u github.com/go/sort` if the package does not exist locally.
  - If the package exists locally, go-get will not be called.
2. Gather `.go` files (skip `_test.go`) in github.com/go/sort
3. Apply AST rewrite to replace Type in those `.go` files with int.
  - Only type that starts with __Type__ can be converted. This enables variable naming like __TypeKey__ or __TypeValue__
  that closely expresses meaning while there is still a namespace for type placeholder.
  - Many rewrite rules are possible: `TypeKey->string TypeValue->int`.
  - We can rewrite non-builtin types with `:`: `Type->github.com/go/types:types.Box`.
4. Type-check results.
5. Save the results as a new package called `int_sort` in `$PWD`.
  - If there is already a dir called `int_sort`, it will first be removed.
  - If the new package starts with `.`, it will save the results in `$PWD`:
      - The package name is set to `$GOPACKAGE` by `go-generate`.
      - All top-level identifiers will have prefixes to prevent conflicts, and their uses will also be updated.
      - Filenames will be renamed to prevent conflicts.

## Example

```go
package queue

type Type string

// TypeQueue represents a queue of Type types.
type TypeQueue struct {
	items []Type
}

// New makes a new empty Type queue.
func New() *TypeQueue {
	return &TypeQueue{items: make([]Type, 0)}
}

// Enq adds an item to the queue.
func (q *TypeQueue) Enq(obj Type) *TypeQueue {
	q.items = append(q.items, obj)
	return q
}

// Deq removes and returns the next item in the queue.
func (q *TypeQueue) Deq() Type {
	obj := q.items[0]
	q.items = q.items[1:]
	return obj
}

// Len gets the current number of Type items in the queue.
func (q *TypeQueue) Len() int {
	return len(q.items)
}
```

```
result Type->int64 TypeQueue->FIFO
```

```go
package result

type FIFO struct {
	items []int64
}

func New() *FIFO {
	return &FIFO{items: make([]int64, 0)}
}

func (q *FIFO) Enq(obj int64) *FIFO {
	q.items = append(q.items, obj)
	return q
}

func (q *FIFO) Deq() int64 {
	obj := q.items[0]
	q.items = q.items[1:]
	return obj
}

func (q *FIFO) Len() int {
	return len(q.items)
}
```

You can find more examples in `fixture/` and their outputs in `output/`.

## FAQ

### Why are type-checking and ast-based replacement important?

Type-checking and ast-based replacement ensure that the tool doesn't generate invalid code even you or the tool make mistakes, and rewrites identifiers in cases that it shouldn't.

### Why is type placeholder designed this way?

`type TypeXXX int32`

 - It provides a namespace for replaceable types.
 - Knowing that this type might be replaced, package creator can still write go-testable code with a concrete type.
 - It can express meaning. For example, `TypeQueue` shows that it is a queue.

### Why does this tool rewrite at package-level instead of file-level?

 - This tool tries NOT to apply any restriction for package creator except that any TypeXXX might be rewritten. Package creator has full flexibility to write normal go code.
 - It is common to distribute go code at package-level.

## LICENSE

The MIT License (MIT)
Copyright (c) 2016 taylorchu.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
