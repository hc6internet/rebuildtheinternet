#!/bin/bash
memcached -d -u memcache
go run server.go
