/*
 * Copyright 2023 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cast

import (
	"encoding/json"
	"log"
	"reflect"
	"sort"
	"strings"
	"sync"
	"unicode"
)

var (
	FAST = &fastEncoding{}
	JSON = &jsonEncoding{}
)

type jsonEncoding struct{}

// Convert converts src to dest using json encoding.
func (e *jsonEncoding) Convert(src, dest any) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}

type fastEncoding struct{}

// Convert converts src to dest using fast encoding.
func (e *fastEncoding) Convert(src, dest any) error {
	srcValue := reflect.ValueOf(src)
	if !srcValue.IsValid() || srcValue.IsNil() {
		return nil
	}
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(destValue)}
	}
	l := newMiddleValueList()
	defer middleValueListPool.Put(l)
	reflectValue(l, 0, srcValue)
	fromMiddleValue(l, l.List[0], destValue)
	return nil
}

type ValueType int

const (
	NilValueType   = ValueType(0)
	ValueValueType = ValueType(1)
	SliceValueType = ValueType(2)
	MapValueType   = ValueType(3)
)

type MiddleValueList struct {
	List []MiddleValue
}

func (b *MiddleValueList) Reset() {
	b.List[0] = MiddleValue{} // root
	b.List = b.List[:1]
}

type MiddleValue struct {
	Type   ValueType
	Name   string
	Value  reflect.Value
	Length int // 孩子数量
	First  int // 首个孩子
}

var middleValueListPool sync.Pool

func newMiddleValueList() *MiddleValueList {
	if v := middleValueListPool.Get(); v != nil {
		e := v.(*MiddleValueList)
		e.Reset()
		return e
	}
	return &MiddleValueList{
		List: make([]MiddleValue, 1, 256),
	}
}

func reflectValue(l *MiddleValueList, current int, v reflect.Value) {
	valueEncoder(v)(l, current, v)
}

func valueEncoder(v reflect.Value) encoderFunc {
	if !v.IsValid() {
		return func(l *MiddleValueList, current int, v reflect.Value) {}
	}
	return typeEncoder(v.Type())
}

var encoderCache sync.Map // map[string]string

type encoderFunc func(l *MiddleValueList, current int, v reflect.Value)

func typeEncoder(t reflect.Type) encoderFunc {
	if fi, ok := encoderCache.Load(t); ok {
		return fi.(encoderFunc)
	}

	// To deal with recursive types, populate the map with an
	// indirect func before we build it. This type waits on the
	// real func (f) to be ready and then calls it. This indirect
	// func is only used for recursive types.
	var (
		wg sync.WaitGroup
		f  encoderFunc
	)
	wg.Add(1)
	fi, loaded := encoderCache.LoadOrStore(t, encoderFunc(func(l *MiddleValueList, current int, v reflect.Value) {
		wg.Wait()
		f(l, current, v)
	}))
	if loaded {
		return fi.(encoderFunc)
	}

	// Compute the real encoder and replace the indirect func with it.
	f = newTypeEncoder(t)
	wg.Done()
	encoderCache.Store(t, f)
	return f
}

func validMapKey(key reflect.Value) (string, bool) {
	if key.Kind() != reflect.String {
		return "", false
	}
	return key.String(), true
}

func newTypeEncoder(t reflect.Type) encoderFunc {
	switch t.Kind() {
	case reflect.Interface:
		return func(l *MiddleValueList, current int, v reflect.Value) {
			if v.IsNil() {
				l.List[current] = MiddleValue{Type: NilValueType}
				return
			}
			reflectValue(l, current, v.Elem())
		}
	case reflect.Pointer:
		toMiddleValue := typeEncoder(t.Elem())
		return func(l *MiddleValueList, current int, v reflect.Value) {
			if v.IsNil() {
				l.List[current] = MiddleValue{Type: NilValueType}
				return
			}
			toMiddleValue(l, current, v.Elem())
		}
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return func(l *MiddleValueList, current int, v reflect.Value) {
			l.List[current] = MiddleValue{Type: ValueValueType, Value: v}
		}
	case reflect.Array, reflect.Slice:
		toMiddleValue := typeEncoder(t.Elem())
		return func(l *MiddleValueList, current int, v reflect.Value) {
			n := v.Len()
			p := &l.List[current]
			p.Type = SliceValueType
			if n == 0 {
				return
			}
			p.Length = n
			end := len(l.List)
			p.First = end
			for i := 0; i < n; i++ {
				l.List = append(l.List, MiddleValue{})
			}
			for i := 0; i < n; i++ {
				toMiddleValue(l, end+i, v.Index(i))
			}
		}
	case reflect.Map:
		toMiddleValue := typeEncoder(t.Elem())
		return func(l *MiddleValueList, current int, v reflect.Value) {
			n := v.Len()
			p := &l.List[current]
			p.Type = MapValueType
			if n == 0 {
				return
			}
			p.Length = n
			end := len(l.List)
			p.First = end
			for i := 0; i < n; i++ {
				l.List = append(l.List, MiddleValue{})
			}
			i := 0
			iter := v.MapRange()
			for iter.Next() {
				strKey, valid := validMapKey(iter.Key())
				if !valid {
					continue
				}
				toMiddleValue(l, end+i, iter.Value())
				l.List[end+i].Name = strKey
				i++
			}
		}
	case reflect.Struct:
		fields := cachedTypeFields(t)
		return func(l *MiddleValueList, current int, v reflect.Value) {
			n := len(fields.list)
			p := &l.List[current]
			p.Type = MapValueType
			if n == 0 {
				return
			}
			p.Length = n
			end := len(l.List)
			p.First = end
			for i := 0; i < n; i++ {
				l.List = append(l.List, MiddleValue{})
			}
			for j := range fields.list {
				f := &fields.list[j]
				fv := v
				breakNil := false
				for _, i := range f.index {
					if fv.Kind() == reflect.Pointer {
						if fv.IsNil() {
							breakNil = true
							break
						}
						fv = fv.Elem()
					}
					fv = fv.Field(i)
				}
				if breakNil {
					l.List[end+j] = MiddleValue{Type: NilValueType, Name: f.name}
					continue
				}
				f.encoder(l, end+j, fv)
				l.List[end+j].Name = f.name
			}
		}
	default:
		return func(l *MiddleValueList, current int, v reflect.Value) {}
	}
}

func fromMiddleValue(l *MiddleValueList, p MiddleValue, destValue reflect.Value) {
	switch p.Type {
	case NilValueType:
		return
	case ValueValueType:
		fromSimple(p.Value, destValue)
	case SliceValueType:
		fromSlice(l, p, destValue)
	case MapValueType:
		fromMap(l, p, destValue)
	default:
		log.Println("should never reach here")
	}
}

func valueInterface(l *MiddleValueList, p MiddleValue) interface{} {
	switch p.Type {
	case NilValueType:
		return nil
	case ValueValueType:
		return p.Value.Interface()
	case SliceValueType:
		data := l.List[p.First : p.First+p.Length]
		return arrayInterface(l, data)
	case MapValueType:
		data := l.List[p.First : p.First+p.Length]
		return objectInterface(l, data)
	default:
		log.Println("should never reach here")
		return nil
	}
}

func arrayInterface(l *MiddleValueList, ps []MiddleValue) []interface{} {
	r := make([]interface{}, len(ps))
	for i, p := range ps {
		r[i] = valueInterface(l, p)
	}
	return r
}

func objectInterface(l *MiddleValueList, pm []MiddleValue) map[string]interface{} {
	r := make(map[string]interface{}, len(pm))
	for _, p := range pm {
		r[p.Name] = valueInterface(l, p)
	}
	return r
}

func fromSimple(pv reflect.Value, destValue reflect.Value) {
	destValue = makeValue(destValue)
	destValue.Set(pv)

	//switch c := item[0]; c {
	//case 't', 'f': // true, false
	//	switch v.Kind() {
	//	case reflect.Bool:
	//		v.SetBool(value)
	//	case reflect.Interface:
	//		if v.NumMethod() == 0 {
	//			v.Set(reflect.ValueOf(value))
	//		} else {
	//			d.saveError(&UnmarshalTypeError{Value: "bool", Type: v.Type(), Offset: int64(d.readIndex())})
	//		}
	//	}
	//
	//case '"': // string
	//	s, ok := unquoteBytes(item)
	//	switch v.Kind() {
	//	case reflect.Slice:
	//		if v.Type().Elem().Kind() != reflect.Uint8 {
	//			d.saveError(&UnmarshalTypeError{Value: "string", Type: v.Type(), Offset: int64(d.readIndex())})
	//			break
	//		}
	//		b := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	//		n, err := base64.StdEncoding.Decode(b, s)
	//		if err != nil {
	//			d.saveError(err)
	//			break
	//		}
	//		v.SetBytes(b[:n])
	//	case reflect.String:
	//		v.SetString(string(s))
	//	case reflect.Interface:
	//		if v.NumMethod() == 0 {
	//			v.Set(reflect.ValueOf(string(s)))
	//		} else {
	//			d.saveError(&UnmarshalTypeError{Value: "string", Type: v.Type(), Offset: int64(d.readIndex())})
	//		}
	//	}
	//
	//default: // number
	//	s := string(item)
	//	switch v.Kind() {
	//	default:
	//		if v.Kind() == reflect.String && v.Type() == numberType {
	//			v.SetString(s)
	//			if !isValidNumber(s) {
	//				return fmt.Errorf("json: invalid number literal, trying to unmarshal %q into Number", item)
	//			}
	//			break
	//		}
	//		if fromQuoted {
	//			return fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal %q into %v", item, v.Type())
	//		}
	//		d.saveError(&UnmarshalTypeError{Value: "number", Type: v.Type(), Offset: int64(d.readIndex())})
	//	case reflect.Interface:
	//		n, err := d.convertNumber(s)
	//		if err != nil {
	//			d.saveError(err)
	//			break
	//		}
	//		if v.NumMethod() != 0 {
	//			d.saveError(&UnmarshalTypeError{Value: "number", Type: v.Type(), Offset: int64(d.readIndex())})
	//			break
	//		}
	//		v.Set(reflect.ValueOf(n))
	//
	//	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	//		n, err := strconv.ParseInt(s, 10, 64)
	//		if err != nil || v.OverflowInt(n) {
	//			d.saveError(&UnmarshalTypeError{Value: "number " + s, Type: v.Type(), Offset: int64(d.readIndex())})
	//			break
	//		}
	//		v.SetInt(n)
	//
	//	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	//		n, err := strconv.ParseUint(s, 10, 64)
	//		if err != nil || v.OverflowUint(n) {
	//			d.saveError(&UnmarshalTypeError{Value: "number " + s, Type: v.Type(), Offset: int64(d.readIndex())})
	//			break
	//		}
	//		v.SetUint(n)
	//
	//	case reflect.Float32, reflect.Float64:
	//		n, err := strconv.ParseFloat(s, v.Type().Bits())
	//		if err != nil || v.OverflowFloat(n) {
	//			d.saveError(&UnmarshalTypeError{Value: "number " + s, Type: v.Type(), Offset: int64(d.readIndex())})
	//			break
	//		}
	//		v.SetFloat(n)
	//	}
	//}
}

func fromSlice(l *MiddleValueList, p MiddleValue, destValue reflect.Value) {
	var data []MiddleValue
	if p.Length > 0 {
		data = l.List[p.First : p.First+p.Length]
	}
	destValue = makeValue(destValue)
	switch destValue.Kind() {
	case reflect.Interface:
		arr := arrayInterface(l, data)
		destValue.Set(reflect.ValueOf(arr))
	default:
		n := len(data)
		if destValue.Kind() == reflect.Slice {
			v := reflect.MakeSlice(destValue.Type(), n, n)
			destValue.Set(v)
		}
		i := 0
		for ; i < n; i++ {
			if i < destValue.Len() {
				fromMiddleValue(l, l.List[p.First+i], destValue.Index(i))
			}
		}
		if i < destValue.Len() {
			if destValue.Kind() == reflect.Array {
				z := reflect.Zero(destValue.Type().Elem())
				for ; i < destValue.Len(); i++ {
					destValue.Index(i).Set(z)
				}
			} else {
				destValue.SetLen(i)
			}
		}
	}
}

func fromMap(l *MiddleValueList, p MiddleValue, destValue reflect.Value) {
	destValue = makeValue(destValue)
	dstType := destValue.Type()
	switch destValue.Kind() {
	case reflect.Interface:
		var data []MiddleValue
		if p.Length > 0 {
			data = l.List[p.First : p.First+p.Length]
		}
		oi := objectInterface(l, data)
		destValue.Set(reflect.ValueOf(oi))
	case reflect.Map:
		if dstType.Key().Kind() != reflect.String {
			return
		}
		if destValue.IsNil() {
			destValue.Set(reflect.MakeMap(dstType))
		}
		fromMapToMap(l, p, destValue, dstType)
	case reflect.Struct:
		fromMapToStruct(l, p, destValue, dstType)
	}
}

func fromMapToMap(l *MiddleValueList, p MiddleValue, destValue reflect.Value, dstType reflect.Type) {
	elemType := dstType.Elem()
	for i := 0; i < p.Length; i++ {
		elemValue := reflect.New(elemType).Elem()
		fromMiddleValue(l, l.List[p.First+i], elemValue)
		keyValue := reflect.ValueOf(p.Name)
		destValue.SetMapIndex(keyValue, elemValue)
	}
}

func fromMapToStruct(l *MiddleValueList, p MiddleValue, destValue reflect.Value, dstType reflect.Type) {
	fields := cachedTypeFields(dstType)
	for i := 0; i < p.Length; i++ {
		e := l.List[p.First+i]
		f, ok := fields.byExactName[e.Name]
		if !ok {
			continue
		}
		subValue := destValue
		for _, j := range f.index {
			if subValue.Kind() == reflect.Ptr {
				if subValue.IsNil() {
					if !subValue.CanSet() {
						subValue = reflect.Value{}
						break
					}
					subValue.Set(reflect.New(subValue.Type().Elem()))
				}
				subValue = subValue.Elem()
			}
			subValue = subValue.Field(j)
		}
		fromMiddleValue(l, e, subValue)
	}
}

func makeValue(v reflect.Value) reflect.Value {
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && e.Elem().Kind() == reflect.Ptr {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		// Prevent infinite loop if v is an interface pointing to its own address:
		//     var v interface{}
		//     v = &v
		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		v = v.Elem()
	}
	return v
}

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) string {
	tag, _, _ = strings.Cut(tag, ",")
	return tag
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}

// A field represents a single field found in a struct.
type field struct {
	name    string
	tag     bool
	index   []int
	typ     reflect.Type
	encoder encoderFunc
}

type structFields struct {
	list        []field
	byExactName map[string]*field
}

// byIndex sorts field by index sequence.
type byIndex []field

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].index {
		if k >= len(x[j].index) {
			return false
		}
		if xik != x[j].index[k] {
			return xik < x[j].index[k]
		}
	}
	return len(x[i].index) < len(x[j].index)
}

var fieldCache sync.Map // map[reflect.Type]structFields

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) structFields {
	if f, ok := fieldCache.Load(t); ok {
		return f.(structFields)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.(structFields)
}

// typeFields returns a list of fields that JSON should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
func typeFields(t reflect.Type) structFields {
	// Anonymous fields to explore at the current level and the next.
	var current []field
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	var count, nextCount map[reflect.Type]int

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				if sf.Anonymous {
					t := sf.Type
					if t.Kind() == reflect.Pointer {
						t = t.Elem()
					}
					if !sf.IsExported() && t.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-struct types.
						continue
					}
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
				} else if !sf.IsExported() {
					// Ignore unexported non-embedded fields.
					continue
				}
				tag := sf.Tag.Get("json")
				if tag == "-" {
					continue
				}
				name := parseTag(tag)
				if !isValidTag(name) {
					name = ""
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Pointer {
					// Follow pointer.
					ft = ft.Elem()
				}

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					field := field{
						name:  name,
						tag:   tagged,
						index: index,
						typ:   ft,
					}

					fields = append(fields, field)
					if count[f.typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{name: ft.Name(), index: index, typ: ft})
				}
			}
		}
	}

	sort.Slice(fields, func(i, j int) bool {
		x := fields
		// sort field by name, breaking ties with depth, then
		// breaking ties with "name came from json tag", then
		// breaking ties with index sequence.
		if x[i].name != x[j].name {
			return x[i].name < x[j].name
		}
		if len(x[i].index) != len(x[j].index) {
			return len(x[i].index) < len(x[j].index)
		}
		if x[i].tag != x[j].tag {
			return x[i].tag
		}
		return byIndex(x).Less(i, j)
	})

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with JSON tags are promoted.

	// The fields are sorted in primary order of name, secondary order
	// of field index length. Loop over names; for each name, delete
	// hidden fields by choosing the one dominant field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		// One iteration per name.
		// Find the sequence of fields with the name of this first field.
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.name != name {
				break
			}
		}
		if advance == 1 { // Only one field with this name
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Sort(byIndex(fields))

	for i := range fields {
		f := &fields[i]
		f.encoder = typeEncoder(typeByIndex(t, f.index))
	}

	exactNameIndex := make(map[string]*field, len(fields))
	for i, field := range fields {
		exactNameIndex[field.name] = &fields[i]
	}
	return structFields{fields, exactNameIndex}
}

func typeByIndex(t reflect.Type, index []int) reflect.Type {
	for _, i := range index {
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		t = t.Field(i).Type
	}
	return t
}

// dominantField looks through the fields, all of which are known to
// have the same name, to find the single field that dominates the
// others using Go's embedding rules, modified by the presence of
// JSON tags. If there are multiple top-level fields, the boolean
// will be false: This condition is an error in Go and we skip all
// the fields.
func dominantField(fields []field) (field, bool) {
	// The fields are sorted in increasing index-length order, then by presence of tag.
	// That means that the first field is the dominant one. We need only check
	// for error cases: two fields at top level, either both tagged or neither tagged.
	if len(fields) > 1 && len(fields[0].index) == len(fields[1].index) && fields[0].tag == fields[1].tag {
		return field{}, false
	}
	return fields[0], true
}
