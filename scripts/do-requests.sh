#!/bin/bash

trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

function _go_metrics() {
    while true; do
        curl -s -o /dev/null http://localhost:8000/
        sleep 5
    done
}

function _python_metrics() {
    while true; do
        curl -s -o /dev/null http://localhost:8001/send-email
        sleep 5
    done
}

function _node_metrics() {
    while true; do
        curl -s -o /dev/null http://localhost:8003
        sleep 5
    done
}

function _ruby_metrics() {
    local temp_file="/tmp/poke.txt"
    local temp_image="/tmp/poke.png"

    cat scripts/pokemons.txt | shuf >$temp_file

    while IFS= read -r pokemon; do
        clear -x

        curl -s -o /dev/null http://localhost:8002/pokemon/$pokemon/details
        curl -s http://localhost:8002/pokemon/$pokemon | jq '.'
        curl -s http://localhost:8002/pokemon/$pokemon/image --output $temp_image
        catimg -t -w100 $temp_image

        sleep 5
    done <$temp_file

}

_go_metrics &
_python_metrics &
_node_metrics &
_ruby_metrics
