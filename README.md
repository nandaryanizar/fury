# Fury

Fury is an ORM for PostgreSQL database that use functional option pattern design.

## Requirements

To use this library, we have to install a few prerequisites:

* Go 1.12 or newer
* PostgreSQL 9.5 (newer version maybe compatible)
* Docker 18.09 or newer (optional, use for testing only)

## Installation

```
go get github.com/nandaryanizar/fury
```

## How to Use

Fury support basic query such as SELECT, INSERT, UPDATE, DELETE. Here's a few example on how to use this library.

### Initialization

To use this library we should have PostgreSQL up and running. First, we need to create environment variable file, for example:

```yaml
# Environment variable example
# database.yaml

# Required environment variables
DATABASE_USERNAME: postgres
DATABASE_PASSWORD: pgadmin123
DATABASE_HOST: postgres_db
DATABASE_PORT: 5432
DATABASE_NAME: testdb

# Optional environment variables
DATABASE_SSLMODE=false
DATABASE_MAXRETRIES=1
DATABASE_MAXIDLECONNS=2
DATABASE_MAXOPENCONNS=0
DATABASE_CONNMAXXLIFETIME=0
```

To connect to the database, we can use `Connect` function. The function will return the DB struct consists of connection pool, configuration, and query struct. Only at the end of the program we need to close the DB connection by calling `Close` method.

```go
// Create DB connection pool
db, err := fury.Connect("database.yaml")

// Close DB connection
defer db.Close()
```

This return the result of SELECT query to pointer(s) to struct. For INSERT, UPDATE and DELETE, this library will generate the query and arguments based on the passed pointer(s) to struct. So for this example, supposed we have below struct for example:

```go
type Account struct {
	UserID    int `fury:"primary_key,auto_increment"`
	Username  string
	Password  string
	Email     string
	CreatedOn time.Time
	LastLogin time.Time
}
```

As this readme file is written, this library neither support converting struct name to another convention (e.g. snake case) nor using tag. It assume database table column name is the lowercase version of the struct field name (e.g. UserID field will be mapped to userid column).

#### Tags

Fury also support some tags, currently `primary_key` and `auto_increment`. These tags are useful when generating query. Field with tag `primary_key` will be used as where condition if the value is not zero value of the type. It will also be ignored in `UPDATE` query when the value is zero value of the type. In `INSERT` query, `auto_increment` tagged field will be ignored as well.

### SELECT Query

Currently, there's two main method to generate a SELECT query. The first is `Find` method. Supposed we want to generate query as simple as `SELECT * FROM account`, we can do this:

```go
// Create an empty slice of pointer to Account struct
accounts := []*Account{}

// Then pass the address of the slice as parameter like below example
db.Find(&accounts)
```

The results will be appended to the slice. If the slice is not empty, it will use the available element to scan the result and create new and append the remaining results if the number of results is more than the number of element in the slice.

We also can use pointer to struct to be passed to the method. But only the first record will be scanned to the struct. Below is the example:

```go
// Create an empty Account struct
account := Account{}

// Generate the same `SELECT * FROM account` query, but only the first record is scanned to the struct
db.Find(&account)
```

### SELECT Query with WHERE Conditions

If we want to add condition to our query we can use `Where` query option. Supposed we want to add condition from previous query to return only the record with `Username` equals to `nandaryanizar`, we can do as below:

```go
// Call find method with addition query option parameter
db.Find(&account, fury.Where(fury.IsEqualsTo("username", "nandaryanizar")))

// or

db.Find(&account, fury.Where("username = nandaryanizar"))
```

The above call will generate query SELECT * FROM account WHERE username = 'nandaryanizar'. The `Where` query option takes `Expression`, `LogicalExpression`, and `string` as parameters.

When there are `Where` query option passed as parameter, it will be treated with AND query condition.

```go
// Call find method with addition query option parameter
db.Find(&account,
    fury.Where(fury.IsEqualsTo("username", "nandaryanizar")),
    fury.Where(fury.IsEqualsTo("password", "test")),
)

// or

db.Find(&account,
    fury.Where(
        And(
            fury.IsEqualsTo("username", "nandaryanizar"),
            fury.IsEqualsTo("password", "test")),
    ),
)

// Generate the same `SELECT * FROM acccount WHERE username = 'nandaryanizar' AND password = 'test'`
```

To create condition with OR query condition, we can create it like below:

```go
db.Find(&account,
    fury.Where(
        Or(
            fury.IsEqualsTo("username", "nandaryanizar"),
            fury.IsEqualsTo("email", "some@test.com")),
    ),
)
````

If we want to query that use the primary key as the condition, we only need to fill the field in the struct without having to explicitly add `Where` query option.

```go
// Create struct with non-zero field tagged with `primary_key`
account := Account{
    UserID: 1,
}

// This method will generate `SELECT * FROM account WHERE account.userid = 1`
db.Find(&account)

// If slice of pointer to struct have non-zero field tagged with `primary_key`
accounts := []*Account{
    &Account{
        UserID: 1,
    },
    &Account{
        UserID: 2,
    },
}

// This method will generate `SELECT * FROM account WHERE account.userid = 1 OR account.userid = 2`
db.Find(&accounts)
```

There are several query option other than `Where`, they can be used like this:

```go
// This will select only column username and email
// Generate `SELECT username, email FROM account`
db.Find(&account, Select("username", "email"))

// This method will set the table name to query on, this method also disable the query to generate the primary key condition based on the pointer to struct passed
// This will generate `SELECT * FROM acct`
db.Find(&account, Table("acct"))

