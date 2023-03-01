/*******************************************************************************
 * Copyright 2021 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package models

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/KarimElghamry/alvarium-sdk-go/pkg/config"
	"io/ioutil"
	"math/rand"
	"time"
)

const alphanumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type SampleData struct {
	Description string `json:"description,omitempty"`
	Seed        string `json:"seed,omitempty"`
	Signature   string `json:"signature,omitempty"`
}

func NewSampleData(cfg config.KeyInfo) (SampleData, error) {
	key, err := ioutil.ReadFile(cfg.Path)
	if err != nil {
		return SampleData{}, err
	}

	x := SampleData{
		Description: factoryRandomFixedLengthString(128, alphanumericCharset),
		Seed:        factoryRandomFixedLengthString(64, alphanumericCharset),
	}

	keyDecoded := make([]byte, hex.DecodedLen(len(key)))
	hex.Decode(keyDecoded, key)
	signed := ed25519.Sign(keyDecoded, []byte(x.Seed))
	x.Signature = fmt.Sprintf("%x", signed)
	return x, nil
}

func factoryRandomFixedLengthString(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
