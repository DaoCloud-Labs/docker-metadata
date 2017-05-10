curl -s -X GET --unix-socket /var/run/dce-metadata/dce-metadata.sock http:/containers/$(hostname)/json |grep -o '"IPAddress":"[^"]*' | grep -o '[^"]*$' | head -1
