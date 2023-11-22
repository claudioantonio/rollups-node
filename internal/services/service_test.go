// (c) Cartesi and individual authors (see AUTHORS)
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/cartesi/rollups-node/internal/logger"
	"github.com/stretchr/testify/suite"
)

type ServicesTestSuite struct {
	suite.Suite
	tmpDir      string
	servicePort int
}

func (s *ServicesTestSuite) SetupSuite() {
	logger.Init("warning", false)
	s.buildFakeService()
	s.servicePort = 55555
}

func (s *ServicesTestSuite) TearDownSuite() {
	err := os.RemoveAll(s.tmpDir)
	if err != nil {
		panic(err)
	}
}

func (s *ServicesTestSuite) SetupTest() {
	s.servicePort++
	serviceAdress := "0.0.0.0:" + fmt.Sprint(s.servicePort)
	os.Setenv("SERVICE_ADDRESS", serviceAdress)
}

// Service should stop when context is cancelled
func (s *ServicesTestSuite) TestServiceStops() {
	service := Service{
		name:            "fake-service",
		binaryName:      "fake-service",
		healthcheckPort: fmt.Sprint(s.servicePort),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start service in goroutine
	result := make(chan error)
	go func() {
		result <- service.Start(ctx)
	}()

	time.Sleep(100 * time.Millisecond)

	// shutdown
	cancel()
	err := <-result
	s.Require().Nil(err, "service exited for the wrong reason: %v", err)
}

// Service should stop if timeout is reached and it isn't ready yet
func (s *ServicesTestSuite) TestServiceTimeout() {
	service := Service{
		name:            "fake-service",
		binaryName:      "fake-service",
		healthcheckPort: "0000", // wrong port
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start service in goroutine
	result := make(chan error, 1)
	go func() {
		result <- service.Start(ctx)
	}()

	// expect timeout because of wrong port
	err := service.Ready(ctx, 500*time.Millisecond)
	s.NotNil(err, "expected service to timeout")

	// shutdown
	cancel()
	s.Nil(<-result, "service exited for the wrong reason: %v", err)
}

// Service should be ready soon after starting
func (s *ServicesTestSuite) TestServiceReady() {
	service := Service{
		name:            "fake-service",
		binaryName:      "fake-service",
		healthcheckPort: fmt.Sprint(s.servicePort),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start service in goroutine
	result := make(chan error)
	go func() {
		result <- service.Start(ctx)
	}()

	// wait for service to be ready
	err := service.Ready(ctx, 500*time.Millisecond)
	s.Nil(err, "service timed out")

	// shutdown
	cancel()
	s.Nil(<-result, "service exited for the wrong reason: %v", err)
}

// Builds the fake-service binary and adds it to PATH
func (s *ServicesTestSuite) buildFakeService() {
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		panic(err)
	}

	s.tmpDir, err = os.MkdirTemp("", "")
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(
		"go",
		"build",
		"-o",
		filepath.Join(s.tmpDir, "fake-service"),
		"internal/services/fakeservice/main.go",
	)
	cmd.Dir = rootDir
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	os.Setenv("PATH", os.Getenv("PATH")+":"+s.tmpDir)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServicesTestSuite))
}
