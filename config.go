package main

import (
	"time"

	"github.com/spf13/viper"
)

func init() {
	// Set the configuration options.
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Set the default path locations for the configuration file. These paths
	// are searched in the order that they are added.
	viper.AddConfigPath("/etc/cyclops")
	viper.AddConfigPath("$HOME/.cyclops")
	viper.AddConfigPath(".")

	// Set defaults for CaptureConfig
	viper.SetDefault("capture",
		map[string]interface{}{
			"device":   "/dev/video0",
			"height":   240,
			"width":    320,
			"interval": time.Millisecond * 100,
		},
	)

	// Set defaults for ThresholdConfig
	viper.SetDefault("threshold",
		map[string]interface{}{
			"min_hue": 0,
			"max_hue": 180,
			"min_sat": 0,
			"max_sat": 255,
			"min_val": 0,
			"max_val": 255,
		},
	)

	// Set defaults for MorphConfig
	viper.SetDefault("morph",
		map[string]interface{}{
			"blur":    0,
			"erosion": 5,
		},
	)

	// Set defaults for TransformConfig
	viper.SetDefault("transform",
		map[string]interface{}{
			"dp":         1.0,
			"min_dist":   1.0,
			"param_1":    1.0,
			"param_2":    1.0,
			"min_radius": 0.0,
			"max_radius": 0.0,
		},
	)
}

// Config is the top level configuration element.
type Config struct {
	Capture   CaptureConfig   `mapstructure:"capture"`
	Threshold ThresholdConfig `mapstructure:"threshold"`
	Morph     MorphConfig     `mapstructure:"morph"`
	Transform TransformConfig `mapstructure:"transform"`
}

// CaptureConfig defines all of the capture related configuration options.
type CaptureConfig struct {
	Device   string        `mapstructure:"device"`
	Height   int           `mapstructure:"height"`
	Width    int           `mapstructure:"width"`
	Interval time.Duration `mapstructure:"interval"`
}

// ThresholdConfig defines all of the thresholding related configuration
// options.
type ThresholdConfig struct {
	MinHue        float64 `mapstructure:"min_hue:"`
	MaxHue        float64 `mapstructure:"max_hue"`
	MinSaturation float64 `mapstructure:"min_sat"`
	MaxSaturation float64 `mapstructure:"max_sat"`
	MinValue      float64 `mapstructure:"min_val"`
	MaxValue      float64 `mapstructure:"max_val"`
}

// MorphConfig defines all of the morphology related configuration options that
// are applied to the thresholding image before the object detection process is
// performed.
type MorphConfig struct {
	Blur    int `mapstructure:"blur"`
	Erosion int `mapstructure:"erosion"`
}

// TransformConfig defines all of the configuration options that are required
// for the object detection process.
type TransformConfig struct {
	DP        float64 `mapstructure:"dp"`
	MinDist   float64 `mapstructure:"min_dist"`
	Param1    float64 `mapstructure:"param_1"`
	Param2    float64 `mapstructure:"param_2"`
	MinRadius int     `mapstructure:"min_radius"`
	MaxRadius int     `mapstructure:"max_radius"`
}

// LoadConfig attempts to the load the configuration file that is specified at
// the given path. If the path is empty, then it will try to load the
// configuration file from one of the default locations.
func LoadConfig(configPath string) (*Config, error) {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	c := &Config{}

	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}

	return c, nil
}
