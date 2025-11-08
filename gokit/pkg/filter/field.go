package filter

import "github.com/mojo-lang/lang/go/pkg/mojo/lang"

type Field struct {
    Name string
    Path []string
}

func (f *Field) Parse(expr *lang.Expression) error {
    switch expr.Expression.(type) {
    case *lang.Expression_StringLiteralExpr:
        f.Name = expr.GetStringLiteralExpr().GetValue()
    case *lang.Expression_IdentifierExpr:
        f.Name = expr.GetIdentifierExpr().GetName()
    case *lang.Expression_ExplicitMemberExpr:
        f.Name = expr.GetExplicitMemberExpr().Member
    }
    return nil
}

type Value struct {
    Val interface{}
}

func (v *Value) Parse(expr *lang.Expression) error {
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
    }
    return nil
}
