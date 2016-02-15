#!/bin/bash

hash=$(git rev-parse HEAD)

sed -i s/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx/$hash/ version.go


