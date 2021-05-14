/*
	Copyright 2019 NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package zitilib_runlevel_5_operation

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/fablab/kernel/model"
	"github.com/openziti/fablab/zitilib/cli"
	"github.com/openziti/fabric/pb/mgmt_pb"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"time"
)

func NewClientMetrics(service string, closer <-chan struct{}) *ClientMetrics {
	return NewClientMetricsWithIdMapper(service, closer, func(id string) string {
		return "#" + id
	})
}

func NewClientMetricsWithIdMapper(service string, closer <-chan struct{}, f func(string) string) *ClientMetrics {
	return &ClientMetrics{
		service:            service,
		closer:             closer,
		idToSelectorMapper: f,
	}
}

type ClientMetrics struct {
	service            string
	listener           net.Listener
	closer             <-chan struct{}
	model              *model.Model
	idToSelectorMapper func(string) string
}

func (metrics *ClientMetrics) Activate(run model.Run) error {
	if _, err := cli.Exec(run.GetModel(), "edge", "delete", "identity", "metrics-host"); err != nil {
		return err
	}

	jwtFilePath := run.GetLabel().GetFilePath("metrics-host.jwt")
	if _, err := cli.Exec(run.GetModel(), "edge", "create", "identity", "service", "metrics-host", "-a", "metrics-host", "-o", jwtFilePath); err != nil {
		return err
	}

	identityConfigPath := run.GetLabel().GetFilePath("metrics-host.json")
	if _, err := cli.Exec(run.GetModel(), "edge", "enroll", jwtFilePath, "-o", identityConfigPath); err != nil {
		return err
	}

	return nil
}

func (metrics *ClientMetrics) Operate(run model.Run) error {
	identityConfigPath := run.GetLabel().GetFilePath("metrics-host.json")

	sdkConfig, err := config.NewFromFile(identityConfigPath)
	if err != nil {
		return err
	}

	metrics.model = run.GetModel()

	context := ziti.NewContextWithConfig(sdkConfig)
	listener, err := context.Listen(metrics.service)
	if err != nil {
		return err
	}

	metrics.listener = listener

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				pfxlog.Logger().WithError(err).Info("metrics listener closed, returning")
				return
			}
			go metrics.HandleMetricsConn(conn)
		}
	}()

	go metrics.runMetrics()

	return nil
}

func (metrics *ClientMetrics) HandleMetricsConn(conn net.Conn) {
	defer func() { _ = conn.Close() }()

	log := pfxlog.Logger()
	log.Infof("new client metrics connection established from: %v", conn.RemoteAddr().String())
	lenBuf := make([]byte, 4)
	msgBuf := make([]byte, 4*1024)
	for {
		if _, err := io.ReadFull(conn, lenBuf); err != nil {
			log.WithError(err).Info("metrics conn closed, exiting read loop")
			return
		}
		msgLen := int(binary.LittleEndian.Uint32(lenBuf))
		if msgLen > 1024*16 {
			log.Errorf("got invalid metrics message len: %v, closing connection", msgLen)
			return
		}

		if msgLen > len(msgBuf) {
			msgBuf = make([]byte, msgLen)
		}

		if _, err := io.ReadFull(conn, msgBuf[:msgLen]); err != nil {
			pfxlog.Logger().WithError(err).Info("metrics conn closed, exiting read loop")
			return
		}

		event := &mgmt_pb.StreamMetricsEvent{}
		err := proto.Unmarshal(msgBuf[:msgLen], event)
		if err != nil {
			logrus.Errorf("error handling metrics receive (%w), exiting", err)
			return
		}

		hostSelector := metrics.idToSelectorMapper(event.SourceId)
		host, err := metrics.model.SelectHost(hostSelector)
		if err == nil {
			modelEvent := metrics.toClientMetricsEvent(event)
			metrics.model.AcceptHostMetrics(host, modelEvent)
			logrus.Infof("<$= [%s]", event.SourceId)
		} else {
			logrus.Errorf("unable to find host (%v)", err)
		}
	}
}

func (metrics *ClientMetrics) runMetrics() {
	logrus.Infof("starting")
	defer logrus.Infof("exiting")

	<-metrics.closer
	_ = metrics.listener.Close()
}

func (metrics *ClientMetrics) toClientMetricsEvent(fabricEvent *mgmt_pb.StreamMetricsEvent) *model.MetricsEvent {
	modelEvent := &model.MetricsEvent{
		Timestamp: time.Unix(fabricEvent.Timestamp.Seconds, int64(fabricEvent.Timestamp.Nanos)),
		Metrics:   model.MetricSet{},
	}

	for name, val := range fabricEvent.IntMetrics {
		group := fabricEvent.MetricGroup[name]
		modelEvent.Metrics.AddGroupedMetric(group, name, val)
	}

	for name, val := range fabricEvent.FloatMetrics {
		group := fabricEvent.MetricGroup[name]
		modelEvent.Metrics.AddGroupedMetric(group, name, val)
	}

	return modelEvent
}
