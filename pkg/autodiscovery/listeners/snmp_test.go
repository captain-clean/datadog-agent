// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2020 Datadog, Inc.

package listeners

import (
	"testing"

	"github.com/DataDog/datadog-agent/pkg/autodiscovery/integration"

	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestSNMPListener(t *testing.T) {
	newSvc := make(chan Service, 10)
	delSvc := make(chan Service, 10)
	testChan := make(chan snmpJob, 10)

	snmpConfig := util.SNMPConfig{
		Network:   "192.168.0.0/24",
		Community: "public",
	}
	listenerConfig := util.SNMPListenerConfig{
		Configs: []util.SNMPConfig{snmpConfig},
		Workers: 1,
	}

	mockConfig := config.Mock()
	mockConfig.Set("snmp_listener", listenerConfig)

	worker = func(l *SNMPListener, jobs <-chan snmpJob) {
		for {
			job := <-jobs
			testChan <- job
		}
	}

	l, err := NewSNMPListener()
	assert.Equal(t, nil, err)
	l.Listen(newSvc, delSvc)

	job := <-testChan

	assert.Equal(t, "snmp", job.subnet.adIdentifier)
	assert.Equal(t, "192.168.0.0", job.currentIP.String())
	assert.Equal(t, "192.168.0.0", job.subnet.startingIP.String())
	assert.Equal(t, "192.168.0.0/24", job.subnet.network.String())
	assert.Equal(t, "public", job.subnet.config.Community)
	assert.Equal(t, "public", job.subnet.defaultParams.Community)

	job = <-testChan
	assert.Equal(t, "192.168.0.1", job.currentIP.String())
	assert.Equal(t, "192.168.0.0", job.subnet.startingIP.String())
}

func TestSNMPListenerIgnoredAdresses(t *testing.T) {
	newSvc := make(chan Service, 10)
	delSvc := make(chan Service, 10)
	testChan := make(chan snmpJob, 10)

	snmpConfig := util.SNMPConfig{
		Network:            "192.168.0.0/24",
		Community:          "public",
		IgnoredIPAddresses: []string{"192.168.0.0"},
	}
	listenerConfig := util.SNMPListenerConfig{
		Configs: []util.SNMPConfig{snmpConfig},
		Workers: 1,
	}

	mockConfig := config.Mock()
	mockConfig.Set("snmp_listener", listenerConfig)

	worker = func(l *SNMPListener, jobs <-chan snmpJob) {
		for {
			job := <-jobs
			testChan <- job
		}
	}

	l, err := NewSNMPListener()
	assert.Equal(t, nil, err)
	l.Listen(newSvc, delSvc)

	job := <-testChan

	assert.Equal(t, "snmp", job.subnet.adIdentifier)
	assert.Equal(t, "192.168.0.1", job.currentIP.String())
	assert.Equal(t, "192.168.0.0", job.subnet.startingIP.String())

	job = <-testChan
	assert.Equal(t, "192.168.0.2", job.currentIP.String())
	assert.Equal(t, "192.168.0.0", job.subnet.startingIP.String())
}

func TestSNMPInfo(t *testing.T) {
	snmpConfig := util.SNMPConfig{
		Network:   "192.168.0.0/24",
		Community: "public",
		Timeout:   5,
		Retries:   2,
	}

	svc := SNMPService{
		adIdentifier: "snmp",
		entityID:     "id",
		deviceIP:     "192.168.0.1",
		creationTime: integration.Before,
		config:       snmpConfig,
	}

	info, err := svc.GetSNMPInfo("community")
	assert.Equal(t, nil, err)
	assert.Equal(t, "public", info)

	info, err = svc.GetSNMPInfo("timeout")
	assert.Equal(t, nil, err)
	assert.Equal(t, "5", info)

	info, err = svc.GetSNMPInfo("retries")
	assert.Equal(t, nil, err)
	assert.Equal(t, "2", info)
}

func TestSNMPInfov3(t *testing.T) {
	snmpConfig := util.SNMPConfig{
		Network:      "192.168.0.0/24",
		User:         "admin",
		AuthKey:      "secret",
		AuthProtocol: "SHA",
		PrivKey:      "private",
		PrivProtocol: "DES",
	}

	svc := SNMPService{
		adIdentifier: "snmp",
		entityID:     "id",
		deviceIP:     "192.168.0.1",
		creationTime: integration.Before,
		config:       snmpConfig,
	}

	info, err := svc.GetSNMPInfo("user")
	assert.Equal(t, nil, err)
	assert.Equal(t, "admin", info)

	info, err = svc.GetSNMPInfo("auth_key")
	assert.Equal(t, nil, err)
	assert.Equal(t, "secret", info)

	info, err = svc.GetSNMPInfo("auth_protocol")
	assert.Equal(t, nil, err)
	assert.Equal(t, "usmHMACSHAAuthProtocol", info)

	info, err = svc.GetSNMPInfo("priv_key")
	assert.Equal(t, nil, err)
	assert.Equal(t, "private", info)

	info, err = svc.GetSNMPInfo("priv_protocol")
	assert.Equal(t, nil, err)
	assert.Equal(t, "usmDESPrivProtocol", info)
}
