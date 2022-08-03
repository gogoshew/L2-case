package main

import "L2_case/pattern"

func main() {
	// 1. Facade
	//fmt.Println(pattern.Man.ToDoFacade(pattern.Man{}))

	// 2. Builder
	//pattern.ToDoBuilderPattern()

	// 3. Visitor
	//pattern.VisitorPattern()
	// 4. Command
	pattern.CommandPattern()
	// 6. Factory method
	//glock, _ := pattern.GetGun("Glock")
	//mp, _ := pattern.GetGun("MP5")
	//pattern.FactoryPrint(glock)
	//pattern.FactoryPrint(mp)
}
