package main

import (
	"fmt"
	"simpleCode/machine"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	//fmt.Println(machine.IsMatch("special", "aaa", "a"))
	//fmt.Println(machine.IsMatch("special", "aaa", "*"))
	//fmt.Println(machine.IsMatch("special", "adceb", "*a*b"))
	//fmt.Println(machine.IsMatch("special", "acdcb", "a*c?b"))
	//fmt.Println(machine.IsMatch("common", "aa", "a*"))
	//fmt.Println(machine.IsMatch("common", "ab", ".*"))
	//fmt.Println(machine.IsMatch("common", "aab", "a*b+c*d?e"))
	fmt.Println(machine.IsMatch("common", "aab", "a*abac"))
	//fmt.Println(machine.IsMatch("common", "a", "ccc.*..a*"))
	//fmt.Println(machine.IsMatch("common", "aaba", "ab*a*c*a"))
}
