#!/bin/bash

rm -rf /var/run/pgbouncer/*
pgbouncer pgbouncer.ini -u nobody