#Build the go application  
FROM golang:alpine AS builder

#Folder where app code will be shared 
WORKDIR /app

#Copy the source code from the hostmachine to the container 
COPY . .

#Download the required dependencies 
RUN go mod download

#Compile the application code
RUN go build -o bulk-parser .

#Run the application 
CMD ["./bulk-parser"]







