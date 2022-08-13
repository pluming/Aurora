package test

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {

	//getenv := os.Getenv("PATH")
	//
	//l := strings.Split(getenv, ";")
	//
	//l = append(l,"D:\\soft\\Redis-x64-5.0.14.1")
	//
	//err := os.Setenv("PATH", strings.Join(l, ";"))
	//if err != nil {
	//	t.Fatal(err)
	//}
	getenv2 := os.Getenv("PATH")

	t.Log(getenv2)

}
