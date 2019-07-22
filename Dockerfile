FROM debian
RUN apt-get update && apt-get install -y ca-certificates
COPY minter-test-task minter-test-task
RUN chmod +x minter-test-task
COPY wait-for-it.sh wait-for-it.sh
RUN chmod +x wait-for-it.sh
CMD ./minter-test-task -continueParsing