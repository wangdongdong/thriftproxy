#!/bin/bash

set -x -e

thrift --gen go:ignore_initialisms=true -out . adder.thrift
