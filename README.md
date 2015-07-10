# qs

Package qs encodes structs into URL query strings.

It uses reflection and Go struct tags (like encoding/json package)
to identify field names and encoding options.

It's handy when systematically constructing URL query parameters.
It only has one function exposed that returns `net/url.Values`.

Browse documentation on [GoDoc](https://godoc.org/github.com/ahmetalpbalkan/qs): [![GoDoc](https://godoc.org/github.com/ahmetalpbalkan/qs?status.svg)](https://godoc.org/github.com/ahmetalpbalkan/qs)

# Example

```go

import "github.com/ahmetalpbalkan/qs"

type SearchParams struct {
	Query string `qs:"q"` 
	Count int `qs:"num,omitempty"` 
	Next  int `qs:"n,omitempty"` 
}

params := SearchParams{
	Query: "apple pie",
	Count: 10}

query := qs.Encode(params).Encode()
// Output: q=apple+pie&num=10
```

# License

This package is distributed under Apache 2.0 License.
See [LICENSE](LICENSE) for more.

# Authors

* Ahmet Alp Balkan ([@ahmetalpbalkan](https://twitter.com/ahmetalpbalkan))
