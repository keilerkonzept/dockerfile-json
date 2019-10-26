package dockerfile

import (
	"encoding/json"
	"fmt"
)

func (c *Command) MarshalJSON() ([]byte, error) {
	rawJSON, err := json.Marshal(c.Command)
	if err != nil {
		return nil, fmt.Errorf("merge json fields: %v", err)
	}
	out := map[string]interface{}{
		"Name": c.Name,
	}
	if err := json.Unmarshal(rawJSON, &out); err != nil {
		return nil, fmt.Errorf("merge json fields: %v", err)
	}
	return json.Marshal(out)
}
