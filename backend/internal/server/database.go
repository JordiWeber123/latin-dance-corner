package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Driver interface {
	OpenDatabase(ctx context.Context) *sql.DB
}

var TABLE_NAME string = "moves"
var COLUMNS [4]string = [4]string{"id", "name", "description", "image"}
const MAX_PAGE_SIZE uint = 100
//Stores a given move into the database
func CreateMove(db *sql.DB , move Move) (int64, error) {
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
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", TABLE_NAME)
	var res Move 
	err := db.QueryRow(query, id).Scan(&res.Id,
		&res.Name, 
		&res.DanceStyle,
		&res.Difficulty, 
		&res.Description, 
		&res.Image,
		&res.CreatedAt)
	
	if errors.Is(err, sql.ErrNoRows) {
		return Move{}, ErrMoveNotFound{}
	}
	return res, err
}


//Gets page_size amount of moves, using offset to skip around pages
func GetMultiPaginated(db *sql.DB, page_size uint, offset uint) ([]Move, error) {
	page_size = min(MAX_PAGE_SIZE, page_size)
	res := make([]Move, 0, page_size)
	query := fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2",
		 TABLE_NAME)
	rows, err := db.Query(query, page_size, offset)
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
		err = rows.Scan(&current.Id, 
			&current.Name, 
			&current.DanceStyle,
			&current.Difficulty, 
			&current.Description, 
			&current.Image,
			&current.CreatedAt)
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


type ErrMoveNotFound struct{}

func (ErrMoveNotFound) Error() string{
	return "move not found"
}
