// Copyright 2021 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package netpoll

import (
	"context"
	"time"
)

// SetNumLoops is used to set the number of pollers, generally do not need to actively set.
// By default, the number of pollers is equal to runtime.GOMAXPROCS(0)/20+1.
// If the number of cores in your service process is less than 20c, theoretically only one poller is needed.
// Otherwise you may need to adjust the number of pollers to achieve the best results.
// Experience recommends assigning a poller every 20c.
func SetNumLoops(numLoops int) error {
	return setNumLoops(numLoops)
}

// LoadBalance sets the load balancing method. Load balancing is always a best effort to attempt
// to distribute the incoming connections between multiple polls.
// This option only works when NumLoops is set.
func SetLoadBalance(lb LoadBalance) error {
	return setLoadBalance(lb)
}

// DisableGopool will remove gopool(the goroutine pool used to run OnRequest),
// which means that OnRequest will be run via `go OnRequest(...)`.
// Usually, OnRequest will cause stack expansion, which can be solved by reusing goroutine.
// But if you can confirm that the OnRequest will not cause stack expansion,
// it is recommended to use DisableGopool to reduce redundancy and improve performance.
func DisableGopool() error {
	return disableGopool()
}

// WithOnPrepare registers the OnPrepare method to EventLoop.
func WithOnPrepare(onPrepare OnPrepare) Option {
	return Option{func(op *options) {
		op.onPrepare = onPrepare
	}}
}

// WithReadTimeout sets the read timeout of connections.
func WithReadTimeout(timeout time.Duration) Option {
	return Option{func(op *options) {
		op.readTimeout = timeout
	}}
}

// WithIdleTimeout sets the idle timeout of connections.
func WithIdleTimeout(timeout time.Duration) Option {
	return Option{func(op *options) {
		op.idleTimeout = timeout
	}}
}

// Option .
type Option struct {
	f func(*options)
}

type options struct {
	onPrepare   OnPrepare
	readTimeout time.Duration
	idleTimeout time.Duration
}

func (opt *options) prepare(onRequest OnRequest) OnPrepare {
	return func(connection Connection) context.Context {
		connection.SetOnRequest(onRequest)
		connection.SetReadTimeout(opt.readTimeout)
		connection.SetIdleTimeout(opt.idleTimeout)
		if opt.onPrepare != nil {
			return opt.onPrepare(connection)
		}
		return context.Background()
	}
}
