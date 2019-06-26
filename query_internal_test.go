package fury

import (
	"reflect"
	"testing"
)

type User struct {
	UserID  int `fury:"primary_key"`
	Counter int `fury:"auto_increment"`
}

func TestNextModel(t *testing.T) {
	cases := []struct {
		have interface{}
	}{
		{&User{}},
		{[]*User{&User{}}},
	}

	for _, tc := range cases {
		q, err := NewQuery(tc.have)
		if err != nil {
			t.Error(err)
		}

		model := q.nextModel()

		if model == nil {
			t.Error("Error: returned model should not be nil")
		}
	}
}

func TestNextOrCreateModel(t *testing.T) {
	cases := []struct {
		have interface{}
	}{
		{&User{}},
		{&[]*User{}},
	}

	for _, tc := range cases {
		q, err := NewQuery(tc.have)
		if err != nil {
			t.Error(err)
		}

		model, err := q.nextOrCreateModel()
		if err != nil {
			t.Error(err)
		}

		if q.modelPtr == nil {
			t.Error("Error: modelPtr should not be nil")
		}

		reflectVal := reflect.ValueOf(tc.have)
		if reflectVal.Kind() == reflect.Ptr {
			reflectVal = reflectVal.Elem()
		}

		if model == nil {
			t.Error("Error: returned model should not be nil")
		}

		if reflectVal.Kind() == reflect.Slice {
			if reflectVal.Len() < 1 {
				t.Error("Error: models should have length more than one")
			}
		}
	}
}

func TestSelectColumns(t *testing.T) {
	cases := []struct {
		have []interface{}
		want string
	}{
		{
			[]interface{}{"user.userid", "user.username"},
			"SELECT user.userid, user.username",
		},
		{
			[]interface{}{"COUNT(user.*)"},
			"SELECT COUNT(user.*)",
		},
		{
			[]interface{}{},
			"SELECT user.*",
		},
	}

	for _, tc := range cases {
		q := &Query{tableName: "user"}

		sel := Select(tc.have...)
		sel(q)

		str, err := q.prepareSelectColumn()
		if err != nil {
			t.Error(err)
		}

		if str != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, str)
		}
	}
}

func TestAddTableName(t *testing.T) {
	cases := []struct {
		want           string
		useModelAsCond bool
	}{
		{"user", true},
		{"", true},
	}

	for _, tc := range cases {
		q := &Query{}

		table := Table(tc.want)
		table(q)

		have := q.tableName

		if have != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, have)
		}
	}
}

func TestUnsupportedWhereConditions(t *testing.T) {
	cases := []struct {
		have interface{}
	}{
		{1},
		{true},
	}

	for _, tc := range cases {
		db := &DB{query: &Query{}}
		whereClause := Where(tc.have)
		if _, err := whereClause(db.query); err == nil {
			t.Error("Expected error found nil")
		}
	}
}

func TestWhereConditions(t *testing.T) {
	cases := []struct {
		have interface{}
		want *Query
	}{
		{
			IsEqualsTo("key", 1),
			&Query{},
		},
		{
			And(IsEqualsTo("key", 1), IsEqualsTo("key", 2)),
			&Query{},
		},
	}

	for _, tc := range cases {
		tc.want.whereConditions = []interface{}{tc.have}
		db := &DB{query: &Query{}}
		whereClause := Where(tc.have)
		if _, err := whereClause(db.query); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(db.query, tc.want) {
			t.Errorf("Error: expected %v, found %v", tc.want, db.query)
		}
	}
}

func TestPKWhereConditions(t *testing.T) {
	cases := []struct {
		have interface{}
		want string
	}{
		{
			&User{},
			" WHERE ",
		},
		{
			&User{UserID: 2},
			" WHERE user.userid = ?",
		},
	}

	for _, tc := range cases {
		q, err := NewQuery(tc.have)
		if err != nil {
			t.Error(err)
		}

		if err := q.addPKWhereConditions(); err != nil {
			t.Error(err)
		}

		str, err := q.prepareWhereQuery()
		if err != nil {
			t.Error(err)
		}

		if str != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, str)
		}
	}
}

func TestAllPKWhereConditions(t *testing.T) {
	cases := []struct {
		have interface{}
		want string
	}{
		{
			&[]*User{
				&User{},
				&User{},
			},
			" WHERE ",
		},
		{
			&[]*User{
				&User{UserID: 2},
				&User{UserID: 3},
			},
			" WHERE (user.userid = ? OR user.userid = ?)",
		},
	}

	for _, tc := range cases {
		q, err := NewQuery(tc.have)
		if err != nil {
			t.Error(err)
		}

		if err := q.addAllPKWhereConditions(); err != nil {
			t.Error(err)
		}

		str, err := q.prepareWhereQuery()
		if err != nil {
			t.Error(err)
		}

		if str != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, str)
		}
	}
}

func TestWhereConditionsToString(t *testing.T) {
	cases := []struct {
		have     []interface{}
		want     string
		wantArgs []interface{}
	}{
		{
			[]interface{}{IsEqualsTo("key", 1), IsEqualsTo("key", 2)},
			" WHERE key = ? AND key = ?",
			[]interface{}{1, 2},
		},
		{
			[]interface{}{And(IsEqualsTo("key", 1), IsEqualsTo("key", 2))},
			" WHERE (key = ? AND key = ?)",
			[]interface{}{1, 2},
		},
	}

	for _, tc := range cases {
		db := &DB{query: &Query{}}

		for _, exps := range tc.have {
			whereClause := Where(exps)
			if _, err := whereClause(db.query); err != nil {
				t.Error(err)
			}
		}

		whereString, err := db.query.prepareWhereQuery()
		if err != nil {
			t.Error(err)
		}

		if whereString != tc.want || !reflect.DeepEqual(tc.wantArgs, db.query.args) {
			t.Errorf("Error: expected %v and %v, found %v and %v", tc.want, tc.wantArgs, whereString, db.query.args)
		}
	}
}

