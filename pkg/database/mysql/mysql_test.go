package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectMySQL(t *testing.T) {
	// Connect to the MySQL database
	db, err := ConnectMySQL()
	if err != nil {
		t.Errorf("Error connecting to MySQL database: %v", err)
	}

	// Make sure the connection is not nil
	assert.NotNil(t, db)

	// Make sure the connection is valid by querying the database
	var result string
	db.Raw("SELECT VERSION()").Scan(&result)
	assert.NotEmpty(t, result)
}
