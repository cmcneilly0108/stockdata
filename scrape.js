var page = require('webpage').create();
page.open('http://financials.morningstar.com/company-profile/c.action?t=MSFT&region=usa&culture=en-US', function () {
    console.log(page.content); //page source
    phantom.exit();
});
