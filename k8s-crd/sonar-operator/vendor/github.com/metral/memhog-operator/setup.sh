#!/bin/bash

# Install dependencies
rm -rf vendor
glide cc
glide install -v -s
