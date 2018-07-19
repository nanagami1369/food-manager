package main

import (
	"fmt"
	"time"
)
func dayonmonth()int {
	var re int
	year := time.Now().Year();
	switch    time := time.Now().Month(); time {
	case 4, 6, 9, 11:
		re = 30;
	case 2:
		if year%4 == 0 && year%100 != 0 || year%400 == 0 {
			re = 29
		} else {
			re = 28
		}
	case 1, 3, 5, 7, 8,10, 12: re = 31
	}
	return re
}
func main() {
	fmt.Println(dayonmonth())
}