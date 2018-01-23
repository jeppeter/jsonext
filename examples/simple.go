package main

import (
	"fmt"
	"github.com/jeppeter/jsonext"
	"os"
)

func main() {
	var vmap map[string]interface{}
	var err error
	msg := `
	{
		"path"  : "new",
		"cc" : "www",
		"zz" : {
			"path" : {
				"hello" : "world"
			},
			"cc" : {
				"bon" : "jour"
			}
		}
	}
	`
	vmap, err = jsonext.GetJsonMap(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse [%s] error[%s]\n", msg, err.Error())
		os.Exit(5)
	}
	for k, v := range vmap {
		switch v.(type) {
		case string:
			fmt.Printf("[%s]=[%s]\n", k, v.(string))
		}
	}
	return
}
