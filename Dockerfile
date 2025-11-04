FROM alpine:latest

# It's a good practice to run containers as a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy the binary from the build context (goreleaser will place it here)
COPY micomido /usr/local/bin/micomido
COPY migrator /usr/local/bin/migrator

# Copy migrations
COPY migrations /migrations

# Expose the port the app runs on
EXPOSE 8080

# The command to run the application
ENTRYPOINT ["/usr/local/bin/micomido"]
