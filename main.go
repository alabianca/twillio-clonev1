package main

func main() {
	client := Client{
		addr: "localhost:16004",
	}

	client.Connect()
}
