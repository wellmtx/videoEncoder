package utils_test

import (
	"testing"
	"videoEncoder/framework/utils"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	json := `{
		"id": "4256435423-4324324312",
		"file_path": "convite.mp4",
		"status": "pending"
	}`

	err := utils.IsJson(json)

	require.Nil(t, err)

	json = `dfndfjondsfojdsfnjo{
		"id": "4256435423-4324324312",
		"file_path": "convite.mp4",
		"status": "pending"
	}`

	err = utils.IsJson(json)
	require.Error(t, err)
}
