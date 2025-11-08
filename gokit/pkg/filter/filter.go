package filter

import (
	"github.com/mojo-lang/lang/go/pkg/mojo/lang"
	"github.com/mojo-lang/mojo/go/pkg/mojo/parser/syntax"
	"gorm.io/gorm/clause"
)

type Filter string

func (f Filter) CompileToSql() (clause.Expression, error) {
	parser := syntax.New(nil)
	if expr, err := parser.ParseExpression(string(f)); err != nil {
		return nil, err
	} else {
		return f.compile(expr)
	}
}

func (f Filter) compile(expr *lang.Expression) (clause.Expression, error) {
	switch expr.Expression.(type) {
	case *lang.Expression_BinaryExpr:
		binary := expr.GetBinaryExpr()
		switch binary.GetOperator().GetSymbol() {
		case "&&", "and":
			if left, err := f.compile(binary.LeftArgument); err != nil {
				return nil, err
			} else {
				if right, err := f.compile(binary.RightArgument); err != nil {
					return nil, err
				} else {
					return clause.And(left, right), nil
				}
			}
		case "||", "or":
			if left, err := f.compile(binary.LeftArgument); err != nil {
				return nil, err
			} else {
				if right, err := f.compile(binary.RightArgument); err != nil {
					return nil, err
				} else {
					return clause.Or(left, right), nil
				}
			}
		case "==", ">=", "<=", ">", "<":
			filed := &Field{}
			if err := filed.Parse(binary.LeftArgument); err != nil {
				return nil, err
			}
			value := &Value{}
			if err := value.Parse(binary.RightArgument); err != nil {
				return nil, err
			}
			switch binary.GetOperator().GetSymbol() {
			case "==":
				return clause.Eq{Column: filed.Name, Value: value.Val}, nil
			case ">=":
				return clause.Gte{Column: filed.Name, Value: value.Val}, nil
			case "<=":
				return clause.Lte{Column: filed.Name, Value: value.Val}, nil
			case ">":
				return clause.Gt{Column: filed.Name, Value: value.Val}, nil
			case "<":
				return clause.Lt{Column: filed.Name, Value: value.Val}, nil
			}
		}
	case *lang.Expression_PrefixUnaryExpr:
		unary := expr.GetPrefixUnaryExpr()
		switch unary.GetOperator().GetSymbol() {
		case "!", "not":
			if expression, err := f.compile(unary.Argument); err != nil {
				return nil, err
			} else {
				return clause.Not(expression), nil
			}
		}
	case *lang.Expression_FunctionCallExpr:
	case *lang.Expression_ParenthesizedExpr:
		return f.compile(expr.GetParenthesizedExpr().Expression)
	}
	return nil, nil
}
