FROM scratch

ADD pkg/NodeWrecker-linux-pi /NodeWrecker

ENTRYPOINT ["/NodeWrecker"]