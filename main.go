package main

func main() {
	client := Client{
		addr:     "localhost:16004",
		user:     "alexlabianca",
		password: "CoffeeBay$77z",
	}

	client.Connect()
}
