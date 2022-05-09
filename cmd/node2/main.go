package main

import "udpforward/forward"

func main() {
	forward.Start("0.0.0.0:11111", "127.0.0.1:1194", "72c5eb83afc3321d1737931e0aa25273718b4591383d13666d49b346eb29178a", "9b0dc2e5a8e562c2c2438e7182449b58")
}
