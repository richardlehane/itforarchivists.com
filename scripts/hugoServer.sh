#!/bin/bash
cd "$( dirname "${BASH_SOURCE[0]}")/.."
hugo -b "http://localhost:8081" -D -w