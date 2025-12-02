# ----------- STAGE 1: BUILD THE GO BINARY -----------------

FROM golang:alpine AS builder

# Create a directory inside the container for building
RUN mkdir /build

# Copy all Go project files into the build directory
ADD . /build/

# Set working directory
WORKDIR /build

# Build a statically compiled Go binary named 'main'
RUN go build -o main .

# ----------- STAGE 2: RUNTIME IMAGE -----------------------

FROM alpine

# Create a non-root user to run the app
RUN adduser -S -H -D appuser

# Switch to that user (for security)
USER appuser

# Create application directory
WORKDIR /app

# Copy app source (optional, only needed if you need views/configs)
COPY . /app

# Copy the compiled Go binary from builder stage
COPY --from=builder /build/main /app/

# Expose port 3000
EXPOSE 3000

# Start the app
CMD ["./main"]
