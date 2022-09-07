#!/bin/bash

SCRIPT_DIR=$(realpath $(dirname $0))

HTTP_METHODS=("GET" "POST" "PUT" "DELETE")
HTTP_VERSIONS=("HTTP/1.0" "HTTP/1.1" "HTTP/2.0")
HTTP_CODES=("200" "201" "202" "203" "204" "300" "302" "304" "307" "308" "400" "401" "403" "404" "408" "500" "501" "502" "503" "504")
HTTP_URIS=("/admin/index.php"
	"/admin/index.php?user=root"
	"/admin/index.php?user=mark"
	"/admin/list.php"
	"/admin/show.php"
	"/admin/show.php?data=value"
	"/admin/show.php?data=value&data2=value2"
	"/admin/logout.php"
	"/admin/login.php"
	"/cart/show.php"
	"/cart/show.php?page=10"
	"/cart/show.php?page=100"
	"/cart/show.php?filter=98765,4356"
	"/cart/checkout.php"
	"/cart/delete.php"
	"/cart/error.php"
	"/cart/check.php"
	"/cart/check.php?items=10,15,18"
	"/cart/check.php?items=10,15,18&data=value&data2=value2"
	"/error.php"
	"/login.php"
	"/logout.php"
	"/index.php"
	"/index.php?filter=12345,3255,888"
	"/index.php?filter=11111,3255,888"
	"/index.php?filter=11111,2222,888"
	"/index.php?filter=11111,2222,333"
	"/index.php?filter=22222,2222,333"
	"/index.php?filter=22222,3333,333"
	"/index.php?filter=22222,4444,333"
	"/index.php?filter=22222,5555,333"
	"/index.php?filter=22222,6666,333"
	"/index.php?filter=22222,7777,333"
	"/index.php?filter=22222,8888,333"
	"/index.php?filter=22222,9999,333"
	"/index.php?filter=33333,1111,111"
	"/index.php?filter=33333,1111,222"
	"/index.php?filter=33333,1111,333"
	"/index.php?filter=33333,1111,444"
	"/index.php?filter=33333,1111,555"
	"/index.php?filter=33333,1111,666"
	"/index.php?filter=33333,1111,777"
	"/index.php?filter=33333,1111,888"
	"/index.php?filter=33333,1111,999"
)
HTTP_IPS=("0.0.0.0"
	"0.0.0.1"
	"0.0.1.1"
	"0.1.1.1"
	"1.1.1.1"
	"1.1.1.2"
	"1.1.2.2"
	"1.2.2.2"
	"2.2.2.2"
	"2.2.2.3"
	"2.2.3.3"
	"2.3.3.3"
	"3.3.3.3"
	"3.3.3.4"
	"3.3.4.4"
	"3.4.4.4"
	"4.4.4.4"
	"4.4.4.5"
	"4.4.5.5"
	"4.5.5.5"
	"5.5.5.5"
	"5.5.5.6"
	"5.5.6.6"
	"5.6.6.6"
	"6.6.6.6"
	"6.6.6.7"
	"6.6.7.7"
	"6.7.7.7"
	"7.7.7.7"
	"7.7.7.8"
	"7.7.8.8"
	"7.8.8.8"
	"8.8.8.8"
	"8.8.8.9"
	"8.8.9.9"
	"8.9.9.9"
	"9.9.9.9"
	"9.9.9.10"
	"9.9.10.10"
	"9.10.10.10"
	"10.10.10.10"
	"10.10.10.11"
	"10.10.11.11"
	"10.11.11.11"
	"11.11.11.11"
	"11.11.11.12"
	"11.11.12.12"
	"11.12.12.12"
	"12.12.12.12"
	"12.12.12.13"
	"12.12.13.13"
	"12.13.13.13"
	"13.13.13.13"
	"13.13.13.14"
	"13.13.14.14"
	"13.14.14.14"
	"14.14.14.14"
	"14.14.14.15"
	"14.14.15.15"
	"14.15.15.15"
	"15.15.15.15"
	"15.15.15.16"
	"15.15.16.16"
	"15.16.16.16"
	"16.16.16.16"
	"16.16.16.17"
	"16.16.17.17"
	"16.17.17.17"
	"17.17.17.17"
	"17.17.17.18"
	"17.17.18.18"
	"17.18.18.18"
	"18.18.18.18"
	"18.18.18.19"
	"18.18.19.19"
	"18.19.19.19"
	"19.19.19.19"
	"19.19.19.20"
	"19.19.20.20"
	"19.20.20.20"
)
HTTP_AGENTS=("Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0" "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)" "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.52 Safari/536.5" "Safari/537.36 OPR/38.0.2220.41" "Opera/9.80 (Macintosh; Intel Mac OS X; U; en)" "Version/13.1.1 Mobile/15E148 Safari/604.1" "curl/7.64.1" "PostmanRuntime/7.26.5")

declare -A ipsMap
declare -A pathsMap
declare -A pathsMapDur

