package jsonlang

import(
	"fmt"
	"strings"
)

// Helper function
func int_num(n interface{}) int {
	return int(n.(float64))
}

// ret_if(val)
// Panics if val is not bool.
func ret_if(a ...interface{}) bool {
	return a[0].(bool) == true
}

// ret_ifn(val)
// Panics if val is not bool.
func ret_ifn(a ...interface{}) bool {
	return a[0].(bool) == false
}

// Creates a label, so we can jump to it with jump_if.
// label("labname")
// Panics if argument is not string.
func label(labels map[string]int, i int, args ...interface{}) {
	labels[args[0].(string)] = i
}

// Assigns a value to a reference.
// set(&varname, value)
func set(a ...interface{}) {
	argCount("set", 2, a...)
	val, ok := a[0].(Ref)
	if !ok { panic("You can not set a nonreference.") }
	val.Set(a[1])
}

// jump_if(val, "labname")
// Special builtin.
func jump_if(i *int, val interface{}, label interface{}, labels map[string]int, negate bool) {
	verdict := val.(bool) == true
	if negate { verdict = !verdict }
	if verdict {
		*i = labels[label.(string)]
	}
}

func prepareEx(s string, a ...interface{}) Ref {
	if len(a) < 2 { panic(fmt.Sprintf("Not enough args for %v exists.", s)) }
	for _, v := range a {
		_, ok := v.(Ref)
		if ! ok { panic(fmt.Sprintf("All args must be Refs at %v exists.", s)) }
	}
	ref, ok := a[0].(Ref)
	if !ok { panic(fmt.Sprintf("Result is not a reference at %v exists.", s)) }
	return ref
}

// exists(&var1, &var2)
// Panics if the arguments is a non-reference.
// Returns true if the references exists.
func exists(a ...interface{}) {
	ref := prepareEx("", a...)
	ref.Set(a[1].(Ref).Exists())
}

// none_exists(&result, &var1, &var2, ...)
// Panics if any of the arguments is a non-reference.
// Returns true if any of the references exists.
func any_exists(a ...interface{}) {
	ref := prepareEx("any", a...)
	for _, v := range a[1:] {
		if v.(Ref).Exists() { ref.Set(true); return }
	}
	ref.Set(false)
}

// none_exists(&result, &var1, &var2, ...)
// Panics if any of the arguments is a non-reference.
func none_exists(a ...interface{}) {
	ref := prepareEx("none", a...)
	for _, v := range a[1:] {
		if v.(Ref).Exists() { ref.Set(false); return }
	}
	ref.Set(true)
}

// all_exists(&result, &var1, &var2, ...)
// Panics if any of the arguments is a non-reference.
func all_exists(a ...interface{}) {
	ref := prepareEx("all", a...)
	for _, v := range a[1:] {
		if !v.(Ref).Exists() { ref.Set(false); return }
	}
	ref.Set(true)
}

// err_if(test, panic_msgs ...interface{})
// Panics only if test is bool and true.
// Panics with converting to and joining to string all panic_msgs.
func err_if(a ...interface{}) {
	verdict, is_bool := a[0].(bool)
	if !is_bool { return }
	if !verdict { return }
	sl := []string{}
	for _, v := range a[1:] {
		sl = append(sl, fmt.Sprint(v))
	}
	panic(strings.Join(sl, " "))
}

func print(a ...interface{}) {
	fmt.Print(a...)
}

// println(anything, "really")
func println(a ...interface{}) {
	fmt.Println(a...)
}

func firstRef(func_name string, a ...interface{}) Ref {
	ref, ok := a[0].(Ref)
	if !ok { panic("First argument must be a Reference at " + func_name) }
	return ref
}

// Helper function.
func minArgs(func_name string, count int, a ...interface{}) {
	if len(a) < count { panic(fmt.Sprint("%v must be called with at least %v args.", func_name, count)) }
}

// Helper function.
func argCount(func_name string, count int, a ...interface{}) {
	if len(a) != count { panic(fmt.Sprint("%v must be called with %v args.", func_name, count)) }
}

// push(&slice, val1, val2, ...)
// Pushes one or more values into an slice.
// Panics if slice is undefined.
// Panics if slice is not []interface{}.
func push(a ...interface{}) {
	minArgs("push", 2)
	ref := firstRef("push", a...)
	sl := ref.DereferStrict().([]interface{})
	sl = append(sl, a[1])
	ref.Set(sl)
}

// set_slice_index(&slice, 3, val)
// Panics if out of bound.
func set_slice_index(a ...interface{}) {
	argCount("set_map_key", 3)
}

// set_map_key(map, "key", val)
// Panics if map is not a map[string]interface{}
func set_map_key(a ...interface{}) {
	argCount("set_map_key", 3)
	ma := a[0].(map[string]interface{})
	key := a[1].(string)
	ma[key] = a[2]
}

// delete_map_key(map, "key")
// panics if map is not map[string]interface{}
func delete_map_key(a ...interface{}) {
	argCount("delete_map_key", 2)
	ma := a[0].(map[string]interface{})
	key := a[1].(string)
	delete(ma, key)
}

// slice(&slice, 1, 3)
// Panics if slice is undefined. Panics if out of bound. Panics if slice is not []interface{}.
func slice(a ...interface{}) {
	ref := firstRef("slice", a...)
	beg := int_num(a[1])
	end := int_num(a[2])
	sl := ref.DereferStrict()
	ref.Set(sl.([]interface{})[beg:end])
}