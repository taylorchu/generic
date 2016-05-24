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
  that closely expresses meaning.
  - In the package, we can define types like `type Type int` to write tests for a specific type. This is why we choose
  to skip `_test.go` because they will not work after rewrite.
  - Many rewrite rules are possible: `TypeKey->string TypeValue->int`.
  - We can rewrite non-builtin types with `:`: `Type->github.com/go/types:types.Box`.
4. Type-check results.
5. Save the results as a new package called `int_sort` in `$PWD`.
  - This prevents conflicting definitions from the package.
  - If there is already a dir called `int_sort`, it will first be removed.
  - If the new package is `.`, it will save the results in the current dir, and use `$GOPACKAGE` set by `go-generate`.
    However, it will not remove any file if an error occurrs.


## Tricky examples that other code generation tools might fail

### Identifier with same name (Type->OtherType)

```go
package p

type Type int

func (Type Type) Type(Type Type) {

}
```

```go
package p

func (Type OtherType) Type(Type OtherType) {

}
```

### Import (Type2->github.com/golang/test:test.OtherType)

```go
package p

type Type3 Type2
```

```go
package p

import "github.com/golang/test"

type Type3 test.OtherType2
```

### Simple type-check (Type2->map[string]int63)

```go
package p

type Type3 Type2
```

```
undeclared name: int63
```

### Type2->github.com/taylorchu/generic:generic.Target1

```go
package p

type Type3 Type2
```

```
Target1 not declared by package generic
```

## LICENSE

The MIT License (MIT)
Copyright (c) 2016 taylorchu.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
