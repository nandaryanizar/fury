package fury_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/nandaryanizar/fury"
)

func TestFirstQuery(t *testing.T) {
	type queryFunc func(out interface{}, opts ...fury.QueryOption) error
	type optionFunc func(conditions interface{}) fury.QueryOption

	cases := []struct {
		have        interface{}
		options     optionFunc
		optionsArgs interface{}
		query       queryFunc
		want        interface{}
	}{
		{
			&Account{
				UserID: 1,
			},
			fury.Where,
			fury.IsEqualsTo("username", "test1"),
			db.First,
			&Account{
				UserID:    1,
				Username:  "test1",
				Password:  "test1",
				Email:     "test1@test.com",
				CreatedOn: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
				LastLogin: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
			},
		},
	}

	for _, tc := range cases {
		if err := tc.query(tc.have, tc.options(tc.optionsArgs)); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(tc.have, tc.want) {
			t.Errorf("Error: expected %v, found %v", tc.want, tc.have)
		}
	}
}

func TestFindQuery(t *testing.T) {
	type queryFunc func(out interface{}, opts ...fury.QueryOption) error
	type optionFunc func(conditions interface{}) fury.QueryOption

	cases := []struct {
		have        interface{}
		options     optionFunc
		optionsArgs interface{}
		query       queryFunc
		want        interface{}
	}{
		{
			&Account{
				UserID: 1,
			},
			fury.Where,
			fury.IsEqualsTo("username", "test1"),
			db.Find,
			&Account{
				UserID:    1,
				Username:  "test1",
				Password:  "test1",
				Email:     "test1@test.com",
				CreatedOn: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
				LastLogin: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
			},
		},
		{
			&Account{},
			fury.Where,
			fury.IsGreaterThan("userid", 0),
			db.Find,
			&Account{
				UserID:    1,
				Username:  "test1",
				Password:  "test1",
				Email:     "test1@test.com",
				CreatedOn: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
				LastLogin: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
			},
		},
		{
			&[]*Account{},
			fury.Where,
			fury.IsGreaterThan("account.userid", 0),
			db.Find,
			&[]*Account{
				&Account{
					UserID:    1,
					Username:  "test1",
					Password:  "test1",
					Email:     "test1@test.com",
					CreatedOn: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
					LastLogin: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
				},
				&Account{
					UserID:    2,
					Username:  "test2",
					Password:  "test2",
					Email:     "test2@test.com",
					CreatedOn: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
					LastLogin: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
				},
			},
		},
	}

	for _, tc := range cases {
		if err := tc.query(tc.have, tc.options(tc.optionsArgs)); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(tc.have, tc.want) {
			t.Errorf("Error: expected %v, found %v", tc.want, tc.have)
		}
	}
}

func TestInsertQuery(t *testing.T) {
	type queryFunc func(out interface{}, opts ...fury.QueryOption) error
	type optionFunc func(conditions interface{}) fury.QueryOption

	cases := []struct {
		have interface{}
		want interface{}
	}{
		{
			&Account{
				UserID: 3,
			},
			&Account{
				UserID:    3,
				Username:  "test3",
				Password:  "test3",
				Email:     "test3@test.com",
				CreatedOn: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
				LastLogin: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
			},
		},
		{
			&[]*Account{
				&Account{
					UserID: 4,
				},
				&Account{
					UserID: 5,
				},
			},
			&[]*Account{
				&Account{
					UserID:    4,
					Username:  "test4",
					Password:  "test4",
					Email:     "test4@test.com",
					CreatedOn: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
					LastLogin: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
				},
				&Account{
					UserID:    5,
					Username:  "test5",
					Password:  "test5",
					Email:     "test5@test.com",
					CreatedOn: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
					LastLogin: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
				},
			},
		},
	}

	for _, tc := range cases {
		if err := db.Insert(tc.want); err != nil {
			t.Error(err)
		}

		if err := db.Find(tc.have); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(tc.want, tc.have) {
			t.Errorf("Error: expected %v, found %v", tc.want, tc.have)
		}
	}
}

func TestUpdateQuery(t *testing.T) {
	type queryFunc func(out interface{}, opts ...fury.QueryOption) error
	type optionFunc func(conditions interface{}) fury.QueryOption

	cases := []struct {
		have interface{}
		want interface{}
	}{
		{
			&[]*Account{
				&Account{
					UserID: 4,
				},
				&Account{
					UserID: 5,
				},
			},
			&[]*Account{
				&Account{
					UserID:    4,
					Username:  "test4",
					Password:  "test4test",
					Email:     "test4@test.com",
					CreatedOn: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
					LastLogin: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
				},
				&Account{
					UserID:    5,
					Username:  "test5test",
					Password:  "test5",
					Email:     "test5@test.com",
					CreatedOn: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
					LastLogin: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
				},
			},
		},
	}

	for _, tc := range cases {
		if err := db.Update(tc.want); err != nil {
			t.Error(err)
		}

		if err := db.Find(tc.have); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(tc.want, tc.have) {
			t.Errorf("Error: expected %v, found %v", tc.want, tc.have)
		}
	}
}

func TestUpdateWhereQuery(t *testing.T) {
	type queryFunc func(out interface{}, opts ...fury.QueryOption) error
	type optionFunc func(conditions interface{}) fury.QueryOption

	cases := []struct {
		have interface{}
		want interface{}
	}{
		{
			&Account{},
			&Account{
				Email:     "update@update.com",
				CreatedOn: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
				LastLogin: time.Date(2016, 06, 22, 19, 10, 25, 0, time.FixedZone("", 0)),
			},
		},
	}

	for _, tc := range cases {
		if err := db.Update(tc.want, fury.Where(fury.IsEqualsTo("account.username", "test2"))); err != nil {
			t.Error(err)
		}

		if err := db.Find(tc.have, fury.Where(fury.IsEqualsTo("account.email", "update@update.com"))); err != nil {
			t.Error(err)
		}

		if accWant, ok := tc.want.(*Account); ok {
			if accHave, haveOk := tc.have.(*Account); haveOk {
				accWant.UserID = accHave.UserID
			}
		}

		if !reflect.DeepEqual(tc.want, tc.have) {
			t.Errorf("Error: expected %v, found %v", tc.want, tc.have)
		}
	}
}

func TestDeleteQuery(t *testing.T) {
	type queryFunc func(out interface{}, opts ...fury.QueryOption) error
	type optionFunc func(conditions interface{}) fury.QueryOption

	cases := []struct {
		have interface{}
	}{
		{
			&[]*Account{
				&Account{
					UserID: 4,
				},
				&Account{
					UserID: 5,
				},
			},
		},
	}

	for _, tc := range cases {
		want := tc.have

		if err := db.Delete(tc.have); err != nil {
			t.Error(err)
		}

		if err := db.Find(tc.have); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(want, tc.have) {
			t.Errorf("Error: expected %v, found %v", want, tc.have)
		}
	}
}

func TestDeleteWhereQuery(t *testing.T) {
	type queryFunc func(out interface{}, opts ...fury.QueryOption) error
	type optionFunc func(conditions interface{}) fury.QueryOption

	cases := []struct {
		have interface{}
		want interface{}
	}{
		{
			&Account{},
			&Account{},
		},
	}

	for _, tc := range cases {
		if err := db.Delete(tc.want, fury.Where(fury.IsEqualsTo("account.username", "test2"))); err != nil {
			t.Error(err)
		}

		if err := db.Find(tc.have, fury.Where(fury.IsEqualsTo("account.username", "test2"))); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(tc.want, tc.have) {
			t.Errorf("Error: expected %v, found %v", tc.want, tc.have)
		}
	}
}
