package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"governance-indexer/pkg/models"
	"strings"
	"time"

	"github.com/lib/pq"
)

type ProposalPostgres struct {
	db *sql.DB
}

func NewProposalPostgres(db *sql.DB) *ProposalPostgres {
	return &ProposalPostgres{db: db}
}

// AddProposal Добавляет новые proposals в БД
func (p ProposalPostgres) AddProposal(proposals []models.Proposals) error {
	if len(proposals) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(proposals))
	args := make([]interface{}, 0, len(proposals)*9)

	i := 1
	for _, t := range proposals {
		placeholders = append(placeholders, fmt.Sprintf(""+
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i, i+1, i+2, i+3, i+4, i+5, i+6, i+7, i+8, i+9, i+10))

		// choices нужно сериализовать в JSON
		choicesJSON, err := json.Marshal(t.Choices)
		if err != nil {
			return err
		}

		args = append(
			args,
			t.ID,
			t.Title,
			t.Author,
			time.Unix(t.Created, 0),
			time.Unix(t.Start, 0),
			time.Unix(t.End, 0),
			t.Snapshot,
			t.State,
			string(choicesJSON),
			t.Space.ID,
			t.Space.Name)
		i += 11
	}

	query := `
      INSERT INTO proposal(hex_id, title, author, created_at, start_at, end_at, snapshot, state, choices, space_id, space_name)
      VALUES ` + strings.Join(placeholders, ", ") + `
      ON CONFLICT (hex_id) DO NOTHING
  `

	_, err := p.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

// FindMissing Возвращает proposals которых нет в БД.
// Читает из БД записи из колонки hex_id, количество которых равно proposals в аргументе,
// сравнивает по Proposals.ID (hex_id) каждый элемент
func (p ProposalPostgres) FindMissing(proposals []models.Proposals) ([]models.Proposals, error) {
	// собираем ID из API
	ids := make([]string, 0, len(proposals))
	for _, p := range proposals {
		ids = append(ids, p.ID)
	}

	// база говорит, каких ID у неё нет
	rows, err := p.db.Query(`
        SELECT unnest($1::text[]) AS hex_id
        EXCEPT
        SELECT hex_id FROM proposal;
    `, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	missingIDs := make(map[string]struct{})
	for rows.Next() {
		var id string
		_ = rows.Scan(&id)
		missingIDs[id] = struct{}{}
	}

	// собираем отсутствующие proposals
	var missing []models.Proposals
	for _, p := range proposals {
		if _, ok := missingIDs[p.ID]; ok {
			missing = append(missing, p)
		}
	}

	return missing, nil
}
