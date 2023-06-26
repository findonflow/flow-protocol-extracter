#/bin/bash

lftp -c 'set net:idle 10
         set net:max-retries 0
         set net:reconnect-interval-base 3
         set net:reconnect-interval-max 3
         pget -n 20 -c "https://storage.googleapis.com/flow-execution-archive/mainnet5/protocol-db.tar.gz" -o mainnet5.tar.gz
	 '
mkdir badger
mv mainnet5.tar.gz badger
cd badger
pigz -dk -p 10 -c mainnet5.tar.gz | pv | tar xf -
