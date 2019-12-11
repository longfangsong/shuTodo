#!/usr/bin/env bash

export DB_ADDRESS="postgresql://localhost:5432/postgres?sslmode=disable"
export JWT_SECRET="test"
export PORT="8001"
gin -p 8000 run main.go