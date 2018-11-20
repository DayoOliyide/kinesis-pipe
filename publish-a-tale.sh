#!/bin/bash

AWS_DEFAULT_REGION=eu-west-1
AWS_ACCESS_KEY_ID=blah
AWS_SECRET_ACCESS_KEY=doubleblah

KINESIS_STREAM=Tales
KINESIS_SHARDS=10

A_TALE="It was the best of times,
it was the worst of times,
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

printf "Creating a kinesis stream called %s with %d shards .." ${KINESIS_STREAM} ${KINESIS_SHARDS}
aws --endpoint-url=http://localhost:4568 kinesis create-stream --stream-name ${KINESIS_STREAM} --shard-count ${KINESIS_SHARDS}
printf ".. Created stream \n" ${KINESIS_STREAM} ${KINESIS_SHARDS}

printf "Publishing a Tale on %s " ${KINESIS_STREAM}
while read -r line
do
    aws --endpoint-url=http://localhost:4568 kinesis put-record --stream-name ${KINESIS_STREAM} --partition-key $(($RANDOM % ${KINESIS_SHARDS})) --data "$line" 1> /dev/null
    printf "."
done < <(printf '%s\n' "$A_TALE")
printf " Finished publishing\n"
