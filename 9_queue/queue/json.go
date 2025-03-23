package queue

import (
	"encoding/json"
	"fmt"
)

func Decode(raw json.RawMessage) (ret any, err error) {
	var meta struct {
		Request RequestType `json:"request"`
	}

	err = json.Unmarshal(raw, &meta)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshall request: %w", err)
	}

	switch meta.Request {
	case RequestGet:
		var resp GetRequest
		err = json.Unmarshal(raw, &resp)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshall get request: %w", err)
		}
		ret = &resp
	case RequestPut:
		var resp PutRequest
		err = json.Unmarshal(raw, &resp)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshall put request: %w", err)
		}
		ret = &resp
	case RequestAbort:
		var resp AbortRequest
		err = json.Unmarshal(raw, &resp)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshall abort request: %w", err)
		}
		ret = &resp
	case RequestDelete:
		var resp DeleteRequest
		err = json.Unmarshal(raw, &resp)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshall delete request: %w", err)
		}
		ret = &resp
	default:
		return nil, fmt.Errorf("can't determine request type: %v", meta)
	}

	return ret, nil
}
