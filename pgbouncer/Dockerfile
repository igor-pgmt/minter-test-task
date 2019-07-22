FROM debian
RUN apt-get update && apt-get install -y pgbouncer
RUN mkdir /var/run/pgbouncer && chmod -R 777 /var/run/pgbouncer
RUN mkdir /workdir && chmod -R 777 /workdir && cd /workdir
WORKDIR /workdir
COPY pgbouncer.ini pgbouncer.ini
COPY userlist.txt userlist.txt
COPY start.sh start.sh
RUN chmod +x start.sh
CMD ./start.sh