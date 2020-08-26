package main

import "github.com/dathan/rabbit-topic/pkg/delivery"

func main() {

	publish := delivery.New("wework.topic.exchange")

	str1 := `{"uid": 1, "offices": ["apple", "peach"]}`
	str2 := `{"uid": 2, "offices": ["mango", "grapes"]}`

	publish.SendTo("wework.production.user", str1)
	publish.SendTo("wework.NOTproduction.user", str2)

}
