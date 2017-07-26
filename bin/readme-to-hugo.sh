#!/bin/bash

# Copy the <middleware>/README.md to middleware/<mddleware>.md, add
# some Hugo meta data end let Hugo render it.

weight=1
tags="middleware"

function header() {
    local middleware=$1
    local title=$2
    local description=$3

    echo "+++"
    echo "title = \"$title\""
    echo "description = \"$description\""
    echo "weight = $weight"
    echo "tags = [ " \"$tags\" "," \"$middleware\" "]"
    echo "categories = [ \"middleware\" ]"
    echo "date = \"2017-07-24T15:25:40+00:00\""
#    echo "date = \"$(date --iso-8601=seconds)\""
    echo "+++"
}

function parse() {
    local file="$1"
    local middleware=$(basename $(dirname $file))
    local state="title"
    local description=""
    local rest=""

    while IFS= read -r line; do
        if [[ $state == "title" ]]; then
            title=${line#\# }
            state="skip"
            continue
        fi
    
        if [[ $state == "skip" ]]; then
            [[ -z "$line" ]] && continue
            state="description"
        fi

        if [[ $state == "description" ]]; then
            if [[ -z "$line" ]]; then
                state="rest"
                continue
            fi
            description+="$line "
        fi

        if [[ $state == "rest" ]]; then
            # stupid shell to-readd newline
            rest+="$line
"
        fi
    done < $file

    header $middleware "$title" "$description"
    echo
    echo "$rest"
}

if [[ -z "$1" ]]; then
    echo $0: Need containing directory of CoreDNS middleware
    exit 1
fi

content="../content/middleware"
if [[ ! -e "$content" ]]; then
    echo $0: Need to be run from the site\'s bin directory
    exit 1
fi

for dir in "$1"/*; do
    file=$dir/README.md
    if [[ ! -e $file ]]; then continue; fi
    middleware=$(basename $(dirname $file))
    
    echo $middleware.md >&2
    echo $file >&2
    parse $file > $content/$middleware.md
    weight=$((weight + 1))
done
