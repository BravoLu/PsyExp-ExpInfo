module github.com/grpc_exp_info

go 1.21.5

replace github.com/BravoLu/grpc_idl => ./api/exp_info

replace github.com/grpc_plugins => ../plugins

require (
	github.com/BravoLu/grpc_idl v0.0.0-00010101000000-000000000000
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc_plugins v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.60.0
	gorm.io/datatypes v1.2.0
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.5
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/prometheus/client_golang v1.17.0 // indirect
	github.com/prometheus/client_model v0.4.1-0.20230718164431-9a2bf3000d16 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/net v0.16.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
