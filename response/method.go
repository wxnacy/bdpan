package response

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/wxnacy/bdpan"
)

func ToInterface[T any](r *http.Response) (*T, error) {
	var i T
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	bdpan.Log.Infof("%s code %d\n", r.Request.URL.Path, r.StatusCode)
	bdpan.Log.Infof("%s response %s\n", r.Request.URL.Path, string(bodyBytes))
	if r.StatusCode == 200 {
		if err := json.Unmarshal(bodyBytes, &i); err != nil {
			return nil, err
		}
	} else {
		var apiErr ApiError
		if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
			return nil, err
		}
		return nil, &apiErr

	}
	return &i, nil
}
