CHAIN_ID = testchain

check:
	@ps -ef | grep ssd
	@ps -ef | grep sscli

clean:
	@rm -rf /root/.ssd /root/.ssd/app.log

stop-daemon:
	@pkill ssd

stop-rest:
	@pkill sscli

start-daemon:
	@nohup ssd start >> ${HOME}/.ssd/app.log  2>&1  &

start-rest:
	@nohup sscli rest-server --chain-id=${CHAIN_ID} --trust-node=true --laddr=tcp://0.0.0.0:1317 >> ${HOME}/.ssd/restServer.log  2>&1  &

# too slow
checksum:
	@linecnt=$$(cat accounts.list|linecount);\
	 sum=0;\
	 for i in $$(python -c "print ' '.join(str(item) for item in range($$linecnt))"); \
	 do \
	 account=$$(cat accounts.list|row $$(($$i+1))|column 1);\
	 balance=$$(sscli query account $$account|row 4|fromstr "         "|endstr satoshi);\
	 echo $$(($$i+1)) $$account $$balance >> accounts_new.list;\
	 sum=$$(($$sum+$$balance));\
	 done;

stopall: stop-daemon stop-rest
restart-daemon: stop-daemon start-daemon
restart-rest: stop-rest start-rest
