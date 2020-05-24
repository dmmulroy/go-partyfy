package main

import (
	"fmt"
	_ "image/gif"
	_ "image/png"
	"os"
)

func main() {
	img, err := os.Open("docker.png")

	if err != nil {
		fmt.Println("Error opening image: ", err)
	}

	defer img.Close()

	if err := partyfy(img); err != nil {
		fmt.Println(err)
	}

}
