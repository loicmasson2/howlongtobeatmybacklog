FROM rust:slim-bookworm

# Add wasm target
RUN rustup target add wasm32-unknown-unknown

# Install wasm-pack for testing
RUN curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh

# Install trunk
ADD https://github.com/thedodd/trunk/releases/download/v0.17.2/trunk-x86_64-unknown-linux-gnu.tar.gz ./tmp
RUN cd /tmp && tar xf trunk-x86_64-unknown-linux-gnu.tar.gz && chmod +x trunk && mv trunk /bin

WORKDIR /usr/src/myapp

# While we’re waiting on a --dependencies-only build options for cargo,
# we can overcome this problem by changing our dev.Dockerfile to have a default
# src/main.rs with which the dependencies are built before we COPY any of our code into the build
COPY Cargo.toml .
COPY Trunk.toml .
# COPY Makefile.toml .
RUN mkdir src
# RUN cargo make

# copy the real source files 
COPY src src
COPY index.html index.html

# COPY ./crates/conduit-wasm .
# COPY .env.example .env

EXPOSE 8000

CMD [ "trunk", "serve" ]