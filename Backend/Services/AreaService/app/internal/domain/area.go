package domain

type InputField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AreaAction struct {
	Active   bool         `json:"active"`
	ID       int          `json:"id"`
	Provider string       `json:"provider"`
	Service  string       `json:"service"`
	Title    string       `json:"title"`
	Type     string       `json:"type"`
	Input    []InputField `json:"input"`
}

type AreaReaction struct {
	ID       int          `json:"id"`
	Provider string       `json:"provider"`
	Service  string       `json:"service"`
	Title    string       `json:"title"`
	Input    []InputField `json:"input"`
}

type Area struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Active    bool           `json:"active"`
	UserID    int            `json:"user_id"`
	Actions   []AreaAction   `json:"actions"`
	Reactions []AreaReaction `json:"reactions"`
}

type AreaRepository interface {
	GetUserAreas(userID int) ([]Area, error)
	GetAreaActions(areaID int) ([]AreaAction, error)
	GetAreaReactions(areaID int) ([]AreaReaction, error)
	SaveArea(area Area) (Area, error)
	SaveActions(areaID int, actions []AreaAction) ([]AreaAction, error)
	SaveReactions(areaID int, reactions []AreaReaction) ([]AreaReaction, error)
	GetAreaFromAction(actionId int) (Area, error)
	GetArea(areaID int) (Area, error)
	ToggleArea(areaID int, isActive bool) error
}
