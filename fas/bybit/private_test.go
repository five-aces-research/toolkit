package bybit

import (
	"fmt"
	"github.com/DawnKosmos/bybit-go5/models"
	"testing"
)

func TestSetOrder(t *testing.T) {
	pr := NewPrivate("ke", "kE5jRrPgSPSuOJikF6", "JsVyekm5hgVl9SPSwxbTTx9nW7XqwAmBEVKT", true)

	res, err := pr.by.GetPositionInfo(models.GetPositionInfoRequest{
		Category: "inverse",
		Symbol:   "null",
	})
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	for _, v := range res.List {
		fmt.Printf("%+v\n", v)
	}

}
