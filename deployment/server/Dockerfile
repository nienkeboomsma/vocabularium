FROM golang:1.24

RUN apt-get update && apt-get install -y \
    git \
    qt5-qmake \
    qtbase5-dev \
    qtbase5-dev-tools \
    libqt5xmlpatterns5-dev \
    build-essential \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /collatinus

RUN git clone --branch Daemon https://github.com/biblissima/collatinus.git .

RUN qmake collatinus.pro && make -j$(nproc)

RUN qmake Client_C11.pro && make -j$(nproc)

RUN mkdir -p /var/log/collatinus && chmod 777 /var/log/collatinus

WORKDIR /server

COPY . .

RUN go build -o server main.go

EXPOSE 5555

WORKDIR /

# The collatinus server logs a bunch of diagnostic info that spams the container logs
CMD bash -c "./collatinus/bin/collatinusd >/dev/null 2>&1 & ./server/server"
