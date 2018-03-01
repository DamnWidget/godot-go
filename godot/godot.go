package godot

import (
	"github.com/pinzolo/casee"
	"github.com/shadowapex/godot-go/gdnative"
	"log"
	"reflect"
	"runtime"
	"strings"
)

var debug = false

// EnableDebug will enable debug logging of the godot library.
func EnableDebug() {
	debug = true
}

// Init is a special Go function that will be called upon library initialization.
func init() {
	// Configure GDNative to use our own NativeScript init function.
	gdnative.SetNativeScriptInit(autoRegisterClasses)
}

// autoRegisterClasses is the script's entrypoint. It is called by Godot
// when a script is loaded. It is responsible for registering all the classes.
func autoRegisterClasses() {
	log.Println("Discovering classes to register with Godot...")

	// Loop through our registered classes and register them with the Godot API.
	for _, constructor := range godotConstructorsToAutoRegister {
		// Use the constructor to build a class to inspect the given structure.
		class := constructor()

		// Get the type of the given struct, and get its name as a string
		classType := reflect.TypeOf(class)
		classString := strings.Replace(classType.String(), "*", "", 1)
		if debug {
			log.Println("Registering class:", classString)
		}

		// Create a registered class structure that will hold information about the
		// cass and its methods.
		regClass := newRegisteredClass(classType)

		// Call the "BaseClass" method on the class to get the base class.
		baseClass := class.(Class).BaseClass()
		if debug {
			log.Println("  Using Base Class:", baseClass)
		}

		// Set up our constructor and destructor function structs.
		createFunc := createConstructor(classString, constructor)
		destroyFunc := createDestructor(classString)

		// Register our class with Godot.
		gdnative.NativeScript.RegisterClass(classString, baseClass, createFunc, destroyFunc)

		// Loop through our class's struct fields. We do this to register properties as well
		// as find the embedded parent struct to ensure we don't register those methods.
		if debug {
			log.Println("  Looking at struct fields:")
			log.Println("    Found", classType.Elem().NumField(), "struct fields.")
		}
		for i := 0; i < classType.Elem().NumField(); i++ {
			classField := classType.Elem().Field(i)
			if debug {
				log.Println("  Found field:", classField.Name)
				log.Println("    Type:", classField.Type.String())
				log.Println("    Anonymous:", classField.Anonymous)
				log.Println("    Package Path:", classField.PkgPath)
			}

			// Look only at anonymously embedded fields
			if !classField.Anonymous {
				continue
			}
		}

		// Loop through our class's methods that are attached to it.
		if debug {
			log.Println("  Looking at methods:")
			log.Println("    Found", classType.NumMethod(), "methods")
		}
		for i := 0; i < classType.NumMethod(); i++ {
			classMethod := classType.Method(i)

			// TODO: For now we are only checking if the given method is embedded or
			// not. If the method comes from an embedded structure, skip it.
			// We need to figure this shit out, so we can allow embedding of non-godot
			// types.
			runMethod := runtime.FuncForPC(classMethod.Func.Pointer())
			filename, _ := runMethod.FileLine(classMethod.Func.Pointer())
			if strings.Contains(filename, "<autogenerated>") {
				continue
			}

			if debug {
				log.Println("  Found method:", classMethod.Name)
				log.Println("    Method in package path:", classMethod.PkgPath)
				log.Println("    Type package path:", classMethod.Type.PkgPath())
				log.Println("    Type:", classMethod.Type.String())
				log.Println("    Kind:", classMethod.Type.Kind().String())
			}

			// Construct a registered method structure that inspects all of the
			// arguments and return types.
			regMethod := newRegisteredMethod(classMethod)
			regClass.addMethod(classMethod.Name, regMethod)
			if debug {
				log.Println("    Method Arguments:", len(regMethod.arguments))
				log.Println("    Method Arguments:", regMethod.arguments)
				log.Println("    Method Returns:", len(regMethod.returns))
				log.Println("    Method Returns:", regMethod.returns)
			}

			// Skip the method if its a Class interface method
			skip := false
			for _, exclude := range []string{"BaseClass", "GetOwner", "SetOwner"} {
				if classMethod.Name == exclude {
					skip = true
				}
			}
			if skip {
				continue
			}

			// Look at the method name to see if it starts with "X_". If it does, we need to
			// replace it with an underscore. This is required because Go method visibility
			// is done through case sensitivity. Since Godot private methods start with an
			// "_" character.
			goMethodName := classMethod.Name
			godotMethodName := toGodotMethodName(goMethodName)

			// Set up our method structure
			method := createMethod(classString, goMethodName)
			attributes := &gdnative.MethodAttributes{
				RPCType: gdnative.MethodRpcModeDisabled,
			}

			// Register the method.
			gdnative.NativeScript.RegisterMethod(classString, godotMethodName, attributes, method)
		}

		// Register our class in our Go registry.
		classRegistry[classString] = regClass

	}
}

