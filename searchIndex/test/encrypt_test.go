package test

import (
	"fmt"
	"searchIndex/utils"
	"testing"
)

func Test_DES(t *testing.T) {

	result, _ := utils.DESDecryptString("9K/zqvMCvqwbe238iPjbGfFN68rnTOZBhaW9SlmUXWrA0nQfP6rEBt5W195V3XsQj+ilx0yLlN2OD+/ZTFavSplELYXLPaxIinAaxd1pjWPjQzJb6yt9lkNwpQYeNxjlQ+KXsNumRts=")

	fmt.Println(result)

}
