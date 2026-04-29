package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	customError "github.com/danila-kuryakin/cryptoServises/controller/errors"
	"github.com/danila-kuryakin/cryptoServises/controller/pkg/models"
)

type ProposalPostgres struct {
	db *sql.DB
}

func NewProposalPostgres(db *sql.DB) *ProposalPostgres {
	return &ProposalPostgres{db: db}
}

func (p *ProposalPostgres) ReadProposalsEvents() ([]models.ProposalEvent, error) {

	query := fmt.Sprintf(`SELECT prop.hex_id, prop.created_at, prop.start_at, prop.end_at				       
				FROM %s AS prop
				LEFT JOIN %s AS evn ON prop.hex_id = evn.hex_id
				WHERE evn.processed_at IS NULL`, proposalsTable, eventOutboxTable)

	rows, err := p.db.Query(query)
	if err != nil {
		log.Println("Query:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Close:", err)
		}
	}(rows)

	var proposals []models.ProposalEvent

	for rows.Next() {
		var p models.ProposalEvent
		if err := rows.Scan(&p.ID, &p.Created, &p.Start, &p.End); err != nil {
			return nil, err
		}
		proposals = append(proposals, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return proposals, nil
}

func (p *ProposalPostgres) ProposalDeliverySuccessful(proposals []models.ProposalEvent) error {
	tx, err := p.db.Begin()
	if err != nil {
		log.Println("Error in Begin:", err)
		return err
	}
	for _, proposal := range proposals {
		query := fmt.Sprintf(`
				UPDATE %s
				SET processed_at = NOW()
				WHERE hex_id = $1;
				`, eventOutboxTable)

		_, err := tx.Exec(query, proposal.ID)
		if err != nil {
			log.Println("Error in Exec:", err)
			return tx.Rollback()
		}
	}

	return tx.Commit()
}

const (
	eventCreate = "create"
	eventStart  = "start"
	eventEnd    = "end"
)

func (p *ProposalPostgres) AddEventScheduler(proposals []models.ProposalEvent) error {
	tx, err := p.db.Begin()
	if err != nil {
		log.Println("Error in Begin:", err)
		return err
	}
	for _, proposal := range proposals {
		query := fmt.Sprintf(`
				INSERT INTO %s (hex_id, event_type, event_at)
				VALUES ($1, $2, $3)
			`, eventSchedulerTable)

		if proposal.Created.Valid {
			if _, err := tx.Exec(query, proposal.ID, eventCreate, proposal.Created); err != nil {
				log.Println("Error in Exec Created:", err)
				return tx.Rollback()
			}
		}
		if proposal.Start.Valid {
			if _, err := tx.Exec(query, proposal.ID, eventStart, proposal.Start); err != nil {
				log.Println("Error in Exec Start:", err)
				return tx.Rollback()
			}
		}
		if proposal.End.Valid {
			if _, err := tx.Exec(query, proposal.ID, eventEnd, proposal.End); err != nil {
				log.Println("Error in Exec End:", err)
				return tx.Rollback()
			}
		}
	}
	return tx.Commit()
}

func (p *ProposalPostgres) GetCurrentEvents(number int64) ([]models.CurrentEvent, error) {

	query := fmt.Sprintf(`
			SELECT evn.hex_id, evn.event_type, evn.event_at, prop.space_id, prop.space_name, prop.title
			FROM %s AS evn
			LEFT JOIN %s AS prop ON evn.hex_id = prop.hex_id
			WHERE evn.processed_at IS NULL
			ORDER BY evn.event_at ASC
			LIMIT $1
		`, eventSchedulerTable, proposalsTable)

	rows, err := p.db.Query(query, number)
	if err != nil {
		log.Println("Query:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Close:", err)
		}
	}(rows)

	var event []models.CurrentEvent

	for rows.Next() {
		var c models.CurrentEvent
		if err := rows.Scan(&c.ID, &c.EventType, &c.EventTime, &c.SpaceID, &c.SpaceName, &c.Title); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, customError.ErrDataNotFound
			}
			return nil, err
		}
		event = append(event, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return event, nil
}

func (p *ProposalPostgres) EventDeliverySuccessful(events []models.CurrentEvent) error {
	tx, err := p.db.Begin()
	if err != nil {
		log.Println("Error in Begin:", err)
		return err
	}
	for _, event := range events {
		query := fmt.Sprintf(`
				UPDATE %s
				SET processed_at = NOW()
				WHERE hex_id = $1 AND event_type = $2;
				`, eventSchedulerTable)

		_, err := tx.Exec(query, event.ID, event.EventType)
		if err != nil {
			log.Println("Error in Exec:", err)
			return tx.Rollback()
		}
	}

	return tx.Commit()
}
