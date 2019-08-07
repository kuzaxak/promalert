## What is PromAlert?
Yet another [Prometheus](https://prometheus.io/) Alertmanager webhook processor inspired by [qvl/promplot](https://github.com/qvl/promplot)

## How it works 
Receive webhook from Alertmanager, draw images from alert expression, upload pictures to S3 bucket, generate public links, send a notification to Slack. 

## Installation

### Helm chart

> Coming soon

### Docker 

Use [env file](https://docs.docker.com/engine/reference/commandline/run/#set-environment-variables--e---env---env-file) to store required params
```bash
docker run -p 8080:8080 --env-file env.list kuzaxak/promalert
```

## Configuration
The following tables list the configurable parameters of the PromAlert and their default values.

You can use a YAML configuration file or env variable. Package [viper](https://github.com/spf13/viper) used to parse them.
Environment variables prefix: `PROMALERT_`

*Required params:*

| Parameter        | Description           | Env variable               |
|:-----------------|:----------------------|:---------------------------|
| `slack_token`    | OAuth bot token       | `PROMALERT_SLACK_TOKEN`    |
| `slack_channel`  | Slack channel to send | `PROMALERT_SLACK_CHANNEL`  |
| `prometheus_url` | Prometheus URL        | `PROMALERT_PROMETHEUS_URL` |
| `s3_bucket`      | S3 bucket name        | `PROMALERT_S3_BUCKET`      |
| `s3_region`      | S3 region             | `PROMALERT_S3_REGION`      |

*Additional params:*

| Parameter           | Description                                  | Default                                          |
|:--------------------|:---------------------------------------------|:-------------------------------------------------|
| `http_port`         | HTTP port                                    | `8080`                                           |
| `metric_resolution` | Amount of point on the graph                 | `100`                                            |
| `debug`             | Verbose log output. Dump HTTP request to log | `false`                                          |
| `message_template`  | Slack message template. Go template syntax   | [`config.example.yaml`](config.example.yaml#L11) |

### AWS
AWS credentials parsed by [aws-go-client](https://github.com/aws/aws-sdk-go) in the following [order](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html):
1. Environment variables.
1. Shared credentials file.
1. If your application is running on an Amazon EC2 instance, IAM role for Amazon EC2.

### Message templating

> Coming soon

## Build 
```
docker build --rm -t promalert .
```

## Bug Reporting

To submit a bug report use the GitHub bug tracker for the project:

[Github Issues](https://github.com/kuzaxak/promalert/issues)

## License

GNU Lesser General Public License v3.0

See [LICENSE](LICENSE) to see the full text.