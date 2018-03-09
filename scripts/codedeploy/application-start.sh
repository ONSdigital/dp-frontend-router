#!/bin/bash

AWS_REGION=
CONFIG_BUCKET=
ECR_REPOSITORY_URI=
GIT_COMMIT=

INSTANCE=$(curl -s http://instance-data/latest/meta-data/instance-id)
CONFIG=$(aws --region $AWS_REGION ec2 describe-tags --filters "Name=resource-id,Values=$INSTANCE" "Name=key,Values=Configuration" --output text | awk '{print $5}')

(aws s3 cp s3://$CONFIG_BUCKET/frontend-router/$CONFIG.asc . && gpg --decrypt $CONFIG.asc > $CONFIG) || exit $?

if [[ $DEPLOYMENT_GROUP_NAME =~ [a-z]+-publishing ]]; then
  DOCKER_NETWORK=publishing
else
  DOCKER_NETWORK=website
fi

source $CONFIG && docker run -d                      \
  --env=BABBAGE_URL=$BABBAGE_URL                     \
  --env=BIND_ADDR=$BIND_ADDR                         \
  --env=RENDERER_URL=$RENDERER_URL                   \
  --env=RESOLVER_URL=$RESOLVER_URL                   \
  --env=DOWNLOADER_URL=$DOWNLOADER_URL               \
  --env=SPLASH_PAGE=$SPLASH_PAGE                     \
  --env=SITE_DOMAIN=$SITE_DOMAIN                     \
  --env=AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID         \
  --env=AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY \
  --env=AWS_REGION=$AWS_REGION                       \
  --env=ANALYTICS_SQS_URL=$ANALYTICS_SQS_URL         \
  --env=REDIRECT_SECRET=$REDIRECT_SECRET             \
  --name=frontend-router                             \
  --net=$DOCKER_NETWORK                              \
  --restart=always                                   \
  $ECR_REPOSITORY_URI/frontend-router:$GIT_COMMIT
