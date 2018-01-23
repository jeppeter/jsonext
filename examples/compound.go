package main

import (
	"fmt"
	"github.com/jeppeter/jsonext"
	"io"
	"os"
)

func DumpAMap(key string, amap []interface{}, level int, writer *io.Writer) error {

}

func DumpVMap(vmap map[string]interface{}, level int, writer *io.Writer) error {
	for k, v := range vmap {
		switch ct := v.(type) {
		case string:
			s := fmt.Sprintf("level[%d][%s]=[%s]\n", level, k, v.(string))
			writer.Write(s)
		case float64:
			s := fmt.Sprintf("level[%d][%s]=[%f]\n", level, k, v.(float64))
			writer.Write(s)
		case int:
			s := fmt.Sprintf("level[%d][%s]=[%d]\n", level, k, v.(int))
			writer.Write(s)
		case int64:
			s := fmt.Sprintf("level[%d][%s]=[%d]\n", level, k, v.(int64))
			writer.Write(s)
		case float32:
			s := fmt.Sprintf("level[%d][%s]=[%f]\n", level, k, v.(float32))
			writer.Write(s)
		case []interface{}:
			DumpAMap(k, v.([]interface{}), level+1, writer)
		case interface{}:
			DumpVMap(v, level+1, writer)
		default:
			s := fmt.Sprintf("level[%d][%s]=[%v]\n", level, k, v)
			writer.Write(s)
		}
	}
	return nil
}

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
		},
		"array" : ["ccc",20 ,333.90]
	}
	`
	vmap, err = jsonext.GetJsonMap(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse [%s] error[%s]\n", msg, err.Error())
		os.Exit(5)
	}

	return
}
