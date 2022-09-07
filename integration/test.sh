#!/bin/bash

BASE_DIR=$(realpath $(dirname $0))
INPUT_BASE_DIR="$BASE_DIR/testdata"
INPUT_ARCHIVE="$INPUT_BASE_DIR/test-input.tar.bz2"
INPUT_FILE_PATTERN="$INPUT_BASE_DIR/test-input_*"
OUTPUT_FILE_PREFIX="$INPUT_BASE_DIR/test-output_"
OUTPUT_FILE_SUFFIX="json"
BINARY="$BASE_DIR/../target/log-parser"
TMP_OUTPUT="$BASE_DIR/.log-parser_$RANDOM"
TMP_OUTPUT_NORMALIZED="$TMP_OUTPUT.json"

tar xf $INPUT_ARCHIVE -C $INPUT_BASE_DIR

for f in $INPUT_FILE_PATTERN
do
	# clean the temp files
	> $TMP_OUTPUT
	> $TMP_OUTPUT_NORMALIZED
	test_suite=$(basename $f .log | cut -f2 -d_)
	expected_output="${OUTPUT_FILE_PREFIX}${test_suite}.$OUTPUT_FILE_SUFFIX"
	MAX_CLIENT_IPS=$(($(cat $expected_output | jq .top_client_ips | wc -l)-2))
	MAX_PATHS=$(($(cat $expected_output | jq .top_path_avg_seconds | wc -l)-2))

	# run the app
	$BINARY --in=$f --out=$TMP_OUTPUT --max-client-ips=$MAX_CLIENT_IPS --max-paths=$MAX_PATHS
	cat $TMP_OUTPUT | jq -S > $TMP_OUTPUT_NORMALIZED

	# Check whether the output equals the expected json
	diff -q $TMP_OUTPUT_NORMALIZED $expected_output
	if (( $? != 0 ))
	then
		echo "Integration test failed: $TMP_OUTPUT_NORMALIZED != $expected_output"
		exit 1
	fi
done
rm -f "$TMP_OUTPUT" "$TMP_OUTPUT_NORMALIZED"