// CreateConstructor will create the InstanceCreateFunc structure with the given class name
// and constructor. This structure can be used when registering a class with Godot.
func createConstructor(classString string, constructor ClassConstructor) *gdnative.InstanceCreateFunc {
	var createFunc gdnative.InstanceCreateFunc
	createFunc.CreateFunc = func(object gdnative.Object, methodData string) string {
		// Create a new instance of the object.
		class := constructor()
		if debug {
			log.Println("Created new object instance:", class, "with instance address:", object.ID())
		}

		// Add the Godot object pointer to the class structure.
		class.SetOwner(object)

		// Add the instance to our instance registry.
		instanceRegistry[object.ID()] = class

		// Return the instance string. This will be passed to the method function as userData, so we
		// can look up the instance in our registry.
		return object.ID()

	}
	createFunc.MethodData = classString
	createFunc.FreeFunc = func(methodData string) {}

	return &createFunc
}

// CreateDestructor will create the InstanceDestroyFunc structure with the given class name.
// This structure can be used when registering a class with Godot.
func createDestructor(classString string) *gdnative.InstanceDestroyFunc {
	var destroyFunc gdnative.InstanceDestroyFunc
	destroyFunc.DestroyFunc = func(object gdnative.Object, className, instanceID string) {
		if debug {
			log.Println("Destroying object instance:", className, "with instance address:", object.ID())
		}

		// Unregister it from our instanceRegistry so it can be garbage collected.
		delete(instanceRegistry, instanceID)
	}
	destroyFunc.MethodData = classString
	destroyFunc.FreeFunc = func(methodData string) {}

	return &destroyFunc
}

// CreateMethod will create the InstanceMethod structure. This will be called whenever
// a Godot method is called.
func createMethod(classString, methodString string) *gdnative.InstanceMethod {
	var methodFunc gdnative.InstanceMethod
	methodFunc.Method = func(object gdnative.Object, classMethod, instanceString string, numArgs int, args []gdnative.Variant) gdnative.Variant {
		var ret gdnative.Variant

		// Get the object instance based on the instance string given in userData.
		class := instanceRegistry[instanceString]
		classValue := reflect.ValueOf(class)

		if debug {
			log.Println("Method was called!")
			log.Println("  godotObject:", object)
			log.Println("  numArgs:", numArgs)
			log.Println("  args:", args)
			log.Println("  instance:", class)
			log.Println("  methodString (methodData):", classMethod)
			log.Println("  instanceString (userData):", instanceString)
		}

		// Create a slice of Godot Variant arguments
		goArgsSlice := []reflect.Value{}

		// If we have arguments, append the first argument.
		for _, arg := range args {
			// Convert the variant into its base type
			goArgsSlice = append(goArgsSlice, variantToGoType(arg))
		}

		// Use the method string to get the class name and method name.
		if debug {
			log.Println("  Getting class name and method name...")
		}
		classMethodSlice := strings.Split(classMethod, "::")
		className := classMethodSlice[0]
		methodName := classMethodSlice[1]
		if debug {
			log.Println("    Class Name: ", className)
			log.Println("    Method Name: ", methodName)
		}

		// Look up the registered class so we can find out how many arguments it takes
		// and their types.
		if debug {
			log.Println("  Look up the registered class and its method")
			log.Println("  Registered classes:", classRegistry)
		}
		regClass := classRegistry[className]
		if regClass == nil {
			log.Fatal("  This class has not been registered! Class name: ", className, " Method name: ", methodName)
		}
		if debug {
			log.Println("  Looked up class:", regClass)
			log.Println("  Methods in class:", regClass.methods)
		}
		regMethod := regClass.methods[methodName]

		if debug {
			log.Println("  Registered method arguments:", regMethod.arguments)
			log.Println("  Arguments to pass:", goArgsSlice)
		}

		// Check to ensure the method has the same number of arguments we expect
		if len(regMethod.arguments)-1 != int(numArgs) {
			gdnative.Log.Error("Invalid number of arguments. Expected ", numArgs, " arguments. (Got ", len(regMethod.arguments), ")")
			panic("Invalid number of arguments.")
		}

		// Get the value of the class, so we can call methods on it.
		method := classValue.MethodByName(methodName)
		rawRet := method.Call(goArgsSlice)
		if debug {
			log.Println("Got raw return value after method call:", rawRet)
		}

		// Check to see if this returns anything.
		if len(rawRet) == 0 {
			nilReturn := gdnative.NewVariantNil()
			return *nilReturn
		}

		// Convert our returned value into a Godot Variant.
		rawRetInterface := rawRet[0].Interface()
		switch regMethod.returns[0].String() {
		case "bool":
			base := gdnative.Bool(rawRetInterface.(bool))
			variant := gdnative.NewVariantBool(base)
			ret = *variant
		case "int64":
			base := gdnative.Int64T(rawRetInterface.(int64))
			variant := gdnative.NewVariantInt(base)
			ret = *variant
		case "int32":
			base := gdnative.Int64T(rawRetInterface.(int32))
			variant := gdnative.NewVariantInt(base)
			ret = *variant
		case "int":
			base := gdnative.Int64T(rawRetInterface.(int))
			variant := gdnative.NewVariantInt(base)
			ret = *variant
		case "uint64":
			base := gdnative.Uint64T(rawRetInterface.(uint64))
			variant := gdnative.NewVariantUint(base)
			ret = *variant
		case "uint32":
			base := gdnative.Uint64T(rawRetInterface.(uint32))
			variant := gdnative.NewVariantUint(base)
			ret = *variant
		case "uint":
			base := gdnative.Uint64T(rawRetInterface.(uint))
			variant := gdnative.NewVariantUint(base)
			ret = *variant
		case "float64":
			base := gdnative.Double(rawRetInterface.(float64))
			variant := gdnative.NewVariantReal(base)
			ret = *variant
		case "string":
			base := gdnative.WcharT(rawRetInterface.(string))
			baseStr := base.AsString()
			variant := gdnative.NewVariantString(*baseStr)
			ret = *variant
		default:
			panic("The return was not valid. Should be Godot Variant or built-in Go type. Received: " + regMethod.returns[0].String())
		}

		return ret
	}
	methodFunc.MethodData = classString + "::" + methodString
	methodFunc.FreeFunc = func(methodData string) {}

	return &methodFunc
}

