## Plugeth-Bor Test Plugin 

In order to run the test first navigate to `cmd/cli` and build the binary with `go build`. From there navigate to `plugins/test-plugin/test`, confirm that `run_test.sh` is executable and run: `BOR=../../../cmd/cli/cli ./run-test.sh`.
 
There are two injections in `./core/rawdb/` not covered by this test: `ModifyAncients()`, `AppendAncient()`. These are covered by stand alone standard go tests which can all be run by navigating to the respective directories and running: `go test -v -run TestPlugethInjections`. 

There are three injections in `./core/` not covered by this test: `NewSideBlock()`, `Reorg()`, `BlockProcessingError()`. These injections are not covered by any test and will need to be tested manually by commenting the hooks and trying to build the project. 

Note: depending on where the script is deployed Geth may complain that the path to the `.ipc` file is too long. Renaming of the directories or moving the project may be necessary.