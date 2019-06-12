package main

import (
    "fmt"
    "reflect"
    "strings"
    ndm "github.com/openebs/maya/pkg/apis/openebs.io/ndm/v1alpha1"
)

/*
 * Just do 'go run genBuilder.go'
 */

/*
 * This is used to avoid infinite loop due to self-referencing structures
 */
var parsedStruct map[string]int

/*
 * Generates predicates based on
 * - Possible value (val)
 * - Struct Field (f) that provides fieldType, fieldName
 * - Path to reach this field (fParentPath)
 * - Main struct for which builder pattern is generated
 */
func generatePredicate(val string, f reflect.StructField, fPPath string, typeName string) {
	var changeBack int
	fType := f.Type.Kind().String()
	fName := f.Name
	fPath := fPPath + "." + f.Name
	predicateString := string(
`
// Is[[Default]] filters the [[StructName]] based on type of the [[FieldName]]
func Is[[Default]]([[SliceIndex]]) Predicate {
        return func(obj *[[StructName]]) bool {
                return obj.Is[[Default]]([[SliceIndex]])
        }
}

// Is[[Default]] returns true if the [[StructName]].[[FieldName]] is of [[Default]]
func (obj *[[StructName]]) Is[[Default]]([[SliceIndex]]) bool {
        return obj.[[FieldPath]] == [[DefaultComparision]]
}
`)
	predicateString = strings.Replace(predicateString, "[[StructName]]", typeName, -1)
	predicateString = strings.Replace(predicateString, "[[Default]]", val, -1)

	/*
	 * Add index 'i' as parameter to predicate in case of slice
	 */
	if fType == "slice" {
		predicateString = strings.Replace(predicateString, "[[SliceIndex]]", "int i", -1)
	}

	/*
	 * Replace value to be compared with based on its type 
	 * mainly, need to consider <double quotes> for string type
	 * TODO: for ptrs
	 */
	changeBack = 0
	if fType == "slice" {
		fType = f.Type.Elem().Kind().String()
		changeBack = 1
	}
	if fType == "string" {
		predicateString = strings.Replace(predicateString, "[[DefaultComparision]]", "\"" + val + "\"", -1)
	} else {
		predicateString = strings.Replace(predicateString, "[[DefaultComparision]]", val, -1)
	}
	if changeBack == 1 {
		fType = "slice"
		changeBack = 0
	}

	/*
	 * Replace FieldPath with what received
	 * Mainly need to consider slice to add indexing
	 * TODO: for ptrs
	 */
	if fType == "slice" {
		predicateString = strings.Replace(predicateString, "[[FieldPath]]", fPath + "[i]", -1)
	} else {
		predicateString = strings.Replace(predicateString, "[[FieldPath]]", fPath, -1)
	}

	fmt.Println(strings.Replace(predicateString, "[[FieldName]]", fName, -1))
}

/*
 * Generate Get op for tagged fields
 * TODO: Need to fix for slice similar to predicates
 */
func generateGet(fType string, fName string, fPath string, typeName string) {
	getString := string(
`
// Get[[FieldName]] returns [[FieldName]] of the [[StructName]]
func (obj *[[StructName]]) Get[[FieldName]]() [[FieldType]]{
        return obj.[[FieldPath]]
}
`)
	getString = strings.Replace(getString, "[[StructName]]", typeName, -1)
	getString = strings.Replace(getString, "[[FieldType]]", fType, -1)
	getString = strings.Replace(getString, "[[FieldPath]]", fPath, -1)
	fmt.Println(strings.Replace(getString, "[[FieldName]]", fName, -1))
}

/*
 * parseStruct goes through the given structure using Reflect
 * Calls generatePredicate and generateGet based on tags on the struct fields
 */
func parseStruct(t reflect.Type, indent int, fPath string, typeName string) {
	indent = indent + 1
//	fmt.Printf("%*sK: %v N: %v fpath: %v\n", indent, " ", t.Kind(), t.Name(), fPath)
	switch(t.Kind()) {
		case reflect.Struct:
			_, ok := parsedStruct[t.Name()]
			if !ok {
				parsedStruct[t.Name()] = 1
			} else {
				break
			}
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				fType := f.Type
				fTag := f.Tag
				possibleValues := fTag.Get("possibleValues")
				ops := fTag.Get("genOps")
//				fmt.Printf("%*sfN: %v fType: %v fPath: %s fDefaults: %v fDefaultOps: %v\n",
//				    indent, " ", f.Name, fType, fPath + "." + f.Name,  possibleValues, ops)
				if possibleValues != "" {
					for _, d := range strings.Split(possibleValues, ",") {
//						fmt.Printf("%*spossibleVal: %v\n", indent + 1, " ", d)
						generatePredicate(d, f, fPath, typeName)
						//generatePredicate(d, fType.Kind().String(), f.Name, fPath + "." + f.Name, typeName)
					}
				}
				if ops != "" {
					for _, d := range strings.Split(ops, ",") {
//						fmt.Printf("%*sopsVal: %v\n", indent + 1, " ", d)
						if d == "get" {
							generateGet(fType.Kind().String(), f.Name, fPath + "." + f.Name, typeName)
						}
					}
				}

				parseStruct(fType, indent, fPath + "." + f.Name, typeName)
			}
			break
		case reflect.Slice:
	// For Elem(), type should be from pointer like reflect.TypeOf(&t).Elem()
	// Or from slice as below
//			eType := t.Elem()
//			parseStruct(eType, indent, fPath + "[i]", typeName)
			break;
	}
}

func main() {
	builderStruct := string(
`// [[StructName]] encapsulates [[StructName]] api object.                                                   
type [[StructName]] struct {                                                                             
        // actual [[StructName]] object                                                                 
        Object *ndm.[[StructName]]
}                                                                                                     
                                                                                                      
// [[StructName]]List holds the list of [[StructName]] api                                                  
type [[StructName]]List struct {                                                                         
        // list of [[StructName]]s
        ObjectList *ndm.[[StructName]]List                                                               
}                                                                                                     
                                                                                                      
// Predicate defines an abstraction to determine conditional checks against the                       
// provided [[StructName]] instance                                                                     
type Predicate func(*[[StructName]]) bool                                                                
                                                                                                      
// predicateList holds the list of Predicates                                                         
type predicateList []Predicate                                                           

// all returns true if all the predicates succeed against the provided block
// device instance.
func (l predicateList) all(c *[[StructName]]) bool {
        for _, pred := range l {
                if !pred(c) {
                        return false
                }
        }
        return true
}

`)
	var t ndm.BlockDevice
	parsedStruct = map[string]int{}
	rType := reflect.TypeOf(t)

	fmt.Println(strings.Replace(builderStruct, "[[StructName]]", rType.Name(), -1))

//	fmt.Printf("K: %v N: %v\n", rType.Kind(), rType.Name())
	parseStruct(rType, 0, "Object", rType.Name())

	var t1 ndm.DeviceDevLink
	rType = reflect.TypeOf(t1)
	parsedStruct = map[string]int{}
//	fmt.Printf("K: %v N: %v\n", rType.Kind(), rType.Name())
	parseStruct(rType, 0, "Object", rType.Name())
}

