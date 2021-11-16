FROM centos:7
ENV REDIS_HOST redis-cn02nwh9l8hdm1gmx.redis.ivolces.com:6379
ENV REDIS_PASSWORD Group12345678
ENV CONFIG_NAME rain
ENV AUTH_USERNAME kxc3esdiu23d
ENV AUTH_PASSWORD xz8cvjs9q3m1
ENV MYSQL_USERNAME group8
ENV MYSQL_PASSWORD Group12345678
ENV MYSQL_HOST rdsmysqlhf6ed4fde2675f947
ENV MYSQL_PORT 3306
ENV MYSQL_DBNAME envelope_rains
ENV GIN_MOD release
WORKDIR /root
COPY envelope_manager ./server
COPY rain.yaml ./rain.yaml
EXPOSE 9090
CMD /root/server

