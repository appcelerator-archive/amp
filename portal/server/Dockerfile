FROM alpine
COPY portal.alpine /portal-server
RUN chmod +x /portal-server && mkdir /public
COPY public/ /public/
CMD [ "/portal-server" ]
