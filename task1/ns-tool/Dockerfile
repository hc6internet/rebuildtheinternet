FROM alpine:3.7

RUN \
    apk update && \
    apk upgrade && \
    apk --update --no-cache add bash iputils python3 py-pip && \

    # add virtualenv
    pip install --upgrade pip && \
    pip install virtualenv && \
    pip install requests && \

    # create virtualenv
    virtualenv -p /usr/bin/python3 py3 && \

    # clean up
    rm -rf /var/cache/apk/* && \
    rm -rf /tmp/* && \
    rm -rf /var/log/*

ADD pdns.py /pdns.py

CMD ["/bin/bash"]
