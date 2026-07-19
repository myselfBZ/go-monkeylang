package evaluator

import (
	"fmt"

	"github.com/myselfBZ/interpreter/internal/ast"
	"github.com/myselfBZ/interpreter/internal/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func boolToBoolOBJ(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
    return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(o object.Object) bool{
    if o != nil{
        return o.Type() == object.ERROR_OBJ
    }
    return false
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
    case *ast.BlockStatement:
        return evalBlock(node, env)
	case *ast.IntLiteral:
		return &object.Integer{Value: int(node.Value)}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.PrefixExpression:
		return evalPrefix(node, node.Operator, env)
	case *ast.InfixExperssion:
		right := Eval(node.Right, env)
        if isError(right){
            return right
        }
		left := Eval(node.Left, env)
		return evalInfix(right, left, node.Operator)
    case *ast.IfExpression:
		newEnv := object.NewEnclosedEnvironment(env)
        return evalIfExp(node, newEnv)
    case *ast.ReturnStatement:
        value := Eval(node.ReturnValue, env)
        if isError(value){
            return value
        }
        return &object.ReturnValue{Value: value}
    case *ast.LetStatement:
        v := Eval(node.Value, env)
        if isError(v){
            return v
        }
        env.Set(node.Name.Value, v)
    case *ast.Identifier:
        v, ok := env.Get(node.Value) 
        if !ok{
            return newError("identifier not found %s", node.Value)
        }
        return v
    // case *ast.FunctionLiteral:
    //     funcEnv := object.NewEnvironment()
    //
	default:
		return NULL
	}
    return nil
}


func evalIdent(node *ast.Identifier, env *object.Environment) object.Object{
    obj, ok := env.Get(node.Value)
    if !ok{
        return newError("identifier not found %s", node.Value)
    }
    return obj
}


func evalProgram(node *ast.Program, env *object.Environment) object.Object {
    var result object.Object
    for _, v := range node.Statements{
        result = Eval(v, env)
        if err, ok := result.(*object.Error); ok{
            return err
        }
        if returnV, ok := result.(*object.ReturnValue); ok{
            return returnV.Value
        }
    }
    return result
}


func evalBlock(node *ast.BlockStatement, env *object.Environment) object.Object {
    var result object.Object
    for _, v := range node.Statements{
        result = Eval(v, env)
        if result != nil && result.Type() == object.ERROR_OBJ{
            return result
        }
        if result != nil && result.Type() == object.RETURN_VALUE{
            return result
        }
    }
    return result
}

func evalPrefix(node *ast.PrefixExpression, op string, env *object.Environment) object.Object {
	v := Eval(node.Right, env)
    if isError(v){
        return v
    }
	switch op {
	case "!":
		return evalBang(v)
	case "-":
		return evalMinus(v)
	default:
		return newError("can't have %s infront of %s", op,v.Type()) 
	}
}

func evalBang(o object.Object) object.Object {
	switch o {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinus(o object.Object) object.Object {
	if o.Type() != object.INTEGER_OBJ {
		return NULL
	}
	val := o.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func evalIntInfix(right object.Object, left object.Object, oprtr string) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch oprtr {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "==":
		return boolToBoolOBJ(leftValue == rightValue)
	case "!=":
		return boolToBoolOBJ(leftValue != rightValue)
	case ">=":
		return boolToBoolOBJ(leftValue >= rightValue)
	case "<=":
		return boolToBoolOBJ(leftValue <= rightValue)
	case ">":
		return boolToBoolOBJ(leftValue > rightValue)
	case "<":
		return boolToBoolOBJ(leftValue < rightValue)
	default:
        return newError("unknown operator: %s%s%s", left.Inspect(), oprtr, right.Inspect())
	}
}

func evalInfix(right object.Object, left object.Object, oprtr string) object.Object {
	if right.Type() == object.INTEGER_OBJ && left.Type() == object.INTEGER_OBJ {
		return evalIntInfix(right, left, oprtr)
	}
	return evalBoolInfix(right, left, oprtr)
}

func compareBool(right object.Object, left object.Object, oprtr string) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value
	switch oprtr {
	case "==":
		return &object.Boolean{Value: leftValue == rightValue}
	case "!=":
		return &object.Boolean{Value: leftValue != rightValue}
	default:
		return newError("unknown operator between booleans %s", oprtr)
	}
}

func evalBoolInfix(right object.Object, left object.Object, oprtr string) object.Object {
	if right.Type() != left.Type() {
        return newError("unknown operation with umatched types") 
	}
	if right.Type() == object.BOOLEAN_OBJ {
		return compareBool(right, left, oprtr)
	}
    return newError("unknown operator for booleans %s", oprtr) 
}

func evalIfExp(node *ast.IfExpression, env *object.Environment) object.Object {
	conditionObj := Eval(node.Condition, env)
	condition, ok := conditionObj.(*object.Boolean)
    if !ok{
        return newError("non-boolean condition in if statement %s", conditionObj.Type())
    }
	if condition.Value {
		return Eval(node.Consequence, env)
	}
    if node.Alternative != nil{
        return Eval(node.Alternative, env)
    }
    return NULL
}



