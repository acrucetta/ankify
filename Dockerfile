FROM golang:1.16-alpine3.14

WORKDIR /app

RUN apk add --no-cache python3
RUN apk add --no-cache wget
RUN wget https://bootstrap.pypa.io/get-pip.py
RUN python3 get-pip.py
RUN pip3 install --no-cache-dir --upgrade pipenv

COPY go.mod go.sum ./
RUN go mod download

COPY Pipfile Pipfile.lock ./

# Install pipenv
RUN pip install pipenv
RUN pipenv install --system --deploy --ignore-pipfile

COPY . .

ENTRYPOINT [ "./ankify" ]
