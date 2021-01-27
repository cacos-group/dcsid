package leaf

import (
	"fmt"
	"testing"
)

func TestGenid_NextId(t *testing.T) {

	leaf := New(&Config{
		DSN:    "user:password@tcp(127.0.0.1:3306)/leaf",
		BizTag: "test",
	})

	go func() {
		for i := 0; i < 100000; i++ {
			id, err := leaf.NextId()
			if err != nil {
				t.Error(err)
				return
			}
			fmt.Println(id)
		}
	}()

	for i := 0; i < 100000; i++ {
		id, err := leaf.NextId()
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(id)
	}

}
