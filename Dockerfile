FROM reg.yl.com/basic/python/al-py-ssd:v6
MAINTAINER ylops
ENV  APP="MODULE"
ADD .  /$APP
WORKDIR /$APP
CMD chmod 777 ./exporter-center && ./exporter-center server