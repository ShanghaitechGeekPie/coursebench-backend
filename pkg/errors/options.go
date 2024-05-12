// Copyright (C) 2021-2024 ShanghaiTech GeekPie
// This file is part of CourseBench Backend.
//
// CourseBench Backend is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// CourseBench Backend is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with CourseBench Backend.  If not, see <http://www.gnu.org/licenses/>.

package errors

type Option interface {
	LogLevel() ErrorLevel
	HideErrorCode() bool
}

type optionImpl struct {
	logLevel      ErrorLevel
	hideErrorCode bool
}

func (o optionImpl) LogLevel() ErrorLevel {
	return o.logLevel
}

func (o optionImpl) HideErrorCode() bool {
	return o.hideErrorCode
}

type OptionBuilder struct {
	opt optionImpl
}

func NewOptionBuilder() *OptionBuilder {
	return &OptionBuilder{
		opt: optionImpl{},
	}
}

func (ob *OptionBuilder) SetLogLevel(level ErrorLevel) *OptionBuilder {
	ob.opt.logLevel = level
	return ob
}

func (ob *OptionBuilder) SetHideErrorCode() *OptionBuilder {
	ob.opt.hideErrorCode = true
	return ob
}

func (ob *OptionBuilder) Build() Option {
	return ob.opt
}
