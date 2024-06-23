SHELL=bash

include .env

dev:
	$(eval export $(sed 's/#.*//g' .env | xargs))
	go run main.go