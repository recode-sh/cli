module github.com/recode-sh/cli

go 1.17

replace github.com/recode-sh/aws-cloud-provider v0.0.0 => ../aws-cloud-provider

replace github.com/recode-sh/recode v0.0.0 => ../recode

replace github.com/recode-sh/agent v0.0.0 => ../agent

require (
	github.com/aws/aws-sdk-go-v2/config v1.13.1
	github.com/briandowns/spinner v1.18.1
	github.com/golang/mock v1.6.0
	github.com/google/go-github/v43 v43.0.0
	github.com/google/wire v0.5.0
	github.com/jwalton/gchalk v1.3.0
	github.com/kevinburke/ssh_config v1.1.0
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8
	github.com/recode-sh/agent v0.0.0
	github.com/recode-sh/aws-cloud-provider v0.0.0
	github.com/recode-sh/recode v0.0.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	golang.org/x/crypto v0.0.0-20220313003712-b769efc7c000
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
	google.golang.org/grpc v1.46.2
)

require (
	github.com/aws/aws-sdk-go-v2 v1.15.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.6.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.10.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.6 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.13.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.11.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.29.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.5.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.9.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.14.0 // indirect
	github.com/aws/smithy-go v1.11.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gosimple/slug v1.12.0 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jsonmaur/aws-regions/v2 v2.3.1 // indirect
	github.com/jwalton/go-supportscolor v1.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/whilp/git-urls v1.0.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220224120231-95c6836cb0e7 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
