FROM scratch
ADD fe fe
ENTRYPOINT ["/fe"]
EXPOSE 8000

HEALTHCHECK --interval=5m --timeout=3s \
  CMD curl -f http://localhost:8000/ || exit 1
