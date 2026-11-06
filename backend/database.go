package main

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

type Driver interface {
	OpenDatabase(ctx context.Context) *sql.DB
}

var TABLE_NAME string = "moves"
var COLUMNS [4]string = [4]string{"id", "name", "description", "image"}

//Stores a given move into the database
func StoreMove(db *sql.DB , move Move) (int64, error) {
	//TODO: Do I have to colon-terminate?
	query := fmt.Sprintf("INSERT INTO %s (name, dance_style, difficulty, description, image) VALUES ($1, $2, $3, $4, $5)", 
		TABLE_NAME)
	result, err := db.Exec(query, move.Name, move.DanceStyle, move.Difficulty, move.Description, move.Image)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//Gets a move from its id on the database
func GetMoveById(db *sql.DB, id int64) (Move, error){
	//TODO: Do I have to colon-terminate?
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", TABLE_NAME)
	rows, err:= db.Query(query, id)
	if err != nil {
		return Move{}, err
	}
	defer rows.Close()
	var res Move

	//rows.Scan(&res)
	ok := rows.Next()
	err = rows.Err()
	if !ok && err == nil {

	}
	err = rows.Scan(&res.Id, &res.Name, &res.Description, &res.Image)
	return res, err
}

func GetNMovesSorted(db *sql.DB, n int, sortBy string) ([]Move, error) {
	if !slices.Contains(COLUMNS[:], strings.ToLower(sortBy)) {
		return make([]Move, 0), fmt.Errorf("no column named: %s", sortBy)
	}
	res := make([]Move, 0, n)
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY $1 LIMIT $2", TABLE_NAME)
	rows, err := db.Query(query, sortBy, n)
	if err != nil {
		return make([]Move, 0), err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return res, err
		}
		var current Move
		err = rows.Scan(&current.Id, &current.Name, &current.Description, &current.Image)
		if err != nil {
			return res, err
		}
		res = append(res, current)
	}
	return res, rows.Err()
}

func DeleteMove(db *sql.DB, id int64) (bool, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", TABLE_NAME)
	affected, err := db.Exec(query, id)
	if err != nil {
		return false, err
	}
	numAffected, err := affected.RowsAffected()
	if err != nil {
		return false, err
	}
	//This error is a little weird, but really no more than 1 delete should happen
	if numAffected > 1 {
		return false, fmt.Errorf("more than one row was deleted")
	}
	//TODO: rollback?

	return numAffected == 1, nil
}
