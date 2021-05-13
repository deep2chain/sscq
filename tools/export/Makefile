# written by junying, 2020-01-15

genesis2accs:
	@printkey accounts genesis_old.json address amount denom|repl "[u'" ""|repl "']" ""|deline stake|rowcnt 2 > accounts.list

checksum:
	@printkey accounts genesis_old.json address amount denom|repl "[u'" ""|repl "']" ""|deline stake|rowcnt 2|column 2|sumup

export:
	@linecnt=$$(cat accounts.list|linecount)
	 for i in $$(python -c "print ' '.join(str(item) for item in range($$linecnt))"); \
	 do \
	 address=$$(cat accounts.list|row $$(($$i+1))|column 1);\
	 balance=$$(cat accounts.list|row $$(($$i+1))|column 2);\
	 ssd add-genesis-account $$address ${balance}satoshi --home ${HOME}/.ssd;\
	 done;