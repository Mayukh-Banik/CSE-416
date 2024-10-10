# Basic Ubuntu Stuff
FROM ubuntu:22.04
RUN apt update && apt upgrade -y && apt install wget -y

# Go Installation
RUN wget https://go.dev/dl/go1.23.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz
# RUN export PATH=$PATH:/usr/local/go/bin && source $HOME/.profile
ENV PATH=$PATH:/usr/local/go/bin

# Installs Kubo IPFS
RUN wget https://dist.ipfs.tech/kubo/v0.30.0/kubo_v0.30.0_linux-amd64.tar.gz && \ 
    tar xvzf kubo_v0.30.0_linux-amd64.tar.gz && \
    mv kubo/ipfs /usr/local/bin/ipfs && \
    rm -rf kubo kubo_v0.30.0_linux-amd64.tar.gz

RUN IPFS_PATH=~/.ipfs ipfs init

RUN ipfs config Addresses.API /ip4/0.0.0.0/tcp/5001 && ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/8080

RUN echo -e "/key/swarm/psk/1.0.0/\n/base16/\n$(tr -dc 'a-f0-9' < /dev/urandom | head -c64)" > $IPFS_PATH/swarm.key

RUN IPFS_PATH=~/.ipfs ipfs bootstrap rm --all
RUN IPFS_PATH=~/.ipfs ipfs bootstrap add /ip4/$(hostname -I | awk '{print $1}')/tcp/4001/ipfs/$(IPFS_PATH=~/.ipfs ipfs config show | grep '"PeerID"' | awk -F'"' '{print $4}')

RUN export LIBP2P_FORCE_PNET=1 

EXPOSE 4001 5001 8080

CMD ["ipfs", "daemon"]

