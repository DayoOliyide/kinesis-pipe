#!/bin/bash

AWS_DEFAULT_REGION=eu-west-1
AWS_ACCESS_KEY_ID=blah
AWS_SECRET_ACCESS_KEY=doubleblah


A_TALE="it was the worst of times,
it was the age of wisdom,
it was the age of foolishness,
it was the epoch of belief,
it was the epoch of incredulity,
it was the season of Light,
it was the season of Darkness,
it was the spring of hope,
it was the winter of despair,
we had everything before us,
we had nothing before us,
we were all going direct to Heaven,
we were all going direct the other way"



while read -r line
do
    aws --endpoint-url=http://localhost:4568 kinesis put-record --stream-name Tales --partition-key $(($RANDOM % 1000)) --data "$line"
done < <(printf '%s\n' "$A_TALE")

