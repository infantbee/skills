package main

import (
	"fmt"

	"github.com/linvon/cuckoo-filter"
)

func main() {
	cf := cuckoo.NewFilter(4, 9, 3900, cuckoo.TableTypePacked)
	fmt.Println(cf.Info())
	fmt.Println(cf.FalsePositiveRate())

	a := []byte("Asssssssssssssssssxxxxxxxxxxxxxxxxxxxxxxxx")
	cf.Add(a)
	fmt.Println(cf.Contain(a))
	fmt.Println(cf.Size())

	b := cf.Encode()
	ncf, _ := cuckoo.Decode(b)
	fmt.Println(ncf.Contain(a))
	fmt.Println("xxxxx: ", cf.Size())

	cf.Delete(a)
	fmt.Println(cf.Size())
}
