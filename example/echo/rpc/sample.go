package service

import (
	"github.com/ronaksoft/rony/config"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/tools"
	"github.com/spf13/cobra"
)

/*
   Creation Time: 2021 - Jul - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

//go:generate protoc -I=. -I=../../.. --go_out=paths=source_relative:. sample.proto
//go:generate protoc -I=. -I=../../.. --gorony_out=paths=source_relative,option=no_edge_dep:. sample.proto
func init() {}

// Sample implements auto-generated service.ISample interface
type Sample struct{}

func (s *Sample) Echo(ctx *edge.RequestCtx, req *EchoRequest, res *EchoResponse) {
	res.ReqID = req.ID
	res.RandomText = tools.RandomID(32)
}

// SampleCli implements service.ISampleCli auto-generated interface
type SampleCli struct{}

func (s *SampleCli) Echo(cli *SampleClient, cmd *cobra.Command, args []string) error {
	req := &EchoRequest{
		ID:            config.GetInt64("id"),
	}
	res, err := cli.Echo(req)
	if err != nil {
		cmd.Println("Receiver Error:", err.Error())
		return err
	}
	cmd.Println("EchoResponse:", res.ReqID, res.RandomText)
	return nil
}
