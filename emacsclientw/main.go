
package main

import (
	"fmt"
	"os"
	"net"
	"bufio"
	"strings"
)

func main() {
	var err error

	serverAddr, authString := get_server_info()	

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Fprint(conn, ` -auth `)
	fmt.Fprint(conn, authString)
	if os.Args[1] != "" {
		fmt.Fprint(conn, ` -file `)
		fmt.Fprint(conn, quote_argument(os.Args[1]))
	}
	fmt.Fprint(conn, "\n")

	// wait closing emacs
	connScanner := bufio.NewScanner(conn)
	for connScanner.Scan() {
		fmt.Println(connScanner.Text())
	}
}

func quote_argument(s string) string {
	r := strings.NewReplacer(
		"&", "&&",
		"-", "&-",
		" ", "&_",
		"\n", "&n")
	return r.Replace(s)
}

func get_server_info() (serverAddr string, authString string) {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("APPDATA")
	}
	
	fp, err := os.Open(home + `\.emacs.d\server\server`)
	if err != nil {
		panic(err)
	}
	defer fp.Close()	

	scanner := bufio.NewScanner(fp)
	if scanner.Scan() {
		// 1st line
		serverAddr = strings.Split(scanner.Text(), " ")[0]
	}
	if scanner.Scan() {
		// 2nd line
		authString = scanner.Text()
	}

	return
}
