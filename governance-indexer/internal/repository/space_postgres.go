package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/danila-kuryakin/cryptoServises/controller/pkg/models"

	"github.com/lib/pq"
)

type SpacePostgres struct {
	db *sql.DB
}

func NewSpacePostgres(db *sql.DB) *SpacePostgres {
	return &SpacePostgres{db: db}
}

func (s SpacePostgres) AddHistory(space []models.Space) error {
	err := s.Add(space, eventHistory)
	if err != nil {
		return err
	}
	return nil
}

func (s SpacePostgres) AddNew(space []models.Space) error {
	err := s.Add(space, eventHistory)
	if err != nil {
		return err
	}
	return nil
}

// AddProposal Добавляет новые proposals в БД
func (s SpacePostgres) Add(space []models.Space, eventType string) error {
	query := fmt.Sprintf(`
				INSERT INTO %s (space_id, name, about, network, symbol, created,
								strategies_name, admins, members,
								filters_min_score, filters_only_members)
				VALUES ($1,$2,$3,$4,$5,$6,
				        $7,$8,$9,
				        $10,$11) ON CONFLICT (space_id) DO NOTHING
				`, spacesTable)

	queryOutbox := fmt.Sprintf(`
				INSERT INTO %s (space_id, event_type, created_at)
				VALUES ($1, $2, now()) ON CONFLICT (space_id) DO NOTHING
				`, spacesOutboxTable)

	if len(space) == 0 {
		log.Println("AddProposal called with no space")
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		log.Println("Error in proposalPostgres.AddSpace:", err)
		return err
	}

	for _, space := range space {
		strategiesJSON, err := json.Marshal(space.Strategies)
		if err != nil {
			log.Println("AddProposal err:", err)
			return err
		}

		adminsJSON, err := json.Marshal(space.Admins)
		if err != nil {
			log.Println("AddProposal err:", err)
			return err
		}

		membersJSON, err := json.Marshal(space.Members)
		if err != nil {
			log.Println("AddProposal err:", err)
			return err
		}

		_, err = tx.Exec(query,
			space.ID,
			space.Name,
			space.About,
			space.Network,
			space.Symbol,
			space.Created,
			string(strategiesJSON),
			string(adminsJSON),
			string(membersJSON),
			space.Filters.MinScore,
			space.Filters.OnlyMembers)
		if err != nil {
			log.Println("Error to exec space :", err)
			err := tx.Rollback()
			if err != nil {
				log.Println("Error to Rollback space :", err)
				return err
			}
			return err
		}

		_, err = tx.Exec(queryOutbox,
			space.ID,
			eventType,
		)
		if err != nil {
			log.Println("Error to exec outbox space:", err)
			err := tx.Rollback()
			if err != nil {
				log.Println("Error to Rollback outbox space:", err)
				return err
			}
			return err
		}
	}

	return tx.Commit()
}

// FindMissing Возвращает spaces которых нет в БД.
// Читает из БД записи из колонки space_id, количество которых равно spaces в аргументе,
// сравнивает по Space.ID (space_id) каждый элемент
func (s SpacePostgres) FindMissing(spaces []models.Space) ([]models.Space, error) {
	// собираем ID из API
	ids := make([]string, 0, len(spaces))
	for _, p := range spaces {
		ids = append(ids, p.ID)
	}
	query := fmt.Sprintf(`
        SELECT unnest($1::text[]) AS space_id
        EXCEPT
        SELECT space_id FROM %s;
    `, spacesTable)
	// база говорит, каких ID у неё нет
	rows, err := s.db.Query(query, pq.Array(ids))
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

	// собираем отсутствующие spaces
	var missing []models.Space
	for _, p := range spaces {
		if _, ok := missingIDs[p.ID]; ok {
			missing = append(missing, p)
		}
	}

	return missing, nil
}
