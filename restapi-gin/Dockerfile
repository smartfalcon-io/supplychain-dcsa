# Step 1: Choose a Base Image with Go installed
FROM golang:1.16

# Step 2: Set the Working Directory inside the container
WORKDIR /home/srivani/Bug/api-code

# Step 3: Copy the Go application files to the container
COPY . .

# Step 4: Build the Go application
RUN go build -o main .

# Step 5: Expose the port the application will run on
EXPOSE 8080

CMD ["api-code"]

