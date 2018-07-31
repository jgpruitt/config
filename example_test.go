package config_test

import (
	"strings"
	"github.com/jgpruitt/config"
	"fmt"
	"time"
	"net"
)

func ExampleRead() {
	var file = `

		number = 1234
		every = 3m20s

		database:
			username = admin
			port=5432

		log:
			path=../out/log.txt
			level=fatal
	`
	cfgs, _ := config.Read(strings.NewReader(file))

	// the "default" config contains key/values occurring before
	// the first named config appears
	def := cfgs[""]

	number, _ := def.IntOrDefault("number", 42)
	fmt.Println("number =", number)

	every, _ := def.DurationOrDefault("every", time.Minute * 9)
	fmt.Println("every =", every)

	db := cfgs["database"]

	username, _ := db.StringOrDefault("username", "not-admin")
	fmt.Println("username =", username)

	port, _ := db.IntOrDefault("port", 8086)
	fmt.Println("port =", port)

	// easily use a default in the case of a missing key/value pair
	ip, _ := db.IPOrDefault("ip", net.ParseIP("127.0.0.1"))
	fmt.Println("ip =", ip)

	log := cfgs["log"]

	path, _ := log.FilePathOrDefault("path", "./log.out")
	fmt.Println("path =", path)

	level, _ := log.StringOrDefault("level", "debug")
	fmt.Println("level =", level)
}
