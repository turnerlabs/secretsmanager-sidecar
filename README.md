secretsmanager-sidecar
=======================

A simple command line program designed to be run as a sidecar container that writes a secret from [AWS SecretsManager](https://aws.amazon.com/secrets-manager/) to a file.

### Example usage

```bash
# name of secrets manager secret
export SECRET_ID=my-secret

# where the secret should be written
export SECRET_FILE=/var/secret/my-secret

# run the program
./secretsmanager-sidecar
```

In this example, the program writes the contents of the decrypted secrets manager secret `my-secret` to `/var/secret/my-secret` and then exits.

### AWS Access

This program uses the [AWS SDK for Go][go-sdk] which looks for credentials in the following locations:

1. [Environment Variables][go-env-vars]

1. [Shared Credentials File][go-shared-credentials-file]

1. [EC2 Instance Profile][go-iam-roles-for-ec2-instances]

For more information see [Specifying Credentials][go-specifying-credentials] in
the AWS SDK for Go documentation.


### Usage with Docker

```bash
docker run -it --rm \
  -v $HOME/.aws/credentials:/root/.aws/credentials:ro \
  -e AWS_PROFILE=default \
  -v $PWD/secret/:/var/secret/ \
  -e SECRET_ID=my-secret \
  -e SECRET_FILE=/var/secret/my-secret \
  quay.io/turner/secretsmanager-sidecar
```

### Usage with ECS/Fargate

The following ECS task definition container definitions configure a sidecar that runs this program and then after it exits with a code of 0, starts the app container which will have access to the secret file.

```json
[
  {
    "name": "app",
    "image": "quay.io/turner/turner-defaultbackend:0.2.0",
    "essential": true,
    "dependsOn": [
      {
        "containerName": "secrets",
        "condition": "SUCCESS"
      }
    ],
    "portMappings": [
      {
        "protocol": "tcp",
        "containerPort": 8080,
        "hostPort": 8080
      }
    ],
    "environment": [
      {
        "name": "SECRET",
        "value": "/var/secret/my-secret"
      }
    ],
    "mountPoints": [
      {
        "readOnly": true,
        "containerPath": "/var/secret",
        "sourceVolume": "secret"
      }
    ],    
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
        "awslogs-group": "/fargate/service/myapp",
        "awslogs-region": "us-east-1",
        "awslogs-stream-prefix": "ecs"
      }
    }
  },
  {
    "name": "secrets",
    "image": "quay.io/turner/secretsmanager-sidecar:0.1.0",
    "essential": false,
    "environment": [
      {
        "name": "SECRET_ID",
        "value": "my-secret"
      },
      {
        "name": "SECRET_FILE",
        "value": "/var/secret/my-secret"
      }
    ],
    "mountPoints": [
      {
        "readOnly": false,
        "containerPath": "/var/secret",
        "sourceVolume": "secret"
      }
    ],    
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
        "awslogs-group": "/fargate/service/myapp",
        "awslogs-region": "us-east-1",
        "awslogs-stream-prefix": "ecs"
      }
    }
  }  
]
```

[go-sdk]: https://aws.amazon.com/documentation/sdk-for-go/
[go-env-vars]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#environment-variables
[go-shared-credentials-file]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#shared-credentials-file
[go-iam-roles-for-ec2-instances]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#iam-roles-for-ec2-instances
[go-specifying-credentials]: http://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials