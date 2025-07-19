FROM postgres:17

COPY init.sh /docker-entrypoint-initdb.d/init.sh
COPY spqr.sql /docker-entrypoint-initdb.d/spqr.sql.disabled

RUN chmod +x /docker-entrypoint-initdb.d/init.sh
