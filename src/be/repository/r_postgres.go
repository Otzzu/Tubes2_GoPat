package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var Db *sqlx.DB

func GetChildrenByParent(parent string) ([]string, error) {
    var results []pq.StringArray
    query := `SELECT children FROM articles WHERE parent = $1`

    err := Db.Select(&results, query, parent)
    if err != nil {
        return nil, fmt.Errorf("error querying children by parent: %v", err)
    }
    if len(results) == 0 {
        return nil, nil  // Return nil explicitly if no children found
    }

    var children []string
    for _, childArray := range results {
        children = append(children, childArray...)
    }

    return children, nil
}

func ExistsParent(parent string) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM articles WHERE parent = $1)"
    err := Db.Get(&exists, query, parent)
    if err != nil {
        return false, err
    }
    return exists, nil
}

func SaveArticleWithChildren(parent string, children []string) error {
    tx, err := Db.BeginTxx(context.Background(), nil)
    if err != nil {
        return err
    }

    var articleID int
    err = tx.QueryRowx("INSERT INTO articles (parent, children) VALUES ($1, $2) RETURNING id", parent, pq.Array(children)).Scan(&articleID)
    if err != nil {
        tx.Rollback()
        var pgErr *pgconn.PgError
        if ok := errors.As(err, &pgErr); ok && pgErr.Code == "23505" {  // Handle unique violation error
            return fmt.Errorf("duplicate entry for parent: %s", parent)
        }
        return fmt.Errorf("failed to insert article, error: %v", err)
    }

    return tx.Commit()
}