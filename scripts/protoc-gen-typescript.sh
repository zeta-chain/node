#!/usr/bin/env bash

cd proto
rm -rf typescript
buf generate --template buf.ts.yaml