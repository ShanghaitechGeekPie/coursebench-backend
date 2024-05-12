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

package config

import (
	"github.com/spf13/viper"
	"time"
)

type FiberConfigType struct {
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	InDevelopment  bool          `mapstructure:"in_development"`
	Listen         string        `mapstructure:"listen"`
	UseXForwardFor bool          `mapstructure:"use_x_forward_for"`
}

var FiberConfig FiberConfigType

func SetupFiberConfig() {
	cfg := viper.Sub("fiber")
	if cfg == nil {
		cfg = viper.New()
	}
	cfg.SetDefault("read_timeout", "10s")
	cfg.SetDefault("write_timeout", "10s")
	cfg.SetDefault("idle_timeout", "1m")
	cfg.SetDefault("in_development", GlobalConf.InDevelopment)
	cfg.SetDefault("listen", "0.0.0.0:10001")
	jsonErr := cfg.Unmarshal(&FiberConfig)
	if jsonErr != nil {
		panic(jsonErr)
	}
}
