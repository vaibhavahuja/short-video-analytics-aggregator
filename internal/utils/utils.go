package utils

import (
	"encoding/binary"
	"encoding/json"
	"github.com/rs/zerolog/log"
)

func MarshalJsonToStruct(in interface{}, obj interface{}) {
	bytes, err := json.Marshal(in)
	if err != nil {
		log.Debug().Err(err).Msg("error in MarshalJsonToStruct, json.Marshal failed")
	}
	if len(bytes) == 0 {
		return
	}

	// skipping error intentionally
	err = json.Unmarshal(bytes, obj)
	if err != nil {
		log.Debug().Str("payload", string(bytes)).Err(err).Msg("error in MarshalJsonToStruct, json.Unmarshal failed")
	}
}

func IntToByte(value int) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(value))
	return b
}
