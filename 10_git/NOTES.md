**TCP** to vcs.protohackers.com:30307

Prompt user with READY

# Basic

## Req

Found requests are:

1. HELP

```
OK usage: HELP|GET|PUT|LIST
```

2. GET

```
ERR usage: GET file [revision]
```

3. PUT

```
ERR usage: PUT file length newline data
```

4. LIST

```
ERR usage: LIST dir
```

## Response

Response are in the format of:

```
OK <msg>
```

or

```
ERR <msg>
```

# More

- dirs are obtained from the file.

- files are unix-like, must begin and be separated by `/`.<br/>
  e.g: `/meong/kucing`.

- increment version number if there is a diff with it's current revision, otherwise don't

- apparently the filesystem is not organized like unix. a directory can contain file too.
  i.e. you can perform these actions:

```
PUT /dir1/dir2/file
PUT /dir1/dir2
```

and store 2 files
