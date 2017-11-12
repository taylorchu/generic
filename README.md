# Generic, a code generation tool to enable generics in go

[![CircleCI](https://circleci.com/gh/taylorchu/generic.svg?style=svg)](https://circleci.com/gh/taylorchu/generic)

__v2__ `go get -u github.com/taylorchu/generic/cmd/gorewrite`

__v1__ _(deprecated)_ `go get -u github.com/taylorchu/generic/cmd/generic`

This is an experiment to enable generics with code generation in the most elegant way.

<img src="http://i.imgur.com/X07XInF.png" width="300">

## Example

You can create a package like this. Note that `Type` and `TypeQueue` are type placeholders.

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

Add rewrite config in `GoRewrite.yaml`, and run `gorewrite`:

```yaml
spec:
  - name: result
    import: github.com/YourName/queue
    typeMap:
      Type:
        expr: int64
      TypeQueue:
        expr: FIFO
```

The output is saved to `$PWD/result/`.

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

## `GoRewrite.yaml` reference

The yaml config contains multiple rewrite specs.

- `spec[*].name` (string): unique identifier the spec. It is a path to the output, and used as package name if the spec is not local.
- `spec[*].local` (bool): true if the spec is local. If the spec is local, the output will be saved in `$PWD` instead of a new package relative to `$PWD`.
  All the top level identifiers and the filename will also be prefixed with `spec[*].name` to avoid conflicts.
- `spec[*].typeMap` (map): type mappings used to replace placeholders. The key is type placeholder. The value `expr` can be any go expression.
  If `expr` references any other packages, all those packages need to be listed in `import`.

```yaml
spec:
  - name: result
    local: true
    import: github.com/YourName/queue
    typeMap:
      Type:
        expr: test.Box
        import:
          - github/YourName/test
```

## FAQ

### What are the existing approaches to generics in go?

This [post](https://appliedgo.net/generics/) summarizes the current state in go to _"simulate"_ generics. 

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

### Why does it remove all the comments?

Comments in go ast are [free-floating](https://github.com/golang/go/issues/20744), so they are hard to work with. Hopefully it is fixed in the near future.

### How do I make sure the rewritten package is not import-able?

The spec name should start with `internal/`. For example, `internal/queue`.

> When the go command sees an import of a package with internal in its path, it verifies that the package doing the import is within the tree rooted at the parent of the internal directory. For example, a package .../a/b/c/internal/d/e/f can be imported only by code in the directory tree rooted at .../a/b/c. It cannot be imported by code in .../a/b/g or in any other repository.

See [Internal packages](https://golang.org/doc/go1.4#internalpackages).

### How do I determine spec name?

Spec name is essentially go package path, and the base name of this package path is a package name.

> Good package names are short and clear. They are lower case, with no under_scores or mixedCaps.

See [Package names](https://blog.golang.org/package-names).
