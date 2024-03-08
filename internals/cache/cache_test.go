package cache

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	value, _ := json.Marshal("value")

	store.Add("key", value)

	data, ok := store.Get("key")

	require.Equal(t, ok, true)

	var retrieved string

	err := json.Unmarshal(data, &retrieved)

	require.NotEmpty(t, data)

	require.NoError(t, err)

	require.Equal(t, "value", retrieved)

}

func TestAddConcurrent(t *testing.T) {

	type data struct {
		Key string `json:"key"`
		Val string `json:"val"`
	}

	slice := []data{{"key1", "val1"}, {"key2", "val2"}, {"key3", "val3"}, {"key4", "val4"}, {"key5", "val5"}}

	channel := make(chan data)

	for i := 0; i < len(slice); i++ {
		obj := slice[i]
		value, _ := json.Marshal(obj.Val)
		go func() {
			store.Add(obj.Key, value)
			channel <- obj
		}()
	}

	for i := 0; i < len(slice); i++ {
		obj := <-channel
		value, ok := store.Get(obj.Key)
		require.NotEmpty(t, value)
		require.Equal(t, ok, true)

		var retrieved string
		err := json.Unmarshal(value, &retrieved)

		require.NoError(t, err)
		require.NotEmpty(t, retrieved)
		require.Equal(t, obj.Val, retrieved)
	}

}