func TestLimitOffset(t *testing.T) {
	cases := []struct {
		haveLimit  int
		haveOffset int
		want       string
	}{
		{1, 0, " LIMIT 1"},
		{0, 2, " OFFSET 2"},
		{2, 3, " LIMIT 2 OFFSET 3"},
	}

	for _, tc := range cases {
		q := &Query{}

		limit := Limit(tc.haveLimit)
		offset := Offset(tc.haveOffset)

		if _, err := limit(q); err != nil {
			t.Error(err)
		}

		if _, err := offset(q); err != nil {
			t.Error(err)
		}

		str := q.prepareLimitOffsetQuery()
		if str != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, str)
		}
	}
}

func TestGroupByColumns(t *testing.T) {
	cases := []struct {
		have []interface{}
		want string
	}{
		{
			[]interface{}{"user.userid", "user.username"},
			" GROUP BY user.userid, user.username",
		},
		{
			[]interface{}{},
			"",
		},
	}

	for _, tc := range cases {
		q := &Query{}

		groupBy := GroupBy(tc.have...)
		groupBy(q)

		str := q.prepareGroupByQuery()

		if str != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, str)
		}
	}
}

func TestOrderByColumns(t *testing.T) {
	cases := []struct {
		have []interface{}
		want string
	}{
		{
			[]interface{}{"user.userid ASC", "user.username DESC"},
			" ORDER BY user.userid ASC, user.username DESC",
		},
		{
			[]interface{}{},
			"",
		},
	}

	for _, tc := range cases {
		q := &Query{}

		orderBy := OrderBy(tc.have...)
		orderBy(q)

		str := q.prepareOrderByQuery()

		if str != tc.want {
			t.Errorf("Error: expected %v, found %v", tc.want, str)
		}
	}
}

func TestReplaceSQLPlaceholders(t *testing.T) {
	cases := []struct {
		have     string
		haveArgs []interface{}
		want     string
	}{
		{
			"user.userid = ?",
			[]interface{}{1},
			"user.userid = $1",
		},
		{
			"user.userid = ? AND user.userid = ?",
			[]interface{}{1, 2},
			"user.userid = $1 AND user.userid = $2",
		},
	}

	for _, tc := range cases {
		q := &Query{
			SQL:  tc.have,
			args: tc.haveArgs,
		}
		q.replaceSQLPlaceholder()

		if tc.want != q.SQL {
			t.Errorf("Error: expected %v, found %v", tc.want, q.SQL)
		}
	}
}

func TestPrepareSelect(t *testing.T) {
	cases := []struct {
		have *Query
		want string
	}{
		{
			&Query{
				tableName:       "user",
				useModelAsCond:  false,
				columns:         []interface{}{"COUNT(user.*)"},
				whereConditions: []interface{}{IsGreaterThan("user.counter", 1)},
				limit:           1,
				offset:          2,
				groups:          []interface{}{"user.counter"},
				orders:          []interface{}{"user.counter DESC"},
			},
			"SELECT COUNT(user.*) FROM user WHERE user.counter > $1 GROUP BY user.counter ORDER BY user.counter DESC LIMIT 1 OFFSET 2;",
		},
	}

	for _, tc := range cases {
		if err := tc.have.prepareSelectQuery(); err != nil {
			t.Error(err)
		}

		if tc.want != tc.have.SQL {
			t.Errorf("Error: expected %s, found %s", tc.want, tc.have.SQL)
		}
	}
}

func TestPrepareInsert(t *testing.T) {
	cases := []struct {
		have interface{}
		want string
	}{
		{
			&User{UserID: 123, Counter: 1},
			"INSERT INTO user(userid) VALUES($1);",
		},
	}

	for _, tc := range cases {
		q, err := NewQuery(tc.have)
		if err != nil {
			t.Error(err)
		}

		if err := q.prepareInsertQuery(); err != nil {
			t.Error(err)
		}

		if tc.want != q.SQL {
			t.Errorf("Error: expected %s, found %s", tc.want, q.SQL)
		}
	}
}

func TestPrepareUpdate(t *testing.T) {
	cases := []struct {
		have interface{}
		want string
	}{
		{
			&User{UserID: 123, Counter: 1},
			"UPDATE user SET (userid,counter) = ($1,$2) WHERE user.userid = $3;",
		},
	}

	for _, tc := range cases {
		q, err := NewQuery(tc.have)
		if err != nil {
			t.Error(err)
		}

		if err := q.prepareUpdateQuery(); err != nil {
			t.Error(err)
		}

		if tc.want != q.SQL {
			t.Errorf("Error: expected %s, found %s", tc.want, q.SQL)
		}
	}
}

func TestPrepareDelete(t *testing.T) {
	cases := []struct {
		have interface{}
		want string
	}{
		{
			&User{UserID: 123, Counter: 1},
			"DELETE FROM user WHERE user.userid = $1;",
		},
	}

	for _, tc := range cases {
		q, err := NewQuery(tc.have)
		if err != nil {
			t.Error(err)
		}

		if err := q.prepareDeleteQuery(); err != nil {
			t.Error(err)
		}

		if tc.want != q.SQL {
			t.Errorf("Error: expected %s, found %s", tc.want, q.SQL)
		}
	}
}
