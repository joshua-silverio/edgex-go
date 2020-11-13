// +build linux

//
// Copyright (c) 2019 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
// in compliance with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under
// the License.
//
// SPDX-License-Identifier: Apache-2.0'
//

package secretstore

import (
	"context"
	"os"
	"testing"

	"github.com/edgexfoundry/edgex-go/internal/security/secretstoreclient"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"

	"github.com/stretchr/testify/assert"
)

// TestCreatesFile only runs on Linux and makes sure the code can
// run a real executable taking real arguments.
func TestCreatesFile(t *testing.T) {
	const testfile = "/tmp/tokenprovider_linux_test.dat"
	config := secretstoreclient.SecretServiceInfo{
		TokenProvider:     "touch",
		TokenProviderType: OneShotProvider,
		TokenProviderArgs: []string{testfile},
	}
	ctx, cancel := context.WithCancel(context.Background())

	err := os.RemoveAll(testfile)
	defer os.RemoveAll(testfile) // cleanup

	p := NewTokenProvider(ctx, logger.MockLogger{}, NewDefaultExecRunner())
	p.SetConfiguration(config)
	assert.NoError(t, err)

	p.Launch()
	defer cancel()

	file, err := os.Open(testfile)
	defer file.Close()
	assert.NoError(t, err) // fails if file wasn't created
}