# badgebadge ![Badge Badge](http://doyouevenbadge.com/github.com/stevvooe/badgebadge)

Do you even badge?

Make sure you have all the badges on your project to make it shine!

Simple place the following markdown in your project to report on whether or not you have sufficient status badging:

```markdown
![Badge Badge](http://doyouevenbadge.com/github.com/stevvooe/badgebadge)
```

Make sure to replace `github.com/stevvooe/badgebadge` with the name of your project.

## Build

You can build this Go server via docker:

```
docker build -t badgebadge .
```

The server exposes port 8080:

```
docker run -p 8080:8080 badgebadge
```


## Server options

Only 2 ENVIRONMENT variables are

* DEBUG - Shows more output in the server logs
* BADGE_COUNT - Minimum amount of badges for it to pass (default: 3)


## Status
[![Status Cake](https://www.statuscake.com/App/button/index.php?Track=0vgZznzlRJ&Days=1&Design=1)](https://www.statuscake.com)