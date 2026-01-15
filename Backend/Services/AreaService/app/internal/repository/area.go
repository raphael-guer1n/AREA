package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/raphael-guer1n/AREA/AreaService/internal/domain"
)

type areaRepository struct {
	db *sql.DB
}

func (a areaRepository) ToggleArea(areaID int, isActive bool) error {
	_, err := a.db.Exec("UPDATE areas SET active = $1 WHERE id = $2", isActive, areaID)
	return err
}

func (a areaRepository) GetArea(areaID int) (domain.Area, error) {
	row, err := a.db.Query("SELECT id, name, active, user_id FROM areas WHERE id = $1", areaID)
	if err != nil {
		return domain.Area{}, err
	}
	var area domain.Area
	row.Next()
	err = row.Scan(&area.ID, &area.Name, &area.Active, &area.UserID)
	row.Close()
	if err != nil {
		return domain.Area{}, err
	}
	area.Actions, err = a.GetAreaActions(areaID)
	if err != nil {
		return domain.Area{}, err
	}
	area.Reactions, err = a.GetAreaReactions(areaID)
	if err != nil {
		return domain.Area{}, err
	}
	return area, nil
}

func (a areaRepository) GetAreaFromAction(actionId int) (domain.Area, error) {
	row, err := a.db.Query("SELECT area_id FROM actions WHERE id = $1", actionId)
	if err != nil {
		return domain.Area{}, err
	}
	var areaID int
	row.Next()
	err = row.Scan(&areaID)
	row.Close()
	if err != nil {
		return domain.Area{}, err
	}
	var area domain.Area
	row, err = a.db.Query("SELECT id, name, active, user_id FROM areas WHERE id = $1", areaID)
	if err != nil {
		return domain.Area{}, err
	}
	row.Next()
	err = row.Scan(&area.ID, &area.Name, &area.Active, &area.UserID)
	row.Close()
	if err != nil {
		return domain.Area{}, err
	}
	area.Actions, err = a.GetAreaActions(areaID)
	if err != nil {
		return domain.Area{}, err
	}
	area.Reactions, err = a.GetAreaReactions(areaID)
	if err != nil {
		return domain.Area{}, err
	}
	return area, nil
}

func (a areaRepository) SaveReactions(areaID int, reactions []domain.AreaReaction) ([]domain.AreaReaction, error) {
	for _, reaction := range reactions {
		inputJSON, err := json.Marshal(reaction.Input)
		if err != nil {
			return reactions, err
		}
		err = a.db.QueryRow(`INSERT INTO reactions (area_id, provider, service, title, inputs) VALUES ($1, $2, $3, $4, $5) RETURNING id`, areaID, reaction.Provider, reaction.Service, reaction.Title, inputJSON).Scan(&reaction.ID)
		if err != nil {
			return reactions, err
		}
	}
	return reactions, nil
}

func (a areaRepository) SaveActions(areaID int, actions []domain.AreaAction) ([]domain.AreaAction, error) {
	for i, action := range actions {
		inputJSON, err := json.Marshal(action.Input)
		if err != nil {
			return actions, err
		}
		err = a.db.QueryRow(`INSERT INTO actions (area_id, provider, service, title, inputs, type) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, areaID, action.Provider, action.Service, action.Title, inputJSON, action.Type).Scan(&actions[i].ID)
		if err != nil {
			return actions, err
		}
	}
	return actions, nil
}

func (a areaRepository) SaveArea(area domain.Area) (domain.Area, error) {
	var areaID int
	err := a.db.QueryRow(`INSERT INTO areas (name, active, user_id) VALUES ($1, $2, $3) RETURNING id`, area.Name, area.Active, area.UserID).Scan(&areaID)
	if err != nil {
		return area, err
	}
	area.ID = areaID
	actions, err := a.SaveActions(areaID, area.Actions)
	if err != nil {
		return area, err
	}
	area.Actions = actions
	reactions, err := a.SaveReactions(areaID, area.Reactions)
	if err != nil {
		return area, err
	}
	area.Reactions = reactions
	return area, nil
}

func (a areaRepository) GetAreaReactions(areaID int) ([]domain.AreaReaction, error) {
	rows, err := a.db.Query("SELECT id, provider, service, title, inputs FROM reactions WHERE area_id = $1", areaID)
	if err != nil {
		return nil, err
	}
	reactions := make([]domain.AreaReaction, 0)
	for rows.Next() {
		var reaction domain.AreaReaction
		var inputJSON []byte
		if err := rows.Scan(&reaction.ID, &reaction.Provider, &reaction.Service, &reaction.Title, &inputJSON); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(inputJSON, &reaction.Input); err != nil {
			return nil, err
		}
		reactions = append(reactions, reaction)
	}
	return reactions, nil
}

func (a areaRepository) GetAreaActions(areaID int) ([]domain.AreaAction, error) {
	rows, err := a.db.Query("SELECT id, provider, service, title, inputs, type FROM actions WHERE area_id = $1", areaID)
	if err != nil {
		return nil, err
	}
	actions := make([]domain.AreaAction, 0)

	for rows.Next() {
		var action domain.AreaAction
		var inputJSON []byte
		if err := rows.Scan(&action.ID, &action.Provider, &action.Service, &action.Title, &inputJSON, &action.Type); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(inputJSON, &action.Input); err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, nil
}

func (a areaRepository) GetUserAreas(userID int) ([]domain.Area, error) {
	rows, err := a.db.Query("SELECT id, name, active FROM areas WHERE user_id = $1", userID)
	areas := make([]domain.Area, 0)

	if err != nil {
		return nil, err
	}
	if rows == nil {
		return areas, nil
	}
	for rows.Next() {
		var area domain.Area
		if err := rows.Scan(&area.ID, &area.Name, &area.Active); err != nil {
			return nil, err
		}
		area.Actions, err = a.GetAreaActions(area.ID)
		if err != nil {
			return nil, err
		}
		area.Reactions, err = a.GetAreaReactions(area.ID)
		if err != nil {
			return nil, err
		}
		areas = append(areas, area)
	}
	return areas, nil
}

func (a areaRepository) DeleteArea(areaID int) error {
	_, err := a.db.Exec("DELETE FROM areas WHERE id = $1", areaID)
	return err
}

func (a areaRepository) DeactivateAreasByProvider(userID int, provider string) (int, error) {
	query := `
		UPDATE areas
		SET active = false
		WHERE user_id = $1
		  AND active = true
		  AND id IN (
			SELECT DISTINCT area_id
			FROM (
			  SELECT area_id FROM actions WHERE provider = $2
			  UNION
			  SELECT area_id FROM reactions WHERE provider = $2
			) AS areas_with_provider
		  )
	`
	result, err := a.db.Exec(query, userID, provider)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}

func NewAreaRepository(db *sql.DB) domain.AreaRepository {
	return &areaRepository{db: db}
}
