package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Deck struct {
	DeckID    int
	Cindex    int
	Cover_img string
	Deck_IDS  []int
	Cards     map[string][]Card
	Commander Card
	Name      string
}

type SaveDeck struct {
	DeckID int    `json:"deckID"`
	CIndex int    `json:"commander"`
	Name   string `json:"deckName"`
	Add    []int  `json:"addToDeck"`
	Remove []int  `json:"removeFromDeck"`
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

func (d *DeckModel) UserDeckCheck(userID int, deckID int) int {
	stmt := "select exists(select * from user_deck where user_id = ? and deck_id = ?);"
	var result int
	err := d.DB.QueryRow(stmt, userID, deckID).Scan(&result)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func (d *DeckModel) GetDeck(id int) (Deck, error) {
	var deck Deck

	var commander_id sql.NullInt64
	var cover_img sql.NullString

	stmt := "select name, commander_id, cover_img from deck where id = ?;"
	err := d.DB.QueryRow(stmt, id).Scan(&deck.Name, &commander_id, &cover_img)
	if err != nil {
		return deck, err
	}
	deck.DeckID = id

	if !commander_id.Valid {
		deck.Cindex = -1
	} else {
		deck.Cindex = int(commander_id.Int64)
	}

	if !cover_img.Valid {
		deck.Cover_img = ""
	} else {
		deck.Cover_img = cover_img.String
	}

	return deck, nil

}

func (d *DeckModel) GetDeckWithCards(id int) (Deck, error) {
	var deck Deck

	var commander_id sql.NullInt64
	var cover_img sql.NullString

	deck.DeckID = id
	stmt := "select name, commander_id, cover_img from deck where id = ?;"
	err := d.DB.QueryRow(stmt, id).Scan(&deck.Name, &commander_id, &cover_img)
	if err != nil {
		return deck, err
	}

	if !commander_id.Valid {
		deck.Cindex = -1
	} else {
		deck.Cindex = int(commander_id.Int64)
	}

	if !cover_img.Valid {
		deck.Cover_img = ""
	} else {
		deck.Cover_img = cover_img.String
	}

	stmt = "select CL.id, CL.name, medimage_url, mana_cost, type_line,power, toughness, ability, cmc, croppedimage_url  from cardlist as CL join deck_card as DC on CL.id = DC.card_id join deck as D on DC.deck_id = D.id where D.id = ?;"
	rows, err := d.DB.Query(stmt, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return deck, err
		} else {
			return deck, err
		}

	}

	cardMap := map[string][]Card{
		"Artifact":     {},
		"Creature":     {},
		"Enchantment":  {},
		"Instant":      {},
		"Land":         {},
		"Planeswalker": {},
		"Sorcery":      {},
	}

	for rows.Next() {
		var card Card
		err = rows.Scan(&card.ID, &card.Name, &card.Image_uri, &card.Mana_cost, &card.Type__line, &card.Power, &card.Toughness, &card.Ability, &card.CMC, &card.Cropped_uri)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return deck, err
			} else {
				return deck, err
			}
		}

		switch {
		case strings.Contains(card.Type__line, "Creature"):
			cardMap["Creature"] = append(cardMap["Creature"], card)
		case strings.Contains(card.Type__line, "Artifact"):
			cardMap["Artifact"] = append(cardMap["Artifact"], card)
		case strings.Contains(card.Type__line, "Enchantment"):
			cardMap["Enchantment"] = append(cardMap["Enchantment"], card)
		case strings.Contains(card.Type__line, "Instant"):
			cardMap["Instant"] = append(cardMap["Instant"], card)
		case strings.Contains(card.Type__line, "Land"):
			cardMap["Land"] = append(cardMap["Land"], card)
		case strings.Contains(card.Type__line, "Planeswalker"):
			cardMap["Planeswalker"] = append(cardMap["Planeswalker"], card)
		case strings.Contains(card.Type__line, "Sorcery"):
			cardMap["Sorcery"] = append(cardMap["Sorcery"], card)

		}
		// cardMap[card.Type__line] = append(cardMap[card.Type__line], card)
	}

	deck.Cards = cardMap

	var card Card
	stmt = "select CL.id, CL.name, medimage_url, mana_cost, type_line,power, toughness, ability, cmc, croppedimage_url from cardlist as CL, deck as D where D.commander_id = CL.id and D.id = ?"
	err = d.DB.QueryRow(stmt, id).Scan(&card.ID, &card.Name, &card.Image_uri, &card.Mana_cost, &card.Type__line, &card.Power, &card.Toughness, &card.Ability, &card.CMC, &card.Cropped_uri)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			deck.Cindex = -1
			return deck, err
		} else {
			return deck, err
		}
	}

	deck.Commander = card
	return deck, nil
}

