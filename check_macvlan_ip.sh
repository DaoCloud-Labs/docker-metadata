# Usage:  ./check_macvlan_ip.sh '172.16.[0-9]\{1,3\}\.[0-9]\{1,3\}'

for loop in 1 2 3 4 5 6
do
    if ifconfig | awk -F "[: ]+" '/inet addr:/ { if ($4 != "127.0.0.1") print $4 }' | grep -o $1 >/dev/null
    then
        echo "online"
        break
    else
        echo "sleep 10s for ip alive...."
        sleep 10s
    fi
done
