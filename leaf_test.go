package leaf

import (
	"fmt"
	"testing"
	"time"
)

func TestGenid_NextId(t *testing.T) {

	leaf := New(&Config{
		DSN:    "",
		BizTag: "test1",
	})

	go func() {
		for i := 0; i < 100000; i++ {
			id, err := leaf.NextId()
			if err != nil {
				t.Error(err)
				return
			}
			fmt.Println(id)
			time.Sleep(time.Millisecond)
		}
	}()

	for i := 0; i < 100000; i++ {
		id, err := leaf.NextId()
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(id)
		time.Sleep(time.Millisecond)
	}

}
