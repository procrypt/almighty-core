package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/almighty/almighty-core/migration"
	"github.com/almighty/almighty-core/models"
	"github.com/almighty/almighty-core/resource"
	"github.com/almighty/almighty-core/transaction"
	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var DB *gorm.DB

func TestMain(m *testing.M) {
	if _, c := os.LookupEnv(resource.Database); c == false {
		fmt.Printf(resource.StSkipReasonNotSet+"\n", resource.Database)
		return
	}

	dbhost := os.Getenv("ALMIGHTY_DB_HOST")
	if "" == dbhost {
		panic("The environment variable ALMIGHTY_DB_HOST is not specified or empty.")
	}
	var err error
	DB, err = gorm.Open("postgres", fmt.Sprintf("host=%s user=postgres password=mysecretpassword sslmode=disable", dbhost))
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	defer DB.Close()
	// Migrate the schema
	ts := models.NewGormTransactionSupport(DB)
	witRepo := models.NewWorkItemTypeRepository(ts)

	if err := transaction.Do(ts, func() error {
		return migration.Perform(context.Background(), ts.TX(), witRepo)
	}); err != nil {
		panic(err.Error())
	}
	os.Exit(m.Run())
}

func TestNewWorkitemController(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)
	assert.Panics(t, func() {
		NewWorkitemController(goa.New("Test service"), nil, nil)
	})
}

func TestParseInts(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)

	integers, err := parseInts(nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, []int{}, integers)

	str := "1, 2, foo"
	_, err = parseInts(&str)
	assert.NotNil(t, err)
}

func TestParseLimit(t *testing.T) {
	t.Parallel()
	resource.Require(t, resource.UnitTest)

	// Test parsing error in parseInts
	str := "1000, foo"
	integers, length, err := parseLimit(&str)
	assert.NotNil(t, err)
	assert.Equal(t, 0, length)
	assert.Nil(t, integers)

	// Test length = 1
	str = "1000"
	integers, length, err = parseLimit(&str)
	assert.Nil(t, err)
	assert.Equal(t, 1000, length)
	assert.Nil(t, integers)

	// Test empty string
	str = ""
	integers, length, err = parseLimit(&str)
	assert.Nil(t, err)
	assert.Equal(t, 100, length)
	assert.Nil(t, integers)
}
