package main

import (
	cpsmpp "cp_sms_gateway/smpp"
)

func main() {
	cpsmpp.StartServer(":8080")

}
// echo -e -n '\x41\x05\x06\x00' > /dev/tcp/localhost/8080
