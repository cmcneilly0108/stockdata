library(quantmod)
library(plyr)
library(ggplot2)
library(scales)
library(lubridate)
# Get historical data for each of my 10 stocks
holdings <- c('ALG','BWLD','DAR','JAZZ','MWIV','FLXS','PCCC','REGI','REX','VOXX','SPY')
#prospects <- c('SPY','KBALB','ISH','FLXS','SMP','WLFC','DCO')
#t2 = c('VOXX')

h2 <- read.csv("pfmDownload.csv",skip=10)
h2 <- filter(h2,!(Symbol %in% c(' ','Cash','YAFFX','PENNX','HFCGX','FAGIX','CHTTX','RZV','DGS')))
h2$Ticker <- as.character(h2$Symbol)
holdings <- h2$Ticker


p2 <- read.csv("results.csv")
p2 <- filter(p2,Score > 2,!(Ticker %in% holdings))
p2$Ticker <- as.character(p2$Ticker)
prospects <- p2$Ticker
prospects <- c(prospects,'SPY')

getPrices <- function(x){
  prices <- as.data.frame(Cl(x))
  prices$date <- as.Date(rownames(prices))
  #print('here')
  #prices$ticker <- factor('SPY') # Need ticker symbol
  #print('h2')
  colnames(prices)[1] <- 'close'
  #print('h3')
  rownames(prices) <- NULL
  #print('Not here')
  prices$date <- ymd(prices$date)
  return(prices)
}

compareStocks <- function(tickers) {
  # calculate start 1 year from today
  #sDate <- today() - years(1)
  sDate <- today() - months(9)
  getSymbols(tickers,src='yahoo',from=as.character(sDate))
  data <- mget(tickers)
  allPrices <- ldply(data,getPrices)
  colnames(allPrices) <- c('ticker','close','date')
  print(head(allPrices))
  
  # Find Start price and append to rows
  iPrice <- allPrices[allPrices$date==min(allPrices$date),]
  iPrice <- subset(iPrice, select = -date)
  colnames(iPrice) <- c('ticker','startPrice')
  allPrices <- join(allPrices,iPrice)  
  
  # Scale the price by start price
  # allPrices$newPrice <- (allPrices$close/allPrices$startPrice) - 1
  allPrices$gain <- (allPrices$close/allPrices$startPrice) - 1
  allPrices$gainP <- percent((allPrices$close/allPrices$startPrice) - 1)
  
  # In order to pick which one to sell
  #last <- max(ymd(allPrices$date[!is.na(allPrices$date)]))
  #first <- min(ymd(allPrices$date[!is.na(allPrices$date)]))
  #m3 <- last - days(90)
  #m6 <- last - days(180)
  # Create df with tickers and purchase date
  # Calculate gain for past year
  # subset(allPrices,date==last | date==m6 | date==m3,select=c(date,ticker,gain))
  results <- subset(allPrices,date==max(allPrices$date),select=c(date,ticker,gain,gainP))
  # Sort results by gain descending
  print(results[order(-results$gain),c('ticker','gainP')])
}

todrop <- compareStocks(holdings)

picks <- compareStocks(prospects)
# Color code compared to SPY for same periods
# Go back to original financial package and see how much of this it can do

# Plot performance over past year
#ggplot(allPrices,aes(x=date,y=close)) + 
#  geom_line() + scale_y_log10() + facet_wrap(~ ticker)

# Try plotting all 10 on same chart
#ggplot(allPrices,aes(x=date,y=close,colour=ticker)) +
#  geom_line() + scale_y_log10()

# plot the log price - better but not easily comparable
#ggplot(allPrices,aes(x=date,y=newPrice,colour=ticker)) + geom_line() + scale_y_log10()
