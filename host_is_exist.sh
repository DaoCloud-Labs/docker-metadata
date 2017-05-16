# Usage:  ./host_is_exist.sh 192.168.1.1

for loop in 1 2 3
do
    if ping -q -c2 $1 >/dev/null
    then
        echo "online"
        break
    else
        echo "sleep 1m for ip alive...."
        sleep 1m
    fi
done