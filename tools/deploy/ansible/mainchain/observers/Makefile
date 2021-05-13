# junying-todo, 2020-02-10
# you should prepare the following files.
# 1. /etc/ansible/hosts
# >>>
# [observers]
# observer1 ansible_host=192.168.5.11 ansible_user=root ansible_ssh_private_key_file=/tmp/192.168.5.11.pem
# observer2 ansible_host=192.168.5.12 ansible_user=root ansible_ssh_private_key_file=/tmp/192.168.5.12.pem
# observer3 ansible_host=192.168.5.13 ansible_user=root ansible_ssh_private_key_file=/tmp/192.168.5.13.pem
# observer4 ansible_host=192.168.5.14 ansible_user=root ansible_ssh_private_key_file=/tmp/192.168.5.14.pem
#
# 2. config.toml, genesis.json, hsd, hscli, Makefile
#################################################################
CONFIG_DIR=/home/ubuntu/config/nosentry/htdf/observation

all: init cp.config start

init:
	@ansible observers -m shell -a "if ! [ -d ' /root/.hsd/config' ];then mkdir -p /root/.hsd/config; fi"

copy: cp.config cp.hsd cp.hscli cp.Makefile
cp.config:
	@ansible observers -m copy -a "src=/home/ubuntu/config/nosentry/htdf/observation/config.toml dest=/root/.hsd/config force=yes"
	@ansible observers -m copy -a "src=/home/ubuntu/config/nosentry/htdf/observation/genesis.json dest=/root/.hsd/config force=yes"
cp.hsd:	
	@ansible observers -m copy -a "src=/home/ubuntu/go/bin/hsd dest=/usr/local/bin force=yes mode=755"
cp.hscli:	
	@ansible observers -m copy -a "src=/home/ubuntu/go/bin/hscli dest=/usr/local/bin force=yes mode=755"
cp.Makefile:
	@ansible observers -m copy -a "src=/home/ubuntu/config/nosentry/htdf/observation/Makefile dest=/root force=yes"	

start:
	@ansible observers -m shell -a "make start-daemon start-rest"

stop:
	@ansible observers -m shell -a "make stop-daemon stop-rest"

clear:
	@@ansible observers -m shell -a "make clean"

# checking part
confirm.validators:
	@ansible observers -m shell -a "curl http://localhost:26657/validators|grep address"

confirm.height:
	@ansible observers -m shell -a "curl http://localhost:26657/status|grep latest_block_height"