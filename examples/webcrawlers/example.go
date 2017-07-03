package main

import (
	"fmt"

	"github.com/PMoneda/flow"
)

func main() {
	basicCrawler()
	crawlerWithError()
}

func basicCrawler() {
	list := flow.NewFlow().From("crawler://www.metalsucks.net/", ".post-title").GetBody()
	fmt.Println(list)
}
func crawlerWithError() {
	pipe := flow.NewFlow().From("crawler://www.metal1sucks.net/", ".post-title")
	ch := pipe.Choice()
	ch = ch.When(flow.Header("error").Exist()).To("print://?msg=Fail to crawler")
	ch = ch.Otherwise().To("print://?msg=${body}")
	pipe.GetBody()
}
