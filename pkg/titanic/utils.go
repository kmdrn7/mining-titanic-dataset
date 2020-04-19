package titanic

import (
	"fmt"
	"os"
)

func Throw(err error, msg string){	
	fmt.Fprintf(os.Stderr, msg+  ". %v", err)
	os.Exit(1)
}

func CheckDuplicateValue(src map[string]int, val string) (bool){
	if _,ok := src[val]; ok {
		return true
	}
	return false
}

func IterateData(){
	//iterator := dataset.ValuesIterator(dataframe.ValuesOptions{0, 1, true})
	//for {
		//row, vals, _ := iterator()
		//if row == nil {
			//break
		//}
		//fmt.Println(vals)
	//}
}
