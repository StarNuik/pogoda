package pgd_test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

var dbUrl = os.Getenv("PG_URL")

func TestOpen(t *testing.T) {
	require := require.New(t)

	db, err := sql.Open("pgd", dbUrl)
	require.Nil(err)

	err = db.Ping()
	require.Nil(err)
}
