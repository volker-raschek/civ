FROM scratch AS build

COPY civ-* /usr/bin/civ

ENTRYPOINT [ "/usr/bin/civ" ]