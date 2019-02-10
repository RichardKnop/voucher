# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Contact maintainer with any issues you encounter
MAINTAINER Richard Knop <risoknop@gmail.com>

# Set environment variables
ENV PATH /go/bin:$PATH

# Create a new unprivileged user
RUN useradd -m -d /home/app --user-group --shell /bin/false app

# Cd into the api code directory
WORKDIR /go/src/github.com/RichardKnop/voucher

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/RichardKnop/voucher

# Chown the application directory to app user
RUN chown -R app:app /go/src/github.com/RichardKnop/voucher

# Use the unprivileged user
USER app

# Set GO111MODULE=on variable to activate module support
ENV GO111MODULE on

# Install the program
RUN go install github.com/RichardKnop/voucher

# Run the voucher command by default when the container starts
ENTRYPOINT ["/go/bin/voucher"]

# Document that the service listens on port 8080
EXPOSE 8080
