package filter

import (
	"fmt"
	"github.com/mojo-lang/mojo/go/pkg/mojo/lang"
)

type Field struct {
	Name string
	Path []string
}

func (f *Field) Parse(expr *lang.Expression) error {
	if expr == nil {
		return fmt.Errorf("missing filter field")
	}
	switch expr.Expression.(type) {
	case *lang.Expression_StringLiteralExpr:
		f.Name = expr.GetStringLiteralExpr().GetValue()
	case *lang.Expression_IdentifierExpr:
		f.Name = expr.GetIdentifierExpr().GetName()
	case *lang.Expression_ExplicitMemberExpr:
		f.Name = expr.GetExplicitMemberExpr().Member
	default:
		return fmt.Errorf("unsupported filter field %T", expr.Expression)
	}
	if f.Name == "" {
		return fmt.Errorf("empty filter field")
	}
	return nil
}

type Value struct {
	Val interface{}
}

func (v *Value) Parse(expr *lang.Expression) error {
	if expr == nil {
		return fmt.Errorf("missing filter value")
	}
	switch expr.Expression.(type) {
	case *lang.Expression_IdentifierExpr:
		v.Val = expr.GetIdentifierExpr().GetName()
	case *lang.Expression_StringLiteralExpr:
		v.Val = expr.GetStringLiteralExpr().EvalValue()
	case *lang.Expression_FloatLiteralExpr:
		v.Val = expr.GetFloatLiteralExpr().EvalValue()
	case *lang.Expression_IntegerLiteralExpr:
		v.Val = expr.GetIntegerLiteralExpr().EvalValue()
	case *lang.Expression_BoolLiteralExpr:
		v.Val = expr.GetBoolLiteralExpr().EvalValue()
	default:
		return fmt.Errorf("unsupported filter value %T", expr.Expression)
	}
	return nil
}
