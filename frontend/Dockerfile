# Choose the Image which has Node installed already
FROM node:lts-alpine

# make the 'client' folder the current working directory
WORKDIR /frontend

# copy project files and folders to the current working directory (i.e. 'app' folder)
COPY . .

# install project dependencies
RUN npm install

ENV HOST=0.0.0.0

EXPOSE 8080
CMD npm run serve
