package server

import (
	"github.com/iPopcorn/investment-manager/types"
)

type StateRepository struct {
	state types.State // Pass by value so state is immutable
}

func StateRepositoryFactory(initialState types.State) *StateRepository {
	return &StateRepository{
		state: initialState,
	}
}

func (r *StateRepository) GetState() (types.State, error) {
	return r.state, nil
}

func (r *StateRepository) Save(newState types.State) error {
	r.state = newState
	return nil
}
