package types

import "encoding/json"

type ZephyrUpdate struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func (z ZephyrUpdate) MarshalBinary() ([]byte, error) {
	return json.Marshal(z)
}

func (z *ZephyrUpdate) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, z)
}
