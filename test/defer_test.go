package test

import "fmt"
import "testing"

func TestDeferFuncArgs(t *testing.T) {
	a()
}

func a() {
	i := 1
	defer fmt.Println(i)
	i++
	return
}

func TestDeferFuncOrder(t *testing.T) {
	b()
}

func b() {
	for i := 0; i < 4; i++ {
		defer fmt.Println(i)
	}
}

func TestDeferReadAndAssignReturingValue(t *testing.T) {
	fmt.Println(c())
}

func c() (i int) {
	defer func() {
		i++
	}()
	return 1
}