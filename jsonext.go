package jsonext

import (
	"encoding/json"
	"fmt"
	//"log"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	EPSILON_VALUE = float64(0.00000001)
)

func parseMessage(msg string) (map[string]interface{}, error) {
	var v map[string]interface{}
	err := json.Unmarshal([]byte(msg), &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func GetJsonMap(valstr string) (retv map[string]interface{}, err error) {
	retv, err = SafeParseMessage(valstr)
	if err != nil {
		return
	}
	err = nil
	return
}

func GetJsonArray(valstr string) (retv []interface{}, err error) {
	dec := json.NewDecoder(strings.NewReader(valstr))
	err = dec.Decode(&retv)
	if err != nil {
		return
	}
	err = nil
	return
}

func SafeParseMessage(fmsg string) (map[string]interface{}, error) {
	v, err := parseMessage(fmsg)
	if err != nil {
		pmsg := `"` + fmsg + `"`
		//pmsg := fmsg
		cmsg, err := strconv.Unquote(pmsg)
		if err != nil {
			cmsg = fmsg
		}
		v, err = parseMessage(cmsg)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func __FormatLevel(level int) string {
	s := ""
	for i := 0; i < level; i++ {
		s += fmt.Sprintf("  ")
	}
	return s
}

type FormatClass interface {
	Format(level int, keyname string, value interface{}) (string, error)
	SupportType() string
}

var (
	suppFormatClass []FormatClass
)

type Float64FormatClass struct {
	FormatClass
}

func (f Float64FormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	s := ""
	ival := int(value.(float64))
	if math.Abs((float64(ival) - (value.(float64)))) < EPSILON_VALUE {
		s += fmt.Sprintf(" %d", ival)
	} else {
		s += fmt.Sprintf(" %f", value.(float64))
	}
	return s, nil
}

func (f Float64FormatClass) SupportType() string {
	return "float64"
}

type Float32FormatClass struct {
	FormatClass
}

func (f Float32FormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	s := ""
	ival := int(value.(float32))
	if math.Abs((float64(ival) - float64((value.(float32))))) < EPSILON_VALUE {
		s += fmt.Sprintf(" %d", ival)
	} else {
		s += fmt.Sprintf(" %f", value.(float32))
	}
	return s, nil
}

func (f Float32FormatClass) SupportType() string {
	return "float32"
}

type ArrayFormatClass struct {
	FormatClass
}

func (f ArrayFormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	s := ""
	s += "["
	for i, aa := range value.([]interface{}) {
		if i != 0 {
			s += ","
		}
		/*we do not format any keyname*/
		arrays, err := __FormatValue(level+1, "", aa)
		if err != nil {
			return "", err
		}
		s += arrays
	}
	s += "]"
	return s, nil
}

func (f ArrayFormatClass) SupportType() string {
	return "[]interface {}"
}

type MapStringFormatClass struct {
	FormatClass
}

func (f MapStringFormatClass) Format(level int, keyname string, v interface{}) (string, error) {
	var curmap map[string]interface{}
	var sortkeys []string
	s := ""

	curmap = v.(map[string]interface{})
	if len(keyname) != 0 {
		s += __FormatName(level, keyname)
	}

	s += "{"

	for k := range curmap {
		sortkeys = append(sortkeys, k)
	}
	sort.Strings(sortkeys)

	for i, kk := range sortkeys {
		kv, ok := curmap[kk]
		if !ok {
			err := fmt.Errorf("can not find key (%s)", kk)
			return "", err
		}
		if i != 0 {
			s += ",\n"
		} else {
			s += "\n"
		}
		ks, err := __FormatJsonValue(level+1, kk, kv)
		if err != nil {
			return "", err
		}
		s += ks
	}
	s += "\n"
	s += __FormatLevel(level)
	s += "}"
	return s, nil
}

func (f MapStringFormatClass) SupportType() string {
	return "map[string]interface {}"
}

type StringFormatClass struct {
	FormatClass
}

func (f StringFormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	s := ""
	s += fmt.Sprintf(" %s", strconv.Quote(value.(string)))
	return s, nil
}

func (f StringFormatClass) SupportType() string {
	return "string"
}

type BoolFormatClass struct {
	FormatClass
}

func (f BoolFormatClass) Format(level int, keyname string, value interface{}) (string, error) {
	s := ""
	if value.(bool) {
		s += "true"
	} else {
		s += "false"
	}
	return s, nil
}

func (f BoolFormatClass) SupportType() string {
	return "bool"
}

var (
	formatMap map[string]FormatClass
)

func init() {
	formatMap = make(map[string]FormatClass)
	arrcls := ArrayFormatClass{}
	formatMap[arrcls.SupportType()] = arrcls
	fl32cls := Float32FormatClass{}
	formatMap[fl32cls.SupportType()] = fl32cls
	fl64cls := Float64FormatClass{}
	formatMap[fl64cls.SupportType()] = fl64cls
	scls := StringFormatClass{}
	formatMap[scls.SupportType()] = scls
	boolcls := BoolFormatClass{}
	formatMap[boolcls.SupportType()] = boolcls
	mapstrcls := MapStringFormatClass{}
	formatMap[mapstrcls.SupportType()] = mapstrcls
}

func __FormatName(level int, keyname string) string {
	s := ""
	s += __FormatLevel(level)
	s += fmt.Sprintf("\"%s\" : ", keyname)
	return s
}

func __FormatValue(level int, keyname string, value interface{}) (string, error) {
	var err error
	var typestr string
	s := ""
	typestr = reflect.TypeOf(value).String()
	fcls, ok := formatMap[typestr]
	if !ok {
		err := fmt.Errorf("(%s) support type %s", keyname, typestr)
		return "", err
	}

	ss, err := fcls.Format(level, keyname, value)
	if err != nil {
		return "", err
	}

	s += ss
	return s, nil
}

func __FormatValueBasic(level int, keyname string, value interface{}) (string, error) {
	var s string
	var err error
	//Debug("[%d].(%s) = %q", level, keyname, value)
	s = ""
	s += __FormatName(level+1, keyname)
	sets, err := __FormatValue(level, keyname, value)
	if err != nil {
		return "", err
	}
	s += sets

	return s, nil
}

func __FormatJsonValue(level int, keyname string, value interface{}) (string, error) {
	var err error
	var s, ss string
	s = ""
	//Debug("[%d].(%s)type(%s) = %q", level, keyname, reflect.TypeOf(value).String(), value)

	switch value.(type) {
	case map[string]interface{}:
		ss, err = FormatJsonValue(level+1, keyname, value.(map[string]interface{}))
		if err != nil {
			return "", err
		}
		s += ss
	default:
		ss, err = __FormatValueBasic(level, keyname, value)
		if err != nil {
			return "", err
		}
		s += ss
	}
	//Debug("end[%d].(%s)type(%s) = %q", level, keyname, reflect.TypeOf(value).String(), value)
	return s, nil
}

func FormatJsonValue(level int, keyname string, v map[string]interface{}) (string, error) {
	var curmap map[string]interface{}
	var sortkeys []string
	s := ""

	curmap = v
	if len(keyname) != 0 {
		s += __FormatName(level, keyname)
	}

	s += "{"

	for k := range curmap {
		sortkeys = append(sortkeys, k)
	}
	sort.Strings(sortkeys)

	for i, kk := range sortkeys {
		kv, ok := curmap[kk]
		if !ok {
			err := fmt.Errorf("can not find key (%s)", kk)
			return "", err
		}
		if i != 0 {
			s += ",\n"
		} else {
			s += "\n"
		}
		ks, err := __FormatJsonValue(level, kk, kv)
		if err != nil {
			return "", err
		}
		s += ks
	}
	s += "\n"
	s += __FormatLevel(level)
	s += "}"
	return s, nil
}

func FormatJsonArray(level int, keyname string, v []interface{}) (string, error) {
	var s string
	var err error
	s, err = __FormatJsonValue(level, keyname, v)
	if err != nil {
		return "", err
	}
	return s, nil
}

func writeToFile(infile string, jsonbytes []byte) error {
	fw, err := os.OpenFile(infile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fw.Close()
	totalw := 0
	for totalw < len(jsonbytes) {
		n, err := fw.Write(jsonbytes[totalw:])
		if err != nil {
			return err
		}
		totalw += n
	}
	return nil
}

func WriteJson(infile string, v map[string]interface{}) error {
	jsonstring, err := FormatJsonValue(0, "", v)
	if err != nil {
		return err
	}
	jsonbytes := []byte(jsonstring)
	return writeToFile(infile, jsonbytes)
}

func WriteJsonString(infile string, val string) error {
	valmap, err := SafeParseMessage(val)
	if err != nil {
		return err
	}

	return WriteJson(infile, valmap)
}

func SetJsonValue(path, typestr, value string, v map[string]interface{}) (map[string]interface{}, error) {
	var pathext []string
	var tmpext []string
	var err error
	var mapv map[string]interface{}
	var arrayv []interface{}

	tmpext = strings.Split(path, "/")
	for _, a := range tmpext {
		if len(a) > 0 {
			/*we fil the path*/
			pathext = append(pathext, a)
		}
	}

	curmap := v
	if len(pathext) == 0 {
		if typestr != "map" {
			err = fmt.Errorf("invalid path(%s) and type(%s)", path, typestr)
			return nil, err
		}
		mapv, err = GetJsonMap(value)
		if err != nil {
			return nil, err
		}
		v = mapv
		return v, nil
	} else {
		for i, curpath := range pathext {
			if i == (len(pathext) - 1) {
				/*this is the last one ,so we set the value*/
				switch typestr {
				case "string":
					curmap[curpath] = value
				case "float64":
					fval, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return nil, err
					}
					curmap[curpath] = fval
				case "map":
					mapv, err = GetJsonMap(value)
					if err != nil {
						return nil, err
					}
					curmap[curpath] = mapv
				case "array":
					arrayv, err = GetJsonArray(value)
					if err != nil {
						return nil, err
					}
					curmap[curpath] = arrayv
				default:
					err = fmt.Errorf("unknown type %s", typestr)
					return nil, err
				}
				return v, nil
			}
			curval, ok := curmap[curpath]
			if !ok {
				/*we make the next use*/
				curmap[curpath] = make(map[string]interface{})
				curmap = curmap[curpath].(map[string]interface{})
			} else {
				switch curval.(type) {
				case map[string]interface{}:
					curmap = curval.(map[string]interface{})
				default:
					/*we make the map string*/
					curval = make(map[string]interface{})
					curmap[curpath] = curval.(map[string]interface{})
					curmap = curval.(map[string]interface{})
				}
			}
		}
	}

	err = fmt.Errorf("unknown path %s", path)
	return nil, err
}

func DeleteJsonValue(path string, v map[string]interface{}, force int) (map[string]interface{}, error) {
	var pathext []string
	var tmpext []string
	var err error

	tmpext = strings.Split(path, "/")
	for _, a := range tmpext {
		if len(a) > 0 {
			pathext = append(pathext, a)
		}
	}

	curmap := v
	if len(pathext) > 0 {
		for i, curpath := range pathext {
			curval, ok := curmap[curpath]
			if !ok {
				if force > 0 {
					return v, nil
				}
				err = fmt.Errorf("can not find (%s) value", path)
				return nil, err
			}
			if i == (len(pathext) - 1) {
				/*this is the last one ,so we set the value*/
				delete(curmap, curpath)
				return v, nil
			}
			switch curval.(type) {
			case map[string]interface{}:
				curmap = curval.(map[string]interface{})
			default:
				if force > 0 {
					delete(curmap, curpath)
					return v, nil
				}
				err = fmt.Errorf("can not handle path %s", curpath)
				return nil, err
			}
		}
	} else {
		/*we set the null for total delete*/
		return nil, nil
	}

	if force > 0 {
		return v, nil
	}

	err = fmt.Errorf("unknown path %s", path)
	return nil, err

}

func GetJson(infile string) (map[string]interface{}, error) {
	var v map[string]interface{}
	fp, err := os.Open(infile)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	dec := json.NewDecoder(fp)

	err = dec.Decode(&v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func __GetJsonValue(path string, v map[string]interface{}) (val string, err error) {
	var pathext []string
	var curmap map[string]interface{}
	var tmpext []string

	val = ""
	err = nil
	tmpext = strings.Split(path, "/")
	for _, a := range tmpext {
		if len(a) > 0 {
			pathext = append(pathext, a)
		}
	}

	curmap = v
	if len(pathext) > 0 {

		for i, curpath := range pathext {
			var curval, cval interface{}
			var ok bool
			curval, ok = curmap[curpath]
			if !ok {
				err = fmt.Errorf("can not find (%s) in %s", curpath, path)
				return
			}

			if i == (len(pathext) - 1) {
				switch curval.(type) {
				case int:
					val = fmt.Sprintf("%d", curval)
				case uint32:
					val = fmt.Sprintf("%d", curval)
				case uint64:
					val = fmt.Sprintf("%d", curval)
				case float64:
					val = fmt.Sprintf("%f", curval)
				case float32:
					val = fmt.Sprintf("%f", curval)
				case map[string]interface{}:
					val, err = FormatJsonValue(0, "", curval.(map[string]interface{}))
					if err != nil {
						return
					}
				case []interface{}:
					val, err = __FormatValue(0, "", curval)
					if err != nil {
						return
					}
				default:
					val = fmt.Sprintf("%s", curval)
				}
				err = nil
				return
			}

			switch curval.(type) {
			case map[string]interface{}:
				cval, ok = curval.(map[string]interface{})
				if !ok {
					err = fmt.Errorf("can not parse in (%s) for path(%s)", curpath, path)
					return
				}
			case []interface{}:
				cval, ok = curval.([]interface{})
				if !ok {
					err = fmt.Errorf("can not parse in (%s) for path(%s)", curpath, path)
					return
				}
			default:
				err = fmt.Errorf("type of (%s) error", path)
				return
			}
			curmap = cval.(map[string]interface{})
		}
	} else {
		/*we format total*/
		val, err = __FormatValue(0, "", v)
		if err != nil {
			return
		}
		err = nil
		return
	}

	err = fmt.Errorf("can not find (%s) all over", path)
	return
}

func GetJsonValue(path string, v map[string]interface{}) (string, error) {
	return __GetJsonValue(path, v)
}

func GetJsonStruct(valstr string, v interface{}) error {
	var err error
	dec := json.NewDecoder(strings.NewReader(valstr))
	err = dec.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

func FormJsonStruct(v interface{}) (valstr string, err error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	valstr = string(b)
	err = nil
	return
}

func GetJsonValueDefault(infile string, path string, defval string) string {
	val := defval
	fp, err := os.Open(infile)
	if err != nil {
		return val
	}
	defer fp.Close()
	dec := json.NewDecoder(fp)

	for {
		var v map[string]interface{}
		err = dec.Decode(&v)
		if err != nil {
			return val
		}

		getval, err := __GetJsonValue(path, v)
		if err == nil {
			val = getval
			return val
		}
	}

	return val
}
