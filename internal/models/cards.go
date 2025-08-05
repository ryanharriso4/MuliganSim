package models

import (
	"database/sql"
)

type Card struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Image_uri   string `json:"imageuri"`
	Cropped_uri string `json:"croppeduri"`
	Mana_cost   string `json:"manacost"`
	Type__line  string `json:"typeline"`
	Power       string `json:"power"`
	Toughness   string `json:"toughness"`
	Ability     string `json:"ability"`
	CMC         int    `json:"cmc"`
}

type CardModel struct {
	DB *sql.DB
}

func (c *CardModel) GetByName(name string) ([]Card, error) {

	var cards []Card

	stmt := "select name, id, image_url from cardlist where name like '%" + name + "%'"
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
