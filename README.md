# ElasticSearch command line tool

`es` is a small command line tool to interact with the
ElasticSearch search engine.

Notice that you could do all of this with `curl` commands, as
seen on the [ElasticSearch API](http://www.elasticsearch.org/guide/reference/api/).
However, you probably save a few keystrokes with `es`.

## Setup

You need to compile yourself currently:

	$ git clone https://github.com/olivere/es
	$ go get
	$ go build
	$ ./es help

## Commands

Before we start, you can always lookup the ElasticSearch API via
the `api` command, like so:

	$ es api indices

The `api` command will open up a browser window with the API page
that matches the specified command. You can find the complete
[ElasticSearch API here](http://www.elasticsearch.org/guide/reference/api/).

Let's get started. First we list existing indices, either all of them
or via a regular expression.

	$ es indices
	master
	marvel
	dummy
	$ es indices 'm.*'
	master
	marvel

Let's create a new index. Use the -f flag to force creation, i.e. it will
not print an error if the index already exists (and won't touch the 
existing index).

	$ es create twitter
	$ es indices
	master
	marvel
	dummy
	twitter
	$ es create twitter
	Error: IndexAlreadyExistsException[[twitter] Already exists] (400)
	$ es create -f twitter

Print some useful information with the status API.

	$ es status twitter
	{ ... }

You can also get the status of all indices. Just leave out the index name.

	$ es status
	{ ... }

Let's remove some indices.

	$ es delete twitter
	$ es indices
	master
	marvel
	dummy
	$ es delete twitter
	Error: IndexMissingException[[twitter] missing] (404)
	$ es delete -f twitter
	$ es indices
	master
	marvel
	dummy

Let's review mappings, and even create mappings from the command line.

	$ es mapping dummy
	{
		"dummy" : {
		}
	}
	$ es mapping nonexistent
	Error: IndexMissingException[[nonexistent] missing] (404)
	$ es create twitter
	$ es put-mapping twitter tweet < tweet-mapping.json
	$ es mapping twitter
	{
	  "twitter" : {
	    "tweet" : {
	      "properties" : {
	        "message" : {
	          "type" : "string",
	          "store" : "yes"
	        }
	      }
	    }
	  }
	}

Templates, oh how I love thee... here's a sample session.

	$ es templates
	dummy-template
	$ es template dummy-template
	{
	}
	$ es create-template another-template < template.json
	$ es templates
	dummy-template
	another-template
	$ es template another-template
	{
	}
	$ es delete-template another-template
	$ es delete-templete -f another-template




## Credits

Thanks a lot for the great folks working hard on
[ElasticSearch](http://www.elasticsearch.org/) and
[Google Go](http://golang.org/).

Also a big thanks to [Blake Mizerany](https://github.com/bmizerany) 
and the [fast heroku client](https://github.com/bmizerany/hk)
for inspiration on how to create a Go-based command line tool.
