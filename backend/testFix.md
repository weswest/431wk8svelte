Data engineers and scientists who use Python and R utilize data frames heavily for fundamental data processing tasks. The data frame provides a structure for organizing data of various types in the same object. For example, an R data frame can include numeric, character, Boolean, and list objects. 

Some might say that R provides the quintessential data frame structure, while Python Pandas' Data Frame offers a cumbersome, less user-friendly version of the R data frame.

Review documentation for data frames in R, Python Pandas, and Go.  Compare and contrast the capabilities and features of data frames across these three languages for data science.

Reflect on your experience with data frames. 

What has been your experience in working with data frames in R and Python?

What do you think of the Go alternatives for data frames?

### Response

Comparing and contrasting data frames across languages:

* R.  Data frames in R are the most intuitive, due to them being in-built into the base R package.  One of the big advantages with R's data frame manipulation is that it's so advanced and baked that it is an intuitive and readable exercise to see how a data frame is being manipulated
* Python.  Since this is the home of the slightly more complex data science work, you see a ton of development to both make data frame manipulation powerful but to also then link those data frames to other applications in a relatively near-native fashion.  This allows for a ton of downstream power from using data frames in python.  I've personally always found the interface for creating new fields to be a little complex, but oh well.
* Go.  As with everything we've discussed in Go, it's a) less popular for data scientists so the ecosystem is weaker; b) designed for performance which biases users away from packages that may generate unnecessary overhead.  The level of verboseness necessary to perform basic filtering or mutating tasks relative to Python is a shock to the system.

For what it's worth, I really struggled to wrap my head around data frames.  I came to Python from a decade-long background in data analysis being limited almost exclusively from a SQL-->Excel or Excel-->Excel pipeline and the data frame seemed overkill.  I really had to retrain my expectations and eventually started thinking about them like mini-SQL tables, which is probably still the wrong mental model but it got me there.

The Go alternatives for data frames seem like they're in their infancy and we haven't seen a "winner" yet.  Investment resources are distributed across a handful of packages which means none of them are as powerful as they could ultimately be.  As I mentioned in my description above, there is a fundamental expectation when using data frames that the ease of use of the tools to manipulate the data will be worth the tradeoff for any possible performance degradation.  Seeing as Go is all about being as optimized as possible, I would need to see some live examples of how fast a data frame is relative to performing calculations in other fashions before I was totally sold on using them in Go.