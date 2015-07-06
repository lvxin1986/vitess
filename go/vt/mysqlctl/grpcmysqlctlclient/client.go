// Copyright 2014, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package grpcmysqlctlclient contains the gRPC1 version of the mysqlctl
// client protocol.
package grpcmysqlctlclient

import (
	"net"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/youtube/vitess/go/vt/mysqlctl/mysqlctlclient"

	pb "github.com/youtube/vitess/go/vt/proto/mysqlctl"
)

type client struct {
	cc *grpc.ClientConn
	c  pb.MysqlCtlClient
}

func factory(network, addr string, dialTimeout time.Duration) (mysqlctlclient.MysqlctlClient, error) {
	// create the RPC client
	cc, err := grpc.Dial(addr, grpc.WithBlock(), grpc.WithTimeout(dialTimeout), grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout(network, addr, timeout)
	}))
	if err != nil {
		return nil, err
	}
	c := pb.NewMysqlCtlClient(cc)

	return &client{
		cc: cc,
		c:  c,
	}, nil
}

// Start is part of the MysqlctlClient interface.
func (c *client) Start(ctx context.Context) error {
	_, err := c.c.Start(ctx, &pb.StartRequest{})
	return err
}

// Shutdown is part of the MysqlctlClient interface.
func (c *client) Shutdown(ctx context.Context, waitForMysqld bool) error {
	_, err := c.c.Shutdown(ctx, &pb.ShutdownRequest{
		WaitForMysqld: waitForMysqld,
	})
	return err
}

// RunMysqlUpgrade is part of the MysqlctlClient interface.
func (c *client) RunMysqlUpgrade(ctx context.Context) error {
	_, err := c.c.RunMysqlUpgrade(ctx, &pb.RunMysqlUpgradeRequest{})
	return err
}

// Close is part of the MysqlctlClient interface.
func (c *client) Close() {
	c.cc.Close()
}

func init() {
	mysqlctlclient.RegisterFactory("grpc", factory)
}
