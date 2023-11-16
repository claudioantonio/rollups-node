ARG NODE_ANVIL_VERSION=devel
FROM cartesi/anvil:${NODE_ANVIL_VERSION}

WORKDIR /usr/share/cartesi
COPY ./entrypoint.sh /
COPY ./anvil_state/state.json .
COPY ./deployments/localhost.json .

ENTRYPOINT ["/entrypoint.sh"]
CMD ["anvil", "--load-state", "/usr/share/cartesi/state.json"]
