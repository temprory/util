package main

import (
	"fmt"
	"github.com/temprory/util"
)

func main() {
	for n := int64(1000000); n < 100000000; n++ {
		s := util.B32ToString(n)
		n2 := util.B32ToNum(s)
		if n%1000000 == 0 {
			fmt.Println("-- 111:", n, s, n2)
		}

		if n != n2 {
			fmt.Println("-- 222:", n, n2)
			panic("n != n2")
		}
	}
}
