# ElasticSearch command line tool

`es` is a small command line tool to interact with the
ElasticSearch search engine.

## Setup

You need to compile yourself currently:

	$ git clone https://github.com/olivere/es
	$ go get
	$ go build
	$ ./es help

## Commands

Lets list indices.

	$ es indices
	master
	marvel
	dummy
	$ es indices 'm.*'
	master
	marvel

And create a new index. This can be useful if you e.g. added templates.

	$ es create twitter
	$ es indices
	master
	marvel
	dummy
	twitter
	$ es create twitter
	Error: IndexAlreadyExistsException[[twitter] Already exists] (400)
	$ es create -f twitter

Now, lets delete indices.

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

Let's review the mapping of an index.

	$ es mapping dummy
	{
		"dummy" : {
		}
	}
	$ es mapping nonexistent
	Error: IndexMissingException[[nonexistent] missing] (404)

## Credits

Thanks a lot for the great folks working hard on
[ElasticSearch](http://www.elasticsearch.org/) and
[Google Go](http://golang.org/).

Also a big thanks to [Blake Mizerany](https://github.com/bmizerany) 
and the [fast heroku client](https://github.com/bmizerany/hk)
for inspiration on how to create a Go-based command line tool.
