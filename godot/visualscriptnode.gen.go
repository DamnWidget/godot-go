package godot

import (
	"github.com/shadowapex/godot-go/gdnative"
)

/*------------------------------------------------------------------------------
//   This code was generated by a tool.
//
//   Changes to this file may cause incorrect behavior and will be lost if
//   the code is regenerated. Any updates should be done in
//   "class.go.tmpl" so they can be included in the generated
//   code.
//----------------------------------------------------------------------------*/

//func NewVisualScriptNodeFromPointer(ptr gdnative.Pointer) VisualScriptNode {
func newVisualScriptNodeFromPointer(ptr gdnative.Pointer) VisualScriptNode {
	owner := gdnative.NewObjectFromPointer(ptr)
	obj := VisualScriptNode{}
	obj.SetBaseObject(owner)

	return obj
}

/*
Undocumented
*/
type VisualScriptNode struct {
	Resource
	owner gdnative.Object
}

func (o *VisualScriptNode) BaseClass() string {
	return "VisualScriptNode"
}

/*
        Undocumented
	Args: [], Returns: Array
*/
func (o *VisualScriptNode) X_GetDefaultInputValues() gdnative.Array {
	//log.Println("Calling VisualScriptNode.X_GetDefaultInputValues()")

	// Build out the method's arguments
	ptrArguments := make([]gdnative.Pointer, 0, 0)

	// Get the method bind
	methodBind := gdnative.NewMethodBind("VisualScriptNode", "_get_default_input_values")

	// Call the parent method.
	// Array
	retPtr := gdnative.NewEmptyArray()
	gdnative.MethodBindPtrCall(methodBind, o.GetBaseObject(), ptrArguments, retPtr)

	// If we have a return type, convert it from a pointer into its actual object.
	ret := gdnative.NewArrayFromPointer(retPtr)
	return ret
}

/*
        Undocumented
	Args: [{ false values Array}], Returns: void
*/
func (o *VisualScriptNode) X_SetDefaultInputValues(values gdnative.Array) {
	//log.Println("Calling VisualScriptNode.X_SetDefaultInputValues()")

	// Build out the method's arguments
	ptrArguments := make([]gdnative.Pointer, 1, 1)
	ptrArguments[0] = gdnative.NewPointerFromArray(values)

	// Get the method bind
	methodBind := gdnative.NewMethodBind("VisualScriptNode", "_set_default_input_values")

	// Call the parent method.
	// void
	retPtr := gdnative.NewEmptyVoid()
	gdnative.MethodBindPtrCall(methodBind, o.GetBaseObject(), ptrArguments, retPtr)

}

/*
        Undocumented
	Args: [{ false port_idx int}], Returns: Variant
*/
func (o *VisualScriptNode) GetDefaultInputValue(portIdx gdnative.Int) gdnative.Variant {
	//log.Println("Calling VisualScriptNode.GetDefaultInputValue()")

	// Build out the method's arguments
	ptrArguments := make([]gdnative.Pointer, 1, 1)
	ptrArguments[0] = gdnative.NewPointerFromInt(portIdx)

	// Get the method bind
	methodBind := gdnative.NewMethodBind("VisualScriptNode", "get_default_input_value")

	// Call the parent method.
	// Variant
	retPtr := gdnative.NewEmptyVariant()
	gdnative.MethodBindPtrCall(methodBind, o.GetBaseObject(), ptrArguments, retPtr)

	// If we have a return type, convert it from a pointer into its actual object.
	ret := gdnative.NewVariantFromPointer(retPtr)
	return ret
}

/*
        Undocumented
	Args: [], Returns: VisualScript
*/
func (o *VisualScriptNode) GetVisualScript() VisualScriptImplementer {
	//log.Println("Calling VisualScriptNode.GetVisualScript()")

	// Build out the method's arguments
	ptrArguments := make([]gdnative.Pointer, 0, 0)

	// Get the method bind
	methodBind := gdnative.NewMethodBind("VisualScriptNode", "get_visual_script")

	// Call the parent method.
	// VisualScript
	retPtr := gdnative.NewEmptyObject()
	gdnative.MethodBindPtrCall(methodBind, o.GetBaseObject(), ptrArguments, retPtr)

	// If we have a return type, convert it from a pointer into its actual object.
	ret := newVisualScriptFromPointer(retPtr)

	// Check to see if we already have an instance of this object in our Go instance registry.
	if instance, ok := InstanceRegistry.Get(ret.GetBaseObject().ID()); ok {
		return instance.(VisualScriptImplementer)
	}

	// Check to see what kind of class this is and create it. This is generally used with
	// GetNode().
	className := ret.GetClass()
	if className != "VisualScript" {
		actualRet := getActualClass(className, ret.GetBaseObject())
		return actualRet.(VisualScriptImplementer)
	}

	return &ret
}

/*
        Undocumented
	Args: [], Returns: void
*/
func (o *VisualScriptNode) PortsChangedNotify() {
	//log.Println("Calling VisualScriptNode.PortsChangedNotify()")

	// Build out the method's arguments
	ptrArguments := make([]gdnative.Pointer, 0, 0)

	// Get the method bind
	methodBind := gdnative.NewMethodBind("VisualScriptNode", "ports_changed_notify")

	// Call the parent method.
	// void
	retPtr := gdnative.NewEmptyVoid()
	gdnative.MethodBindPtrCall(methodBind, o.GetBaseObject(), ptrArguments, retPtr)

}

/*
        Undocumented
	Args: [{ false port_idx int} { false value Variant}], Returns: void
*/
func (o *VisualScriptNode) SetDefaultInputValue(portIdx gdnative.Int, value gdnative.Variant) {
	//log.Println("Calling VisualScriptNode.SetDefaultInputValue()")

	// Build out the method's arguments
	ptrArguments := make([]gdnative.Pointer, 2, 2)
	ptrArguments[0] = gdnative.NewPointerFromInt(portIdx)
	ptrArguments[1] = gdnative.NewPointerFromVariant(value)

	// Get the method bind
	methodBind := gdnative.NewMethodBind("VisualScriptNode", "set_default_input_value")

	// Call the parent method.
	// void
	retPtr := gdnative.NewEmptyVoid()
	gdnative.MethodBindPtrCall(methodBind, o.GetBaseObject(), ptrArguments, retPtr)

}

// VisualScriptNodeImplementer is an interface that implements the methods
// of the VisualScriptNode class.
type VisualScriptNodeImplementer interface {
	ResourceImplementer
	X_GetDefaultInputValues() gdnative.Array
	X_SetDefaultInputValues(values gdnative.Array)
	GetDefaultInputValue(portIdx gdnative.Int) gdnative.Variant
	GetVisualScript() VisualScriptImplementer
	PortsChangedNotify()
	SetDefaultInputValue(portIdx gdnative.Int, value gdnative.Variant)
}
