FROM gitpod/workspace-full

USER gitpod
RUN sh -c "$(curl -sSfL https://release.solana.com/v1.10.32/install)"
ENV PATH="/home/gitpod/.local/share/solana/install/active_release/bin:$PATH"
RUN solana config set --url devnet
RUN cargo install sugar-cli
RUN sudo curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin