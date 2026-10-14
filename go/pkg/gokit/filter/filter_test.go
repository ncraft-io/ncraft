package filter

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm/clause"
)

func TestFilterOperators(t *testing.T) {
	for _, tc := range []struct {
		input      string
		expression clause.Expression
	}{
		{`name == "alice"`, clause.Eq{}},
		{`count != 3`, clause.Neq{}},
		{`count >= 3`, clause.Gte{}},
		{`count <= 3`, clause.Lte{}},
		{`count > 3`, clause.Gt{}},
		{`count < 3`, clause.Lt{}},
		{`count > 1 && count < 3`, clause.AndConditions{}},
		{`count == 1 || count == 3`, clause.OrConditions{}},
		{`!(count == 1)`, clause.NotConditions{}},
		{`(count == 1)`, clause.Eq{}},
	} {
		t.Run(tc.input, func(t *testing.T) {
			expr, err := Filter(tc.input).CompileToSql()
			require.NoError(t, err)
			require.IsType(t, tc.expression, expr)
		})
	}
	expr, err := Filter(`name == "alice"`).CompileToSql()
	require.NoError(t, err)
	require.Equal(t, clause.Eq{Column: "name", Value: "alice"}, expr)
}

func TestFilterRejectsUnsupportedExpressions(t *testing.T) {
	for _, input := range []string{`count + 1`, `fn(count)`, `1 == 2`, `count == fn(2)`, `true`, `count ==`} {
		t.Run(input, func(t *testing.T) {
			expr, err := Filter(input).CompileToSql()
			require.Error(t, err)
			require.Nil(t, expr)
		})
	}
	_, err := Filter("").compile(nil)
	require.Error(t, err)
	require.Error(t, (&Field{}).Parse(nil))
	require.Error(t, (&Value{}).Parse(nil))
}

func TestQueryFilterOnlyAndStableOrdering(t *testing.T) {
	require.Equal(t, Filter(""), FieldFilters{}.Compile(""))
	require.Equal(t, Filter("(count > 1)"), HttpQueryFilter{}.Compile(url.Values{"filter": {"count > 1"}}))
	query := url.Values{"b": {"2"}, "a": {"1"}, "page_size": {"10"}, "empty": {}, "filter": {"count > 1"}}
	for range 10 {
		require.Equal(t, Filter("a == 1 && b == 2 && (count > 1)"), HttpQueryFilter{}.Compile(query))
	}
}
