curl --user runner:$RUNNER_AUTH -X POST -d @../testdata/develop/1dm6kby/14kvrph.json http://localhost:8081/siegfried/logs/develop
curl --user runner:$RUNNER_AUTH -X POST -d @../testdata/develop/1dm6kby/1sq1zaj.json http://localhost:8081/siegfried/logs/develop
curl --user runner:$RUNNER_AUTH -X POST -d @../testdata/develop/1dm6kby/2f36wgx.json http://localhost:8081/siegfried/logs/develop
curl --user runner:$RUNNER_AUTH -X POST -d @../testdata/develop/1dm6kby/37985jf.json http://localhost:8081/siegfried/logs/develop
curl --user runner:$RUNNER_AUTH -X POST -d @../testdata/develop/1dm6kby/x24z03.json http://localhost:8081/siegfried/logs/develop

curl --user runner:$RUNNER_AUTH -X POST -d @../testdata/bench/1dnj7df/1mj4wgx.json http://localhost:8081/siegfried/logs/bench
curl --user runner:$RUNNER_AUTH -X POST -d @../testdata/bench/1dnj7df/1qr7zaj.json http://localhost:8081/siegfried/logs/bench
curl --user runner:$RUNNER_AUTH -X POST -d @../testdata/bench/1dnj7df/3n445jf.json http://localhost:8081/siegfried/logs/bench
curl --user runner:$RUNNER_AUTH -X POST -d @../testdata/bench/1dnj7df/km5zaj.json http://localhost:8081/siegfried/logs/bench









