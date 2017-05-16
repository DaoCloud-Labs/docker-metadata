# Usage:  ./host_is_exist.sh '172.16.[0-9]\{1,3\}\.[0-9]\{1,3\}'

for loop in 1 2 3
do
    if ifconfig | awk -F "[: ]+" '/inet addr:/ { if ($4 != "127.0.0.1") print $4 }' | grep -o $1 >/dev/null
    then
        echo "online"
        break
    else
        echo "sleep 1m for ip alive...."
        sleep 1m
    fi
done