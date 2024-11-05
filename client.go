package bdpan

import (
	sdk "github.com/wxnacy/bdpan/openapi"
)

func init() {
	initClient()
}

var (
	_client *sdk.APIClient
)

func initClient() {
	configuration := sdk.NewConfiguration()
	_client = sdk.NewAPIClient(configuration)
}

func GetClient() *sdk.APIClient {
	return _client
}
