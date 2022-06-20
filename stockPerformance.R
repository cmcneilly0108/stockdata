# TODOs
# Create spreadsheet - 3 tabs
# Auto download files - etrade
# Improve parsing of pfm - remove blank rows at bottom
# Bring back the notes column
# Add 12-month performance report

library(quantmod)
library(dplyr)
library(scales)
library(lubridate)
library(stringr)
library(xlsx)

h2 <- read.csv("pfmDownload.csv",stringsAsFactors=FALSE)
#h2 <- read.csv("pfmDownload.csv",skip=10)
h2 <- filter(h2,!(Symbol %in% c(' ','Cash','YAFFX','PENNX','HFCGX','FAGIX','CHTTX','RZV','DGS','FM','VNQ','FPX')))
h2$Ticker <- as.character(h2$Symbol)
holdings <- h2$Ticker
#holdings <- c(holdings,'EBAY','PYPL')


# update!  run bash shell, load xlsx, convert
# Update data files
fd <- file.info("ss.xlsx")$mtime
cd <- Sys.time()
dt <- difftime(cd, fd, units = "hours")
if (dt > 10) {
  system("./getSS.sh")
}

p <- read.xlsx('ss.xlsx',1,startRow=3)
p <- filter(p,!is.na(Ticker))

#p <- read.csv('shadow_stock_portfolio.csv',skip=2,stringsAsFactors=FALSE)
#p2 <- read.csv("results.csv")
p2 <- filter(p,!(Ticker %in% holdings))
#p2 <- filter(p2,Score > 1,!(Ticker %in% holdings))
p2$Ticker <- as.character(p2$Ticker)
prospects <- p2$Ticker
prospects <- c(prospects,'SPY')

candidates <- filter(h2,!(Ticker %in% prospects))

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
  sDate <- today() - months(4)
  getSymbols(tickers,src='yahoo',from=as.character(sDate))
  data <- mget(tickers)
  allPrices <- ldply(data,getPrices)
  colnames(allPrices) <- c('ticker','close','date')
  #print(head(allPrices))
  
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
print("Collecting Holdings")
todrop <- compareStocks(holdings)

print("Collecting Prospects")
picks <- compareStocks(prospects)

#filter(h2,!(Ticker %in% p$Ticker))

#notes <- filter(p2,str_length(X.2)>5) %>% mutate(Notes = X.2) %>% select(Ticker,Notes) %>% print

aaiiSell <- filter(todrop,!(ticker %in% p$Ticker)) %>% print

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
