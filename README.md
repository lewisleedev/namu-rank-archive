# namu-rank-archive

namu-rank-archive archives Namu wiki's search ranks to SQLite DB. Its docker image also does this automatically every 30 min via cron.

## Usage
Using docker is the recommended way to use `namu-rank-archive`. use `make build` to build the image and run it.

```sh
make build
docker run namu-rank-archive
```
If you wish to change the frequency at which the docker image cron runs, change the crontab value inside `entrypoint.sh`. There is no env var to change it.

or you can use the cli instead, which you can build it yourself with `go build .`:

```sh
namu-rank-archive -d data/rank.db migrate  # Migration is just CREATE TABLE IF NOT EXISTS. Always back things up before you do anything!
namu-rank-archive -d data/rank.db archive 
```
There are two important env vars: `LOG_LEVEL` and `NAMU_RANK_DB`, which are both optional and pretty self-explanatory. You can also change the timezone with `TZ`, which defaults to `Asia/Seoul`.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
