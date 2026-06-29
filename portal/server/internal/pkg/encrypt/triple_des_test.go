package encrypt

import (
	"fmt"
	"testing"
)

func TestTripe(t *testing.T) {

	rs, err := TripleDesCbcDecrypt("be8a5e494489f61a26233cb4f43fd7c7", []byte("CCE31A176862725175bb539f"), []byte("CCE31A176862725175bb539f"[1:9]))
	fmt.Println(string(rs), err)

}
