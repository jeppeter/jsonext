package main

import (
	"fmt"
	"github.com/jeppeter/jsonext"
	"io"
	"os"
)

func DumpAMap(key string, amap []interface{}, level int, writer io.Writer) error {
	var icnt int
	for icnt = 0; icnt < len(amap); icnt++ {
		v := amap[icnt]
		switch v.(type) {
		case string:
			s := fmt.Sprintf("level[%d][%s][%d]=[%s]\n", level, key, icnt, v.(string))
			writer.Write([]byte(s))
		case float64:
			s := fmt.Sprintf("level[%d][%s][%d]=[%f]\n", level, key, icnt, v.(float64))
			writer.Write([]byte(s))
		case int:
			s := fmt.Sprintf("level[%d][%s][%d]=[%d]\n", level, key, icnt, v.(int))
			writer.Write([]byte(s))
		case int64:
			s := fmt.Sprintf("level[%d][%s][%d]=[%d]\n", level, key, icnt, v.(int64))
			writer.Write([]byte(s))
		case float32:
			s := fmt.Sprintf("level[%d][%s][%d]=[%f]\n", level, key, icnt, v.(float32))
			writer.Write([]byte(s))
		case map[string]interface{}:
			DumpVMap(v.(map[string]interface{}), level+1, writer)
		case []interface{}:
			DumpAMap(fmt.Sprintf("[%s][%d]", key, icnt), v.([]interface{}), level+1, writer)
		}
	}
	return nil
}

func DumpVMap(vmap map[string]interface{}, level int, writer io.Writer) error {
	for k, v := range vmap {
		switch v.(type) {
		case string:
			s := fmt.Sprintf("level[%d][%s]=[%s]\n", level, k, v.(string))
			writer.Write([]byte(s))
		case float64:
			s := fmt.Sprintf("level[%d][%s]=[%f]\n", level, k, v.(float64))
			writer.Write([]byte(s))
		case int:
			s := fmt.Sprintf("level[%d][%s]=[%d]\n", level, k, v.(int))
			writer.Write([]byte(s))
		case int64:
			s := fmt.Sprintf("level[%d][%s]=[%d]\n", level, k, v.(int64))
			writer.Write([]byte(s))
		case float32:
			s := fmt.Sprintf("level[%d][%s]=[%f]\n", level, k, v.(float32))
			writer.Write([]byte(s))
		case map[string]interface{}:
			DumpVMap(v.(map[string]interface{}), level+1, writer)
		case []interface{}:
			DumpAMap(k, v.([]interface{}), level+1, writer)
		default:
			s := fmt.Sprintf("level[%d][%s]=[%v]\n", level, k, v)
			writer.Write([]byte(s))
		}
	}
	return nil
}

func main() {
	var vmap map[string]interface{}
	var err error
	var vs string
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
	DumpVMap(vmap, 0, os.Stdout)
	vs, err = jsonext.GetJsonValue("zz/path/non", vmap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not get zz/path/non [%s]\n", err.Error())
	}
	vs, err = jsonext.GetJsonValue("zz/path/hello", vmap)
	fmt.Fprintf(os.Stdout, "get [zz/path/hello]=[%s]\n", vs)
	vs, err = jsonext.GetJsonValue("array", vmap)
	fmt.Fprintf(os.Stdout, "get [array]=[%s]\n", vs)
	return
}
