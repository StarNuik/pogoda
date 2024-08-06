package pogoda_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

var (
	pgUser = os.Getenv("PG_USER")
	pgUrl  = os.Getenv("PG_URL")
	pgPass = os.Getenv("PG_PASSWORD")
	dbUrl  = fmt.Sprintf("%s:%s@%s", pgUser, pgPass, pgUrl)
)

func ctx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func TestQuery(t *testing.T) {
	require := require.New(t)

	pool, err := sql.Open("pogoda", dbUrl)
	require.Nil(err)

	row := pool.QueryRow("select count(*) from users;")
	require.Nil(row.Err())
}