// VariantToGoType will check the given variant type and convert it to its
// actual type. The value is returned as a reflect.Value.
func variantToGoType(variant gdnative.Variant) reflect.Value {
	switch variant.GetType() {
	case gdnative.VariantTypeBool:
		return reflect.ValueOf(variant.AsBool())
	case gdnative.VariantTypeInt:
		return reflect.ValueOf(variant.AsInt())
	case gdnative.VariantTypeReal:
		return reflect.ValueOf(variant.AsReal())
	case gdnative.VariantTypeString:
		return reflect.ValueOf(variant.AsString())
	case gdnative.VariantTypeVector2:
		return reflect.ValueOf(variant.AsVector2())
	case gdnative.VariantTypeRect2:
		return reflect.ValueOf(variant.AsRect2())
	case gdnative.VariantTypeVector3:
		return reflect.ValueOf(variant.AsVector3())
	case gdnative.VariantTypeTransform2D:
		return reflect.ValueOf(variant.AsTransform2D())
	case gdnative.VariantTypePlane:
		return reflect.ValueOf(variant.AsPlane())
	case gdnative.VariantTypeQuat:
		return reflect.ValueOf(variant.AsQuat())
	case gdnative.VariantTypeAabb:
		return reflect.ValueOf(variant.AsAabb())
	case gdnative.VariantTypeBasis:
		return reflect.ValueOf(variant.AsBasis())
	case gdnative.VariantTypeTransform:
		return reflect.ValueOf(variant.AsTransform())
	case gdnative.VariantTypeColor:
		return reflect.ValueOf(variant.AsColor())
	case gdnative.VariantTypeNodePath:
		return reflect.ValueOf(variant.AsNodePath())
	case gdnative.VariantTypeRid:
		return reflect.ValueOf(variant.AsRid())
	case gdnative.VariantTypeObject:
		return reflect.ValueOf(variant.AsObject())
	case gdnative.VariantTypeDictionary:
		return reflect.ValueOf(variant.AsDictionary())
	case gdnative.VariantTypeArray:
		return reflect.ValueOf(variant.AsArray())
	case gdnative.VariantTypePoolByteArray:
		return reflect.ValueOf(variant.AsPoolByteArray())
	case gdnative.VariantTypePoolIntArray:
		return reflect.ValueOf(variant.AsPoolIntArray())
	case gdnative.VariantTypePoolRealArray:
		return reflect.ValueOf(variant.AsPoolRealArray())
	case gdnative.VariantTypePoolStringArray:
		return reflect.ValueOf(variant.AsPoolStringArray())
	case gdnative.VariantTypePoolVector2Array:
		return reflect.ValueOf(variant.AsPoolVector2Array())
	case gdnative.VariantTypePoolVector3Array:
		return reflect.ValueOf(variant.AsPoolVector3Array())
	case gdnative.VariantTypePoolColorArray:
		return reflect.ValueOf(variant.AsPoolColorArray())
	default:
		panic("Unknown type of variant argument.")
	}
	return reflect.Value{}
}

// toGoMethodName will take the given Godot method name in snake_case and convert it
// to a CamelCase method name.
func toGoMethodName(methodName string) string {
	goMethodName := casee.ToPascalCase(methodName)
	if strings.HasPrefix(methodName, "_") {
		goMethodName = "X_" + goMethodName
	}
	if debug {
		log.Println("    Godot method name:", methodName)
		log.Println("    Go mapped method name:", goMethodName)
	}

	return goMethodName
}

// toGodotMethodName will take the given Go method name in CamelCase and convert it
// to a snake_case method name.
func toGodotMethodName(goMethodName string) string {
	methodName := goMethodName
	privatePrefix := ""
	if strings.HasPrefix(goMethodName, "X_") {
		methodName = strings.Replace(methodName, "X_", "_", 1)
		privatePrefix = "_"
	}
	methodName = casee.ToSnakeCase(methodName)
	methodName = privatePrefix + methodName
	if debug {
		log.Println("    Go method name:", goMethodName)
		log.Println("    Godot mapped method name:", methodName)
	}

	return methodName
}
