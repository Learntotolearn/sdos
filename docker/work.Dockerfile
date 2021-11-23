ARG ALPINE_VERSION=3.14
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

FROM --platform=$TARGETPLATFORM alpine:${ALPINE_VERSION}

RUN apk add --update --no-cache bash wireguard-tools tcpdump git make ipset dnsmasq tini curl fping mtr jq tzdata ca-certificates dante-server

RUN git clone https://github.com/kuaifan/wondershaper.git && \
    cd wondershaper/ && \
    make install && \
    touch /var/log/dnsmasq.log

COPY ./conf/dnsmasq.conf /etc/dnsmasq.conf
COPY ./conf/resolv.conf /etc/resolv.conf
COPY ./conf/resolv.dnsmasq.conf /etc/resolv.dnsmasq.conf
COPY ./conf/sysctl.conf /etc/sysctl.conf

COPY ../release/sdos_${TARGETOS}_${TARGETARCH} /usr/bin/sdos
RUN chmod +x /usr/bin/sdos

COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY xray/xray.sh /xray.sh
RUN set -ex \
	&& mkdir -p /var/log/xray /usr/share/xray \
	&& chmod +x /xray.sh \
	&& /xray.sh "${TARGETPLATFORM}" \
	&& rm -fv /xray.sh \
	&& wget -O /usr/share/xray/geosite.dat https://github.com/v2fly/domain-list-community/releases/latest/download/dlc.dat \
	&& wget -O /usr/share/xray/geoip.dat https://github.com/v2fly/geoip/releases/latest/download/geoip.dat

WORKDIR /tmp/.sdwan

ENTRYPOINT ["/entrypoint.sh"]


# docker buildx create --use
# docker buildx build --platform linux/amd64 -t kuaifan/sdwan:work-0.0.1 --push -f ./work.Dockerfile .
# 需要 docker login 到 docker hub, 用户名 (docker id): kuaifan
