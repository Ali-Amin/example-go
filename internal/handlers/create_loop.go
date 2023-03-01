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
package handlers

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/KarimElghamry/alvarium-sdk-go/pkg/config"
	"github.com/KarimElghamry/alvarium-sdk-go/pkg/contracts"
	"github.com/KarimElghamry/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/example-go/internal/models"
	logInterface "github.com/project-alvarium/provider-logging/pkg/interfaces"
	"github.com/project-alvarium/provider-logging/pkg/logging"
)

type CreateLoop struct {
	cfg       config.SdkInfo
	chPublish chan []byte
	logger    logInterface.Logger
	sdk       interfaces.Sdk
}

func NewCreateLoop(sdk interfaces.Sdk, ch chan []byte, cfg config.SdkInfo, logger logInterface.Logger) CreateLoop {
	return CreateLoop{
		cfg:       cfg,
		chPublish: ch,
		logger:    logger,
		sdk:       sdk,
	}
}

func (c *CreateLoop) BootstrapHandler(ctx context.Context, wg *sync.WaitGroup) bool {
	cancelled := false
	wg.Add(1)
	go func() {
		defer wg.Done()

		for !cancelled {
			data, err := models.NewSampleData(c.cfg.Signature.PrivateKey)
			if err != nil {
				c.logger.Error(err.Error())
				continue
			}
			b, _ := json.Marshal(data)

			deviceId := os.Getenv("DEVICEID")
			c.sdk.Create(context.WithValue(context.Background(), contracts.DeviceIdKey, deviceId), b)
			c.chPublish <- b
			time.Sleep(1 * time.Second)
		}
		close(c.chPublish)
		c.logger.Write(logging.DebugLevel, "cancel received")
	}()

	wg.Add(1)
	go func() { // Graceful shutdown
		defer wg.Done()

		<-ctx.Done()
		c.logger.Write(logging.InfoLevel, "shutdown received")
		cancelled = true
	}()
	return true
}
