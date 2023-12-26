package hashutil_test

import (
	"fmt"

	"github.com/gsk148/urlShorteningService/internal/app/hashutil"
)

func ExampleEncode() {
	shortID := hashutil.Encode([]byte("string"))
	shortID = "some id like mtyxZrr"
	fmt.Println(shortID)
	// Output:
	// some id like mtyxZrr
}
