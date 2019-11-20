// Copyright 2019 The Dice Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package log provides utility for logging into files and the console.
package log

import (
	"github.com/sirupsen/logrus"
	"io"
)

// Level is the logging level which decides if a value will be logged or not.
type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warning"
	ErrorLevel Level = "error"
)

// Logger prescribes methods for logging to any io.Writer with different priorities.
// It also provides corresponding formatting methods.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// NewLogger creates a new instance of a Logger implementation that will use a
// io.Writer (such as stdout or a file) for writing the logs.
func NewLogger(level Level, output io.Writer) Logger {
	l := logrus.New()
	l.SetOutput(output)

	return l
}