func (d *DeckModel) SaveDeckChanges(ctx context.Context, deck SaveDeck, userID int) (int, error) {

	result := d.UserDeckCheck(userID, deck.DeckID)

	if result == 0 && deck.DeckID != -1 {
		return -1, ErrInvalidCredentials
	}

	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	if deck.DeckID == -1 {
		stmt, err := tx.PrepareContext(ctx, "insert into deck(name) values(?)")
		if err != nil {
			return -1, err
		}

		result, err := stmt.Exec(deck.Name)
		if err != nil {
			return -1, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return -1, err
		}

		deck.DeckID = int(id)
	} else {
		var lowIndex sql.NullInt64
		var highIndex sql.NullInt64
		err := d.DB.QueryRowContext(ctx, "select min(id), max(id) from deck;").Scan(&lowIndex, &highIndex)
		if err != nil {
			return -1, err
		}

		if (!lowIndex.Valid || !highIndex.Valid) || deck.DeckID < int(lowIndex.Int64) || deck.DeckID > int(highIndex.Int64) {
			return -1, ErrInvalidDeckID
		}

		stmt, err := d.DB.PrepareContext(ctx, "update deck set name = ? where id = ?")
		if err != nil {
			return -1, err
		}

		_, err = stmt.ExecContext(ctx, deck.Name, deck.DeckID)
		if err != nil {
			return -1, err
		}
	}

	var lowComIndex int
	var highComIndex int
	err = tx.QueryRowContext(ctx, "select min(id), max(id) from cardlist;").Scan(&lowComIndex, &highComIndex)
	if err != nil {
		return -1, err
	}

	if deck.CIndex >= lowComIndex && deck.CIndex <= highComIndex {
		stmt, err := tx.PrepareContext(ctx, "update deck set commander_id = ? where id = ?")
		if err != nil {
			return -1, err
		}
		defer stmt.Close()
		_, err = stmt.Exec(deck.CIndex, deck.DeckID)
		if err != nil {
			return -1, err
		}

		stmt2, err := tx.PrepareContext(ctx, "update deck set cover_img = (select croppedimage_url from cardlist where cardlist.id = ?) where id = ?")
		if err != nil {
			return -1, err
		}
		defer stmt2.Close()
		_, err = stmt2.Exec(deck.CIndex, deck.DeckID)
		if err != nil {
			return -1, err
		}
	}

	addStmt, err := tx.PrepareContext(ctx, "insert ignore into deck_card(deck_id, card_id) values(?, ?)")
	if err != nil {
		return -1, err
	}
	defer addStmt.Close()

	for _, value := range deck.Add {
		_, err = addStmt.Exec(deck.DeckID, value)
		if err != nil {
			return -1, err
		}
	}

	removeStmt, err := tx.PrepareContext(ctx, "delete from deck_card where deck_id = ? and card_id = ?")
	if err != nil {
		return -1, err
	}
	defer removeStmt.Close()

	for _, value := range deck.Remove {
		_, err = removeStmt.Exec(deck.DeckID, value)
		if err != nil {
			return -1, err
		}
	}

	linkStmt, err := tx.PrepareContext(ctx, "insert ignore into user_deck(user_id, deck_id) values(?, ?)")
	if err != nil {
		return -1, err
	}
	defer linkStmt.Close()

	_, err = linkStmt.Exec(userID, deck.DeckID)
	if err != nil {
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		return -1, err
	}

	return deck.DeckID, nil
}

func (d *DeckModel) DeleteDeck(deckID int, ctx context.Context) error {

	tx, err := d.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "delete from deck_card where deck_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, deckID)
	if err != nil {
		return err
	}

	stmt2, err := tx.PrepareContext(ctx, "delete from user_deck where deck_id = ?")
	if err != nil {
		return err
	}
	defer stmt2.Close()

	_, err = stmt2.ExecContext(ctx, deckID)
	if err != nil {
		return err
	}

	stmt3, err := tx.PrepareContext(ctx, "delete from deck where id = ?")
	if err != nil {
		return err
	}
	defer stmt3.Close()

	_, err = stmt3.ExecContext(ctx, deckID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
