package sql

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Custom struct{}

var _ sql.Scanner = (*Custom)(nil)

func (c *Custom) Scan(interface{}) error {
	return nil
}

type User struct {
	ID        uint
	Name      string
	OtherInfo string
	OtherName string `db:"real_name"`
	Ignore    string `db:"-"`
	PtrString *string
	Custom    Custom
}

// testRows is a mock version of sql.Rows which can only scan uint and strings
type testRows struct {
	mock.Mock
	columns []string
	values  []interface{}
	count   int
}

func (r *testRows) Scan(dest ...interface{}) error {
	if len(dest) == len(r.values) {
		for i := range r.values {
			if s, ok := dest[i].(sql.Scanner); ok {
				s.Scan(r.values[i])
				continue
			}

			rt := reflect.TypeOf(dest[i])
			if rt.Kind() != reflect.Ptr {
				panic("Not a pointer!")
			}

			switch dest[i].(type) {
			case **uint:
				**(dest[i].(**uint)) = r.values[i].(uint)
			case **string:
				if *(dest[i].(**string)) != nil {
					**(dest[i].(**string)) = r.values[i].(string)
				}
			default:
				// Do nothing.
			}
		}
	}

	args := r.Called()
	return args.Error(0)
}

func (r *testRows) Columns() ([]string, error) {
	args := r.Called()
	return r.columns, args.Error(0)
}

func (r *testRows) Next() bool {
	r.count++
	return r.count == 1
}

func (r *testRows) addValue(c string, v interface{}) {
	r.columns = append(r.columns, c)
	r.values = append(r.values, v)
}

func createRows() *testRows {
	rows := new(testRows)
	rows.addValue("id", uint(10))
	rows.addValue("name", "name")
	rows.addValue("other_info", "other info")
	rows.addValue("real_name", "real name")
	rows.addValue("ignore", "ignore")
	rows.addValue("ptr_string", "ptr string")

	return rows
}

func TestScan(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	user := User{}
	count, err := Scan(&user, rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, User{
		ID:        uint(10),
		Name:      "name",
		OtherInfo: "other info",
		OtherName: "real name",
		Ignore:    "",
		PtrString: nil,
		Custom:    Custom{},
	}, user)
}

func TestScan_columnError(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(errors.New("error"))

	user := User{}
	count, err := Scan(&user, rows)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), count)
}

func TestScan_scanError(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(errors.New("error"))

	user := User{}
	count, err := Scan(&user, rows)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), count)
}

func TestScan_panicWhenNotPointer(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)

	user := User{}
	assert.Panics(t, func() {
		Scan(user, rows)
	})
}

func TestScan_slice(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	users := []User{}
	count, err := Scan(&users, rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, User{
		ID:        uint(10),
		Name:      "name",
		OtherInfo: "other info",
		OtherName: "real name",
		Ignore:    "",
		PtrString: nil,
		Custom:    Custom{},
	}, users[0])
}

func TestScan_scanner(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	custom := Custom{}
	count, err := Scan(&custom, rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}
