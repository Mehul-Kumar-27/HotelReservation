package main

import (
	"fmt"

	"github.com/Mehul-Kumar-27/HotelReservation/types"
	"github.com/brianvoe/gofakeit/v6"
)




func main() {

	var user types.User
	gofakeit.Struct(&user)
	fmt.Println(user.Password)

	// for i := 0; i < 10; i++ {
	// 	response, err := http.Get("http://localhost:8080/api/v1/admin")
	// 	if err != nil {
	// 		fmt.Println("Error:", err)
	// 		return
	// 	}
	// 	defer response.Body.Close()

	// 	// Read the response body
	// 	responseBody, err := ioutil.ReadAll(response.Body)
	// 	if err != nil {
	// 		fmt.Println("Error reading response:", err)
	// 		return
	// 	}

	// 	// Print the response body
	// 	fmt.Println("Response Body:", string(responseBody))
	// }
	// Make a GET request

}
