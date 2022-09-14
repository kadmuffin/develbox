package develbox

import (
	"encoding/json"
)

func parse(bytes []byte) DevSetings {
	var configs DevSetings
	json.Unmarshal(bytes, &configs)

	return configs
}
