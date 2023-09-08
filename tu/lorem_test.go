package tu

import (
	"bytes"
	"fmt"
	"testing"
)

func CountWords(s []byte) int {
	// return len(bytes.FieldsFunc(s, func(c rune)bool{
	// 	return c==' ' || c =='.' || c ==','  
	// }))
	return len(bytes.Fields(s))
}


func Test(t *testing.T) {
	s:= Lorem(7)
	fmt.Println(string(s))
	Equal(t, CountWords(s), 7)
	s= Lorem(1)
	fmt.Println(string(s))
	Equal(t, CountWords(s), 1)
	s= Lorem(700)
	fmt.Println(string(s))
	Equal(t, CountWords(s), 700)

}
