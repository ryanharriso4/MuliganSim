package models

import (
	"database/sql"
)

type Deck struct {
	DeckID    int    `json:"deckID"`
	Cindex    int    `json:"commander"`
	Cover_img string `json:"coverimg"`
	Deck      []int  `json:"decklist"`
	Name      string `json:"deckName"`
}

type DeckModel struct {
	DB *sql.DB
}

func (d *DeckModel) GetUserDecks(id int) ([]Deck, error) {
	var decks []Deck

	stmt := "select deck_id from user_deck where user_id = ?"
	rows, err := d.DB.Query(stmt, id)
	if err != nil {
		return decks, err
	}
	defer rows.Close()

	var deckID int
	for rows.Next() {
		err = rows.Scan(&deckID)
		if err != nil {
			return decks, err
		}

		deck, err := d.GetDeck(deckID)

		if err != nil {
			return decks, err
		}

		decks = append(decks, deck)
	}

	return decks, nil
}

func (d *DeckModel) GetDeck(id int) (Deck, error) {
	var deck Deck

	stmt := "select name, commander_id, cover_img from deck where id = ?;"
	err := d.DB.QueryRow(stmt, id).Scan(&deck.Name, &deck.Cindex, &deck.Cover_img)
	if err != nil {
		return deck, err
	}
	deck.DeckID = id

	return deck, nil

}

func (d *DeckModel) SaveDeckChanges(deck Deck) error {
	var lowIndex int
	err := d.DB.QueryRow("select min(id) from cardlist;").Scan(&lowIndex)
	if err != nil {
		return err
	}

	var highIndex int
	err = d.DB.QueryRow("select max(id) from cardlist").Scan(&highIndex)
	if err != nil {
		return err
	}

	// var deckIndex int64
	// if deck.Cindex >= lowIndex && deck.Cindex <= highIndex {
	// 	result, err := d.DB.Exec("insert into deck(name, commander_id) values(?, ?)", deck.Name, deck.Cindex)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	deckIndex, err = result.RowsAffected()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// for _, value := range deck.Deck {
	// 	fmt.Println()
	// }

	return nil

}
