module github.com/devil-dwj/wms/contrib/log/zap

go 1.16

require (
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/devil-dwj/wms v0.0.0-20220801062859-222b280d665d
	github.com/natefinch/lumberjack v2.0.0+incompatible
	go.uber.org/zap v1.21.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)

replace github.com/devil-dwj/wms => ../../../