// To use order by query, simply use below method
// The method will generate `SELECT * FROM account ORDER BY createdon DESC
db.Find(&account, OrderBy("createdon DESC"))

// The method below will add limit and offset query
// The generated query will look like `SELECT * FROM account LIMIT 1 OFFSET 2`
db.Find(&account, 
    fury.Limit(1),
    fury.Offset(2),
)
```

Although the code for group by query has been added, it is not quite useful right now as this library has not been support aggregate query.

The second method to get the generate query is `First` method. The method will get only the first record, it is equivalent to add `Limit` query option with argument 1 to `Find` method.

```go
// This first method will generate `SELECT * FROM account LIMIT 1`
db.First(&account)

// Equivalent to the above method
db.Find(&account, fury.Limit(1))
```

### INSERT Query

The `Insert` method only takes `Table` as working query option to specify the table name. The `Insert` method will generate insert query, omitting the field with `auto_increment` tag. The `Insert` method can be used as follows:

```go
// Create and initialize the struct that wants to be inserted to database
account := Account{
    Username:  "nandaryanizar",
    Password:  "test",
    Email:     "some@test.com",
    CreatedOn: time.Now(),
    LastLogin: time.Now(),
}

// Simply use the method like this
// And it will generate query equivalent to INSERT INTO account(username, password, email, createdon, lastlogin) VALUES("nandaryanizar", "test", "some@test.com", "2019-06-28T02:26:00+07.000", "2019-06-28T02:26:00+07.000")
db.Insert(&account)

// Insert also support pointer to slice of pointer to struct
accounts := []*Account{
    &Account{
        Username:  "nandaryanizar",
        Password:  "test",
        Email:     "some@test.com",
        CreatedOn: time.Now(),
        LastLogin: time.Now(),
    },
    &Account{
        Username:  "otheruser",
        Password:  "test",
        Email:     "someotheruser@test.com",
        CreatedOn: time.Now(),
        LastLogin: time.Now(),
    }
}

// This method will insert above struct values to database in separate insert query for each pointer to struct
db.Insert(&accounts)
```

### UPDATE Query

In current implementation, `Update` method will generate UPDATE query based on all passed struct field, whether it is zero value or non-zero value, except field with tag `primary-key`, which will be omitted if it contains zero value for the field type. This implementation need to be reviewed and enhance or change, so the zero value can somehow bet omitted, either using tag or omitted by default.

`Update` method can be used like this:

```go
// Create and initialize the struct that wants to be updated to database
account := Account{
    Username:  "nandaryanizar",
    Password:  "test123",
    Email:     "some@test.com",
    CreatedOn: time.Now(),
    LastLogin: time.Now(),
}

// Simply use the method like this
// And it will generate query equivalent to UPDATE account SET (username, password, email, createdon, lastlogin) = ("nandaryanizar", "test", "some@test.com", "2019-06-28T02:26:00+07.000", "2019-06-28T02:26:00+07.000")
db.Update(&account, fury.IsEqualsTo("username", "nandaryanizar"))

// To update multiple record
accounts := []*Account{
    &Account{
        UserID: 1,
        Username:  "nandaryanizar",
        Password:  "test",
        Email:     "some@test.com",
        CreatedOn: time.Now(),
        LastLogin: time.Now(),
    },
    &Account{
        UserID: 2,
        Username:  "otheruser",
        Password:  "test",
        Email:     "someotheruser@test.com",
        CreatedOn: time.Now(),
        LastLogin: time.Now(),
    }
}

// Although one method call, the query will be run separately for different struct, each struct has its own query context. This method will generate two query that will be executed in separate connection.
db.Update(&accounts)

// Query 1: UPDATE account SET (userid, username, password, email, createdon, lastlogin) = (1, "nandaryanizar", "test", "some@test.com", "2019-06-28T02:26:00+07.000", "2019-06-28T02:26:00+07.000") WHERE userid = 1

// Query 2: UPDATE account SET (userid, username, password, email, createdon, lastlogin) = (1, "otheruser", "test", "someotheruser@test.com", "2019-06-28T02:26:00+07.000", "2019-06-28T02:26:00+07.000") WHERE userid = 1
```

### DELETE Query

The `Delete` method pretty much works the same way as `Update` method, but current implementation will prevent to run the method without `Where` condition specified from `primary_key` tag or `Where` query option itself.

We can use `Delete` method like below:

```go
// To update multiple record
accounts := []*Account{
    &Account{
        UserID: 1,
        Username:  "nandaryanizar",
        Password:  "test",
        Email:     "some@test.com",
        CreatedOn: time.Now(),
        LastLogin: time.Now(),
    },
    &Account{
        UserID: 2,
        Username:  "otheruser",
        Password:  "test",
        Email:     "someotheruser@test.com",
        CreatedOn: time.Now(),
        LastLogin: time.Now(),
    }
}

// Call delete method
db.Update(&accounts)

// The above method will generate two query
// Query 1: DELETE FROM account WHERE UserID = 1
// Query 2: DELETE FROM account WHERE UserID = 2

// Create and initialize the struct that wants to be updated to database without field with `primary_key` intialized
account := Account{
    Username:  "nandaryanizar",
    Password:  "test123",
    Email:     "some@test.com",
    CreatedOn: time.Now(),
    LastLogin: time.Now(),
}

// This will return error as there is not any single where condition in the query
db.Update(&account)
```

Above are some example usage of this library. This library still need improvements to better suit the real cases.