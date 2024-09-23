# changed here bc air updated
FROM golang:1.23-alpine


WORKDIR /usr/src/app

RUN apk add --no-cache python3 py3-pip

# Create a virtual environment !!! imp ---> if here not working, go docker shell and install manually
RUN python3 -m venv /usr/src/app/venv

# Upgrade pip and install required packages inside the virtual environment !! important
RUN /usr/src/app/venv/bin/pip install --no-cache-dir -U pip \
    && /usr/src/app/venv/bin/pip install --no-cache-dir python-dotenv ipython google google-generativeai

RUN go install github.com/air-verse/air@latest


COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN go mod tidy
