package state

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iPopcorn/investment-manager/types"
	"github.com/iPopcorn/investment-manager/util"
)

type StateRepository struct {
	filename string
}

func StateRepositoryFactory(filename string) *StateRepository {
	defaultName := "/server/state/state.json"

	if filename != "" {
		return &StateRepository{
			filename: filename,
		}
	}

	return &StateRepository{
		filename: defaultName,
	}
}

func (r *StateRepository) GetState() (*types.State, error) {
	location := "StateRepository.GetState()\n"
	filepath, err := util.GetPathToFile("/server/state", r.filename)

	if err != nil {
		fmt.Printf(location+"Failed to get path to file\n%v\n", err)
		return nil, err
	}

	data, err := os.ReadFile(filepath)

	if err != nil {
		fmt.Printf(location+"Failed to read file\n%v\n", err)
		return nil, err
	}

	var state types.State

	err = json.Unmarshal(data, &state)

	return &state, nil
}

func (r *StateRepository) Save(newState types.State) error {
	location := "StateRepository.Save()\n"
	filepath, err := util.GetPathToFile("/server/state", r.filename)

	if err != nil {
		fmt.Printf(location+"Failed to get path to file\n%v\n", err)
		return err
	}

	data, err := json.Marshal(newState)

	if err != nil {
		fmt.Printf(location + "Failed to marshal state into []byte")
		return err
	}

	return os.WriteFile(filepath, data, 0666)
}
