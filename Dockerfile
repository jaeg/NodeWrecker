FROM scratch
ARG binary
ENV env_var_name=$binary
ADD pkg/$binary /NodeWrecker

ENTRYPOINT ["/NodeWrecker"]