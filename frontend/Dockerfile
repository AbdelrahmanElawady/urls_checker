# Choose the Image which has Node installed already
FROM node:lts-alpine

# install simple http server for serving static content
RUN npm install -g http-server

# make the 'client' folder the current working directory
WORKDIR /client

# copy project files and folders to the current working directory (i.e. 'app' folder)
COPY . .

RUN chmod +x /client/start.sh
RUN apk add openssh-server && ssh-keygen -A

# install project dependencies
RUN npm install

# run
RUN npm run build

CMD npm run serve -- --port 8080

EXPOSE 8080
