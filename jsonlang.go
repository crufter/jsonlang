package jsonlang

import (
	"encoding/json"
	"fmt"
	"github.com/opesun/jsonp"
	"reflect"
	"regexp"
	"strings"
)

const (
	Marker = "."
	Parens = true // println($x) vs println, $x
)

type node struct {
	Text   string
	Inside bool // Inside regexp find, means IsString here.
}

// TODO: Very similar to what is in the package opesun/require, maybe we could use that here.
func split(src string, pos [][]int) []node {
	sl := []node{}
	last := 0
	for _, v := range pos {
		sl = append(sl, node{src[last:v[0]], false})
		sl = append(sl, node{src[v[0]:v[1]], true})
		last = v[1]
	}
	sl = append(sl, node{src[last:], false})
	return sl
}

func quote(src []node) error {
	for i, v := range src {
		if !v.Inside {
			continue
		}
		non_string_const := false
		if i == 0 {
			non_string_const = true
		} else {
			prev_node := src[i-1]
			last_char := prev_node.Text[len(prev_node.Text)-1]
			if last_char != '"' {
				non_string_const = true
			}
		}
		if non_string_const {
			src[i].Text = "\"" + Marker + v.Text + "\""
			continue
		} else if string(v.Text[0]) == Marker {
			return fmt.Errorf("String constant \"%v\" starts with the Marker: \"%v\".", v.Text, Marker)
		}
	}
	return nil
}

func join(src []node) string {
	s := []string{}
	for _, v := range src {
		s = append(s, v.Text)
	}
	return strings.Join(s, "")
}

func Compile(src string) ([]interface{}, error) {
	r := regexp.MustCompile("[$&a-zA-Z._]+")
	pos := r.FindAllIndex([]byte(src), -1)
	s := split(src, pos)
	err := quote(s)
	if err != nil {
		return nil, err
	}
	src = join(s)
	src = strings.Replace(src, ";", "],[", -1)
	if Parens {
		src = strings.Replace(src, "(", ", ", -1)
		src = strings.Replace(src, ")", "", -1)
	}
	src = "[[" + src + "]]"
	//fmt.Println(src)
	var v interface{}
	err = json.Unmarshal([]byte(src), &v)
	if err != nil {
		return nil, err
	}
	sl := v.([]interface{})
	l := len(sl)
	if len(sl[l-1].([]interface{})) == 0 {
		sl = sl[:l-1]
	}
	return sl, nil
}

type Ref struct {
	vars map[string]interface{}
	name string
}

func (r Ref) Derefer() interface{} {
	val, _ := jsonp.Get(r.vars, r.name)
	return val
}

func (r Ref) DereferStrict() interface{} {
	val, ok := jsonp.Get(r.vars, r.name)
	if !ok {
		panic(r.name + " is undefined.")
	}
	return val
}

func (r Ref) Set(a interface{}) {
	// TODO: this will not handle []s, set should be implemented in opesun/jsonp.
	sl := strings.Split(r.name, ".")
	l := len(sl)
	if l == 1 {
		r.vars[r.name] = a
	} else {
		// TODO: This is not generic yet!
		val, _ := jsonp.Get(r.vars, strings.Join(sl[:l-1], "."))
		name := sl[l-1:][0]
		val.(map[string]interface{})[name] = val
	}
}

//func (r Ref) Pointer() interface{} {
//	return nil
//}

func (r Ref) Type() reflect.Type {
	val, _ := jsonp.Get(r.vars, r.name)
	return reflect.TypeOf(val)
}

func (r Ref) Exists() bool {
	_, ok := jsonp.Get(r.vars, r.name)
	return ok
}

func eval_rec(i interface{}, vars map[string]interface{}, nested bool) interface{} {
	switch val := i.(type) {
	case string:
		if string(val[0]) == Marker {
			if string(val[1]) == "&" { // Reference
				if nested {
					panic("Can't interpret reference in map or slice.")
				} else {
					return Ref{vars, string(val[2:])}
				}
			} else { // Variable
				val, _ := jsonp.Get(vars, val[1:])
				return val
			}
		} else {
			return i
		}
	case map[string]interface{}:
		for i, v := range val {
			val[i] = eval_rec(v, vars, true)
		}
	case []interface{}:
		for i, v := range val {
			val[i] = eval_rec(v, vars, true)
		}
	default:
		return i
	}
	return i
}

func eval(i interface{}, vars map[string]interface{}) interface{} {
	return eval_rec(i, vars, false)
}

func evalArgs(a []interface{}, vars map[string]interface{}) []interface{} {
	ret := []interface{}{}
	for _, v := range a {
		ret = append(ret, eval(v, vars))
	}
	return ret
}

func Interpret(src []interface{}, vars map[string]interface{}, funcs map[string]func(...interface{})) (err error) {
	if vars == nil {
		vars = map[string]interface{}{}
	}
	if funcs == nil {
		funcs = map[string]func(...interface{}){}
	}
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf(fmt.Sprint(r))
		}
	}()
	i := 0
	labels := map[string]int{}
	l := len(src)
	for i < l {
		v := src[i]
		val := v.([]interface{})
		if len(val) == 0 {
			panic("Empty operation.")
		}
		func_name := val[0].(string)[1:]
		args := evalArgs(val[1:], vars)
		switch func_name {
		case "ret_if":
			if ret_if(args...) {
				return nil
			}
		case "ret_ifn":
			if ret_ifn(args...) {
				return nil
			}
		case "label":
			label(labels, i, args...)
		case "jump_if":
			jump_if(&i, args[0], args[1], labels, false)
		case "jump_ifn":
			jump_if(&i, args[0], args[1], labels, true)
		// Following functions need no special treatment, they could be implemented as user supplied functions too.
		// You can delete them out if you want, they are included for convenience only.
		case "exists":
			exists(args...)
		case "all_exists":
			all_exists(args...)
		case "any_exists":
			any_exists(args...)
		case "none_exists":
			none_exists(args...)
		case "set":
			set(args...)
		case "print":
			print(args...)
		case "println":
			println(args...)
		case "err_if":
			err_if(args...)
		case "push":
			push(args...)
		case "set_slice_index":
			set_slice_index(args...)
		case "set_map_key":
			set_map_key(args...)
		case "delete_map_key":
			delete_map_key(args...)
		case "slice":
			slice(args...)
		default:
			if function, has := funcs[func_name]; has {
				function(args...)
			} else {
				panic("Unkown function: " + func_name)
			}
		}
		i++
	}
	return
}
