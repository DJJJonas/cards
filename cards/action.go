package cards

import "github.com/google/uuid"

var (
	actionTypes = []string{Play, Heropower, Attack, EndTurn}
)

func IsActionValid(action *Action) bool {
	if !isAnActionType(action.Type) {
		return false
	}
	switch action.Type {
	case Play:
		// this action also requires position but if not specified, default to 0
		_, err := uuid.Parse(action.SourceId)
		if err != nil {
			return false
		}
	case Attack:
		_, err := uuid.Parse(action.SourceId)
		if err != nil {
			return false
		}
		_, err = uuid.Parse(action.TargetId)
		if err != nil {
			return false
		}
	}
	return true
}

func isAnActionType(action string) bool {
	for _, at := range actionTypes {
		if at == action {
			return true
		}
	}
	return false
}

type Action struct {
	Type     string
	SourceId string
	TargetId string
	Position int
}
