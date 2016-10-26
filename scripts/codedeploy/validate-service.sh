#!/bin/bash

if [[ $(docker inspect --format="{{ .State.Running}}" frontend-router) == "false" ]]; then
  exit 1;
fi