increment_ip_counters() {
	ip=$1
	if [[ -z "${ipsMap[$ip]}" ]]
	then
		ipsMap[$ip]=0
	fi
	ipsMap[$ip]=$((${ipsMap[$ip]}+1))
}

increment_path_counters() {
	path=$1
	dur=$2
	path_key=$(echo "$path" | cut -f1 -d\?)
	if [[ -z "${pathsMap[$path_key]}" ]]
	then
		pathsMap[$path_key]=0
		pathsMapDur[$path_key]=0
	fi
	pathsMap[$path_key]=$((${pathsMap[$path_key]}+1))
	pathsMapDur[$path_key]=$((${pathsMapDur[$path_key]}+$dur))
}

random_alnum() {
	len=$1
	cat /dev/urandom | tr -dc 'A-Za-z0-9' | head -c $len
}
random_alnum_space() {
	len=$1
	cat /dev/urandom | tr -dc 'A-Za-z0-9 ' | head -c $len
}
random_string() {
	len=$1
	cat /dev/urandom | tr -dc A-Za-z._- | head -c $len
}

generate_ip() {
	is_rand_ip=$(($RANDOM%10))
	if (( $is_rand_ip > 0 ))
	then
		index=$(($RANDOM%${#HTTP_IPS[@]}))
		ip=${HTTP_IPS[$index]}
	else
		byte1=$(($RANDOM%256))
		byte2=$(($RANDOM%256))
		byte3=$(($RANDOM%256))
		byte4=$(($RANDOM%256))
		ip="$byte1.$byte2.$byte3.$byte4"
	fi
	echo "$ip"
}

generate_user() {
	username=$(random_alnum 10)
	echo $username
}

generate_date() {
	seconds_back=$(od -A n -t dL -N 4 /dev/urandom)
	date -d "$seconds_back seconds ago" -u "+%d/%b/%Y:%H:%M:%S %z"
}

generate_method() {
	index=$(($RANDOM%${#HTTP_METHODS[@]}))
	echo ${HTTP_METHODS[$index]}
}

generate_random_uri() {
	path_len=$(($RANDOM%3+1))
	uri=""
	i=0
	while (($i < $path_len))
	do
		p=$(random_string 10)
		uri=$uri/$p
		((i++))
	done
	is_query=$(($RANDOM%2))
	if (($is_query > 0))
	then
		query_len=$(($RANDOM%3+1))
		query="?"
		j=0
		while (($j < $query_len))
		do
			q=$(random_alnum 10 | sed -e 's/ /+/g')
			v=$(random_alnum_space 10 | sed -e 's/ /+/g')
			query="${query}${q}=${v}"
			if (($j < $query_len-1))
			then
				query="${query}&"
			fi
			((j++))
		done
		uri="${uri}${query}"
	fi
	echo "$uri"
}
generate_uri() {
	is_rand_uri=$(($RANDOM%10))
	if (( $is_rand_uri > 0 ))
	then
		index=$(($RANDOM%${#HTTP_URIS[@]}))
		uri=${HTTP_URIS[$index]}
	else
		uri=$(generate_random_uri)
	fi
	echo "$uri"
}

generate_version() {
	index=$(($RANDOM%${#HTTP_VERSIONS[@]}))
	echo ${HTTP_VERSIONS[$index]}
}

generate_status() {
	index=$(($RANDOM%${#HTTP_CODES[@]}))
	echo ${HTTP_CODES[$index]}
}

generate_response_time() {
	ms=$(($RANDOM%10000+30))
	echo $ms
}

generate_user_agent() {
	index=$(($RANDOM%${#HTTP_AGENTS[@]}))
	echo ${HTTP_AGENTS[$index]}
}

## Main

NUM_OF_LINES=$1
if [[ -z "$NUM_OF_LINES" ]]
then
	NUM_OF_LINES=1000
fi

i=0
while (($i < $NUM_OF_LINES))
do
	ip=$(generate_ip)
	usr=$(generate_user)
	d=$(generate_date)
	m=$(generate_method)
	uri=$(generate_uri)
	v=$(generate_version)
	code=$(generate_status)
	ms=$(generate_response_time)
	ua=$(generate_user_agent)

	increment_ip_counters $ip
	increment_path_counters "$uri" "$ms"

	echo "$ip - $usr [$d] \"$m $uri $v\" $code $ms \"$ua\""
	((i++))
done

echo "=================================================="
echo "== Log Stats"

echo "== Stats for IPs"
for key in ${!ipsMap[@]}
do
	echo "   $key: ${ipsMap[$key]}"
done | sort -n -k2 -t:

echo
echo "== Stats for URIs"
for key in "${!pathsMap[@]}"
do
	#echo "DEBUG: $key: ${pathsMapDur[$key]} / ${pathsMap[$key]}"
	#echo "   $key: ${pathsMap[$key]}"
	avg=$(echo "scale=2; ${pathsMapDur[$key]}/${pathsMap[$key]}" | bc -l)
	echo "   $key: $avg"
done | sort -n -k2 -t:
