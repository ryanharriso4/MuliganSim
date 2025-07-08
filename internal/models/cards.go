package models

import (
	"database/sql"
)

type Card struct {
	ID         int
	Name       string
	Image_uri  string
	Mana_cost  string
	Type__line string
	Power      string
	Toughness  string
	Ability    string
	CMC        int
}

type CardModel struct {
	DB *sql.DB
}

func (c *CardModel) GetByName(name string) ([]Card, error) {

	var cards []Card

	stmt := "select name, id, image_url from cardlist where name like '" + name + "%'"
	rows, err := c.DB.Query(stmt)

	if err != nil {
		return cards, err
	}

	for rows.Next() {
		var card Card
		err := rows.Scan(&card.Name, &card.ID, &card.Image_uri)
		if err != nil {
			return cards, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}
