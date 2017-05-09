FROM alpine
COPY amp-ui.alpine /amp-ui-server
RUN chmod +x /amp-ui-server && mkdir /public
COPY public/ /public/
CMD [ "/amp-ui-server" ]
