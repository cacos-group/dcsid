package leaf

import (
	"fmt"
	"testing"
	"time"
)

func TestGenid_NextId(t *testing.T) {

	leaf := New(&Config{
		DSN:    "souti_growth:yC3f4NTLLTS8OUHV@tcp(rm-2ze335994i6a5ii8mfm.mysql.rds.aliyuncs.com:3306)/leaf",
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
