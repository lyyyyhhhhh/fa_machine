package main

import (
	"fa_machine/service"
	"fmt"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	//fmt.Println(service.IsMatch("special", "aaa", "a"))
	//fmt.Println(service.IsMatch("special", "aaa", "*"))
	//fmt.Println(service.IsMatch("special", "adceb", "*a*b"))
	//fmt.Println(service.IsMatch("special", "acdcb", "a*c?b"))
	//fmt.Println(service.IsMatch("common", "aa", "a*"))
	//fmt.Println(service.IsMatch("common", "ab", ".*"))
	//fmt.Println(service.IsMatch("common", "aab", "a*b+c*d?e"))
	fmt.Println(service.IsMatch("common_self", "aab", "a*abac"))
	//fmt.Println(service.IsMatch("common", "a", "ccc.*..a*"))
	//fmt.Println(service.IsMatch("common", "aaba", "ab*a*c*a"))
}
