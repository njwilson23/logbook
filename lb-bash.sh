#! /usr/bin/bash

lb() {
    echo "here"
    if [[ $@ == "yesterday" ]]; then
        command lb -days -1
    else
        command lb "$@"
    fi
}
