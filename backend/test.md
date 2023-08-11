Data engineers and scientists who use Python and R utilize data frames heavily for fundamental data processing tasks. The data frame provides a structure for organizing data of various types in the same object. For example, an R data frame can include numeric, character, Boolean, and list objects. 

Some might say that R provides the quintessential data frame structure, while Python Pandas' DataFrame offers a cumbersome, less user-friendly version of the R data frame.

The Go language has a [DataFrameLinks to an external site.](https://pkg.go.dev/github.com/go-gota/gota/dataframe) package, offering a Python-Pandas-like [DataFrameLinks to an external site.](https://pandas.pydata.org/docs/reference/api/pandas.DataFrame.html). There are, in fact, several Go data frame libraries: [DataFrames in Go with gota, qframe, and dataframe-go.Links to an external site.](https://mungingdata.com/go/dataframes-gota-qframe/)  

Review documentation for data frames in R, Python Pandas, and Go.  Compare and contrast the capabilities and features of data frames across these three languages for data science.

Reflect on your experience with data frames. 

What has been your experience in working with data frames in R and Python?

What do you think of the Go alternatives for data frames?

### Response

Comparing and contrasting dataframes across languages:

* R.  Dataframes in R are the most intuitive, due to them being in-built into the base R package.  One of the big advantages with R's dataframe manipulation is that it's so advanced and baked that it is an intuitive and readable exercise to see how a dataframe is being manipulated
* Python.  Since this is the home of the slightly more complex data science work, you see a ton of development to both make dataframe manipulation powerful but to also then link those dataframes to other applications in a relatively near-native fashion.  This allows for a ton of downstream power from using dataframes in python.  I've personally always found the interface for creating new fields to be a little complex, but c'est la vie
* Go.  As with everything we've discussed in Go, it's a) less popular for data scientists so the ecosystem is weaker; b) designed for performance which biases users away from packages that may generate unnecessary overhead.  The level of verboseness necessary to perform basic filtering or mutating tasks relative to Python is a shock to the system.

For what it's worth, I really struggled to wrap my head around dataframes.  I came to Python from a decade-long background in data analysis being limited almost exclusively from a SQL-->Excel or Excel-->Excel pipeline and the dataframe seemed overkill.  I really had to retrain my expectations and eventually started thinking about them like mini-SQL tables, which is probably still the wrong mental model but it got me there.

The Go alternatives for dataframes seem like they're in their infancy and we haven't seen a "winner" yet.  Investment resources are distributed across a handful of packages which means none of them are as powerful as they could ultimately be.  As I mentioned in my description above, there is a fundamental expectation when using dataframes that the ease of use of the tools to manipulate the data will be worth the tradeoff for any possible performance degradation.  Seeing as Go is all about being as optimized as possible, I would need to see some live examples of how fast a dataframe is relative to performing calculations in other fashions before I was totally sold on using them in Go.