package qs

import (
	"fmt"
)

func ExampleEncode() {
	type QueryParams struct {
		Query string   `qs:"q"`
		Count int      `qs:"num,omitempty"`
		Opt   []string `qs:"opt"`
	}

	p := QueryParams{
		Query: "apple pie",
		Count: 10,
		Opt:   []string{"safe", "localized"},
	}

	fmt.Println(Encode(p).Encode())
	// Output:
	// num=10&opt=safe&opt=localized&q=apple+pie
}
