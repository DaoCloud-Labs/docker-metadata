# Usage:  ./check_macvlan_ip.sh '172.16.[0-9]\{1,3\}\.[0-9]\{1,3\}'

for loop in 1 2 3 4 5 6
do
    if ip addr show | grep -Po 'inet \K[\d.]+' | grep -o $1 >/dev/null
    then
        MAC_VLAN_SET_IP=`ip addr show | grep -Po 'inet \K[\d.]+' | grep -o $1 `
        echo $MAC_VLAN_SET_IP
        break
    else
        # echo -e "sleep 10s for ip alive....\n"
        sleep 10s
    fi
done

if [ ! $MAC_VLAN_SET_IP ]
then
    echo "please check you macvlan....."
    exit -1
fi