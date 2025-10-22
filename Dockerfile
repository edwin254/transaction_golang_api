FROM golang:1.25.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download  && go mod verify

COPY . .

RUN ls -alh

RUN CGO_ENABLED=0 go build -v -o /app/main ./cmd/api

# ----------------------------------------------------------------------
# Stage 2: Final Image
# Use a minimal, non-OS image for the final deployment for security and size.
# ----------------------------------------------------------------------
FROM scratch

# Copy the built executable from the builder stage
# It's now located at /app/main in the scratch environment
COPY --from=builder /app/main /main

# Define the user to run the application as
USER 1000:1000

# Set the entry point to run the application
ENTRYPOINT ["/main"]