package main

//
//import (
//	"fmt"
//	"math"
//	"time"
//)
//
//func main() {
//	z := 50000000
//	a := make(map[int]int, z)
//	b := make([]int, z)
//
//	for i := 0; i < z; i++ {
//		a[i] = i
//		b[i] = i
//	}
//
//	t0 := time.Now()
//	for key, value := range a {
//		if key != value { // never happens
//			fmt.Println("a", key, value)
//		}
//	}
//	d0 := time.Now().Sub(t0)
//
//	t1 := time.Now()
//	for key, value := range b {
//		if key != value { // never happens
//			fmt.Println("b", key, value)
//		}
//	}
//	d1 := time.Now().Sub(t1)
//
//	fmt.Println(
//		"a:", d0,
//		"b:", d1,
//		"diff:", math.Max(float64(d0), float64(d1))/math.Min(float64(d0), float64(d1)),
//	)
//}
